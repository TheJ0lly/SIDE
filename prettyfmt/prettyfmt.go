package prettyfmt

import (
	"bufio"
	"fmt"
	"os"
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

func Print(text string, color COLOR) {

	if color == NO_COLOR {
		fmt.Print(text)
	} else {
		fmt.Print(color)
		fmt.Print(text)
		fmt.Print(RESET_COLORS)
	}
}

func Printf(format string, color COLOR, a ...any) {

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

func ErrorF(format string, a ...any) {
	Printf(format, RED, a)
}

func CSprintf(format string, color COLOR, a ...any) string {
	var string_to_return string

	if color != NO_COLOR {
		string_to_return += string(color)
	}

	string_to_return += Sprintf(format, a)

	if color != NO_COLOR {
		string_to_return += string(RESET_COLORS)
	}
	return string_to_return
}

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

func Scanln(str *string) {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	*str += s.Text()
}
