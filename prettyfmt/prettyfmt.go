package prettyfmt

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/osspecifics"
)

type COLOR string

const (
	NoColor     COLOR = ""
	RED         COLOR = "\x1b[31m"
	GREEN       COLOR = "\x1b[32m"
	YELLOW      COLOR = "\x1b[33m"
	BLUE        COLOR = "\x1b[34m"
	WHITE       COLOR = "\x1b[37m"
	ResetColors COLOR = "\x1b[0m"
)

// Print the text.
func Print(text string) {
	fmt.Print(text)
}

// CPrint - Print the text with a selected color.
func CPrint(text string, color COLOR) {

	if color != NoColor {
		fmt.Print(color)
	}

	fmt.Print(text)

	if color != NoColor {
		fmt.Print(ResetColors)
	}
}

// Printf - Exact same use as the fmt.Printf.
func Printf(format string, a ...any) {
	var argsIndex = 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[argsIndex])
			argsIndex++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}
}

// CPrintf - Exact same use as the fmt.Printf, just add the color for the text as the second paramater, then continue as usual.
func CPrintf(format string, color COLOR, a ...any) {

	if color != NoColor {
		fmt.Print(color)
	}

	var argsIndex = 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[argsIndex])
			argsIndex++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}

	if color != NoColor {
		fmt.Print(ResetColors)
	}
}

// ErrorF - prints the text in a RED color.
func ErrorF(format string, a ...any) {
	fmt.Print(RED)
	var argsIndex = 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[argsIndex])
			argsIndex++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}
	fmt.Print(ResetColors)
}

// CSprintf - Exact same use as the fmt.Sprintf, just add the color for the text as the second parametr, then continue as usual.
func CSprintf(format string, color COLOR, a ...any) string {
	var stringToReturn string

	if color != NoColor {
		stringToReturn += string(color)
	}

	var argsIndex = 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			stringToReturn += fmt.Sprintf("%s", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			stringToReturn += fmt.Sprintf("%d", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'X' {
			stringToReturn += fmt.Sprintf("%X", a[argsIndex])
			argsIndex++
			i += 1
		} else {
			stringToReturn += string(format[i])
		}
	}

	if color != NoColor {
		stringToReturn += string(ResetColors)
	}
	return stringToReturn
}

// Sprintf - Exact same use as the fmt.Sprintf.
func Sprintf(format string, a ...any) string {
	var stringToReturn string
	var argsIndex = 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			stringToReturn += fmt.Sprintf("%s", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			stringToReturn += fmt.Sprintf("%d", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'X' {
			stringToReturn += fmt.Sprintf("%X", a[argsIndex])
			argsIndex++
			i += 1
		} else {
			stringToReturn += string(format[i])
		}
	}

	return stringToReturn
}

// Scanln - Exact same use as the fmt.Scanln.
func Scanln(str *string) {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	*str += s.Text()
}

func WarningF(format string, a ...any) {
	fmt.Print(YELLOW)
	var argsIndex = 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[argsIndex])
			argsIndex++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[argsIndex])
			argsIndex++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}
	fmt.Print(ResetColors)
}

func SPathF(format ...string) string {
	var stringToReturn string

	for i := 0; i < len(format); i++ {
		stringToReturn += format[i]

		if i != len(format)-1 {
			stringToReturn += osspecifics.PathSep
		}
	}
	return stringToReturn
}
