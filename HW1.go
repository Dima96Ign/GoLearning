package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Функция для задачи №1
func countWords(s string) map[string]int {
	words := strings.Fields(s)
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}
	return wordCounts
}

// Функция для задачи №2
func AreAnagrams(s1, s2 string) bool {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	s1 = strings.ReplaceAll(s1, " ", "")
	s2 = strings.ReplaceAll(s2, " ", "")

	if len(s1) != len(s2) {
		return false
	}

	runes1 := []rune(s1)
	runes2 := []rune(s2)

	sort.Slice(runes1, func(i, j int) bool { return runes1[i] < runes1[j] })
	sort.Slice(runes2, func(i, j int) bool { return runes2[i] < runes2[j] })

	return string(runes1) == string(runes2)
}

// Функция для задачи №3
func FirstUnique(s string) rune {
	charCounts := make(map[rune]int)
	for _, char := range s {
		charCounts[char]++
	}
	for _, char := range s {
		if charCounts[char] == 1 {
			return char
		}
	}
	return 0
}

// Функция для задачи №4
func RemoveDuplicates(nums []int) []int {
	numsCount := make(map[int]bool)
	var result []int
	for _, num := range nums {
		if !numsCount[num] {
			result = append(result, num)
			numsCount[num] = true
		}
	}
	return result
}

// Функция для задачи №5
func RemoveElement(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, errors.New("index out of range")
	}
	result := append(nums[:index], nums[index+1:]...)

	return result, nil
}

// Функция для задачи №6
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	var filteredRunes []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			filteredRunes = append(filteredRunes, r)
		}
	}
	left, right := 0, len(filteredRunes)-1
	for left < right {
		if filteredRunes[left] != filteredRunes[right] {
			return false
		}
		left++
		right--
	}
	return true
}

// Функция для задачи №7
func DrawChessBoard(size int) {
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if (row+col)%2 == 0 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func main() {

	//Задача №1
	text := "go is fun go Go"
	wordCount := countWords(text)
	for word, count := range wordCount {
		fmt.Printf("%s: %d\n", word, count)
	}

	//Задача №2
	fmt.Println(AreAnagrams("Lis ten", "Silent"))

	//Задача №3
	fmt.Printf("%c\n", FirstUnique("abcabcd"))

	//Задача №4
	fmt.Println(RemoveDuplicates([]int{1, 2, 3, 2, 1, 4, 5, 5, 6}))

	//Задача №5
	nums := []int{1, 2, 2, 3, 1}
	newNums, err := RemoveElement(nums, 5)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", newNums)
	}

	//Задача №6
	fmt.Println(IsPalindrome("А роза упала на лапу Азора"))

	//Задача №7
	DrawChessBoard(8)
}
