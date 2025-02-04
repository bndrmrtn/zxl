package builtin

import (
	"fmt"
	"html"

	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// HtmlModule is a struct responsible for generating HTML content using an existing HttpModule.
type HtmlModule struct {
	httpModule *HttpModule
}

// NewHtmlModule creates and returns a new instance of HtmlModule using the provided HttpModule.
func NewHtmlModule(hm *HttpModule) *HtmlModule {
	return &HtmlModule{
		httpModule: hm,
	}
}

// Access handles the retrieval of variables (not used in this case).
func (hm *HtmlModule) Access(variable string) (*Variable, error) {
	return nil, fmt.Errorf("variable %s not exists", variable)
}

// Execute performs the requested HTML function like 'doc', 'head', 'body', etc.
func (hm *HtmlModule) Execute(fn string, args []*Variable) (*FuncReturn, error) {
	switch fn {
	default:
		return nil, fmt.Errorf("function %s not exists in html module", fn)
	case "doc":
		return hm.doc(args)
	case "head":
		return hm.head(args)
	case "body":
		return hm.body(args)
	case "title":
		return hm.title(args)
	case "h1":
		return hm.h1(args)
	case "p":
		return hm.p(args)
	case "escape":
		return hm.escape(args)
	}
}

// doc generates the entire HTML document by wrapping the provided content.
func (hm *HtmlModule) doc(args []*Variable) (*FuncReturn, error) {
	var content string
	// Concatenate all the arguments into one string.
	for _, arg := range args {
		content += fmt.Sprintf("%v", arg.Value)
	}
	// Wrap the content in <html> tags.
	htmlContent := fmt.Sprintf("<!DOCTYPE html><html>%s</html>", content)
	hm.httpModule.Execute("html", nil)
	_, err := hm.httpModule.Body.WriteString(htmlContent)
	return nil, err
}

// head generates the <head> section of the HTML document.
func (hm *HtmlModule) head(args []*Variable) (*FuncReturn, error) {
	var content string
	// Concatenate all the arguments into one string.
	for _, arg := range args {
		content += fmt.Sprintf("%v", arg.Value)
	}
	return &FuncReturn{
		Type:  tokens.StringVariable,
		Value: fmt.Sprintf("<head>%s</head>", content),
	}, nil
}

// body generates the <body> section of the HTML document.
func (hm *HtmlModule) body(args []*Variable) (*FuncReturn, error) {
	var content string
	// Concatenate all the arguments into one string.
	for _, arg := range args {
		content += fmt.Sprintf("%v", arg.Value)
	}
	return &FuncReturn{
		Type:  tokens.StringVariable,
		Value: fmt.Sprintf("<body>%s</body>", content),
	}, nil
}

// title generates the <title> tag inside the <head> section.
func (hm *HtmlModule) title(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("title function expects 1 argument, got %d", len(args))
	}
	// Ensure the argument is a string.
	if str, ok := args[0].Value.(string); ok {
		return &FuncReturn{
			Type:  tokens.StringVariable,
			Value: fmt.Sprintf("<title>%s</title>", str),
		}, nil
	}
	return nil, fmt.Errorf("title function expects a string argument, got %T", args[0].Value)
}

// h1 generates the <h1> header tag.
func (hm *HtmlModule) h1(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("h1 function expects 1 argument, got %d", len(args))
	}
	// Ensure the argument is a string.
	if str, ok := args[0].Value.(string); ok {
		return &FuncReturn{
			Type:  tokens.StringVariable,
			Value: fmt.Sprintf("<h1>%s</h1>", str),
		}, nil
	}
	return nil, fmt.Errorf("h1 function expects a string argument, got %T", args[0].Value)
}

// p generates the <p> paragraph tag.
func (hm *HtmlModule) p(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("p function expects 1 argument, got %d", len(args))
	}
	// Ensure the argument is a string.
	if str, ok := args[0].Value.(string); ok {
		return &FuncReturn{
			Type:  tokens.StringVariable,
			Value: fmt.Sprintf("<p>%s</p>", str),
		}, nil
	}
	return nil, fmt.Errorf("p function expects a string argument, got %T", args[0].Value)
}

// escape
func (hm *HtmlModule) escape(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("escape function expects 1 argument, got %d", len(args))
	}
	// Ensure the argument is a string.
	if str, ok := args[0].Value.(string); ok {
		return &FuncReturn{
			Type:  tokens.StringVariable,
			Value: html.EscapeString(fmt.Sprintf("%s", str)),
		}, nil
	}
	return nil, fmt.Errorf("escape function expects a string argument, got %T", args[0].Value)
}
