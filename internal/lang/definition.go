package lang

import "github.com/bndrmrtn/zexlang/internal/models"

type Definition struct {
	Object

	name string

	variables map[string]Object
	functions map[string]Method

	debug *models.Debug
}
