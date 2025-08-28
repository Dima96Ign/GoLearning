package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	brokenWeapon   = errors.New("оружие сломано")
	ErrPotionEmpty = errors.New("зелье пусто")
	MissedItem     = errors.New("предмет отсутствует")
)

type Weapon struct {
	Name       string
	Damage     int
	Durability int
}

type Armor struct {
	Name    string
	Defense int
	Weight  float64
}

type Potion struct {
	Name    string
	Effect  string
	Charges int
}

func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		w.Durability--
		return "", fmt.Errorf("%s: %w", w.Name, brokenWeapon)
	}
	w.Durability--
	return fmt.Sprintf("Атаковали %s (%d урона)", w.Name, w.Damage), nil
}

func (a *Armor) Use() (string, error) {
	if a.Defense <= 0 {
		return "", errors.New(a.Name + " Защита отсутствует!")
	}
	return fmt.Sprintf("Вы носите %s. Защита: %d, Вес: %.2f", a.Name, a.Defense, a.Weight), nil
}

func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", fmt.Errorf("%s: %w", p.Name, ErrPotionEmpty)
	}
	p.Charges--
	return fmt.Sprintf("Использовали %s (%s)", p.Name, p.Effect), nil
}

//- Методы интерфейса Item

type Item interface {
	GetName() string
	GetWeight() float64
	Use() (string, error)
}

// Реализация методов интерфейса Item для Weapon
func (w Weapon) GetName() string {
	return w.Name
}

func (w Weapon) GetWeight() float64 {
	return 0.0
}

func (w *Weapon) Serialize(wr io.Writer) error {
	_, err := fmt.Fprintf(wr, "Armor|%s|%d|%f", w.Name, w.Damage, w.Durability)
	return err
}

func (w *Weapon) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	parts := strings.Split(string(data), "|")
	if len(parts) != 4 || parts[0] != "Weapon" {
		return fmt.Errorf("неверный формат данных")
	}

	w.Name = parts[1]

	w.Damage, err = strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("ошибка преобразования урона: %w", err)
	}

	w.Durability, err = strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("ошибка преобразования прочности: %w", err)
	}

	return nil
}

func (a *Armor) Serialize(wr io.Writer) error {
	_, err := fmt.Fprintf(wr, "Armor|%s|%d|%f", a.Name, a.Defense, a.Weight)
	return err
}

// Метод Deserialize для структуры Armor
func (a *Armor) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	parts := strings.Split(string(data), "|")
	if len(parts) != 4 || parts[0] != "Armor" {
		return fmt.Errorf("неверный формат данных")
	}

	a.Name = parts[1]

	a.Defense, err = strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("ошибка преобразования защиты: %w", err)
	}

	a.Weight, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return fmt.Errorf("ошибка преобразования веса: %w", err)
	}

	return nil
}

// Реализация методов интерфейса Item для Armor
func (a Armor) GetName() string {
	return a.Name
}

func (a Armor) GetWeight() float64 {
	return a.Weight
}

// Реализация методов интерфейса Item для Potion
func (p Potion) GetName() string {
	return p.Name
}

func (p Potion) GetWeight() float64 {
	return 0.0
}

// Реализация функции
func DescribeItem(i Item) (string, error) {
	if i == nil {
		return "", MissedItem
	}
	return fmt.Sprintf("%s (вес: %.2f)", i.GetName(), i.GetWeight()), nil
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	var result []R
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) error {
	if item == nil {
		return errors.New("невозможно добавить отсутствующий предмет")
	}
	inv.Items = append(inv.Items, item)
	return nil
}

func (inv *Inventory) GetItemNames() []string {
	return Map(inv.Items, func(item Item) string {
		return item.GetName()
	})
}

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (inv *Inventory) Save(w io.Writer) {
	for _, item := range inv.Items {
		if storable, ok := item.(Storable); ok {
			storable.Serialize(w)

			_, _ = fmt.Fprintln(w)
		}
	}
}

func (inv *Inventory) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Weapon") {
			var w Weapon
			r := strings.NewReader(line)
			if err := w.Deserialize(r); err != nil {
				return fmt.Errorf("ошибка десериализации оружия: %w", err)
			}
			if err := inv.AddItem(&w); err != nil {
				return fmt.Errorf("ошибка добавления оружия: %w", err)
			}
		} else if strings.HasPrefix(line, "Armor") {
			var a Armor
			r := strings.NewReader(line)
			if err := a.Deserialize(r); err != nil {
				return fmt.Errorf("ошибка десериализации брони: %w", err)
			}
			if err := inv.AddItem(&a); err != nil {
				return fmt.Errorf("ошибка добавления брони: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка сканирования: %w", err)
	}

	return nil
}

func SafeUse(item Item) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("паника: %v", r)
		}
	}()

	if item.GetName() == "Ящик Пандоры" {
		panic("открытие Ящика Пандоры!")
	}

	return item.Use()
}

func main() {
	inv := Inventory{}

	sword := &Weapon{Name: "Меч", Damage: 10, Durability: 5}
	healthPotion := &Potion{Name: "Лечебное", Effect: "+50 HP", Charges: 3}
	pandoraBox := &Weapon{Name: "Ящик Пандоры", Damage: 100, Durability: 1}

	// Добавление предметов в инвентарь
	if err := inv.AddItem(sword); err != nil {
		fmt.Println("Ошибка:", err)
	}
	if err := inv.AddItem(healthPotion); err != nil {
		fmt.Println("Ошибка:", err)
	}
	if err := inv.AddItem(pandoraBox); err != nil {
		fmt.Println("Ошибка:", err)
	}
	if err := inv.AddItem(nil); err != nil {
		fmt.Println("Ошибка:", err)
	}

	// Описание предметов
	if desc, err := DescribeItem(sword); err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(desc)
	}
	if desc, err := DescribeItem(nil); err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(desc)
	}

	// Использование предмета до потери прочности
	fmt.Println("Используем меч до поломки:")
	for {
		result, err := SafeUse(sword)
		if err != nil {
			fmt.Println("Ошибка:", err)
			break
		}
		fmt.Println(result)
	}

	// Сохранение в файл
	fmt.Println("Сохраняем в файл")
	file, err := os.OpenFile("homework_solved.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Ломаем файл
	fmt.Println("Ломаем файл")
	fmt.Fprintf(file, "Weapon||")

	// Загрузка из файла
	fmt.Println("Загружаем из файла")
	file, err = os.Open("homework_solved.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	inv = Inventory{}
	if err := inv.Load(file); err != nil {
		fmt.Println("Ошибка при загрузке:", err)
	}

	// Получение имен предметов
	names := inv.GetItemNames()
	fmt.Println("Имена предметов:", names)

	for _, item := range inv.Items {
		if desc, err := DescribeItem(item); err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			fmt.Println("-", desc)
		}
	}

	// Обработка паники для "Ящика Пандоры"
	fmt.Println("Используем Ящик Пандоры:")
	if result, err := SafeUse(pandoraBox); err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(result)
	}
}
