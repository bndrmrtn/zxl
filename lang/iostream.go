package lang

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type IOStream struct {
	Base
	reader io.Reader
}

func NewIOStream(name string, r io.Reader) Object {
	return &IOStream{
		Base:   NewBase(name, nil),
		reader: r,
	}
}

func (f *IOStream) Type() ObjType {
	return TIOStream
}

func (f *IOStream) Name() string {
	return f.name
}

func (f *IOStream) Value() any {
	return f
}

func (i *IOStream) Method(name string) Method {
	switch name {
	case "readLine":
		return NewFunction(func(_ []Object) (Object, error) {
			reader := bufio.NewReader(i.reader)
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			return NewString("line", line, i.debug), nil
		})
	case "readLines":
		return NewFunction(func(_ []Object) (Object, error) {
			reader := bufio.NewReader(i.reader)
			lines := strings.Builder{}
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}
				lines.WriteString(line)
			}

			result := lines.String()
			if strings.HasSuffix(result, "\n") {
				result = result[:len(result)-1]
			}

			return NewString("lines", result, i.debug), nil
		})
	case "close":
		return NewFunction(func(_ []Object) (Object, error) {
			if rc, ok := i.reader.(io.Closer); ok {
				return nil, rc.Close()
			}
			return nil, nil
		})
	}

	return nil
}

func (i *IOStream) Methods() []string {
	return []string{"readLine", "readLines", "close"}
}

func (i *IOStream) Variable(name string) Object {
	switch name {
	default:
		return nil
	case "$addr":
		return addr(i)
	}
}

func (*IOStream) Variables() []string {
	return []string{"$addr"}
}

func (i *IOStream) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (i *IOStream) String() string {
	return fmt.Sprintf("<IOStream %s>", addr(i))
}

func (i *IOStream) Copy() Object {
	return &IOStream{
		Base:   NewBase(i.name, nil),
		reader: i.reader,
	}
}
