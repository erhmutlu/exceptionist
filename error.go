package exceptionist

import (
	"fmt"
)

type ObservedError struct {
	Key                 string
	Args                []interface{}
	InternalErrorDetail *string
	RevealError         bool
}

type TranslatedError struct {
	ErrorCode                int    `json:"errorCode"`
	ErrorMessage             string `json:"-"`
	UserFriendlyErrorMessage string `json:"errorMessage"`
	InternalErrorDetail      string `json:"-"`
}

func (err ObservedError) Error() string {
	return fmt.Sprint(err.Key, "with args", err.Args)
}

func (err TranslatedError) Error() string {
	return err.ErrorMessage
}

func NewObservedError(key string, revealError bool, args []interface{}) ObservedError {
	return ObservedError{
		Key:         key,
		Args:        args,
		RevealError: revealError,
	}
}

func NewObservedErrorWithErrorDetail(key string, errorDetail string, revealError bool, args []interface{}) ObservedError {
	return ObservedError{
		Key:                 key,
		Args:                args,
		InternalErrorDetail: &errorDetail,
		RevealError:         revealError,
	}
}

func newTranslatedError(errorCode int, errorMessage string, userFriendlyErrorMessage string, errorDetail *string) TranslatedError {
	var internalErrorDetail string
	if errorDetail == nil {
		internalErrorDetail = "Error occurred: " + string(errorCode)
	} else {
		internalErrorDetail = *errorDetail
	}

	return TranslatedError{
		ErrorCode:                errorCode,
		ErrorMessage:             errorMessage,
		UserFriendlyErrorMessage: userFriendlyErrorMessage,
		InternalErrorDetail:      internalErrorDetail,
	}
}
