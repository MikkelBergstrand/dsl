package color

import (
	"fmt"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Magenta = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"

func Println(color string, str ...any) {
	fmt.Print(color)
	fmt.Println(str...)
	fmt.Print(Reset)
}
