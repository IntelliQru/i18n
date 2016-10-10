package i18n

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
)

var instance *i18n

type (
	language map[string]interface{}

	i18n struct {
		Translations map[string]language
		pluralsSets PluralsSets
	}

	translation struct {
		Id          string
		Translation interface{}
	}

	TFuncHandler func(translationId string)
)

func Tfunc(lang string) func(translationId string, args ...interface{}) string {
	return func(translationId string, args ...interface{}) string {
		return translate(lang, translationId, args...)
	}
}

func LoadTranslationFile(filename string) (err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	data := []translation{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return
	}

	file := path.Base(filename)
	lang := strings.Split(file, ".")[0]

	if _, ok := instance.Translations[lang]; !ok {
		instance.Translations[lang] = language{}
	}

	for _, item := range data {
		instance.Translations[lang][item.Id] = item.Translation
	}

	return
}

func translate(lang, translationId string, args... interface{}) (result string) {

	var data interface{}
	var count interface{}

	if _, ok := instance.Translations[lang]; !ok {
		return translationId
	}

	if argc := len(args); argc > 0 {
		if isNumber(args[0]) {
			count = args[0]

			if argc > 1 {
				data = args[1]
			}

		} else {
			data = args[0]
		}
	}

	if data == nil {
		data = map[string]interface{}{}
	} else {
		data = toMap(data)
	}

	var ok bool
	var dataMap map[string]interface{}
	lang_phrase := ""

	if count != nil {
		p, err := plural("ru", count)
		if err != nil {
			return translationId
		}

		dataMap, ok = instance.Translations[lang][translationId].(map[string]interface{})
		if !ok {
			return translationId
		}

		lang_phrase, ok = dataMap[p].(string)
		if !ok {
			return translationId
		}
	} else {
		lang_phrase, ok = instance.Translations[lang][translationId].(string)
		if !ok {
			return translationId
		}
	}

	tmpl, err := template.New(translationId).Parse(lang_phrase)
	if err != nil {
		return translationId
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = tmpl.Execute(writer, data)
	if err != nil {
		return translationId
	}

	writer.Flush()

	result = b.String()

	return b.String()
}

func init() {
	instance = &i18n{
		Translations: make(map[string]language),
		pluralsSets: PluralsSets{},
	}

	setPluralSets([]string{"bm", "bo", "dz", "id", "ig", "ii", "in", "ja", "jbo", "jv", "jw", "kde", "kea", "km", "ko", "lkt", "lo", "ms", "my", "nqo", "root", "sah", "ses", "sg", "th", "to", "vi", "wo", "yo", "zh"}, func(ops *operands)string {
		return PLURAL_OTHER
	})
	setPluralSets([]string{"am", "as", "bn", "fa", "gu", "hi", "kn", "mr", "zu"}, func(ops *operands)string {
		// i = 0 or n = 1
		if intEqualsAny(ops.I, 0) ||
			ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ff", "fr", "hy", "kab"}, func(ops *operands)string {
		// i = 0,1
		if intEqualsAny(ops.I, 0, 1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ast", "ca", "de", "en", "et", "fi", "fy", "gl", "it", "ji", "nl", "sv", "sw", "ur", "yi"}, func(ops *operands)string {
		// i = 1 and v = 0
		if intEqualsAny(ops.I, 1) && intEqualsAny(ops.V, 0) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"si"}, func(ops *operands)string {
		// n = 0,1 or i = 0 and f = 1
		if ops.NequalsAny(0, 1) ||
			intEqualsAny(ops.I, 0) && intEqualsAny(ops.F, 1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ak", "bh", "guw", "ln", "mg", "nso", "pa", "ti", "wa"}, func(ops *operands)string {
		// n = 0..1
		if ops.NinRange(0, 1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"tzm"}, func(ops *operands)string {
		// n = 0..1 or n = 11..99
		if ops.NinRange(0, 1) ||
			ops.NinRange(11, 99) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"pt"}, func(ops *operands)string {
		// n = 0..2 and n != 2
		if ops.NinRange(0, 2) && !ops.NequalsAny(2) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"af", "asa", "az", "bem", "bez", "bg", "brx", "ce", "cgg", "chr", "ckb", "dv", "ee", "el", "eo", "es", "eu", "fo", "fur", "gsw", "ha", "haw", "hu", "jgo", "jmc", "ka", "kaj", "kcg", "kk", "kkj", "kl", "ks", "ksb", "ku", "ky", "lb", "lg", "mas", "mgo", "ml", "mn", "nah", "nb", "nd", "ne", "nn", "nnh", "no", "nr", "ny", "nyn", "om", "or", "os", "pap", "ps", "rm", "rof", "rwk", "saq", "sdh", "seh", "sn", "so", "sq", "ss", "ssy", "st", "syr", "ta", "te", "teo", "tig", "tk", "tn", "tr", "ts", "ug", "uz", "ve", "vo", "vun", "wae", "xh", "xog"}, func(ops *operands)string {
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"pt_PT"}, func(ops *operands)string {
		// n = 1 and v = 0
		if ops.NequalsAny(1) && intEqualsAny(ops.V, 0) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"da"}, func(ops *operands)string {
		// n = 1 or t != 0 and i = 0,1
		if ops.NequalsAny(1) ||
			!intEqualsAny(ops.T, 0) && intEqualsAny(ops.I, 0, 1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"is"}, func(ops *operands)string {
		// t = 0 and i % 10 = 1 and i % 100 != 11 or t != 0
		if intEqualsAny(ops.T, 0) && intEqualsAny(ops.I%10, 1) && !intEqualsAny(ops.I%100, 11) ||
			!intEqualsAny(ops.T, 0) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"mk"}, func(ops *operands)string {
		// v = 0 and i % 10 = 1 or f % 10 = 1
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%10, 1) ||
			intEqualsAny(ops.F%10, 1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"fil", "tl"}, func(ops *operands)string {
		// v = 0 and i = 1,2,3 or v = 0 and i % 10 != 4,6,9 or v != 0 and f % 10 != 4,6,9
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I, 1, 2, 3) ||
			intEqualsAny(ops.V, 0) && !intEqualsAny(ops.I%10, 4, 6, 9) ||
			!intEqualsAny(ops.V, 0) && !intEqualsAny(ops.F%10, 4, 6, 9) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"lv", "prg"}, func(ops *operands)string {
		// n % 10 = 0 or n % 100 = 11..19 or v = 2 and f % 100 = 11..19
		if ops.NmodEqualsAny(10, 0) ||
			ops.NmodInRange(100, 11, 19) ||
			intEqualsAny(ops.V, 2) && intInRange(ops.F%100, 11, 19) {
			return PLURAL_ZERO
		}
		// n % 10 = 1 and n % 100 != 11 or v = 2 and f % 10 = 1 and f % 100 != 11 or v != 2 and f % 10 = 1
		if ops.NmodEqualsAny(10, 1) && !ops.NmodEqualsAny(100, 11) ||
			intEqualsAny(ops.V, 2) && intEqualsAny(ops.F%10, 1) && !intEqualsAny(ops.F%100, 11) ||
			!intEqualsAny(ops.V, 2) && intEqualsAny(ops.F%10, 1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"lag"}, func(ops *operands)string {
		// n = 0
		if ops.NequalsAny(0) {
			return PLURAL_ZERO
		}
		// i = 0,1 and n != 0
		if intEqualsAny(ops.I, 0, 1) && !ops.NequalsAny(0) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ksh"}, func(ops *operands)string {
		// n = 0
		if ops.NequalsAny(0) {
			return PLURAL_ZERO
		}
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"iu", "kw", "naq", "se", "sma", "smi", "smj", "smn", "sms"}, func(ops *operands)string {
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		// n = 2
		if ops.NequalsAny(2) {
			return PLURAL_TWO
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"shi"}, func(ops *operands)string {
		// i = 0 or n = 1
		if intEqualsAny(ops.I, 0) ||
			ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		// n = 2..10
		if ops.NinRange(2, 10) {
			return PLURAL_FEW
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"mo", "ro"}, func(ops *operands)string {
		// i = 1 and v = 0
		if intEqualsAny(ops.I, 1) && intEqualsAny(ops.V, 0) {
			return PLURAL_ONE
		}
		// v != 0 or n = 0 or n != 1 and n % 100 = 1..19
		if !intEqualsAny(ops.V, 0) ||
			ops.NequalsAny(0) ||
			!ops.NequalsAny(1) && ops.NmodInRange(100, 1, 19) {
			return PLURAL_FEW
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"bs", "hr", "sh", "sr"}, func(ops *operands)string {
		// v = 0 and i % 10 = 1 and i % 100 != 11 or f % 10 = 1 and f % 100 != 11
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%10, 1) && !intEqualsAny(ops.I%100, 11) ||
			intEqualsAny(ops.F%10, 1) && !intEqualsAny(ops.F%100, 11) {
			return PLURAL_ONE
		}
		// v = 0 and i % 10 = 2..4 and i % 100 != 12..14 or f % 10 = 2..4 and f % 100 != 12..14
		if intEqualsAny(ops.V, 0) && intInRange(ops.I%10, 2, 4) && !intInRange(ops.I%100, 12, 14) ||
			intInRange(ops.F%10, 2, 4) && !intInRange(ops.F%100, 12, 14) {
			return PLURAL_FEW
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"gd"}, func(ops *operands)string {
		// n = 1,11
		if ops.NequalsAny(1, 11) {
			return PLURAL_ONE
		}
		// n = 2,12
		if ops.NequalsAny(2, 12) {
			return PLURAL_TWO
		}
		// n = 3..10,13..19
		if ops.NinRange(3, 10) || ops.NinRange(13, 19) {
			return PLURAL_FEW
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"sl"}, func(ops *operands)string {
		// v = 0 and i % 100 = 1
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%100, 1) {
			return PLURAL_ONE
		}
		// v = 0 and i % 100 = 2
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%100, 2) {
			return PLURAL_TWO
		}
		// v = 0 and i % 100 = 3..4 or v != 0
		if intEqualsAny(ops.V, 0) && intInRange(ops.I%100, 3, 4) ||
			!intEqualsAny(ops.V, 0) {
			return PLURAL_FEW
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"dsb", "hsb"}, func(ops *operands)string {
		// v = 0 and i % 100 = 1 or f % 100 = 1
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%100, 1) ||
			intEqualsAny(ops.F%100, 1) {
			return PLURAL_ONE
		}
		// v = 0 and i % 100 = 2 or f % 100 = 2
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%100, 2) ||
			intEqualsAny(ops.F%100, 2) {
			return PLURAL_TWO
		}
		// v = 0 and i % 100 = 3..4 or f % 100 = 3..4
		if intEqualsAny(ops.V, 0) && intInRange(ops.I%100, 3, 4) ||
			intInRange(ops.F%100, 3, 4) {
			return PLURAL_FEW
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"he", "iw"}, func(ops *operands)string {
		// i = 1 and v = 0
		if intEqualsAny(ops.I, 1) && intEqualsAny(ops.V, 0) {
			return PLURAL_ONE
		}
		// i = 2 and v = 0
		if intEqualsAny(ops.I, 2) && intEqualsAny(ops.V, 0) {
			return PLURAL_TWO
		}
		// v = 0 and n != 0..10 and n % 10 = 0
		if intEqualsAny(ops.V, 0) && !ops.NinRange(0, 10) && ops.NmodEqualsAny(10, 0) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"cs", "sk"}, func(ops *operands)string {
		// i = 1 and v = 0
		if intEqualsAny(ops.I, 1) && intEqualsAny(ops.V, 0) {
			return PLURAL_ONE
		}
		// i = 2..4 and v = 0
		if intInRange(ops.I, 2, 4) && intEqualsAny(ops.V, 0) {
			return PLURAL_FEW
		}
		// v != 0
		if !intEqualsAny(ops.V, 0) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"pl"}, func(ops *operands)string {
		// i = 1 and v = 0
		if intEqualsAny(ops.I, 1) && intEqualsAny(ops.V, 0) {
			return PLURAL_ONE
		}
		// v = 0 and i % 10 = 2..4 and i % 100 != 12..14
		if intEqualsAny(ops.V, 0) && intInRange(ops.I%10, 2, 4) && !intInRange(ops.I%100, 12, 14) {
			return PLURAL_FEW
		}
		// v = 0 and i != 1 and i % 10 = 0..1 or v = 0 and i % 10 = 5..9 or v = 0 and i % 100 = 12..14
		if intEqualsAny(ops.V, 0) && !intEqualsAny(ops.I, 1) && intInRange(ops.I%10, 0, 1) ||
			intEqualsAny(ops.V, 0) && intInRange(ops.I%10, 5, 9) ||
			intEqualsAny(ops.V, 0) && intInRange(ops.I%100, 12, 14) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"be"}, func(ops *operands)string {
		// n % 10 = 1 and n % 100 != 11
		if ops.NmodEqualsAny(10, 1) && !ops.NmodEqualsAny(100, 11) {
			return PLURAL_ONE
		}
		// n % 10 = 2..4 and n % 100 != 12..14
		if ops.NmodInRange(10, 2, 4) && !ops.NmodInRange(100, 12, 14) {
			return PLURAL_FEW
		}
		// n % 10 = 0 or n % 10 = 5..9 or n % 100 = 11..14
		if ops.NmodEqualsAny(10, 0) ||
			ops.NmodInRange(10, 5, 9) ||
			ops.NmodInRange(100, 11, 14) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"lt"}, func(ops *operands)string {
		// n % 10 = 1 and n % 100 != 11..19
		if ops.NmodEqualsAny(10, 1) && !ops.NmodInRange(100, 11, 19) {
			return PLURAL_ONE
		}
		// n % 10 = 2..9 and n % 100 != 11..19
		if ops.NmodInRange(10, 2, 9) && !ops.NmodInRange(100, 11, 19) {
			return PLURAL_FEW
		}
		// f != 0
		if !intEqualsAny(ops.F, 0) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"mt"}, func(ops *operands)string {
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		// n = 0 or n % 100 = 2..10
		if ops.NequalsAny(0) ||
			ops.NmodInRange(100, 2, 10) {
			return PLURAL_FEW
		}
		// n % 100 = 11..19
		if ops.NmodInRange(100, 11, 19) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ru", "uk"}, func(ops *operands)string {
		// v = 0 and i % 10 = 1 and i % 100 != 11
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%10, 1) && !intEqualsAny(ops.I%100, 11) {
			return PLURAL_ONE
		}
		// v = 0 and i % 10 = 2..4 and i % 100 != 12..14
		if intEqualsAny(ops.V, 0) && intInRange(ops.I%10, 2, 4) && !intInRange(ops.I%100, 12, 14) {
			return PLURAL_FEW
		}
		// v = 0 and i % 10 = 0 or v = 0 and i % 10 = 5..9 or v = 0 and i % 100 = 11..14
		/*if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%10, 0) ||
			intEqualsAny(ops.V, 0) && intInRange(ops.I%10, 5, 9) ||
			intEqualsAny(ops.V, 0) && intInRange(ops.I%100, 11, 14) {
			return PLURAL_MANY
		}*/
		return PLURAL_MANY

	})
	setPluralSets([]string{"br"}, func(ops *operands)string {
		// n % 10 = 1 and n % 100 != 11,71,91
		if ops.NmodEqualsAny(10, 1) && !ops.NmodEqualsAny(100, 11, 71, 91) {
			return PLURAL_ONE
		}
		// n % 10 = 2 and n % 100 != 12,72,92
		if ops.NmodEqualsAny(10, 2) && !ops.NmodEqualsAny(100, 12, 72, 92) {
			return PLURAL_TWO
		}
		// n % 10 = 3..4,9 and n % 100 != 10..19,70..79,90..99
		if (ops.NmodInRange(10, 3, 4) || ops.NmodEqualsAny(10, 9)) && !(ops.NmodInRange(100, 10, 19) || ops.NmodInRange(100, 70, 79) || ops.NmodInRange(100, 90, 99)) {
			return PLURAL_FEW
		}
		// n != 0 and n % 1000000 = 0
		if !ops.NequalsAny(0) && ops.NmodEqualsAny(1000000, 0) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ga"}, func(ops *operands)string {
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		// n = 2
		if ops.NequalsAny(2) {
			return PLURAL_TWO
		}
		// n = 3..6
		if ops.NinRange(3, 6) {
			return PLURAL_FEW
		}
		// n = 7..10
		if ops.NinRange(7, 10) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"gv"}, func(ops *operands)string {
		// v = 0 and i % 10 = 1
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%10, 1) {
			return PLURAL_ONE
		}
		// v = 0 and i % 10 = 2
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%10, 2) {
			return PLURAL_TWO
		}
		// v = 0 and i % 100 = 0,20,40,60,80
		if intEqualsAny(ops.V, 0) && intEqualsAny(ops.I%100, 0, 20, 40, 60, 80) {
			return PLURAL_FEW
		}
		// v != 0
		if !intEqualsAny(ops.V, 0) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"ar"}, func(ops *operands)string {
		// n = 0
		if ops.NequalsAny(0) {
			return PLURAL_ZERO
		}
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		// n = 2
		if ops.NequalsAny(2) {
			return PLURAL_TWO
		}
		// n % 100 = 3..10
		if ops.NmodInRange(100, 3, 10) {
			return PLURAL_FEW
		}
		// n % 100 = 11..99
		if ops.NmodInRange(100, 11, 99) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})
	setPluralSets([]string{"cy"}, func(ops *operands)string {
		// n = 0
		if ops.NequalsAny(0) {
			return PLURAL_ZERO
		}
		// n = 1
		if ops.NequalsAny(1) {
			return PLURAL_ONE
		}
		// n = 2
		if ops.NequalsAny(2) {
			return PLURAL_TWO
		}
		// n = 3
		if ops.NequalsAny(3) {
			return PLURAL_FEW
		}
		// n = 6
		if ops.NequalsAny(6) {
			return PLURAL_MANY
		}
		return PLURAL_OTHER
	})

}




