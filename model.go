package exceptionist

import (
	"bytes"
	"text/template"
)

//Language
type Language string

const (
	TR Language = "tr"
	EN Language = "en"
)

//Bucket
type bucket struct {
	rows             map[string]row
	messageTemplates *template.Template
}
type row struct {
	errorCode    int
	templateName string
}

var defaultTranslationRow = newRow(1000, "default")

func newBucket() bucket {
	return bucket{
		rows:             map[string]row{},
		messageTemplates: template.New("tmpl"),
	}
}

func (bucket bucket) addRow(key string, errorCode int, template string) {
	bucket.messageTemplates.New(key).Parse(template)
	bucket.rows[key] = newRow(errorCode, key)
}

func (bucket bucket) findRow(key string) row {
	if row, ok := bucket.rows[key]; ok {
		return row
	}

	return defaultTranslationRow
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
	dir    *string
	prefix *string
}

func NewConfig(dir *string, prefix *string) Config {
	if dir == nil {
		panic("messages dir is required field.")
	}

	if prefix == nil {
		var defaultPrefix = "messages"
		prefix = &defaultPrefix
	}

	return Config{
		dir:    dir,
		prefix: prefix,
	}
}
