package exceptionist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewTranslator(t *testing.T) {
	//given
	config := prepareConfig()

	//when
	translator := NewTranslator(config)

	//then
	assert.NotNil(t, translator)
	assert.NotNil(t, translator.reader)
	assert.Equal(t, config, translator.config)
	assert.NotNil(t, translator.translations)
}

func TestErrorTranslator_AddLanguageSupport(t *testing.T) {
	//given
	translator := prepareTranslatorWithMockReader()

	//when
	translator.AddLanguageSupport(TR)

	//then
	translations := *translator.translations
	assert.Contains(t, translations, TR, "TR translations should be added")

	bucket := translations[TR]
	rows := bucket.rows
	assert.Contains(t, rows, "default", "Bucket should be initialized with the `default row")
	assert.Equal(t, 4, len(rows))

	assert.Equal(t, 1000, rows["key1"].errorCode)
	assert.Equal(t, "key1", rows["key1"].templateName)

	assert.Equal(t, 1001, rows["key2"].errorCode)
	assert.Equal(t, "key2", rows["key2"].templateName)

	assert.Equal(t, 1002, rows["key3"].errorCode)
	assert.Equal(t, "key3", rows["key3"].templateName)
}

func TestErrorTranslator_AddLanguageSupport_Callable_More_Then_One(t *testing.T) {
	//given
	translator := prepareTranslatorWithMockReader()
	translator.AddLanguageSupport(TR)

	//when
	translator.AddLanguageSupport(EN)

	//then
	translations := *translator.translations
	assert.Contains(t, translations, TR, "TR translations should be added")
	assert.Contains(t, translations, EN, "EN translations should be added")
}

func TestErrorTranslator_Translate(t *testing.T) {
	//given
	translator := prepareTranslatorWithMockReader()
	translator.AddLanguageSupport(TR)

	observedError := prepareObservedError()

	//when
	translatedError := translator.Translate(observedError, TR)

	//then
	assert.NotNil(t, translatedError)
	assert.Equal(t, 1001, translatedError.ErrorCode)
	assert.Equal(t, "errorTemplate2", translatedError.ErrorMessage)
}

func TestErrorTranslator_Translate_Should_Return_DefaultErrorResponse_When_LangIsTR_And_KeyIsNotFound(t *testing.T) {
	//given
	translator := prepareTranslatorWithMockReader()
	translator.AddLanguageSupport(TR)

	observedError := prepareObservedError()
	observedError.Key = "notExistingKey"

	//when
	translatedError := translator.Translate(observedError, TR)

	//then
	assert.NotNil(t, translatedError)
	assert.Equal(t, 100, translatedError.ErrorCode)
	assert.Equal(t, "İşleminizi şuanda gerçekleştiremiyoruz.", translatedError.ErrorMessage)
}

func TestErrorTranslator_Translate_Should_Return_DefaultErrorResponse_When_LangIsEN_And_KeyIsNotFound(t *testing.T) {
	//given
	translator := prepareTranslatorWithMockReader()
	translator.AddLanguageSupport(EN)

	observedError := prepareObservedError()
	observedError.Key = "notExistingKey"

	//when
	translatedError := translator.Translate(observedError, EN)

	//then
	assert.NotNil(t, translatedError)
	assert.Equal(t, 100, translatedError.ErrorCode)
	assert.Equal(t, "We are currently unable to complete your transaction.", translatedError.ErrorMessage)
}

func TestErrorTranslator_Translate_Should_Return_TRDefaultErrorResponse_When_GivenLangIsNotSupported(t *testing.T) {
	//given
	translator := prepareTranslatorWithMockReader()
	translator.AddLanguageSupport(EN)

	observedError := prepareObservedError()

	//when
	translatedError := translator.Translate(observedError, TR)

	//then
	assert.NotNil(t, translatedError)
	assert.Equal(t, 100, translatedError.ErrorCode)
	assert.Equal(t, "İşleminizi şuanda gerçekleştiremiyoruz.", translatedError.ErrorMessage)
}

func prepareConfig() Config {
	dir := "myDir"
	prefix := "myPrefix"
	config := Config{dir: &dir, prefix: &prefix}
	return config
}

func prepareObservedError() ObservedError {
	return ObservedError{Key: "key2"}
}

func prepareTranslatorWithMockReader() ErrorTranslator {
	config := prepareConfig()
	emptyTranslations := make(map[Language]bucket)
	return ErrorTranslator{
		config:       config,
		reader:       mockFileReader{},
		translations: &emptyTranslations,
	}
}

type mockFileReader struct {
}

func (reader mockFileReader) Read(filePath string) map[string]string {
	mockData := make(map[string]string)
	mockData["key1"] = "1000;errorTemplate1"
	mockData["key2"] = "1001;errorTemplate2"
	mockData["key3"] = "1002;errorTemplate3"
	return mockData
}
