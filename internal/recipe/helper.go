package recipe

//go:generate pigeon -o parser.go recipe.peg

import (
	"log"
	"reflect"
	"strings"
)

type pair struct {
	key   string
	value Buildable
}

func asString(val any) string {
	builder := strings.Builder{}
	for _, chars := range val.([]any) {
		builder.Write(chars.([]byte))
	}
	return builder.String()
}

func makeString(val any) *RecipeString {
	builder := strings.Builder{}
	result := make([]Buildable, 0)
	for _, content := range val.([]any) {
		switch element := content.(type) {
		case []byte:
			builder.Write(element)
		case Buildable:
			if builder.Len() > 0 {
				result = append(result, &RecipeStringLiteral{builder.String()}, element)
				builder.Reset()
			} else {
				result = append(result, element)
			}
		default:
			log.Panicf("unexpected element: %v\n", reflect.TypeOf(element))
		}
	}
	if builder.Len() > 0 {
		result = append(result, &RecipeStringLiteral{builder.String()})
	}
	return &RecipeString{result}
}

// Combine head and tail into a single slice
func headTail[T any](head any, tail []any) []T {
	result := make([]T, 0, 1+len(tail))
	result = append(result, head.(T))
	for _, t := range tail {
		result = append(result, t.(T))
	}
	return result
}

func toAnySlice[T any](input []any) []T {
	result := make([]T, len(input))
	for i, e := range input {
		if e != nil {
			result[i] = e.(T)
		}
	}
	return result
}

func collectPairs(pairs []pair) map[string]Buildable {
	result := make(map[string]Buildable)
	for _, keyvalue := range pairs {
		result[keyvalue.key] = keyvalue.value
	}
	return result
}
