package str

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHumanizeInt(t *testing.T) {
	assert.Equal(t, "123.46M", HumanizeInt(123456789))
	assert.Equal(t, "123.46k", HumanizeInt(123456))
	assert.Equal(t, "123", HumanizeInt(123))
	assert.Equal(t, "12.35T", HumanizeInt(12345678900000))
}

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "camel_case", ToSnakeCase("CamelCase"))
	assert.Equal(t, "camel_camel_case", ToSnakeCase("CamelCamelCase"))
	assert.Equal(t, "camel__camel__case", ToSnakeCase("Camel_Camel_Case"))
	assert.Equal(t, "camel_case", ToSnakeCase("camelCase"))
	assert.Equal(t, "camel_camel_case", ToSnakeCase("camelCamelCase"))
	assert.Equal(t, "", ToSnakeCase(""))
	assert.Equal(t, "case", ToSnakeCase("Case"))
	assert.Equal(t, "case", ToSnakeCase("case"))
}
