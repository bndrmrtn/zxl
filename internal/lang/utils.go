package lang

import "fmt"

func addr(s any) Object {
	return NewString("$addr", fmt.Sprintf("%p", s), nil)
}
