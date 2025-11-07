package utils

import (
	mathRand "math/rand"
	"strings"
	"time"
)

func init() {
	mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
}

var MAP_NUMBER = []rune("1234567890")
var MAP_CHARACTERS = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

func MakePinCode(max int) string {
	characterSize := len(MAP_NUMBER)
	var sb strings.Builder

	for i := 0; i < max; i++ {
		ch := MAP_NUMBER[mathRand.Intn(characterSize)]
		sb.WriteRune(ch)
	}

	return sb.String()
}

func MakeRandString(max int) string {
	characterSize := len(MAP_CHARACTERS)
	var sb strings.Builder

	for i := 0; i < max; i++ {
		ch := MAP_CHARACTERS[mathRand.Intn(characterSize)]
		sb.WriteRune(ch)
	}

	return sb.String()
}
