package settings

var (
	HTMLBlocks                   = []string{"div", "p", "head", "body", "footer", "header", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "li", "nav", "form", "html", "address", "article", "aside", "blockquote", "canvas", "dd", "dl", "dt", "fieldset", "figcaption", "figure", "hr", "main", "ol", "pre", "section", "table", "tfoot", "ul", "video", "audio", "iframe"}
	DontAddSpaceWhereSymbolRight = []string{".", ",", ";", ":", "?", "!", ")", "]", "}"}
	DontAddSpaceWhereSymbolLeft  = []string{"(", "[", "{"}
	NotAllowedTags               = []string{"script", "noscript", "style", "video", "source", "audio"}
	ImageExtensions              = []string{".apng", ".avif", ".gif", ".jpg", ".jpeg", ".jfif", ".pjpeg", ".pjp", ".png", ".svg", ".webp", ".bmp", ".ico", ".cur", ".tif", ".tiff"}
	FontsExtensions              = []string{".ttf", ".otf", ".woff", ".woff2", ".eot"}
	FilesExtensions              = []string{".css", ".xml", ".xlsx", ".csv", ".pdf", ".txt"}
	VideoExtensions              = []string{".mp4", ".mov", ".wmv", ".avi", ".avchd", ".flv", ".f4v", ".swf", ".mkv", ".mpeg-2", ".webm", ".mts"}
	AudioExtensions              = []string{".mp3", ".m4a", ".aac", ".oga", ".ogg", ".flac", ".pcm", ".wav", ".aiff"}
)
