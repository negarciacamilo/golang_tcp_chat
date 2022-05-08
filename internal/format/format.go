package format

import (
	"fmt"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func ToByte(format string, a ...any) []byte {
	return []byte(fmt.Sprintf(format, a...))
}

func CyanMessage(format string, a ...any) []byte {
	str := Cyan + format + Reset
	return []byte(fmt.Sprintf(str, a...))
}

func YellowMessage(format string, a ...any) []byte {
	str := Yellow + format + Reset
	return []byte(fmt.Sprintf(str, a...))
}

func RedMessage(format string, a ...any) []byte {
	str := Red + format + Reset
	return []byte(fmt.Sprintf(str, a...))
}

func GrayMessage(format string, a ...any) []byte {
	str := Gray + format + Reset
	return []byte(fmt.Sprintf(str, a...))
}

func PurpleMessage(format string, a ...any) []byte {
	str := Purple + format + Reset
	return []byte(fmt.Sprintf(str, a...))
}
