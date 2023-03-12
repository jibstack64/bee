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

	// Sections
	SECTION_GOTO = '^'
	SECTION_DEF  = '#'
	BLOCKING     = '!'

	// Commenting
	COMMENT_START = "<bu"
	COMMENT_END   = "zz>"

	// Disposable, result
	DISPOSABLE = string(BUILT_IN) + ":disposable"
	RESULT     = string(BUILT_IN) + ":result"

	// Symbols that cannot be used in object names
	INVALID_SYMBOLS = "{}[]-=+/\\,.<>:;()!\"£$%^&*|?~'@ "

	// Define, set, func_in
	SET, FUNC_IN = '<', ':'

	// Tuple, func
	TUPLE_B, FUNC_B = "()", "[]"
)

var (
	// Arguments provided through the command-line
	fileName       string
	codeString     string
	generateGlobal bool

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
		for r := 0; r < len(parsed); r++ {
			if err, pointer := Run(parsed[r]); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			} else if pointer == -1 {
				for _, g := range global {
					switch c := g.Data.(type) {
					case int:
						if c == -1 {
							g.Data = r
						}
					}
				}
			} else if pointer == -2 {
				continue
			} else {
				r = pointer
			}
		}
	}
	if generateGlobal {
		for id, object := range global {
			fmt.Printf("ID: %d.   SYMBOL: '%s'.   DATA: '%+s'\n", id, object.Symbol, object.Data)
		}
	}
}

func init() {
	// parse flags
	flag.StringVar(&fileName, "filename", "", "Specifies the input file.")
	flag.StringVar(&codeString, "code", "", "When used, runs the code provided.")
	flag.BoolVar(&generateGlobal, "final", false, "Generates a list of all final objects - useful for debugging.")
	flag.Parse()
}
