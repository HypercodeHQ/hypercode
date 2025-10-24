package validator

import (
	"net/http"

	"github.com/hypercommithq/hypercommit/validator/rule"
)

type Validator interface {
	Validate(values map[string]string) map[string][]string
	ValidateForm(r *http.Request) map[string][]string
}

type FieldDef struct {
	name  string
	rules []rule.Rule
}

type validator struct {
	fields []FieldDef
}

func Field(name string, rules ...rule.Rule) FieldDef {
	return FieldDef{
		name:  name,
		rules: rules,
	}
}

func New(fields ...FieldDef) Validator {
	return &validator{
		fields: fields,
	}
}

func (v *validator) Validate(values map[string]string) map[string][]string {
	errors := make(map[string][]string)

	for _, field := range v.fields {
		value, exists := values[field.name]
		if !exists {
			continue
		}

		var fieldErrors []string
		for _, r := range field.rules {
			valid, errMsg := r.Validate(field.name, value)
			if !valid {
				fieldErrors = append(fieldErrors, errMsg)
			}
		}

		if len(fieldErrors) > 0 {
			errors[field.name] = fieldErrors
		}
	}

	return errors
}

func (v *validator) ValidateForm(r *http.Request) map[string][]string {
	if err := r.ParseForm(); err != nil {
		return map[string][]string{
			"_form": {err.Error()},
		}
	}

	values := make(map[string]string)
	for _, field := range v.fields {
		values[field.name] = r.FormValue(field.name)
	}

	return v.Validate(values)
}
