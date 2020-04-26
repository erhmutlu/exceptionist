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

type translation struct {
	errorCode string
	errorMessage string
}

var defaultTranslation = translation{
	errorCode: "0",
	errorMessage: "Default.",
}

var t1 = translation{
	errorCode: "1",
	errorMessage: "Geçersiz değer.",
}

var t2 = translation{
	errorCode: "2",
	errorMessage: "Hata.",
}

var translations = map[string]translation{
	"invalid.value": t1,
	"error": t2,
}

func New(key string, args []interface{}) ObservedError {
	return ObservedError{
		Key: key,
		Args: args,
	}
}

func new(t translation) TranslatedError {
	return TranslatedError{
		errorCode:    t.errorCode,
		errorMessage: t.errorMessage,
	}
}

func (err ObservedError) Translate() TranslatedError{
	if translation, ok := translations[err.Key]; ok {
		return new(translation)
	}

	return new(defaultTranslation)
}