package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	Values url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		url.Values(data),
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, filed := range fields {
		value := f.Values.Get(filed)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(filed, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, maxLen int) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > maxLen {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", maxLen))
	}
}

func (f *Form) MinLength(field string, minLen int) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < minLen {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", minLen))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Values.Get(field)

	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MatchesPattern(field string) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		f.Errors.Add(field, "This field is invalid")
	}
}

// func valid(email string) bool {
// 	_, err := mail.ParseAddress(email)
// 	return err == nil
// }

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
