package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

func Top10(text string) []string {
	if len(text) > 0 {
		re := regexp.MustCompile(`\S+`)
		words := re.FindAllString(text, -1)
		return WordCount10(words)
	}
	return nil
}

func WordCount10(words []string) []string {
	wordFiltered := make(map[string]int)
	for _, word := range words {
		wordFiltered[word]++
	}
	// Создаем структуру для определения слова и количества повторений
	type wordItem struct {
		word  string
		count int
	}
	// Создаем слайс структур
	wordItems := make([]wordItem, 0, len(wordFiltered))
	for text, count := range wordFiltered {
		wordItems = append(wordItems, wordItem{text, count})
	}

	// Сортируем слова
	sort.Slice(wordItems, func(i, j int) bool {
		if wordItems[i].count == wordItems[j].count {
			return wordItems[i].word < wordItems[j].word
		}
		return wordItems[i].count > wordItems[j].count
	})

	countItems := len(wordItems)
	if len(wordItems) >= 10 {
		countItems = 10
	}
	// Формируем массив строк
	result := make([]string, len(wordItems[:countItems]))
	for i, wordItem := range wordItems[:countItems] {
		result[i] = wordItem.word
	}
	return result
}
