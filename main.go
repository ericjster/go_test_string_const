package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// This project tries to find a way to determine if a string is a const, in an
// efficient way.
//
// The simple idea is to create a long string which is a concatenation of all
// strings. Then we create the individual strings from it.
//
// In practice we would want to use `go generate` to create the separate
// strings, very similar to `stringer`. See this blog post by Rob Pike for
// details on `go generate` and `stringer`: https://blog.golang.org/generate
//
// Then to determine if a string is a const, we check if the memory for the
// string is inside a larger string. We expect that can be done in a few machine
// instructions.
//
// Example data "Pill" from Rob Pike, in  https://blog.golang.org/generate
//
// Run:
//   go generate
//   go build && ./try_string_const
//   go tool objdump ./try_string_const >x.x

// This directive tells "go generate" what to do.
//go:generate stringer -type=Pill

// Pill is a sample type for stringer to generate.
type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen
)

// These are our constants, in case we want a variable name to hold one of the strings.
// Strings are just an array, with a pointer and length.
var constAll string
var constPlacebo string
var constAspirin string
var constIbuprofen string
var constParacetamol string
var constAcetaminophen string

// Example function to create the constant strings.
// Later this could be done by some tool similar to stringer.
// Each string is created with:
//    constPlacebo = _Pill_name[_Pill_index[Placebo]:_Pill_index[Placebo+1]]
func createConstStrings() {
	constAll = _Pill_name
	constPlacebo = Placebo.String()
	constAspirin = Aspirin.String()
	constIbuprofen = Ibuprofen.String()
	constParacetamol = Paracetamol.String()
	constAcetaminophen = Acetaminophen.String()
}

// IsStringInside returns true if one string is contained inside another, in
// terms of memory layout. This is the inner loop and generates only a few
// machine instructions.
//go:noinline
func IsStringInside(outer, inner string) bool {
	outerStringHeader := (*reflect.StringHeader)(unsafe.Pointer(&outer))
	innerStringHeader := (*reflect.StringHeader)(unsafe.Pointer(&inner))

	// fmt.Printf("  outer: %d, len:%d\n", outerStringHeader.Data, outerStringHeader.Len)
	// fmt.Printf("  inner: %d, len:%d\n", innerStringHeader.Data, innerStringHeader.Len)

	return outerStringHeader.Data <= innerStringHeader.Data &&
		innerStringHeader.Data+uintptr(innerStringHeader.Len) <= outerStringHeader.Data+uintptr(outerStringHeader.Len)
}

func main() {

	createConstStrings()

	strlist := [...]string{
		"Penn",
		"Teller",
		constAll,
		constPlacebo,
		constAspirin,
		constIbuprofen,
		constParacetamol,
		constAcetaminophen,
	}
	for _, s := range strlist {
		fmt.Printf("%-15s: %v\n", s, IsStringInside(constAll, s))
	}
}
