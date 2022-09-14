package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var ValueMustNotBeEmpty = "Dieses Feld darf nicht leer sein"
var ValueInvalidEmail = "Das Format der eingegebenen Email wird nicht erkannt"
var InvalidISBN = "Fehlerhafte ISBN-Nummer"

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

func IsValidISBN10(isbn string) bool {
	// Check if ISBN-10 is empty
	if isbn == "" {
		return true
	}

	var digits []int

	// Split string into digits
	for _, r := range isbn {
		digit, err := strconv.Atoi(string(r))
		if err == nil {
			digits = append(digits, digit)
		}
	}

	// Check if 10 digits were supplied
	if len(digits) == 10 {
		// Calculate checksum which should be zero
		sum := 0
		for i := 1; i <= len(digits); i++ {
			sum = sum + (i * digits[i-1])
		}
		checksum := sum % 11
		return checksum == 0
	} else {
		return false
	}
}

func IsValidISBN13(isbn string) bool {
	// Check if ISBN-13 is empty
	if isbn == "" {
		return true
	}

	var digits []int

	// Split string into digits
	for _, r := range isbn {
		digit, err := strconv.Atoi(string(r))
		if err == nil {
			digits = append(digits, digit)
		}
	}

	// Check if 13 digits were supplied
	if len(digits) == 13 {
		// Calculate checksum which should be zero
		checksum := (digits[0] + digits[2] + digits[4] + digits[6] + digits[8] + digits[10] + digits[12] + 3*(digits[1]+digits[3]+digits[5]+digits[7]+digits[9]+digits[11])) % 10
		return checksum == 0
	} else {
		return false
	}
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
