package prettycode

func (p *PrettyCode) highlightKeyword(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#ffd230\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightIdentifier(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		if keyword == "this" {
			return "<span style=\"color:#ff6467\">" + keyword + "</span>"
		}
		return "<span style=\"color:#a684ff\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightFunction(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		if keyword == "import" {
			return p.highlightKeyword(mode, keyword)
		}
		return "<span style=\"color:#8ec5ff\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightIdentifierWithDot(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		if keyword == "this" {
			return "<span style=\"color:#ff6467\">" + keyword + "</span>"
		}
		return "<span style=\"color:#96f7e4\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightString(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#f0b100\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightUnknown(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#99a1af\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightBracket(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#f3f4f6\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightNumber(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#fe9a00\">" + keyword + "</span>"
	}
	return ""
}

func (p *PrettyCode) highlightOperator(mode Mode, keyword string) string {
	switch mode {
	case HtmlMode:
		return "<span style=\"color:#d1d5dc\">" + keyword + "</span>"
	}
	return ""
}
