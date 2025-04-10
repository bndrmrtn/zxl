package prettycode

import "github.com/fatih/color"

func (p *PrettyCode) highlightKeyword(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#ffd230\">" + keyword + "</span>"
	case ConsoleMode:
		return color.New(color.FgHiYellow, color.Bold).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightIdentifier(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		if keyword == "this" {
			return "<b style=\"color:#ff6467\">" + keyword + "</b>"
		}
		return "<span style=\"color:#a684ff\">" + keyword + "</span>"
	case ConsoleMode:
		if keyword == "this" {
			return color.New(color.FgHiRed, color.Bold).Sprint(keyword)
		}
		return color.New(color.FgMagenta).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightFunction(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		if keyword == "import" || keyword == "type" {
			return p.highlightKeyword(mode, keyword)
		}
		return "<span style=\"color:#8ec5ff\">" + keyword + "</span>"
	case ConsoleMode:
		if keyword == "import" {
			return color.New(color.FgHiYellow, color.Bold).Sprint(keyword)
		}
		return color.New(color.FgHiBlue).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightIdentifierWithDot(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		if keyword == "this" {
			return "<b style=\"color:#ff6467\">" + keyword + "</b>"
		}
		return "<span style=\"color:#96f7e4\">" + keyword + "</span>"
	case ConsoleMode:
		if keyword == "this" {
			return color.New(color.FgHiRed, color.Bold).Sprint(keyword)
		}
		return color.New(color.FgHiGreen).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightString(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#f0b100\">" + keyword + "</span>"
	case ConsoleMode:
		return color.New(color.FgYellow).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightUnknown(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#99a1af\">" + keyword + "</span>"
	case ConsoleMode:
		return color.New(color.FgHiBlack).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightBracket(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#f3f4f6\">" + keyword + "</span>"
	case ConsoleMode:
		return color.New(color.FgHiWhite).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightNumber(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#fe9a00\">" + keyword + "</span>"
	case ConsoleMode:
		return color.New(color.FgHiYellow).Sprint(keyword)
	}
	return ""
}

func (p *PrettyCode) highlightOperator(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#d1d5dc\">" + keyword + "</span>"
	case ConsoleMode:
		return color.New(color.FgHiBlack).Sprint(keyword)
	}
	return ""
}
