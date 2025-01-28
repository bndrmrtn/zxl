package runtime

type RuntimeMode int

const (
	EntryPoint RuntimeMode = iota
	CodeBlock
)
