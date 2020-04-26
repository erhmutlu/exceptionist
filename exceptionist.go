package exceptionist

import "fmt"

type TranslatedError struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type ObservedError struct {
	Key  string
	Args []interface{}
}

func (err ObservedError) Error() string {
	return fmt.Sprint(err.Key, "with args", err.Args)
}

func (err TranslatedError) Error() string {
	return err.ErrorMessage
}

func NewObservedError(key string, args []interface{}) ObservedError {
	return ObservedError{
		Key:  key,
		Args: args,
	}
}

func newTranslatedError(errorCode int, errorMessage string) TranslatedError {
	return TranslatedError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}
