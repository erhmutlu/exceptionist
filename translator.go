package exceptionist

import (
	"fmt"
	"github.com/erhmutlu/g-exceptionist/internal/reader"
	"strconv"
	"strings"
)

type ErrorTranslator struct {
	config       Config
	reader       reader.TranslationFileReader
	translations *map[Language]bucket
}

func NewTranslator(config Config) ErrorTranslator {
	translations := make(map[Language]bucket)

	return ErrorTranslator{
		config:       config,
		reader:       reader.NewPropertiesFileReader(),
		translations: &translations,
	}
}

func (errorTranslator *ErrorTranslator) AddLanguageSupport(lang Language) {
	translations := *errorTranslator.translations
	if _, ok := translations[lang]; !ok {
		translations[lang] = errorTranslator.prepareTranslationBucket(lang)
	}
}

func (errorTranslator ErrorTranslator) Translate(err ObservedError, lang Language) TranslatedError {
	translations := *errorTranslator.translations
	if _, ok := translations[lang]; !ok {
		return TR.toDefaultTranslatedError()
	}

	bucket := translations[lang]
	row := bucket.findRow(err.Key)
	errorMessage := bucket.formatToErrorMessage(row, err.Args)

	if err.RevealError {
		return newTranslatedError(row.errorCode, errorMessage, errorMessage, err.InternalErrorDetail)
	} else {
		return newTranslatedError(row.errorCode, errorMessage, lang.defaultErrorMessage, err.InternalErrorDetail)
	}
}

func (errorTranslator *ErrorTranslator) prepareTranslationBucket(lang Language) bucket {
	bucket := newBucket(lang)

	filepath := *errorTranslator.config.dir + "/" + *errorTranslator.config.prefix + "_" + lang.symbol + ".properties"
	props := errorTranslator.reader.Read(filepath)
	for key, val := range props {
		if semiColon := strings.Index(val, ";"); semiColon >= 0 {
			errorCode, err := strconv.Atoi(val[:semiColon])
			if err != nil {
				fmt.Println("invalid errorCode in the props file for language:", lang.symbol)
				continue
			}

			errorMessageTemplate := val[semiColon+1:]
			bucket.addRow(key, errorCode, errorMessageTemplate)
		}
	}

	return bucket
}
