# g-exceptionist

Error translation tool. Translates errors into errorCode and errorMessage according to given error templates.

## Installation
To install this package, you need to setup your Go workspace. The simplest way to install the library is to run:

```$ go get require github.com/erhmutlu/g-exceptionist@{version}```

## Configuration

#### Template Files
```properties
# cat .../{dir}/{prefix}_en.properties

error.key1=10001;Error message template without arg.
error.key2=10002;Error message template with arg {{index . 0}}, {{index . 1}}, {{index . 2}}.

# cat .../{dir}/{prefix}_tr.properties

error.key1=10001;Argümansız hata mesajı taslağı.
error.key2=10002;Argümanlı hata mesajı taslağı {{index . 0}}, {{index . 1}}, {{index . 2}}.
```

- Type of the template files must be ```.properties```.
- Every supported language should have its own message template file.
- Each row represents an error with specific code and message.
- Each row has 3 parts.
```
errorKey=errorCode;errorMessageTemplate
```
- errorKey must be given to Translator while translation for indicating a specific error.
- Right side of the equation represents an error data, and it must be separated into two parts with `;` delimiter.
- errorMessageTemplate may contain dynamic parts for args:```{{index . 0}}``` represents for 0th index of args slice.

#### Translator
```golang
import ex "github.com/erhmutlu/g-exceptionist"

var dir = "/path/to/templates/dir"
var prefix = "template-files-prefix"

config := ex.NewConfig(&dir, $prefix)
var errorTranslator = ex.NewTranslator(config)
errorTranslator.AddLanguageSupport(ex.TR)
errorTranslator.AddLanguageSupport(ex.EN)
```

- At least one language should be supported, otherwise Translator will always translate to default ```Turkish Error Translation```.

#### Usage
```golang
import ex "github.com/erhmutlu/g-exceptionist"

projectError := YourCustomProjectError{
			ObservedError: ex.NewObservedError("error.key2", []interface{}{"0", "1", "2}),
			//your other fields here
		}

translatedError := errorTranslator.Translate(projectError.ObservedError, ex.EN)
fmt.Println(translatedError.ErrorCode) //10002
fmt.Println(translatedError.ErrorMessage) //Error message template with arg 0, 1, 2
```

- If given language is not supported, Translator will produce default ```Turkish Error Translation```.
- If given errorKey is not found in supported language errors, Translator will produce default ```Error Translation``` for corresponding language.


#### Default Errors
- en
```json
{
  "errorCode": 100,
  "errorMessage": "We are currently unable to complete your transaction."
}
```

- tr
```json
{
  "errorCode": 100,
  "errorMessage": "İşleminizi şuanda gerçekleştiremiyoruz."
}
```