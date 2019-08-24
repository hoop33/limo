# Entrevista

> A question-and-answer command-line prompt for go

## Table of Contents

* [Introduction](#introduction)
* [Installation](#installation)
* [Usage](#usage)
  * [Questions](#questions)
  * [Interviews](#interviews)
  * [Color](#color)
* [Contributing](#contributing)
* [Credits](#credits)
* [License](#license)

## Introduction

Entrevista is a Go library for conducting a question-and-answer interview from the command line. It supports validation, different return types, hidden entries, and default responses.

## Installation

```sh
$ go get -u github.com/hoop33/entrevista
```

## Usage

To conduct an interview, you create an interview, add some questions, and run the interview, like this:

```golang
import (
	"fmt"

	"github.com/hoop33/entrevista"
)

interview := entrevista.NewInterview()
interview.Questions = []entrevista.Question{
	{
		Key:      "name",
		Text:     "Enter your name",
		Required: true,
	},
	*entrevista.NewBoolQuestion("like", "Do you like entrevista?"),
}
answers := interview.Run()
fmt.Println(answers["name"].(string))
fmt.Println(strconv.FormatBool(answers["like"].(bool)))
```

### Questions

Questions have the following fields:

* `Key` -- `string` -- the key to use in the answer map. Required.
* `Text` -- `string` -- the text of the question. Required.
* `AnswerKind` -- `reflect.Kind` -- The type of the expected answer.
* `Required` -- `bool` -- Whether an answer is required.
* `DefaultAnswer` -- `string` -- The default answer.
* `RegularExpression` -- `*regexp.Regexp` -- The regular expression used to validate the answer.
* `ValidateMinMax` -- `bool` -- whether to validate minimum and maximum.
* `Minimum` -- `int` -- the minimum (length for a string, value for a number)
* `Maximum` -- `int` -- the maximum (length for a string, value for a number)
* `RequiredMessage` -- `string` -- the error message to display if the answer is required and not supplied.
* `InvalidMessage` -- `string` -- the error message to display if the answer is invalid.
* `Hidden` -- `bool` -- whether to hide the answer.

The minimum required fields are `Key` and `Text`.

You can also create questions using a few convenience functions:

* `NewQuestion(key string, text string)`
* `NewStringQuestion(key string, text string, minLength int, maxLength int)`
* `NewBoolQuestion(key string, text string)`
* `NewNumberQuestion(key string, text string)`
* `NewNumberInRangeQuestion(key string, text string, min int, max int)`

### Interviews

You can customize a few things on interviews as well:

* `PromptTerminator` -- `string` -- the string to show at the end of questions.
* `RequiredMessage` -- `string` -- the message to display if an answer is required and not supplied. Note that a question's `RequiredMessage` will override this.
* `InvalidMessage` -- `string` -- the message to display if an answer is invalid. Note that a question's `InvalidMessage` will override this.
* `ShowOutput` -- `func(message string)` -- the function to use for normal output.
* `ShowError` -- `func(message string)` -- the function to use for error output.
* `Questions` -- `[]Question` -- the questions.
* `QuitOnInvalidAnswer` -- `bool` -- whether to stop the interview on an invalid answer.
* `ReadAnswer` -- `func(question *Question) (string, error)` -- the method to use to read an answer. Useful for testing or automation.

### Color

Although the default interview displays in black and white, you can easily change the output to use color using the color library of your choosing. Say, for example, you wanted to use (https://github.com/mgutz/ansi)[https://github.com/mgutz/ansi] to display questions in green and errors in bold red. You would do something like this:

```golang
import (
	"fmt"

	"github.com/hoop33/entrevista"
	"github.com/mgutz/ansi"
)

interview := entrevista.NewInterview()
interview.ShowOutput = func(message string) {
	fmt.Print(ansi.ColorFunc("green")(message))
}
interview.ShowError = func(message string) {
	fmt.Println(ansi.ColorFunc("red+b")(message))
}
```

## Contributing

Please note that this project is released with a [Contributor Code of Conduct](http://contributor-covenant.org/). By participating in this project you agree to abide by its terms. See [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md) file.

Contributions are welcome! Please open pull requests with code that passes all the checks.

## Credits

Entrevista uses the following open source libraries -- thank you!

* [gopass](https://github.com/howeyc/gopass)
* [Testify](https://github.com/stretchr/testify)

Apologies if I've inadvertently omitted any library or any contributor.

## License

Copyright &copy; 2015 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)
