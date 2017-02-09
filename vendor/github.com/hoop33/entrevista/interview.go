package entrevista

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
)

// Interview contains all the questions and configuration to conduct an interview
type Interview struct {
	// The string to show at the end of questions. Default is ":".
	PromptTerminator string
	// The error message to display if the answer is required and not supplied.
	RequiredMessage string
	// The error message to display if the answer is invalid.
	InvalidMessage string
	// The function to use for normal output
	ShowOutput func(message string)
	// The function to use for error output
	ShowError func(message string)
	// The questions in the interview.
	Questions []Question
	// Whether to quit on an invalid answer
	QuitOnInvalidAnswer bool
	// The method to read an answer. Used for testing.
	ReadAnswer func(question *Question) (string, error)
}

func showOutput(message string) {
	fmt.Print(message)
}

func showError(message string) {
	fmt.Println(message)
}

func (interview *Interview) displayPrompt(question *Question) {
	interview.ShowOutput(question.Text)
	if question.DefaultAnswer != "" {
		interview.ShowOutput(fmt.Sprintf(" (%s)", question.DefaultAnswer))
	}
	interview.ShowOutput(interview.PromptTerminator)
}

func isValid(value interface{}, text string, question *Question) bool {
	if question.AnswerKind == reflect.Bool {
		return true
	}
	if question.AnswerKind == reflect.String {
		length := len(text)
		if length < question.Minimum || (question.Maximum != 0 && length > question.Maximum) {
			return false
		}
		if question.RegularExpression == nil {
			return true
		}
		return question.RegularExpression.MatchString(text)
	}
	if question.AnswerKind == reflect.Int {
		num := value.(int)
		if num < question.Minimum || (question.Maximum != 0 && num > question.Maximum) {
			return false
		}
		return true
	}
	return false
}

func readAnswer(question *Question) (string, error) {
	if question.Hidden {
		password, err := gopass.GetPasswd()
		fmt.Println()
		return string(password), err
	}

	var answer string
	fmt.Scanln(&answer)
	return answer, nil
}

func convertAnswer(answer string, kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.String:
		return answer, nil
	case reflect.Bool:
		return strings.HasPrefix(strings.ToUpper(answer), "Y"), nil
	case reflect.Int:
		return strconv.Atoi(answer)
	default:
		return answer, fmt.Errorf("The answer type %v is not supported", kind)
	}
}

func answerOrDefault(answer string, defaultAnswer string) string {
	if answer == "" && defaultAnswer != "" {
		return defaultAnswer
	}
	return answer
}

func getErrorMessage(qMessage string, iMessage string) string {
	if qMessage != "" {
		return qMessage
	}
	return iMessage
}

func (interview *Interview) getAnswer(question *Question) (interface{}, error) {
	for {
		interview.displayPrompt(question)
		answer, err := interview.ReadAnswer(question)
		if err != nil {
			return answer, err
		}

		// If they left answer blank and there's a default, set to default
		answer = answerOrDefault(answer, question.DefaultAnswer)

		// If it's still blank and it's required, show an error
		if answer == "" && question.Required {
			interview.ShowError(getErrorMessage(question.RequiredMessage, interview.RequiredMessage))
		} else {
			// Convert the answer to the appropriate type
			converted, err := convertAnswer(answer, question.AnswerKind)
			if err != nil {
				return converted, err
			}

			if !isValid(converted, answer, question) {
				// If answer isn't valid, show an error
				interview.ShowError(getErrorMessage(question.InvalidMessage, interview.InvalidMessage))
			} else {
				// We have a valid answer; return it
				return converted, nil
			}
		}
		// Loop if configured to do so
		if interview.QuitOnInvalidAnswer {
			return answer, err
		}
	}
}

// NewInterview creates a new interview with sane defaults
func NewInterview() *Interview {
	return &Interview{
		PromptTerminator: ": ",
		RequiredMessage:  "You must provide an answer to this question.",
		InvalidMessage:   "Your answer is not valid.",
		ShowOutput:       showOutput,
		ShowError:        showError,
		ReadAnswer:       readAnswer,
	}
}

// Run conducts an interview
func (interview *Interview) Run() (map[string]interface{}, error) {
	answers := make(map[string]interface{}, len(interview.Questions))
	for index, question := range interview.Questions {
		// If they haven't set the answer type, set it to String
		if question.AnswerKind == reflect.Invalid {
			question.AnswerKind = reflect.String
		}

		// If they haven't set a key, return an error
		if question.Key == "" {
			return nil, fmt.Errorf("Question %d has no key", index)
		}

		// If they haven't set the text for a question, return an error
		if question.Text == "" {
			return nil, fmt.Errorf("Question %d has no text", index)
		}

		answer, err := interview.getAnswer(&question)
		if err == nil {
			answers[question.Key] = answer
		} else {
			return answers, err
		}
	}
	return answers, nil
}
