package formatter

type Scope int

const (
	ScopeGlobal Scope = iota
	ScopeDefinition
	ScopeFunction
)
