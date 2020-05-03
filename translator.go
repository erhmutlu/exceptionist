package exceptionist

import (
	"fmt"
	"github.com/magiconair/properties"
	"strconv"
	"strings"
)

type TranslationService struct {
	config       Config
	translations *map[Language]bucket
}

func NewTranslationService(config Config, supportedLanguages ...Language) TranslationService {
	translations := make(map[Language]bucket)

	translationService := TranslationService{
		config:       config,
		translations: &translations,
	}
	for _, language := range supportedLanguages {
		translationService = translationService.addLanguageSupport(language)
	}

	return translationService
}

func (translationService TranslationService) Translate(err ObservedError, lang Language) TranslatedError {
	translations := *translationService.translations
	bucket := translations[lang]
	row := bucket.findRow(err.Key)
	errorMessage := bucket.formatToErrorMessage(row, err.Args)
	return newTranslatedError(row.errorCode, errorMessage)
}

func (translationService TranslationService) addLanguageSupport(lang Language) TranslationService {
	translations := *translationService.translations
	if _, ok := translations[lang]; !ok {
		filepath := *translationService.config.dir + "/messages_" + lang.symbol + ".properties"
		translations[lang] = readTranslations(filepath, lang)
	}
	return translationService
}

func readTranslations(filepath string, lang Language) bucket {
	props := properties.MustLoadFile(filepath, properties.UTF8)

	bucket := newBucket(lang)

	for _, key := range props.Keys() {
		val := props.MustGet(key)
		if semiColon := strings.Index(val, ";"); semiColon >= 0 {
			errorCode, err := strconv.Atoi(val[:semiColon])
			if err != nil {
				fmt.Println("invalid errorCode in the props file:", filepath)
				continue
			}

			errorMessageTemplate := val[semiColon+1:]
			bucket.addRow(key, errorCode, errorMessageTemplate)
		}
	}
	return bucket
}
