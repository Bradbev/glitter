package imguix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	p := newTagParser()
	details := &numericTagDetail{}
	name, err := p.ParseInto("slider:{min:0,max:10}", details)
	assert.NoError(t, err)
	assert.Equal(t, "slider", name)
	assert.Equal(t, &numericTagDetail{
		Min: 0, Max: 1,
	}, details)
}

func TestParseSimple(t *testing.T) {
	p := newTagParser()
	details := &numericTagDetail{}
	name, err := p.ParseInto("slider", details)
	assert.NoError(t, err)
	assert.Equal(t, "slider", name)
}
