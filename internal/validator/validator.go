package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var ValueMustNotBeEmpty = "Dieses Feld darf nicht leer sein"
var ValueInvalidEmail = "Das Format der eingegebenen Email wird nicht erkannt"

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, regex *regexp.Regexp) bool {
	return regex.MatchString(value)
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
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

func ValueMustBeLongerThan(n int) string {
	str := fmt.Sprintf("Der Wert muss mehr als %d Zeichen lang sein", n)
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
