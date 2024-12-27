package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func main() {
	// コードポイント文字列の配列。それぞれ、9, ¢, あ, 𠀋, 🚩(絵文字の旗)
	codepoints := []string{"39", "00A2", "3042", "2000B", "1F6A9"}

	// byte配列を、16進数の文字列に変換する
	bytesToStr := func(bytes []byte) string {
		var str string
		for _, b := range bytes {
			str += fmt.Sprintf("%02X ", b)
		}
		return strings.TrimSuffix(str, " ")
	}

	// 2byte配列を、16進数の文字列に変換する
	wordsToStr := func(words []uint16) string {
		var str string
		for _, w := range words {
			str += fmt.Sprintf("%04X ", w)
		}
		return strings.TrimSuffix(str, " ")
	}

	for _, code := range codepoints {
		char, _ := strconv.ParseUint(code, 16, 32)
		r := rune(char)

		// codepoint -> utf8へ
		bytes := make([]byte, 4)
		size := utf8.EncodeRune(bytes, r)

		// codepoint -> utf16へ
		words := utf16.Encode([]rune{r})

		// byte配列を、16進数文字列化して表示
		fmt.Println("char =", string(r), ", utf8 =", bytesToStr(bytes[:size]), ", utf16 =", wordsToStr(words))
	}
}
