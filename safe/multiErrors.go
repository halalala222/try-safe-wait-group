package safe

import (
	"errors"
	"strings"
)

type MultiErrors []error

var _ error = &MultiErrors{}

func (m MultiErrors) Error() string {
	var errorStringBuilder strings.Builder

	for _, err := range m {
		errorStringBuilder.WriteString(err.Error())
	}

	return errorStringBuilder.String()
}

func (m MultiErrors) ErrorOrNil() error {
	if len(m) == 0 {
		return nil
	}
	return m
}

func (m MultiErrors) IsIn(targetErr error) bool {
	for _, err := range m {
		if errors.Is(targetErr, err) {
			return true
		}
	}
	return false
}

func (m MultiErrors) AsIn(targetError error) bool {
	for _, err := range m {
		if errors.As(err, &targetError) {
			return true
		}
	}
	return false
}

func (m MultiErrors) MultiErrorsIs(multiErrors MultiErrors) bool {
	for _, err := range m {
		if !multiErrors.IsIn(err) {
			return false
		}
	}
	return true
}

func (m MultiErrors) MultiErrorsAs(multiErrors MultiErrors) bool {
	for _, err := range m {
		if !multiErrors.AsIn(err) {
			return false
		}
	}
	return true
}
