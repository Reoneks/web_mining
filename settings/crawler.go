package settings

var (
	HTMLBlocks                   = []string{"div", "p", "head", "body", "footer", "header", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "li", "nav", "form", "html", "address", "article", "aside", "blockquote", "canvas", "dd", "dl", "dt", "fieldset", "figcaption", "figure", "hr", "main", "ol", "pre", "section", "table", "tfoot", "ul", "video", "audio", "iframe"}
	DontAddSpaceWhereSymbolRight = []string{".", ",", ";", ":", "?", "!", ")", "]", "}"}
	DontAddSpaceWhereSymbolLeft  = []string{"(", "[", "{"}
	NotAllowedTags               = []string{"script", "noscript", "style", "video", "source", "audio"}
)
