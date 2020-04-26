package exceptionist

type Language string

const (
	TR Language = "tr"
	EN Language = "en"
)

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

type bucket map[string]translation
type TranslationService struct {
	translations *map[Language]bucket
}

func NewTranslationService() TranslationService {
	translations := make(map[Language]bucket)
	return TranslationService{translations: &translations}
}

func (translationService TranslationService) Add(lang Language, filepath string) TranslationService {
	translations := *translationService.translations
	if _, ok := translations[lang]; !ok {
		bucket := parseFile(filepath)
		//var bucket = bucket{
		//	"invalid.value": t1,
		//	"error":         t2,
		//}
		translations[lang] = bucket
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
