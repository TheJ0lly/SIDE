package prettyfmt

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/osspecifics"
)

type COLOR string

const (
	NO_COLOR     COLOR = ""
	RED          COLOR = "\x1b[31m"
	GREEN        COLOR = "\x1b[32m"
	YELLOW       COLOR = "\x1b[33m"
	BLUE         COLOR = "\x1b[34m"
	WHITE        COLOR = "\x1b[37m"
	RESET_COLORS COLOR = "\x1b[0m"
)

// Print the text.
func Print(text string) {
	fmt.Print(text)
}

// Print the text with a selected color.
func CPrint(text string, color COLOR) {

	if color != NO_COLOR {
		fmt.Print(color)
	}

	fmt.Print(text)

	if color != NO_COLOR {
		fmt.Print(RESET_COLORS)
	}
}

// Exact same use as the fmt.Printf.
func Printf(format string, a ...any) {
	var args_index int = 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[args_index])
			args_index++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}
}

// Exact same use as the fmt.Printf, just add the color for the text as the second paramater, then continue as usual.
func CPrintf(format string, color COLOR, a ...any) {

	if color != NO_COLOR {
		fmt.Print(color)
	}

	var args_index int = 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[args_index])
			args_index++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}

	if color != NO_COLOR {
		fmt.Print(RESET_COLORS)
	}
}

// It uses CPrintf under the hood, and just prints the text in a RED color.
func ErrorF(format string, a ...any) {
	fmt.Print(RED)
	var args_index int = 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[args_index])
			args_index++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}
	fmt.Print(RESET_COLORS)
}

// Exact same use as the fmt.Sprintf, just add the color for the text as the second parametr, then continue as usual.
func CSprintf(format string, color COLOR, a ...any) string {
	var string_to_return string

	if color != NO_COLOR {
		string_to_return += string(color)
	}

	var args_index int = 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			string_to_return += fmt.Sprintf("%s", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			string_to_return += fmt.Sprintf("%d", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'X' {
			string_to_return += fmt.Sprintf("%X", a[args_index])
			args_index++
			i += 1
		} else {
			string_to_return += string(format[i])
		}
	}

	if color != NO_COLOR {
		string_to_return += string(RESET_COLORS)
	}
	return string_to_return
}

// Exact same use as the fmt.Sprintf.
func Sprintf(format string, a ...any) string {
	var string_to_return string
	var args_index int = 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			string_to_return += fmt.Sprintf("%s", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			string_to_return += fmt.Sprintf("%d", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'X' {
			string_to_return += fmt.Sprintf("%X", a[args_index])
			args_index++
			i += 1
		} else {
			string_to_return += string(format[i])
		}
	}

	return string_to_return
}

// Exact same use as the fmt.Scanln.
func Scanln(str *string) {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	*str += s.Text()
}

func WarningF(format string, a ...any) {
	fmt.Print(YELLOW)
	var args_index int = 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			fmt.Printf("%s", a[args_index])
			args_index++
			i += 1
		} else if format[i] == '%' && format[i+1] == 'd' {
			fmt.Printf("%d", a[args_index])
			args_index++
			i += 1
		} else {
			fmt.Printf("%c", format[i])
		}
	}
	fmt.Print(RESET_COLORS)
}

func SPathF(format ...string) string {
	var string_to_return string

	for i := 0; i < len(format); i++ {
		string_to_return += format[i]

		if i != len(format)-1 {
			string_to_return += osspecifics.PATH_SEP
		}
	}
	return string_to_return
}
