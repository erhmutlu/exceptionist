package exceptionist

const (
	requiredTag = "required"
	defaultTag  = "default"
)

const defaultPrefix string = "messages"

type Config struct {
	Dir    *string
	Prefix *string
}

type bucket map[string]translation

func (config Config) ensure() Config{
	if config.Dir == nil {
		panic("messages dir is required field.")
	}

	if config.Prefix == nil {
		var prefix = defaultPrefix
		config.Prefix = &prefix
	}

	return config
}
