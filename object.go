package main

import (
	"strings"
)

var (
	global = Stack{}
)

type Stack map[int]*Object

// Contains an object's id, symbol, type and data.
type Object struct {
	Id     uint
	Symbol string
	Data   interface{}
	Result *Object
}

// Deletes the parent object.
func (ob *Object) Delete() {
	delete(global, int(ob.Id))
}

// Creates a new object with the given values.
// Generates and handles id's automatically.
// If `nil` is returned, something dastardly has gone wrong.
func NewObject(symbol string, data interface{}) *Object {
	ob := &Object{
		Id:     uint(len(global)) + 1,
		Symbol: symbol,
		Data:   data,
	}
	// Double check
	for _, ok := global[int(ob.Id)]; ok; ob.Id++ {
		continue
	}
	global[int(ob.Id)] = ob
	return ob
}

// Fetches an object by it's symbol.
// Returns `nil` if the object is not found.
func FetchObject(symbol string) *Object {
	symbol = strings.ReplaceAll(symbol, " ", "")
	for _, v := range global {
		if v.Symbol == symbol {
			return v
		}
	}
	return nil
}
