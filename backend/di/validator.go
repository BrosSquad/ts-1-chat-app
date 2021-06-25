package di

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/sr_Latn"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"

	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
)

func (c *container) GetValidator() validators.Validator {
	if c.validator == (validators.Validator{}) {
		locale := c.viper.GetString("locale")

		log.Trace().Str("locale", locale).Msg("Validation messages locale")

		v, translator := createValidatorAndTranslator(locale)

		log.Trace().Msg("Registering default translations")

		if err := en_translations.RegisterDefaultTranslations(v, translator); err != nil {
			log.Fatal().Err(err).Msg("Error while registering english translations")
		}

		log.Trace().Msg("Registering custom translations")
		if err := validators.Register(v, translator); err != nil {
			log.Fatal().Err(err).Msg("Error while registering custom validators")
		}

		c.validator = validators.Validator{
			Validator:  v,
			Translator: translator,
		}

		log.Trace().Msg("Validator created")
	}

	return c.validator
}

func createValidatorAndTranslator(l string) (*validator.Validate, ut.Translator) {
	var translations locales.Translator
	switch l {
	case "en":
		translations = en.New()
	case "sr":
		translations = sr_Latn.New()
	default:
		log.Fatal().
			Strs("supported", []string{"en"}).
			Str("locale", l).
			Msg("Translation locale is not supported")
	}
	uni := ut.New(translations, translations)

	log.Trace().Msg("Creating translator")
	translator, _ := uni.GetTranslator(l)

	log.Trace().Msg("Creating validator")
	v := validator.New()

	return v, translator
}
