package exceptionist

import (
	"bytes"
	"fmt"
	"text/template"
)

//Language
type Language struct {
	symbol              string
	defaultErrorCode    int
	defaultErrorMessage string
}

var (
	TR = Language{
		symbol:              "tr",
		defaultErrorCode:    100,
		defaultErrorMessage: "İşleminizi şuanda gerçekleştiremiyoruz.",
	}
	EN = Language{
		symbol:              "en",
		defaultErrorCode:    100,
		defaultErrorMessage: "We are currently unable to complete your transaction.",
	}
)

func (lang Language) toDefaultTranslatedError() TranslatedError {
	return newTranslatedError(lang.defaultErrorCode, lang.defaultErrorMessage, lang.defaultErrorMessage, nil)
}

//Bucket
type bucket struct {
	rows             map[string]row
	messageTemplates *template.Template
}

type row struct {
	errorCode    int
	templateName string
}

func newBucket(lang Language) bucket {
	tmpl := template.Must(template.New("default").Parse(lang.defaultErrorMessage))
	rows := map[string]row{
		"default": newRow(lang.defaultErrorCode, "default"),
	}
	return bucket{
		rows:             rows,
		messageTemplates: tmpl,
	}
}

func (bucket bucket) addRow(key string, errorCode int, errorMessageTemplate string) {
	if _, ok := bucket.rows[key]; ok {
		fmt.Println("Key:", key, "is duplicated on translation file")
	}

	bucket.rows[key] = newRow(errorCode, key)
	bucket.messageTemplates.New(key).Parse(errorMessageTemplate)
}

func (bucket bucket) findRow(key string) row {
	if row, ok := bucket.rows[key]; ok {
		return row
	}

	return bucket.rows["default"]
}

func (bucket bucket) formatToErrorMessage(row row, args []interface{}) string {
	buf := &bytes.Buffer{}
	if err := bucket.messageTemplates.ExecuteTemplate(buf, row.templateName, args); err != nil {
		panic(err)
	}

	return buf.String()
}

func newRow(errorCode int, templateName string) row {
	return row{
		errorCode:    errorCode,
		templateName: templateName,
	}
}

//Config
type Config struct {
	dir    string
	prefix string
}

func NewConfig(dir string, prefix string) Config {
	if dir == "" {
		panic("messages dir is required field.")
	}

	if prefix == "" {
		var defaultPrefix = "messages"
		prefix = defaultPrefix
	}

	return Config{
		dir:    dir,
		prefix: prefix,
	}
}
