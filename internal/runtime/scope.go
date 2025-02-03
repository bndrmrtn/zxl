package runtime

// ExecuterScope is the scope of the executer
type ExecuterScope int

const (
	// ExecuterScopeGlobal is the global scope
	ExecuterScopeGlobal ExecuterScope = iota
	// ExecuterScopeBlock is the block scope (if, else, for, etc.)
	ExecuterScopeBlock
	// ExecuterScopeFunction is the function scope
	ExecuterScopeFunction
	// ExecuterScopeDefinition is the definition scope
	ExecuterScopeDefinition
)

func (e ExecuterScope) String() string {
	switch e {
	default:
		return "@exec:global"
	case ExecuterScopeBlock:
		return "@exec:block"
	case ExecuterScopeFunction:
		return "@exec:function"
	case ExecuterScopeDefinition:
		return "@exec:definition"
	}
}
