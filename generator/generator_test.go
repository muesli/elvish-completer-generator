package generator

import (
	"strings"
	"testing"

	"github.com/dcu/elvish-completer-generator/manpage"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	c := require.New(t)

	p := manpage.New("../manpage/test_files/cp.1.gz")
	err := p.Parse()
	c.Nil(err)

	buffer := &strings.Builder{}
	g := New(p.Name, p.Flags, p.SubCommands)
	err = g.Render(buffer)
	c.Nil(err)

	c.True(buffer.Len() > 0)
}
