package validators

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
)

type (
	Validator struct {
		Validator  *validator.Validate
		Translator ut.Translator
	}

	ValidationError struct {
		ValidationErrors map[string]string
		InnerError       error
	}
)

func (v ValidationError) Error() string {
	if v.InnerError != nil {
		return v.InnerError.Error()
	}

	return "Validation errors"
}

func (v *Validator) Struct(data interface{}) (err error) {
	err = conform.Strings(data)

	if err != nil {
		return ValidationError{
			ValidationErrors: nil,
			InnerError:       err,
		}
	}

	if err = v.Validator.Struct(data); err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return ValidationError{
			ValidationErrors: nil,
			InnerError:       err,
		}
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		return ValidationError{
			ValidationErrors: err.Translate(v.Translator),
			InnerError:       nil,
		}
	}

	return
}

func Register(v *validator.Validate, trans ut.Translator) error {
	return AlphaNumericUnicodeSpaceTranslationRegister(v, trans)
}
