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
	COMMENT_START = "buzz"

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

func read(fn string) string {
	if raw, err := os.ReadFile(fn); err != nil {
		fmt.Printf("'%s' does not exist or missing access.\n", fn)
		os.Exit(1)
		return ""
	} else {
		return string(raw)
	}
}

func start(parsed [][]*Token) {
	for r := 0; r < len(parsed); r++ {
		if pointer, err := Run(parsed[r]); err != nil {
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
			// regenerate (otherwise you encounter some.. strange stuff)
			parsed, _ = GenerateTokens(program)
			r = pointer
		}
	}
}

func main() {
	// Get the code
	if codeString != "" {
		program = codeString
	} else {
		if program == "" {
			if fileName == "" {
				fmt.Println("No input provided.")
				os.Exit(1)
			} else {
				program = read(fileName)
			}
		}
	}

	// Parse code
	if parsed, err := GenerateTokens(program); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		start(parsed)
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
