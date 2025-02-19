package t

import (
	"io"
	"io/fs"
	"strings"

	"github.com/fighterlyt/log"
	"github.com/fighterlyt/t/locale"
	"github.com/fighterlyt/t/translator"
)

const (
	extPo = ".po"
	extMo = ".mo"
)

// trNoop is a no-op Translation
var trNoop = NewTranslation("", nil)

// Translation can provide different language translation of a domain
// tr. [翻译域]包含各个语言的翻译
type Translation struct {
	domain string
	logger log.Logger            // 日志器
	langs  map[string]Translator // key is language
}

// NewTranslation create a new Translation
func NewTranslation(domain string, logger log.Logger, translators ...Translator) *Translation {
	tr := &Translation{
		domain: domain,
		langs:  make(map[string]translator.Translator),
		logger: logger,
	}

	for _, tor := range translators {
		tr.AddOrReplace(tor)
	}

	return tr
}

// AddOrReplace add a translator and return the previous translator of this language
func (tr *Translation) AddOrReplace(tor Translator) Translator {
	lang := tor.Lang()
	if lang == "" {
		return nil
	}

	lang = locale.Normalize(lang)

	if pre, ok := tr.langs[lang]; ok {
		tr.langs[lang] = tor
		return pre
	}

	tr.langs[lang] = tor

	return nil
}

// Get get the Translator of the specified lang
func (tr *Translation) Get(lang string) (Translator, bool) {
	tor, ok := tr.langs[lang]
	return tor, ok
}

// GetOrNoop return the Translator of the specified language
// 获取指定语言的翻译
func (tr *Translation) GetOrNoop(lang string) Translator {
	if tor, ok := tr.langs[lang]; ok {
		return tor
	}
	return noopTranslator
}

// LoadFS load a translator from file system
func (tr *Translation) LoadFS(f fs.FS) bool {
	var loaded = false
	fn := func(ext string) func(path string, d fs.DirEntry, err error) error {
		return func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d != nil && !d.IsDir() {
				if strings.HasSuffix(d.Name(), ext) { // 这里应该使用 d.Name;
					of, err := f.Open(path) // 这里应该使用 path: file asFS 时 path=. d.Name=file name
					if err == nil {

						defer func() {
							_ = of.Close()
						}()

						if err := tr.LoadFile(of); err == nil {
							loaded = true
						}
					}
				}
			}
			return nil
		}
	}

	_ = fs.WalkDir(f, ".", fn(extMo))
	_ = fs.WalkDir(f, ".", fn(extPo))

	return loaded
}

// LoadFile load a translator from a file
func (tr *Translation) LoadFile(file fs.File) error {
	fi, err := file.Stat()

	if err != nil {
		return err
	}

	fileName := fi.Name()

	var (
		content []byte
	)

	if content, err = io.ReadAll(file); err != nil {
		return err
	}

	if strings.HasSuffix(fileName, extPo) {
		err = tr.LoadPo(content)
	} else if strings.HasSuffix(fileName, extMo) {
		err = tr.LoadMo(content)
	}

	return err
}

// LoadPo load po file
func (tr *Translation) LoadPo(content []byte) error {
	poFile, err := translator.ReadPo(content, tr.logger)
	if err != nil {
		return err
	}

	tr.AddOrReplace(poFile)

	return nil
}

// LoadMo load mo file
func (tr *Translation) LoadMo(content []byte) error {
	moFile, err := translator.ReadMo(content, tr.logger)

	if err != nil {
		return err
	}

	tr.AddOrReplace(moFile)

	return nil
}
