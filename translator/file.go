package translator

import (
	"regexp"
	"sort"
	"strings"

	"github.com/fighterlyt/t/format"
)

const (
	// HeaderPluralForms 表明该语言的复数形式
	HeaderPluralForms = "Plural-Forms"
	// HeaderLanguage 表明该文件是什么语言
	HeaderLanguage = "Language"
)

var _ Translator = (*File)(nil) // 触发编译检查，是否实现接口
var reHeader = regexp.MustCompile(`(.*?): (.*)`)

// File 一个翻译文件
type File struct {
	entries map[string]*Entry
	headers map[string]string
	plural  *plural
}

// Lang get this translations' language
func (f *File) Lang() string {
	lang, _ := f.GetHeader(HeaderLanguage)
	return lang
}

// X is a short name for p.gettext
func (f *File) X(msgCtxt, msgID string, args ...interface{}) string {
	entry, ok := f.entries[key(msgCtxt, msgID)]
	if !ok || entry.MsgStr == "" {
		return format.Format(msgID, args...)
	}

	return format.Format(entry.MsgStr, args...)
}

// XN64 is a short name for np.gettext
func (f *File) XN64(msgCtxt, msgID, msgIDPlural string, n int64, args ...interface{}) string {
	entry, ok := f.entries[key(msgCtxt, msgID)]
	if !ok {
		return format.DefaultPlural(msgID, msgIDPlural, n, args...)
	}
	plural := f.getPlural()
	if plural.totalForms <= 0 || plural.fn == nil {
		return format.DefaultPlural(msgID, msgIDPlural, n, args...)
	}
	index := plural.fn(n)

	if index < 0 || index >= int(plural.totalForms) || index > len(entry.MsgStrN) || entry.MsgStrN[index] == "" {
		// 超出范围
		return format.DefaultPlural(msgID, msgIDPlural, n, args...)
	}

	return format.Format(entry.MsgStrN[index], args...)
}

// SortedEntry sort entry by key
func (f *File) SortedEntry() (entries []*Entry) {
	for _, e := range f.entries {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		left := entries[i]
		right := entries[j]
		return left.Key() < right.Key()
	})
	return
}

// AddEntry adds a Entry
func (f *File) AddEntry(e *Entry) {
	if f.entries == nil {
		f.entries = map[string]*Entry{}
	}
	f.entries[e.Key()] = e
	if e.isHeader() {
		f.initHeader()
		f.initPlural()
	}
}

// GetHeader get header value by key
func (f *File) GetHeader(key string) (value string, ok bool) {
	f.initHeader()
	value, ok = f.headers[key]
	return
}

func (f *File) initHeader() {
	if f.headers == nil {
		headers := make(map[string]string)
		if headerEntry, ok := f.entries[key("", "")]; ok {
			kvs := strings.Split(headerEntry.MsgStr, "\n")
			for _, kv := range kvs {
				if kv == "" {
					continue
				}
				find := reHeader.FindAllStringSubmatch(kv, -1)
				if len(find) != 1 || len(find[0]) != 3 {
					continue
				}
				kv := find[0]
				k := strings.TrimSpace(kv[1])
				v := strings.TrimSpace(kv[2])
				headers[k] = v
			}
		}
		f.headers = headers
	}
}

func (f *File) getPlural() *plural {
	f.initPlural()
	return f.plural
}

func (f *File) initPlural() {
	if f.plural == nil {
		forms, _ := f.GetHeader(HeaderPluralForms)
		f.plural = parsePlural(forms)
	}
}
