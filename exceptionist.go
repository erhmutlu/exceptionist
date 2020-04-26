package exceptionist

import "fmt"

type TranslatedError struct {
	errorCode string
	errorMessage string
}

type ObservedError struct {
	Key string
	Args []interface{}
}

func (err ObservedError) Error() string{
	return fmt.Sprint(err.Key, "with args", err.Args)
}

func (err TranslatedError) Error() string{
	return err.errorMessage
}

func NewObservedError(key string, args []interface{}) ObservedError {
	return ObservedError{
		Key: key,
		Args: args,
	}
}

func newTranslatedError(errorCode string, errorMessage string) TranslatedError {
	return TranslatedError{
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}