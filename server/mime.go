package std

import "path/filepath"

var mimeTypes = map[string]string{
	".aac":   "audio/aac",
	".avi":   "video-x-msvideo",
	".bz":    "application/x-bzip",
	".bz2":   "application/x-bzip2",
	".csh":   "application/x-csh",
	".css":   "text/css",
	".doc":   "application/msword",
	".docx":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".gif":   "image/gif",
	".html":  "text/html",
	".htm":   "text/html",
	".ico":   "image/vnd.microsoft.icon",
	".jar":   "application/java-archive",
	".js":    "text/javascript",
	".json":  "application/json",
	".mjs":   "text/javascript",
	".mp3":   "audio/mpeg",
	".mp4":   "video/mp4",
	".mpeg":  "video/mpeg",
	".odp":   "application/vnd.oasis.opendocument.presentation",
	".ods":   "application/vnd.oasis.opendocument.spreadsheet",
	".odt":   "application/vnd.oasis.opendocument.text",
	".oga":   "audio/ogg",
	".ogv":   "video/ogg",
	".txt":   "text/plain",
	".otf":   "font/otf",
	".png":   "image/png",
	".pdf":   "application/pdf",
	".php":   "application/x-httpd-php",
	".ppt":   "application/vnd.ms-powerpoint",
	".pptx":  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".rar":   "application/vnd.rar",
	".rtf":   "application/rtf",
	".sh":    "application/x-sh",
	".svg":   "image/svg+xml",
	".tar":   "application/x-tar",
	".tiff":  "image/tiff",
	".tf":    "image/tiff",
	".ttf":   "font/ttf",
	".wav":   "audio/wav",
	".weba":  "audio/webm",
	".webm":  "video/webm",
	".webp":  "image/webp",
	".xhtml": "application/xhtml+xml", //Why tf would anyone use this
	".xls":   "application/vnd.ms-excel",
	".xlsx":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xml":   "application/xml",
	".xul":   "application/vnd.mozilla.xul+xml",
	".zip":   "application/zip",
	".7z":    "application/x-7z-compressed",
}

func GetMimeByExt(p string) string {
	if mime, ok := mimeTypes[filepath.Ext(p)]; ok {
		return mime
	} else {
		return "application/octet-stream"
	}
}
