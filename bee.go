package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	// The symbol for built-ins
	BUILT_IN = '@'

	// Splits 'lines'
	SPLITTER = ';'

	// Symbols that cannot be used in object names
	INVALID_SYMBOLS = "{}[]-=+/\\,.<>:;()!\"£$%^&*|?~'@ "

	// Define, set, func_in
	SET, FUNC_IN = '<', ':'

	// Tuple, func
	TUPLE_B, FUNC_B = "()", "[]"
)

var (
	// Arguments provided through the command-line
	fileName   string
	codeString string

	// All code
	program string
)

func main() {
	// Get the code
	if codeString != "" {
		program = codeString
	} else {
		if raw, err := os.ReadFile(fileName); err != nil {
			fmt.Printf("'%s' does not exist.\n", fileName)
			os.Exit(1)
		} else {
			program = string(raw)
		}
	}

	// Parse code
	if parsed, err := GenerateTokens(program); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		for _, p := range parsed {
			if err = Run(p); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}

func init() {
	// parse flags
	flag.StringVar(&fileName, "filename", "", "Specifies the input file.")
	flag.StringVar(&codeString, "code", "", "When used, runs the code provided.")
	flag.Parse()
}
