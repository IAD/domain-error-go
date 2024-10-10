package log2

import (
	"fmt"
	"strings"
)

const (
	DetailedInternalServiceError        = "DetailedInternalServiceError"
	DetailedNumericInternalServiceError = "DetailedNumericInternalServiceError"
	InternalServiceError                = "InternalServiceError"
)

type DomainError struct {
	Text       string                `json:"text"`
	Code       string                `json:"code"`
	Attributes DomainErrorAttributes `json:"attributes"`
	Error      string                `json:"error"`
}

func (de DomainError) String() string {
	text := de.Text
	for key, value := range de.Attributes {
		text = strings.ReplaceAll(text, fmt.Sprintf("{{%s}}", key), value)
		text = strings.ReplaceAll(text, fmt.Sprintf("{{ %s }}", key), value)
	}
	if de.Error != "" {
		text += ": " + de.Error
	}

	return text
}

type DomainErrorAttributes map[string]string

func toErrStr(errs []error) string {
	str := ""

	for i, err := range errs {
		if err != nil {
			if i > 0 {
				str += ": "
			}
			str += err.Error()
		}
	}

	return str
}

// NewDetailedInternalServiceError generated from the code "DetailedInternalServiceError".
func NewDetailedInternalServiceError(reason string, errs ...error) *DomainError {
	if len(errs) > 0 && errs[0] == nil {
		return nil
	}

	return &DomainError{
		Text: "internal service error. Reason {{reason}}",
		Code: DetailedInternalServiceError,
		Attributes: DomainErrorAttributes{
			"reason": fmt.Sprintf("%s", reason),
		},
		Error: toErrStr(errs),
	}
}

// NewDetailedNumericInternalServiceError generated from the code "DetailedNumericInternalServiceError".
func NewDetailedNumericInternalServiceError(reason int64, errs ...error) *DomainError {
	if len(errs) > 0 && errs[0] == nil {
		return nil
	}

	return &DomainError{
		Text: "internal service error. Reason {{ reason }}",
		Code: DetailedNumericInternalServiceError,
		Attributes: DomainErrorAttributes{
			"reason": fmt.Sprintf("%d", reason),
		},
		Error: toErrStr(errs),
	}
}

// NewInternalServiceError generated from the code "InternalServiceError".
func NewInternalServiceError(errs ...error) *DomainError {
	if len(errs) > 0 && errs[0] == nil {
		return nil
	}

	return &DomainError{
		Text:       "internal service error",
		Code:       InternalServiceError,
		Attributes: DomainErrorAttributes{},
		Error:      toErrStr(errs),
	}
}

type DomainErrorDeclaration struct {
	ErrorCode  string
	Text       string
	Attributes []string
}

func GetDomainErrors() []DomainErrorDeclaration {
	return []DomainErrorDeclaration{
		{
			ErrorCode:  "DetailedInternalServiceError",
			Text:       "internal service error. Reason {{reason}}",
			Attributes: []string{"reason"},
		},
		{
			ErrorCode:  "DetailedNumericInternalServiceError",
			Text:       "internal service error. Reason {{ reason }}",
			Attributes: []string{"reason"},
		},
		{
			ErrorCode:  "InternalServiceError",
			Text:       "internal service error",
			Attributes: []string{},
		},
	}
}
