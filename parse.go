package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Holds a token (symbol + it's relative object/data).
type Token struct {
	Origin string
	Object *Object
}

func NewToken(origin string, object *Object) *Token {
	return &Token{
		Origin: origin,
		Object: object,
	}
}

// Regex
var (
	stringMatcher = GenerateRegexp("^\"[^\"]*\"$")
	decMatcher    = GenerateRegexp("^-?\\d+(?:\\.\\d+)?$")
	intMatcher    = GenerateRegexp("^-?\\d+$")
)

// Generates a regex, ignoring the error.
func GenerateRegexp(regex string) *regexp.Regexp {
	if r, err := regexp.Compile(regex); err != nil {
		panic(err)
	} else {
		return r
	}
}

// Converts the provided symbol to it's relative object.
// If one is not found, then it is parsed as a string or another valid type.
// If it cannot be parsed into a type, then an error is returned.
func Objectify(symbol string) (*Object, *Error) {
	ob := FetchObject(symbol) // trim whitespace
	if ob == nil {
		if stringMatcher.Find([]byte(symbol)) != nil {
			ob = NewObject("", symbol[1:]+symbol[:len(symbol)-2])
		} else if decMatcher.Find([]byte(symbol)) != nil {
			c, _ := strconv.ParseFloat(symbol, 64)
			ob = NewObject("", c)
		} else if intMatcher.Find([]byte(symbol)) != nil {
			c, _ := strconv.ParseInt(symbol, 10, 64)
			ob = NewObject("", c)
		} else {
			return nil, InvalidObjectError
		}
		return ob, nil
	} else {
		return ob, nil
	}
}

// Generates tokens from code.
func GenerateTokens(code string) ([][]*Token, error) {
	allTokens := [][]*Token{}
	for l, line := range strings.Split(code, string(SPLITTER)) {

		allTokens = append(allTokens, []*Token{})
		tmp := ""

		for c := 0; c < len(line); c++ {
			// set/define/func
			if string(line[c]) == string(SET) {
				allTokens[l] = append(allTokens[l], NewToken(strings.TrimSpace(tmp), nil))
				tmp = "" // reset
				if len(line)-1 == c {
					return nil, SyntaxError.Format(string(line[c]) + "' in '" + line)
				} else {
					if string(line[c+1]) == string(FUNC_IN) {
						allTokens[l] = append(allTokens[l], NewToken(string(SET)+string(FUNC_IN), nil))
						c++
					} else {
						allTokens[l] = append(allTokens[l], NewToken(string(SET), nil))
					}
				}
			} else {
				// add to tmp
				tmp += string(line[c])
			}
			if c == len(line)-1 {
				allTokens[l] = append(allTokens[l], NewToken(strings.TrimSpace(tmp), nil))
				tmp = "" // reset
			}
		}
	}

	// return tokens
	return allTokens, nil
}

// Parses and runs the tokens provided.
func Run(tokens []*Token) error {
	// parse objects, strings, etc.
	toRemove := []int{}
	for i := len(tokens) - 1; i > -1; i-- {
		token := tokens[i]
		if token.Origin == "" {
			toRemove = append(toRemove, i)
			continue
		}
		// attempt to get object
		if token.Origin != string(SET) && token.Origin != string(SET)+string(FUNC_IN) {
			if o, err := Objectify(token.Origin); err != nil {
				if i < len(tokens)-1 {
					if tokens[i+1].Origin == string(SET) {
						// ensure no invalid
						for _, c1 := range INVALID_SYMBOLS {
							for _, c2 := range token.Origin {
								if c1 == c2 {
									return InvalidObjectError.Format(token.Origin)
								}
							}
						}
						token.Object = NewObject(token.Origin, nil)
					}
				}
				return err.Format(token.Origin)
			} else {
				token.Object = o
			}
		}
	}

	// remove
	for _, i := range toRemove {
		tokens = append(tokens[:i-(len(tokens)-len(toRemove))], tokens[i-(len(tokens)-len(toRemove))+1:]...)
	}

	// run!
	for i, token := range tokens {
		// set
		if token.Origin[0] == SET {
			if i == 0 || i == len(tokens)-1 {
				return NoObjectError.Format(token.Origin)
			} else {
				if token.Origin == string(SET)+string(FUNC_IN) {
					c, err := tokens[i-1].Object.Data.(func(ob *Object, v ...*Object) (*Object, error))(tokens[i-1].Object, tokens[i+1].Object)
					if err != nil {
						return err
					} else {
						tokens[i-1].Object.Result = c
					}
				} else if token.Origin == string(SET) {
					if tokens[i+1].Object.Result != nil {
						tokens[i-1].Object.Data = tokens[i+1].Object.Result.Data
					} else {
						tokens[i-1].Object.Data = tokens[i+1].Object.Data
					}
				}
			}
		}
	}

	return nil
}
