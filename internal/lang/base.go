package lang

import "github.com/bndrmrtn/zexlang/internal/models"

// Base is the base object
type Base struct {
	name    string
	mutable bool
	debug   *models.Debug
}

func NewBase(name string, debug *models.Debug) Base {
	return Base{
		name:    name,
		mutable: true,
		debug:   debug,
	}
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Rename(s string) {
	b.name = s
}

func (b *Base) IsMutable() bool {
	return b.mutable
}

func (b *Base) Immute() {
	b.mutable = false
}

func (b *Base) Debug() *models.Debug {
	return b.debug
}
