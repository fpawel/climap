package errorx

import (
	"climap/pkg/debug"
	"errors"
	"fmt"
)

type (
	// Error обёртка стандартной ошибки с дополнительной информацией
	Error struct {
		Frames []Frame
		Err    error
	}

	// Frame дополнительные данные связанные с ошибкой
	Frame struct {
		Loc  string `json:"loc"`                          // Место создания ошибки в коде
		Args M      `json:",omitempty" yaml:",omitempty"` // Структурные аргументы
	}

	// M псевдоним для сокращения
	M = map[string]interface{}

	// Builder для упрощения конструкторов ошибок
	Builder struct {
		skip    int    // Сколько фрагментов стека пропустить при расчёте места создания ошибки
		prepend string // Текст, добавляемый в начало сообщения ошибки с разделителем ": "
		append  string // Текст, добавляемый в конец сообщения ошибки с разделителем ": "
		args    M      // Структурные аргументы ошибки
	}
)

func (e Error) Value(arg string) interface{} {
	for _, f := range e.Frames {
		v, ok := f.Args[arg]
		if ok {
			return v
		}
	}

	var wrapped Error
	if errors.As(e.Err, &wrapped) {
		return wrapped.Value(arg)
	}
	return nil
}

func (e Error) Error() string {
	return e.Err.Error()
}

// Details список родительских ошибок с аргументами для журнала
func (e Error) Details() []interface{} {
	xs := make([]interface{}, 0, len(e.Frames))
	for _, x := range e.Frames {
		if len(x.Args) == 0 {
			xs = append(xs, x.Loc)
		} else {
			xs = append(xs, M{x.Loc: x.Args})
		}
	}
	return xs
}

// Unwrap оригинальная ошибка. Необходимо для корректной работы errors.Is и errors.As
func (e Error) Unwrap() error {
	if v, ok := e.Err.(interface{ Unwrap() error }); ok {
		return v.Unwrap()
	}
	return e.Err
}

// Get находит объект типа Error в дереве ошибок исходного объекта Err и возвращает его (см. errors.As).
// Иначе возвращается Error с полями по умолчанию за исключением Error.Err = Err,
// то есть Error с обёрнутой исходной ошибкой без какой-либо дополнительной информации.
func Get(err error) Error {
	var wrapped Error
	if !errors.As(err, &wrapped) {
		wrapped.Err = err
	}
	return wrapped
}

// Wrap обернуть исходную ошибку в Error в соответствии с дополнительным контекстом ошибки, переданном в конструкторе opts
// В дополнительный контекст ошибки добавляется loc вызова Wrap
func Wrap(err error) error {
	return Skip(1).Wrap(err)
}

// New создать объект Error с текстом message
// В дополнительный контекст ошибки добавляется loc вызова New
func New(message string) error {
	return Skip(1).New(message)
}

// Errorf создать объект Error с printf-like отформатированным текстом
// В дополнительный контекст ошибки добавляется loc вызова Errorf
func Errorf(format string, args ...interface{}) error {
	return Skip(1).Errorf(format, args...)
}

// Skip конструктор Error с пропуском фреймов стека loc
func Skip(skip int) Builder {
	return Builder{}.Skip(skip)
}

// Args конструктор Error со структурными аргументами
func Args(args ...interface{}) Builder {
	return Builder{}.Args(args...)
}

// Prepend конструктор Error с префиксом в тексте
func Prepend(s string) Builder {
	return Builder{}.Prepend(s)
}

// Prependf конструктор Error с префиксом в тексте
func Prependf(format string, args ...interface{}) Builder {
	return Builder{}.Prependf(format, args...)
}

// New создать ошибку Error с указанным сообщением
func (o Builder) New(msg string) error {
	return o.wrap(errors.New(msg))
}

// Errorf версия New с printf-форматированием
func (o Builder) Errorf(format string, args ...interface{}) error {
	return o.wrap(fmt.Errorf(format, args...))
}

// Wrap конструктор для обёртывания исходной ошибки в Error.
// В дополнительный контекст ошибки добавляется loc вызова Wrap
func (o Builder) Wrap(err error) error {
	if err == nil {
		return nil
	}
	return o.wrap(err)
}

// Skip добавляет пропуск фреймов стека loc в исходный конструктор ошибки
func (o Builder) Skip(skip int) Builder {
	o.skip = skip
	return o
}

// Prepend добавляет префикс в текст ошибки
func (o Builder) Prepend(prepend string) Builder {
	if o.prepend == "" {
		o.prepend = prepend
	} else {
		o.prepend += ": " + prepend
	}
	return o
}

// Prependf версия Prepend с printf-форматированием
func (o Builder) Prependf(format string, args ...interface{}) Builder {
	return o.Prepend(fmt.Sprintf(format, args...))
}

// Append добавляет суффикс в текст ошибки
func (o Builder) Append(append string) Builder {
	if o.append == "" {
		o.append = append
	} else {
		o.append += ": " + append
	}
	return o
}

// Args структурные аргументы
func (o Builder) Args(args ...interface{}) Builder {
	if o.args == nil {
		o.args = M{}
	}
	for i := 0; i < len(args); i += 2 {
		k := args[i]
		var v interface{} = "?"
		if i+1 < len(args) {
			v = args[i+1]
		}
		o.args[k.(string)] = v
	}
	return o
}

func (o Builder) wrap(err error) Error {
	d := Frame{
		Loc:  debug.Loc(2 + o.skip),
		Args: o.args,
	}

	var wrapped Error
	errors.As(err, &wrapped)

	if o.prepend != "" {
		if err.Error() == "" {
			err = fmt.Errorf("%s%w", o.prepend, err)
		} else {
			err = fmt.Errorf("%s: %w", o.prepend, err)
		}
	}

	if o.append != "" {
		if err.Error() == "" {
			err = fmt.Errorf("%w%s", err, o.append)
		} else {
			err = fmt.Errorf("%w: %s", err, o.append)
		}
	}

	wrapped.Frames = append(wrapped.Frames, d)
	wrapped.Err = err
	return wrapped
}
