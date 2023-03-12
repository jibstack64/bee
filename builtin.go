package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// A built-in object.
type BuiltIn struct {
	*Object

	Function func(ob *Object, v ...*Object) (*Object, error)
}

// Creates a built-in from the provided values.
func NewBuiltIn(symbol string, f func(ob *Object, v ...*Object) (*Object, error)) *BuiltIn {
	bi := &BuiltIn{
		Object:   NewObject(string(BUILT_IN)+symbol, nil),
		Function: f,
	}
	bi.Data = bi.Function
	return bi
}

// Returns the appropriate error for too many/little args.
func ArgCheck(v int, min int, max int) *Error {
	if v < min {
		return ValuesTooLittleError
	} else if v > max {
		return ValuesTooManyError
	} else {
		return nil
	}
}

// Forms a converter for the given type.
func Converter[T any](f func(v interface{}) (T, error)) func(ob *Object, v ...*Object) (*Object, error) {
	return func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, 1); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		a, err := f(v[0].Data)
		if err != nil {
			return nil, err
		}
		return NewObject(RESULT, a), nil
	}
}

var (
	// Types
	STRING = NewBuiltIn("string", Converter(StringConvert))
	BOOL   = NewBuiltIn("bool", Converter(BooleanConvert))
	NUM    = NewBuiltIn("num", Converter(Float64Convert))
	NIL    = NewBuiltIn("nil", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 0, 0); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		return nil, nil
	})

	// Time
	SLEEP = NewBuiltIn("sleep", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, 1); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		x, err := Float64Convert(v[0].Data)
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(x * float64(time.Second)))
		return nil, nil
	})

	// Destructors
	DEL = NewBuiltIn("del", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, len(global)); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		for _, o := range v {
			o.Delete()
		}
		return nil, nil
	})

	// Constructors
	LINK = NewBuiltIn("link", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, len(global)); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		return NewObject(RESULT, v), nil
	})

	// I/O
	IN = NewBuiltIn("in", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 0, 1); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		if len(v) > 0 {
			s, _ := StringConvert(v[0].Data)
			fmt.Print(s)
		}
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		} else {
			return NewObject(RESULT, strings.ReplaceAll(input, "\n", "")), nil
		}
	})
	OUT = NewBuiltIn("out", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, len(global)); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		for _, o := range v {
			s, _ := StringConvert(o.Data)
			fmt.Print(s)
		}
		return nil, nil
	})

	// Holder
	TMP = NewBuiltIn("tmp", nil)
)
