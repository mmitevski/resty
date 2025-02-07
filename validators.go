package resty

import (
	"fmt"
)

var (
	validators map[string][]ValidationFunk = make(map[string][]ValidationFunk)
)

func toActionKey(f ActionFunc) string {
	return fmt.Sprintf("%#v", f)
}

func getValidators(action ActionFunc) []ValidationFunk {
	key := toActionKey(action)
	if validationHandlerFunk, ok := validators[key]; ok {
		return validationHandlerFunk
	}
	return nil
}

func AddValidator(action ActionFunc, validationFunk ValidationFunk) {
	key := toActionKey(action)
	if _, ok := validators[key]; !ok {
		validators[key] = make([]ValidationFunk, 0)
	}
	validators[key] = append(validators[key], validationFunk)
}
