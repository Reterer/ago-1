package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	const sep = "\t* "
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errors)))

	for _, err := range e.errors {
		sb.WriteString(sep)
		sb.WriteString(err.Error())
	}

	sb.WriteRune('\n')

	return sb.String()
}

func Append(err error, errs ...error) *MultiError {
	notNilErrors := 0

	if err != nil {
		notNilErrors++
	}

	for _, err := range errs {
		if err != nil {
			notNilErrors++
		}
	}

	if notNilErrors == 0 {
		return nil
	}

	me, ok := err.(*MultiError)
	if !ok {
		me = &MultiError{
			errors: make([]error, 0, notNilErrors),
		}

		if err != nil {
			me.errors = append(me.errors, err)
		}
	} else {
		notNilErrors-- // Так как мы подсчитали me
		me.errors = slices.Grow(me.errors, notNilErrors)
	}

	for _, err := range errs {
		if err != nil {
			me.errors = append(me.errors, errs...)
		}
	}

	return me
}

func (e *MultiError) Unwrap() []error {
	return e.errors
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
