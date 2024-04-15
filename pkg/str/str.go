// Функции для обработки срок

package str

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/constraints"
	"math"
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// Jsonify форматировать объект в json.
// Полезно для
//   - форматирования объекта для testify.JSONEq
//   - включения малых объектов непосредственно в текст ошибки.
func Jsonify(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// HumanizeInt returns a string with default formatting of integer number, without trailing zeros and with unit suffix.
// Usefully for debugging. It relies on github.com/dustin/go-humanize.ComputeSI
// example:
//
//	HumanizeInt(506) --> 506
//	HumanizeInt(293456) --> 293.46k
//	HumanizeInt(293456789) --> 293.46M
func HumanizeInt[T constraints.Integer | constraints.Float](n T) string {
	v, unit := humanize.ComputeSI(float64(n))
	return fmt.Sprintf("%v%s", math.Round(v*100)/100, unit)
}

// ToSnakeCase "snake_cased" string
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// Truncate string.
func Truncate(str string, length int) string {
	if length <= 0 {
		return ""
	}
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		if count++; count >= length {
			break
		}
	}
	return truncated
}

func ObjId(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}
