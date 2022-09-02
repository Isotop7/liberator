package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

var ValueMustNotBeEmpty = "Dieses Feld darf nicht leer sein"

type Validator struct {
	ValueErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.ValueErrors) == 0
}

func (v *Validator) AddValueError(key, message string) {
	if v.ValueErrors == nil {
		v.ValueErrors = make(map[string]string)
	}

	if _, exists := v.ValueErrors[key]; !exists {
		v.ValueErrors[key] = message
	}
}

func (v *Validator) CheckValue(ok bool, key, message string) {
	if !ok {
		v.AddValueError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func GreaterThan(value int, lbound int) bool {
	return value > lbound
}

func InBounds(value int, lbound int, ubound int) bool {
	return (value >= lbound && value <= ubound)
}

func ValueMustNotBeLongerThan(n int) string {
	str := fmt.Sprintf("Der Titel darf nicht mehr als %d Zeichen lang sein", n)
	return str
}

func ValueMustBeGreaterThan(n int) string {
	str := fmt.Sprintf("Der Wert muss größer als %d sein", n)
	return str
}

func ValueMustBeInRange(l int, u int) string {
	str := fmt.Sprintf("Der Wert muss größer oder gleich %d und gleichzeitig kleiner oder gleich %d sein", l, u)
	return str
}
