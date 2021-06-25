package validators

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
)

type Validator struct {
	Validator *validator.Validate
	Translator ut.Translator
}

func (v *Validator) Struct(data interface{}) (translations map[string]string, err error) {
	err = conform.Strings(data)

	if err != nil {
		return nil, err
	}

	if err = v.Validator.Struct(data); err == nil {
		return nil, nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return nil, err
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		translations = err.Translate(v.Translator)
		return translations,  nil
	}

	return
}

func Register(v *validator.Validate, trans ut.Translator) error {
	return AlphaNumericUnicodeSpaceTranslationRegister(v, trans)
}

