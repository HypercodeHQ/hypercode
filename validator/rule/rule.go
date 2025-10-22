package rule

import (
	"fmt"
	"regexp"
	"strconv"
)

type Rule interface {
	Validate(field string, value string) (bool, string)
	WithMessage(message string) Rule
}

type requiredRule struct {
	customMessage string
}

func Required() Rule {
	return &requiredRule{}
}

func (r *requiredRule) Validate(field string, value string) (bool, string) {
	if value == "" {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s is required", field)
		}
		return false, msg
	}
	return true, ""
}

func (r *requiredRule) WithMessage(message string) Rule {
	r.customMessage = message
	return r
}

type emailRule struct {
	customMessage string
}

func Email() Rule {
	return &emailRule{}
}

func (r *emailRule) Validate(field string, value string) (bool, string) {
	if value == "" {
		return true, ""
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s must be a valid email address", field)
		}
		return false, msg
	}
	return true, ""
}

func (r *emailRule) WithMessage(message string) Rule {
	r.customMessage = message
	return r
}

type minRule struct {
	min           int
	customMessage string
}

func Min(min int) Rule {
	return &minRule{min: min}
}

func (r *minRule) Validate(field string, value string) (bool, string) {
	if value == "" {
		return true, ""
	}

	if len(value) < r.min {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s must be at least %d characters", field, r.min)
		}
		return false, msg
	}
	return true, ""
}

func (r *minRule) WithMessage(message string) Rule {
	r.customMessage = message
	return r
}

type maxRule struct {
	max           int
	customMessage string
}

func Max(max int) Rule {
	return &maxRule{max: max}
}

func (r *maxRule) Validate(field string, value string) (bool, string) {
	if value == "" {
		return true, ""
	}

	if len(value) > r.max {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s must be at most %d characters", field, r.max)
		}
		return false, msg
	}
	return true, ""
}

func (r *maxRule) WithMessage(message string) Rule {
	r.customMessage = message
	return r
}

type numericRule struct {
	customMessage string
}

func Numeric() Rule {
	return &numericRule{}
}

func (r *numericRule) Validate(field string, value string) (bool, string) {
	if value == "" {
		return true, ""
	}

	if _, err := strconv.ParseFloat(value, 64); err != nil {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s must be a number", field)
		}
		return false, msg
	}
	return true, ""
}

func (r *numericRule) WithMessage(message string) Rule {
	r.customMessage = message
	return r
}

type gteRule struct {
	min           int
	customMessage string
}

func Gte(min int) Rule {
	return &gteRule{min: min}
}

func (r *gteRule) Validate(field string, value string) (bool, string) {
	if value == "" {
		return true, ""
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s must be a number", field)
		}
		return false, msg
	}

	if num < r.min {
		msg := r.customMessage
		if msg == "" {
			msg = fmt.Sprintf("%s must be greater than or equal to %d", field, r.min)
		}
		return false, msg
	}
	return true, ""
}

func (r *gteRule) WithMessage(message string) Rule {
	r.customMessage = message
	return r
}
