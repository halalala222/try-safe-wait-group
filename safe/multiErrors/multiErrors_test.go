package multiErrors

import (
	"errors"
	"testing"
)

var multiErrors MultiErrors

var testErrors = []error{
	errors.New("test1"),
	errors.New("test2"),
	errors.New("test3"),
	errors.New("test4"),
}

func TestMultiErrors_Error(t *testing.T) {
	multiErrors = append(multiErrors, testErrors...)
	t.Log(multiErrors.Error())
}

func TestMultiErrors_ErrorOrNil(t *testing.T) {
	if err := multiErrors.ErrorOrNil(); err != nil {
		t.Errorf("expect nil , but get : [%v]", err)
	}
	multiErrors = append(multiErrors, testErrors...)
	if err := multiErrors.ErrorOrNil(); err == nil {
		t.Error("expect not nil , but get nil")
	}
}

func TestMultiErrors_IsIn(t *testing.T) {
	var notInError = errors.New("not in error")
	if multiErrors.IsIn(notInError) {
		t.Errorf("expect [%v] is not in, but get in", notInError)
	}
	multiErrors = append(multiErrors, testErrors...)
	if multiErrors.IsIn(notInError) {
		t.Errorf("expect [%v] is not in, but get in", notInError)
	}
	if !multiErrors.IsIn(testErrors[1]) {
		t.Errorf("expect [%v] is in,but get not in", testErrors[1])
	}
}

func TestMultiErrors_Cap(t *testing.T) {
	var (
		multiError = Cap(10)
		newError   = errors.New("new error")
	)
	multiError = append(multiError, newError)
	t.Log(multiError)
}
