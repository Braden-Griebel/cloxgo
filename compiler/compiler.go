package compiler

import "fmt"

func Compile(source string) {
	scanner := initScanner(&source)
	line := -1

	for {
		token := scanner.scanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Print("   | ")
		}
		var tokenRepr string
		if token.tokenType == TOKEN_EOF {
			tokenRepr = "EOF"
		} else {
			tokenRepr = string(scanner.code[token.start : token.start+token.length])
		}

		fmt.Printf("%s '%s'\n", token.tokenType, tokenRepr)

		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}
