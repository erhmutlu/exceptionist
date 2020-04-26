package exceptionist

import "fmt"

type ObservedError struct {
	Key string
	Args []interface{}
}

func (c ObservedError) Error() string{
	return c.Key
}

func New(key string, args []interface{}) ObservedError {
	return ObservedError{
		Key: key,
		Args: args,
	}
}

func DoPrint(a string) {
	fmt.Print("exceptionist" + a)
}
