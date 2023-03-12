package main

import (
	"bufio"
	"fmt"
	"os"
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
func Converter[T any]() func(ob *Object, v ...*Object) (*Object, error) {
	return func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, 1); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		if r, ok := v[0].Data.(T); !ok {
			return nil, ConversionError.Format(ob.Symbol)
		} else {
			return NewObject("", r), nil
		}
	}
}

var (
	// Types
	STRING = NewBuiltIn("string", Converter[string]())
	BOOL   = NewBuiltIn("bool", Converter[bool]())
	NUM    = NewBuiltIn("num", Converter[float64]())
	NIL    = NewBuiltIn("nil", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 0, 0); e != nil {
			return nil, e.Format(ob.Symbol)
		}
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
		return NewObject("", v), nil
	})

	// I/O
	IN = NewBuiltIn("in", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 0, 1); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		if len(v) > 0 {
			fmt.Print(StringConvert(v[0].Data))
		}
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		} else {
			return NewObject("", input), nil
		}
	})
	OUT = NewBuiltIn("out", func(ob *Object, v ...*Object) (*Object, error) {
		if e := ArgCheck(len(v), 1, len(global)); e != nil {
			return nil, e.Format(ob.Symbol)
		}
		for _, o := range v {
			fmt.Print(StringConvert(o.Data))
		}
		return nil, nil
	})

	// Holder
	TMP = NewBuiltIn("tmp", nil)
)
