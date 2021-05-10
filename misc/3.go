// https://gist.github.com/hvoecking/10772475

// Traverses an arbitrary struct and translates all stings it encounters
//
// I haven't seen an example for reflection traversing an arbitrary struct, so
// I want to share this with you. If you encounter interface{} bugs or want to see
// another example please comment.
//
// The MIT License (MIT)
//
// Copyright (c) 2014 Heye Vöcking
//
// Permission is hereby granted, free of charge, to interface{} person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"reflect"
)

var dict = map[string]string{
	"Hello!":                 "Hallo!",
	"What's up?":             "Was geht?",
	"translate this":         "übersetze dies",
	"point here":             "zeige hier her",
	"translate this as well": "übersetze dies auch...",
	"and one more":           "und noch eins",
	"deep":                   "tief",
}

type I interface{}

type A struct {
	Greeting string
	Message  string
	Pi       float64
}

type B struct {
	Struct    A
	Ptr       *A
	Answer    int
	Map       map[string]string
	StructMap map[string]interface{}
	Slice     []string
}

func create() I {
	// The type C is actually hidden, but reflection allows us to look inside it
	type C struct {
		String string
	}

	return B{
		Struct: A{
			Greeting: "Hello!",
			Message:  "translate this",
			Pi:       3.14,
		},
		Ptr: &A{
			Greeting: "What's up?",
			Message:  "point here",
			Pi:       3.14,
		},
		Map: map[string]string{
			"Test": "translate this as well",
		},
		StructMap: map[string]interface{}{
			"C": C{
				String: "deep",
			},
		},
		Slice: []string{
			"and one more",
		},
		Answer: 42,
	}
}

func main() {
	// Some example test cases so you can mess around and see if it's working
	// To check if it's correct look at the output, no automated checking here

	// Test the simple cases
	{
		fmt.Println("Test with nil pointer to struct:")
		var original *B
		translated := translate(original)
		fmt.Println("original:  ", original)
		fmt.Println("translated:", translated)
		fmt.Println()
	}
	{
		fmt.Println("Test with nil pointer to interface:")
		var original *I
		translated := translate(original)
		fmt.Println("original:  ", original)
		fmt.Println("translated:", translated)
		fmt.Println()
	}
	{
		fmt.Println("Test with struct that has no elements:")
		type E struct {
		}
		var original E
		translated := translate(original)
		fmt.Println("original:  ", original)
		fmt.Println("translated:", translated)
		fmt.Println()
	}
	{
		fmt.Println("Test with empty struct:")
		var original B
		translated := translate(original)
		fmt.Println("original:  ", original, "->", original.Ptr)
		fmt.Println("translated:", translated, "->", translated.(B).Ptr)
		fmt.Println()
	}

	// Imagine we have no influence on the value returned by create()
	created := create()
	{
		// Assume we know that `created` is of type B
		fmt.Println("Translating a struct:")
		original := created.(B)
		translated := translate(original)
		fmt.Println("original:  ", original, "->", original.Ptr)
		fmt.Println("translated:", translated, "->", translated.(B).Ptr)
		fmt.Println()
	}
	{
		// Assume we don't know created's type
		fmt.Println("Translating a struct wrapped in an interface:")
		original := created
		translated := translate(original)
		fmt.Println("original:  ", original, "->", original.(B).Ptr)
		fmt.Println("translated:", translated, "->", translated.(B).Ptr)
		fmt.Println()
	}
	{
		// Assume we don't know B's type and want to pass a pointer
		fmt.Println("Translating a pointer to a struct wrapped in an interface:")
		original := &created
		translated := translate(original)
		fmt.Println("original:  ", (*original), "->", (*original).(B).Ptr)
		fmt.Println("translated:", (*translated.(*I)), "->", (*translated.(*I)).(B).Ptr)
		fmt.Println()
	}
	{
		// Assume we have a struct that contains an interface of an unknown type
		fmt.Println("Translating a struct containing a pointer to a struct wrapped in an interface:")
		type D struct {
			Payload *I
		}
		original := D{
			Payload: &created,
		}
		translated := translate(original)
		fmt.Println("original:  ", original, "->", (*original.Payload), "->", (*original.Payload).(B).Ptr)
		fmt.Println("translated:", translated, "->", (*translated.(D).Payload), "->", (*(translated.(D).Payload)).(B).Ptr)
		fmt.Println()
	}
}

func translate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	translateRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		translatedString := dict[original.Interface().(string)]
		copy.SetString(translatedString)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}
