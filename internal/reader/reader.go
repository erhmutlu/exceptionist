package reader

import "github.com/magiconair/properties"

type Read func(filePath string) map[string]string

type TranslationFileReader interface {
	Read(filePath string) map[string]string
}

type PropertiesFileReader struct {
}

func NewPropertiesFileReader() PropertiesFileReader {
	return PropertiesFileReader{}
}

func (reader PropertiesFileReader) Read(filePath string) map[string]string {
	props := properties.MustLoadFile(filePath, properties.UTF8)
	result := make(map[string]string)

	for _, key := range props.Keys() {
		val := props.MustGet(key)
		result[key] = val
	}

	return result
}
