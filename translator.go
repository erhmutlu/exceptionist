package exceptionist

type translation struct {
	errorCode    int
	errorMessage string
}

var defaultTranslation = translation{
	errorCode:    10001,
	errorMessage: "Default.",
}

var t1 = translation{
	errorCode:    1,
	errorMessage: "Geçersiz değer.",
}

var t2 = translation{
	errorCode:    2,
	errorMessage: "Hata.",
}

type TranslationService struct {
	translations *map[string]translation
}

func InitializeTranslationService() TranslationService {
	var translations = map[string]translation{
		"invalid.value": t1,
		"error":         t2,
	}

	return TranslationService{translations: &translations}
}

func (translationService TranslationService) Translate(err ObservedError) TranslatedError {
	translations := *translationService.translations

	if translation, ok := translations[err.Key]; ok {
		return newTranslatedError(translation.errorCode, translation.errorMessage)
	}

	return newTranslatedError(defaultTranslation.errorCode, defaultTranslation.errorMessage)
}
