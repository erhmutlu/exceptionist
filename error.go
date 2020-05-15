package exceptionist

import "fmt"

type ObservedError interface {
	Error() string
	translate(bucket bucket) TranslatedError
}

type BusinessError struct {
	key         string
	args        []interface{}
	revealError bool
}

func NewError(key string, revealError bool, args []interface{}) BusinessError {
	return BusinessError{
		key:         key,
		args:        args,
		revealError: revealError,
	}
}

func (err BusinessError) Error() string {
	return fmt.Sprint(err.key, "with args", err.args)
}

func (err BusinessError) translate(bucket bucket) TranslatedError {
	row := bucket.findRow(err.key)
	if !err.revealError {
		return newTranslatedError(row.errorCode, bucket.lang.defaultErrorMessage, string(row.errorCode))
	}

	errorMessage := bucket.formatToErrorMessage(row, err.args)

	return newTranslatedError(row.errorCode, errorMessage, string(row.errorCode))
}

type InternalError struct {
	error       error
	revealError bool
}

func WrapInternalError(error error, revealError bool) InternalError {
	return InternalError{
		error:       error,
		revealError: revealError,
	}
}

func (err InternalError) Error() string {
	return fmt.Sprint("Error: ", err.error)
}

func (err InternalError) translate(bucket bucket) TranslatedError {
	lang := bucket.lang
	context := err.error.Error()
	if !err.revealError {
		return newTranslatedError(lang.defaultErrorCode, lang.defaultErrorMessage, context)
	}

	return newTranslatedError(lang.defaultErrorCode, context, context)
}

type TranslatedError struct {
	ErrorCode           int    `json:"errorCode"`
	ErrorMessage        string `json:"errorMessage"`
	InternalErrorDetail string `json:"-"`
}

func (err TranslatedError) Error() string {
	return err.ErrorMessage
}

func newTranslatedError(errorCode int, userFriendlyErrorMessage string, errorDetail string) TranslatedError {
	if errorDetail == "" {
		errorDetail = "Error occurred: " + string(errorCode)
	}

	return TranslatedError{
		ErrorCode:           errorCode,
		ErrorMessage:        userFriendlyErrorMessage,
		InternalErrorDetail: errorDetail,
	}
}
