package debug

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"path/filepath"
	"runtime"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// PrintJson распечатать объект в виде json. Полезно для отладки.
func PrintJson(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "\t")
	fmt.Printf("%s\n", b)
}

func Loc(skip int) string {
	return formatFrame(skip+3, func(frame runtime.Frame) string {
		function := filepath.Base(frame.Function)
		for i, ch := range function {
			if string(ch) == "." {
				function = function[i:]
				break
			}
		}
		return fmt.Sprintf("%s:%d%s", filepath.Base(frame.File), frame.Line, function)
	})
}

func formatFrame(skip int, f func(frame runtime.Frame) string) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return f(frame)
}
