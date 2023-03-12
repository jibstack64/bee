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
	numMatcher    = GenerateRegexp("^-?\\d+(?:\\.\\d+)?$")
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
// This also handles addition, multiplication, division, etc.
func Objectify(s string) (*Object, *Error) {
	ob := FetchObject(s)
	if ob == nil {
		if c := stringMatcher.Find([]byte(s)); c != nil {
			ob = NewObject(DISPOSABLE, strings.ReplaceAll(string(c[1:len(s)-1]), "\\n", "\n"))
		} else if numMatcher.Find([]byte(s)) != nil {
			c, _ := strconv.ParseFloat(s, 64)
			ob = NewObject(DISPOSABLE, c)
		} else if s == "true" {
			ob = NewObject(DISPOSABLE, true)
		} else if s == "false" {
			ob = NewObject(DISPOSABLE, false)
		} else if s == "nil" {
			ob = NewObject(DISPOSABLE, nil)
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
					switch string(line[c+1]) {
					case string(FUNC_IN):
						allTokens[l] = append(allTokens[l], NewToken(string(SET)+string(FUNC_IN), nil))
						c++
					case "+":
						allTokens[l] = append(allTokens[l], NewToken(string(SET)+"+", nil))
						c++
					case "-":
						allTokens[l] = append(allTokens[l], NewToken(string(SET)+"-", nil))
						c++
					case "*":
						allTokens[l] = append(allTokens[l], NewToken(string(SET)+"*", nil))
						c++
					case "/":
						allTokens[l] = append(allTokens[l], NewToken(string(SET)+"/", nil))
						c++
					default:
						allTokens[l] = append(allTokens[l], NewToken(string(SET), nil))
					}
				}
			} else if line[c] == FUNC_B[0] || line[c] == FUNC_B[1] {
				allTokens[l] = append(allTokens[l], NewToken(string(line[c]), nil))
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

// Converts v to a string.
func StringConvert(v interface{}) (string, error) {
	switch c := v.(type) {
	case float64:
		return strconv.FormatFloat(c, 'g', -1, 64), nil
	case bool:
		return strconv.FormatBool(c), nil
	case nil:
		return "nil", nil
	default:
		return c.(string), nil
	}
}

// Converts v to a float64.
func Float64Convert(v interface{}) (float64, error) {
	switch c := v.(type) {
	case string:
		return strconv.ParseFloat(c, 64)
	case bool:
		if c {
			return 1, nil
		} else {
			return 0, nil
		}
	case nil:
		return 0, nil
	default:
		return c.(float64), nil
	}
}

// Converts v to a boolean.
func BooleanConvert(v interface{}) (bool, error) {
	switch c := v.(type) {
	case float64:
		if c > 1 {
			return true, nil
		} else {
			return false, nil
		}
	case string:
		if c == "true" {
			return true, nil
		} else {
			return false, nil
		}
	case nil:
		return false, nil
	default:
		return c.(bool), nil
	}
}

// Parses and runs the tokens provided.
func Run(tokens []*Token) error {
	// parse objects, strings, etc.
	for t := len(tokens) - 1; t > -1; t-- {
		token := tokens[t]
		if token.Origin[0] == SET {
			var or string
			if len(token.Origin) > 1 {
				or = token.Origin
			} else {
				or = string(SET) + "#"
			}
			if or[1] == FUNC_IN {
				// function
				if tokens[t-1].Object != nil {
					if tokens[t-1].Object.Symbol == "" {
						return NotCallableError.Format(tokens[t-1].Origin)
					}
				}
				if o := FetchObject(tokens[t-1].Origin); o != nil {
					tokens[t-1].Object = o // set obj
					var r *Object
					var err error
					if tokens[t+1].Object.Result != nil {
						r, err = tokens[t-1].Object.Data.(func(ob *Object, v ...*Object) (*Object, error))(tokens[t-1].Object, tokens[t+1].Object.Result)
					} else {
						r, err = tokens[t-1].Object.Data.(func(ob *Object, v ...*Object) (*Object, error))(tokens[t-1].Object, tokens[t+1].Object)
					}
					if err != nil {
						return err
					} else {
						tokens[t-1].Object.Result = r
					}
				}
			} else {
				// set value
				if tokens[t-1].Object == nil {
					if o, err := Objectify(tokens[t-1].Origin); err != nil {
						for _, c1 := range INVALID_SYMBOLS {
							for _, c2 := range tokens[t-1].Origin {
								if c1 == c2 {
									return SyntaxError.Format(tokens[t-1].Origin)
								}
							}
						}
						tokens[t-1].Object = NewObject(tokens[t-1].Origin, nil)
					} else {
						tokens[t-1].Object = o
					}
				}
				// no no!
				if tokens[t-1].Origin[0] == BUILT_IN {
					return NotAssignableError.Format(tokens[t-1].Origin)
				}
				// func funcs!
				if tokens[t+1].Object.Result == nil && tokens[t+1].Origin[0] == '@' {
					c, err := tokens[t+1].Object.Data.(func(ob *Object, v ...*Object) (*Object, error))(tokens[t+1].Object)
					if err != nil {
						return err
					} else {
						tokens[t+1].Object.Result = c
					}
				}
				if or[1] == '#' {
					if tokens[t+1].Object.Result != nil {
						tokens[t-1].Object.Data = tokens[t+1].Object.Result.Data
					} else {
						tokens[t-1].Object.Data = tokens[t+1].Object.Data
					}
				} else {
					// STOP OBJECT BEING OVERWRITTEN!!!!
					ob := tokens[t-1].Object
					tokens[t-1].Object = NewObject(DISPOSABLE, ob.Data)
					// conversion to int/string, etc.
					var v2 interface{}
					if tokens[t+1].Object.Result != nil {
						v2 = tokens[t+1].Object.Result.Data
					} else {
						v2 = tokens[t+1].Object.Data
					}
					switch v1 := tokens[t-1].Object.Data.(type) {
					case float64:
						switch or[1] {
						case '+':
							tokens[t-1].Object.Data = v1 + v2.(float64)
						case '-':
							tokens[t-1].Object.Data = v1 - v2.(float64)
						case '*':
							tokens[t-1].Object.Data = v1 * v2.(float64)
						case '/':
							tokens[t-1].Object.Data = v1 / v2.(float64)
						}
					case string:
						switch or[1] {
						case '+':
							tokens[t-1].Object.Data = v1 + v2.(string)
						case '-':
							return NumericsError.Format(tokens[t-1].Origin)
						case '*':
							return NumericsError.Format(tokens[t-1].Origin)
						case '/':
							return NumericsError.Format(tokens[t-1].Origin)
						}
					case bool:
						return NumericsError.Format(tokens[t-1].Origin)
					case nil:
						return NumericsError.Format(tokens[t-1].Origin)
					}
				}
			}
		} else {
			// object
			if token.Object == nil {
				if o, err := Objectify(token.Origin); err != nil {
					return err.Format(token.Origin)
				} else {
					token.Object = o
				}
			}
		}
	}

	// clear garbage
	DisposeGarbage()

	return nil
}
