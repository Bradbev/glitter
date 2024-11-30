package imguix

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bzick/tokenizer"
)

func MustParseInto(toParse string, dest any) string {
	p := newTagParser()
	if toParse == "" {
		return ""
	}
	kind, err := p.ParseInto(toParse, dest)
	if err != nil {
		panic(err)
	}
	return kind
}

const (
	tSeparator = iota + 1
	tBraceOpen
	tBraceClose
	tQuote
	tComma
	tNegative
)

type tagParser struct {
	tokenizer  *tokenizer.Tokenizer
	stream     *tokenizer.Stream
	parseError string
	err        error

	name string
}

func (t *tagParser) ParseInto(toParse string, dest any) (name string, err error) {
	t.stream = t.tokenizer.ParseString(toParse)
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				t.err = e
			} else {
				panic(r)
			}
		}
		t.stream.Close()
		err = t.err
	}()
	t.name = t.expectKeyword()
	if t.stream.IsValid() {
		t.expect(tSeparator)
		t.expectDetails(dest)
	}
	return t.name, nil
}

func (t *tagParser) endWithError(e error) {
	remaining := ""
	for t.stream.IsValid() {
		remaining += t.stream.CurrentToken().ValueString() + " "
		t.stream.GoNext()
	}
	panic(fmt.Errorf("%w with remaining\n%s", e, remaining))
}

func (t *tagParser) expectKeyword() string {
	defer t.stream.GoNext()
	tok := t.stream.CurrentToken()
	if !tok.Is(tokenizer.TokenKeyword) {
		t.endWithError(fmt.Errorf("expected a keyword, but got %v", tok.String()))
	}
	return tok.ValueString()
}

func (t *tagParser) expectFloat() float64 {
	defer t.stream.GoNext()
	sign := 1.0
	if t.maybe(tNegative) {
		sign = -1.0
	}
	tok := t.stream.CurrentToken()

	if tok.Is(tokenizer.TokenInteger) {
		return sign * float64(tok.ValueInt64())
	}
	if !tok.Is(tokenizer.TokenFloat) {
		t.endWithError(fmt.Errorf("expected a float for the name, but got %v", tok))
	}
	return sign * tok.ValueFloat64()
}

func (t *tagParser) peek(key tokenizer.TokenKey) bool {
	return t.stream.CurrentToken().Key() == key
}

// maybe will peek for key, and if the next token is key, advance the stream
func (t *tagParser) maybe(key tokenizer.TokenKey) bool {
	ret := t.peek(key)
	if ret {
		t.stream.GoNext()
	}
	return ret
}

func (t *tagParser) expect(key tokenizer.TokenKey) {
	defer t.stream.GoNext()
	tok := t.stream.CurrentToken()
	if tok.Key() != key {
		t.endWithError(fmt.Errorf("expected tokenkey %v, but have %v", key, tok.Key()))
	}
}

func (t *tagParser) makeFieldMap(dest any) map[string]reflect.Value {
	result := map[string]reflect.Value{}
	typ := reflect.TypeOf(dest).Elem()
	val := reflect.ValueOf(dest).Elem()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		result[strings.ToLower(field.Name)] = val.Field(i)
	}
	return result
}

func (t *tagParser) expectDetails(dest any) {
	defer t.stream.GoNext()
	t.expect(tBraceOpen)
	fields := t.makeFieldMap(dest)
	for {
		key := t.expectKeyword()
		t.expect(tSeparator)
		val := t.expectFloat()
		field, ok := fields[strings.ToLower(key)]
		if !ok {
			t.endWithError(fmt.Errorf("Unexpected field %s", key))
		}
		field.SetFloat(val)

		if !t.maybe(tComma) {
			break
		}
	}
	t.expect(tBraceClose)
}

func newTagParser() *tagParser {
	t := *tokenizer.New()
	t.DefineTokens(tSeparator, []string{":"})
	t.DefineTokens(tComma, []string{","})
	t.DefineTokens(tNegative, []string{"-"})
	t.DefineTokens(tBraceOpen, []string{"{"})
	t.DefineTokens(tBraceClose, []string{"}"})
	t.DefineStringToken(tQuote, `"`, `"`).SetEscapeSymbol(tokenizer.BackSlash)
	return &tagParser{
		tokenizer: &t,
	}
}
