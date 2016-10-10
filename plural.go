package i18n

import "errors"

const (
	PLURAL_ZERO  = "zero"
	PLURAL_ONE   = "one"
	PLURAL_TWO	 = "two"
	PLURAL_FEW   = "few"
	PLURAL_MANY  = "many"
	PLURAL_OTHER = "other"
)

type PluralsSets map[string]func(*operands)string

func plural(lang string, count interface{}) (result string, err error) {
	ops, err := newOperands(count)
	if err != nil {
		return
	}

	if _, ok := instance.pluralsSets[lang]; !ok {
		err = errors.New("Unknown language")
		return
	}

	result = instance.pluralsSets[lang](ops)
	return
}

func setPluralSets(langs []string, handler func(number *operands)string) {
	for _, lang := range langs {
		instance.pluralsSets[lang] = handler
	}
}
