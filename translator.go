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

	if translation, ok := bucket[err.Key]; ok {
		return newTranslatedError(translation.errorCode, translation.errorMessage)
	}

	return newTranslatedError(defaultTranslation.errorCode, defaultTranslation.errorMessage)
}

func (translationService TranslationService) addLanguageSupport(lang Language) TranslationService {
	translations := *translationService.translations
	if _, ok := translations[lang]; !ok {
		filepath := *translationService.config.dir + "/messages_" + string(lang) + ".properties"
		bucket := readTranslations(filepath)
		translations[lang] = bucket
	}
	return translationService
}

func readTranslations(filepath string) bucket {
	props := properties.MustLoadFile(filepath, properties.UTF8)

	var bucket bucket = map[string]translation{}

	for _, key := range props.Keys() {
		val := props.MustGet(key)
		if semiColon := strings.Index(val, ";"); semiColon >= 0 {
			errorCode, err := strconv.Atoi(val[:semiColon])
			if err != nil {
				fmt.Println("invalid errorCode in the props file:", filepath)
				continue
			}

			errorMessage := val[semiColon+1:]
			bucket[key] = translation{
				errorCode:    errorCode,
				errorMessage: errorMessage,
			}
		}
	}

	return bucket
}
