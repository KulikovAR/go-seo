package usecases

import "errors"

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func IsDomainError(err error) bool {
	var domainErr *DomainError
	return errors.As(err, &domainErr)
}

func GetDomainErrorCode(err error) string {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Code
	}
	return ""
}
