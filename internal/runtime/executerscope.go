package runtime

// ExecuterScope is the scope of the executer
type ExecuterScope int

const (
	ExecuterScopeGlobal ExecuterScope = iota
	ExecuterScopeFunction
	ExecuterScopeDefinition
)
