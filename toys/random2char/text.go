package random2char

import (
	"math/rand"
	"strings"
	"time"
)

const maximumCount = 140

func MakeText() string {
	rand.Seed(time.Now().Unix())
	OneChars := []string{
		"あ", "い", "う", "え", "お",
		"か", "き", "く", "け", "こ",
		"さ", "い", "す", "せ", "そ",
		"た", "ち", "つ", "て", "と",
		"な", "に", "ぬ", "ね", "の",
		"は", "ひ", "ふ", "へ", "ほ",
		"ま", "み", "む", "め", "も",
		"ら", "り", "る", "れ", "ろ",
		"が", "ぎ", "ぐ", "げ", "ご",
		"ざ", "じ", "ず", "ぜ", "ぞ",
		"だ", "ぢ", "づ", "で", "ど",
		"ば", "び", "ぶ", "べ", "ぼ",
		"ぱ", "ぴ", "ぷ", "ぺ", "ぽ",
		"や", "ゆ", "よ", "わ",
	}

	OneCharsSecond := []string{}
	OneCharsSecond = append(OneCharsSecond, OneChars...)
	OneCharsSecond = append(OneCharsSecond, "ん", "っ", "ー")

	TwoChars := []string{
		"きゃ", "きゅ", "きょ",
		"しゃ", "しゅ", "しょ",
		"ちゃ", "ちゅ", "ちょ",
		"にゃ", "にゅ", "にょ",
		"ひゃ", "ひゅ", "ひょ",
		"みゃ", "みゅ", "みょ",
		"ぎゃ", "ぎゅ", "ぎょ",
		"じゃ", "じゅ", "じょ",
		"びゃ", "びゅ", "びょ",
		"ぴゃ", "ぴゅ", "ぴょ",
	}

	TwoCharsSecond := []string{}
	TwoCharsSecond = append(TwoCharsSecond, TwoChars...)
	TwoCharsSecond = append(TwoCharsSecond, "ーん", "ーっ", "んっ")

	ThreeCharsSecond := []string{"ーんっ"}

	word := ""
	count := 0
	// select first character
	if idx := rand.Intn(len(OneChars) + len(TwoChars)); idx < len(OneChars) {
		word += OneChars[idx]
		count += 1
	} else {
		word += TwoChars[idx-len(OneChars)]
		count += 2
	}

	// select second character
	if idx := rand.Intn(len(OneCharsSecond) + len(TwoCharsSecond) + len(ThreeCharsSecond)); idx < len(OneCharsSecond) {
		word += OneCharsSecond[idx]
		count += 1
	} else if idx -= len(OneCharsSecond); idx < len(TwoCharsSecond) {
		word += TwoCharsSecond[idx]
		count += 2
	} else {
		word += ThreeCharsSecond[idx-len(TwoCharsSecond)]
		count += 3
	}

	looping := maximumCount / count
	word = strings.Repeat(word, looping)
	return word
}
