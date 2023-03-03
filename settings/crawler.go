package settings

var (
	HTMLBlocks                   = []string{"div", "p", "head", "body", "footer", "header", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "li", "nav", "form", "html", "address", "article", "aside", "blockquote", "canvas", "dd", "dl", "dt", "fieldset", "figcaption", "figure", "hr", "main", "ol", "pre", "section", "table", "tfoot", "ul", "video", "audio", "iframe"}
	DontAddSpaceWhereSymbolRight = []string{".", ",", ";", ":", "?", "!", ")", "]", "}"}
	DontAddSpaceWhereSymbolLeft  = []string{"(", "[", "{"}
	NotAllowedTags               = []string{"script", "noscript", "style", "video", "source", "audio"}
	ImageExtensions              = []string{".apng", ".avif", ".gif", ".jpg", ".jpeg", ".jfif", ".pjpeg", ".pjp", ".png", ".svg", ".webp", ".bmp", ".ico", ".cur", ".tif", ".tiff", ".ai", ".ps", ".psd", ".icns"}
	FontsExtensions              = []string{".ttf", ".otf", ".woff", ".woff2", ".eot", ".fnt", ".fon"}
	FilesExtensions              = []string{".css", ".xml", ".xlsx", ".csv", ".pdf", ".txt", ".asp", ".cer", ".cfm", ".cgi", ".pl", ".js", ".jsp", ".part", ".py", ".rss", ".apk", ".bat", ".bin", ".torrent", ".com", ".exe", ".jar", ".msi", ".wsf", ".dat", ".db", ".dbf", ".log", ".sql", ".tar", ".bin", ".dmg", ".iso", ".toast", ".vcd", ".7z", ".arj", ".deb", ".pkg", ".rar", ".rpm", ".tar.gz", ".z", ".zip", ".magnet", ".key", ".odp", ".pps", ".ppt", ".pptx", ".ods", ".xls", ".xlsm", ".bak", ".cfg", ".cur", ".dll", ".dmp", ".drv", ".ini", ".lnk", ".tmp", ".doc", ".docx", ".odt", ".rtf", ".tex", ".wpd"}
	VideoExtensions              = []string{".mp4", ".mov", ".wmv", ".avi", ".avchd", ".flv", ".f4v", ".swf", ".mkv", ".mpeg-2", ".webm", ".mts", ".3g2", ".3gp", ".h264", ".m4v", ".mpg", ".mpeg", ".rm", ".swf", ".vob", ".wmv"}
	AudioExtensions              = []string{".mp3", ".m4a", ".aac", ".oga", ".ogg", ".flac", ".pcm", ".wav", ".aiff", ".aif", ".cda", ".mid", ".midi", ".mpa", ".wma", ".wpl"}
)
