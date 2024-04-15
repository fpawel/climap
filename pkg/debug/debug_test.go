package debug

import (
	"climap/pkg/str"
	"github.com/elliotchance/pie/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoc(t *testing.T) {
	xs := pie.Map([]int{0, -1, -2}, func(x int) string {
		return Loc(x)
	})
	const expected = `[
    "debug_test.go:12.TestLoc.func1",
    "debug.go:19.Loc",
    "debug.go:33.formatFrame"
]`
	assert.JSONEq(t, expected, str.Jsonify(xs))
}
