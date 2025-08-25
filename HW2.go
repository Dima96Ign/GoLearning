package main

import (
	"fmt"
	"io"
)

// TODO: Реализуйте структуры:
// TODO: - Weapon: Name (string), Damage (int), Durability (int)
type Weapon struct {
	Name       string
	Damage     int
	Durability int
}

// TODO: - Armor: Name (string), Defense (int), Weight (float64)
type Armor struct {
	Name    string
	Defense int
	Weight  float64
}

// TODO: - Potion: Name (string), Effect (string), Charges (int)
type Potion struct {
	Name    string
	Effect  string
	Charges int
}

// TODO:
// TODO: Можете добавить свои структуры :)
// TODO:
// TODO: Для каждой структуры реализуйте:
// TODO: - Метод Use() string (описание использования, например "Используется <имя>", и изменение Durability или Charges и т.д.)
func (w *Weapon) Use() string {
	if w.Durability > 0 {
		w.Durability--
		return fmt.Sprintf("Используется %s. Остаток прочности: %d", w.Name, w.Durability)
	}
	return fmt.Sprintf("%s сломано и не может быть использовано", w.Name)
}

func (a *Armor) Use() string {
	return fmt.Sprintf("Вы носите %s. Защита: %d, Вес: %.2f", a.Name, a.Defense, a.Weight)
}

func (p *Potion) Use() string {
	if p.Charges > 0 {
		p.Charges--
		return fmt.Sprintf("Используется %s. Эффект: %s. Осталось зарядов: %d", p.Name, p.Effect, p.Charges)
	}
	return fmt.Sprintf("%s закончились заряды и не может быть использовано.", p.Name)
}

func (inv *Inventory) UseItem(index int) {
	if index < 0 || index >= len(inv.Items) {
		fmt.Println("Неверный индекс предмета")
		return
	}
	item := inv.Items[index]
	fmt.Println(item.Use())
}

// TODO: - Методы интерфейса Item

type Item interface {
	GetName() string
	GetWeight() float64
	Use() string
}

// Реализация методов интерфейса Item для Weapon
func (w Weapon) GetName() string {
	return w.Name
}

func (w Weapon) GetWeight() float64 {
	return 0.0
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

// TODO: Реализуйте функцию
func DescribeItem(i Item) string {
	if i == nil {
		return "Предмет отсутствует"
	}
	return fmt.Sprintf("%s (вес: %.2f)", i.GetName(), i.GetWeight())
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	// TODO: Верните новый слайс только с элементами, для которых predicate вернул true
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	// TODO: Примените transform к каждому элементу и верните слайс с результатами
	var result []R
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	// TODO: Найдите первый элемент, удовлетворяющий condition и верните элемент и true или zero value и false
	var zero T
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	return zero, false
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) GetWeapons() []*Weapon {
	// TODO: Используйте Filter для отбора Weapon, затем Map для преобразования Item -> *Weapon
	weapons := Filter(inv.Items, func(item Item) bool {
		_, ok := item.(*Weapon)
		return ok
	})
	// Используем Map для преобразования Item в *Weapon
	return Map(weapons, func(item Item) *Weapon {
		weapon, _ := item.(*Weapon)
		return weapon
	})
}

func (inv *Inventory) GetBrokenItems() []Item {
	return Filter(inv.Items, func(item Item) bool {
		switch v := item.(type) {
		case *Weapon:
			return v.Durability <= 0
		case *Potion:
			return v.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) GetItemNames() []string {
	return Map(inv.Items, func(item Item) string {
		return item.GetName()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(item Item) bool {
		return item.GetName() == name
	})
}

// TODO: Бонус: реализуйте интефейс Storable для Weapon и Armor:
// TODO: - Weapon: формат "Weapon|Name|Damage|Durability"
// TODO: - Armor: формат "Armor|Name|Defense|Weight"

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (inv *Inventory) Save(w io.Writer) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
}

func (inv *Inventory) Load(r io.Reader) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
}

func main() {
	//TODO: Создайте инвентарь и добавьте:
	//TODO: - Оружие: "Меч" (урон 10, прочность 5)
	sword := &Weapon{
		Name:       "Aerondight",
		Damage:     10,
		Durability: 5,
	}

	// TODO: - Броню: "Щит" (защита 5, вес 4.5)
	shield := &Armor{
		Name:    "Rivian Shield",
		Defense: 5,
		Weight:  4.5,
	}

	// TODO: - Зелье: "Лечебное" (+50 HP, 3 заряда)
	healingPotion := &Potion{
		Name:    "Healing Potion",
		Effect:  "Restores 50 HP",
		Charges: 3,
	}

	// TODO: - Оружие: "Сломанный лук" (урон 5, прочность 0)
	brokenBow := &Weapon{
		Name:       "Broken Bow",
		Damage:     5,
		Durability: 0,
	}
	// TODO: Реализуйте логику/вызовы:
	// TODO: 1. Use предмета с выводом в консоль
	inventory := Inventory{
		Items: []Item{sword, shield, healingPotion, brokenBow},
	}
	inventory.UseItem(0)
	inventory.UseItem(1)
	inventory.UseItem(2)
	inventory.UseItem(3)

	// TODO: 2. DescribeItem с предметом и с nil, так же с выводом в консоль
	fmt.Println(DescribeItem(sword))
	fmt.Println(DescribeItem(shield))
	fmt.Println(DescribeItem(healingPotion))
	fmt.Println(DescribeItem(brokenBow))
	fmt.Println(DescribeItem(nil))
	// TODO: 3. Вывести в консоль результат вызова GetWeapons (должны вернуться только меч и лук)
	weapons := inventory.GetWeapons()
	fmt.Println("Оружие в инвентаре:")
	for _, weapon := range weapons {
		fmt.Printf("- %s (Урон: %d, Прочность: %d)\n", weapon.Name, weapon.Damage, weapon.Durability)
	}
	// TODO: 4. Вывести в консоль результат вызова GetBrokenItems (должен вернуть сломанный лук)
	brokenItems := inventory.GetBrokenItems()
	fmt.Println("Сломанные предметы в инвентаре:")
	for _, item := range brokenItems {
		fmt.Printf("- %s\n", item.GetName())
	}
	// TODO: 5. Вывести в консоль результат вызова GetItemNames (все названия)
	itemNames := inventory.GetItemNames()
	fmt.Println("Названия всех предметов в инвентаре:")
	for _, name := range itemNames {
		fmt.Println("-", name)
	}
	// TODO: 6. Вывести в консоль результат вызова FindItemByName (поиск "Щит")
	if item, found := inventory.FindItemByName("Rivian Shield"); found {
		fmt.Printf("Найден предмет: %s\n", item.GetName())
	} else {
		fmt.Println("Предмет 'Rivian Shield' не найден")
	}

	// TODO: Бонус: сделайте сохранение инвентаря в файл и загрузку инвентаря из файла
}
