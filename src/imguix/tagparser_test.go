package imguix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	p := newTagParser()
	details := &numericTagDetail{}
	name, err := p.ParseInto("slider:{min:-1,max:10}", details)
	assert.NoError(t, err)
	assert.Equal(t, "slider", name)
	assert.Equal(t, &numericTagDetail{
		Min: -1, Max: 10,
	}, details)
}

func TestParseSimple(t *testing.T) {
	p := newTagParser()
	details := &numericTagDetail{}
	name, err := p.ParseInto("slider", details)
	assert.NoError(t, err)
	assert.Equal(t, "slider", name)
}
