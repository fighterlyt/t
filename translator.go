package t

import (
	"github.com/fighterlyt/t/format"
	"github.com/fighterlyt/t/translator"
)

// Translator  翻译接口
type Translator = translator.Translator

// NoopTranslator return a no-op Translator
func NoopTranslator() Translator { return noopTranslator }

var noopTranslator Translator = (*nooptor)(nil)

type nooptor struct{}

func (tor *nooptor) Lang() string { return "" }

func (tor *nooptor) X(msgCtxt, msgID string, args ...interface{}) string {
	return format.Format(msgID, args...)
}

func (tor *nooptor) XN64(msgCtxt, msgID, msgIDPlural string, n int64, args ...interface{}) string {
	return format.DefaultPlural(msgID, msgIDPlural, n, args...)
}
