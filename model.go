package exceptionist

type Language string

const (
	TR Language = "tr"
	EN Language = "en"
)

type bucket map[string]translation
type translation struct {
	errorCode    int
	errorMessage string
}

var defaultTranslation = translation{
	errorCode:    10001,
	errorMessage: "Default.",
}

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
