package uz_Cyrl

import "github.com/theplant/cldr"

var Locale = &cldr.Locale{
	Locale: "uz_Cyrl",
	Number: cldr.Number{
		Symbols:    symbols,
		Formats:    formats,
		Currencies: currencies,
	},
	Calendar:   calendar,
	PluralRule: pluralRule,
}

func init() {
	cldr.RegisterLocale(Locale)
}
