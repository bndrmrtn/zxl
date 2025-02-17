package lang

import "github.com/bndrmrtn/zxl/internal/models"

// Function represents a function method
type Function struct {
	Method

	args []string
	exec func(args []Object) (Object, error)

	debug *models.Debug
}

// NewFunction creates a new function method
func NewFunction(args []string, exec func(args []Object) (Object, error), debug *models.Debug) Method {
	return &Function{
		args:  args,
		exec:  exec,
		debug: debug,
	}
}

func (f *Function) Args() []string {
	return f.args
}

func (f *Function) Execute(args []Object) (Object, error) {
	return f.exec(args)
}

func (f *Function) Debug() *models.Debug {
	return f.debug
}
