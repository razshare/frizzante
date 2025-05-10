package frizzante

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
)

// existsInEmbeddedFileSystem checks if file exists.
func existsInEmbeddedFileSystem(embeddedFileSystem embed.FS, fileName string) bool {
	return isEmbeddedFile(embeddedFileSystem, fileName) || isEmbeddedDirectory(embeddedFileSystem, fileName)
}

// isEmbeddedFile check if file exists and is a file.
func isEmbeddedFile(embeddedFileSystem embed.FS, fileName string) bool {
	_, err := embeddedFileSystem.ReadFile(fileName)
	if err != nil {
		return false
	}
	return true
}

// isEmbeddedDirectory checks if file exists and is a directory.
func isEmbeddedDirectory(embeddedFileSystem embed.FS, fileName string) bool {
	_, err := embeddedFileSystem.ReadDir(fileName)
	if err != nil {
		return false
	}
	return true
}

// exists checks if file exists.
func exists(fileName string) bool {
	_, statError := os.Stat(fileName)
	return nil == statError || !errors.Is(statError, os.ErrNotExist)
}

// deleteFile deletes a file.
func deleteFile(fileName string) bool {
	removeError := os.Remove(fileName)
	return nil == removeError || !errors.Is(removeError, os.ErrNotExist)
}

// isFile check if file exists and is a file.
func isFile(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return !stat.IsDir()
}

// isDirectory checks if file exists and is a directory.
func isDirectory(fileName string) bool {
	stat, statError := os.Stat(fileName)
	if statError != nil {
		return !errors.Is(statError, os.ErrNotExist)
	}
	return stat.IsDir()
}

var mimes = map[string]string{
	".html":  "text/html",
	".css":   "text/css",
	".txt":   "text/plain",
	".ttf":   "font/ttf",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ico":   "image/x-icon",
	".jpeg":  "image/jpeg",
	".jpg":   "image/jpeg",
	".png":   "image/png",
	".gif":   "image/gif",
	".bmp":   "image/bmp",
	".svg":   "image/svg+xml",
	".tif":   "image/tiff",
	".tiff":  "image/tiff",
	".js":    "text/javascript",
	".json":  "application/json",
	".pdf":   "application/pdf",
	".avi":   "video/x-msvideo",
	".mp4":   "video/mp4",
	".mpeg":  "video/mpeg",
	".ogv":   "video/ogg",
	".webm":  "video/webm",
	".jpgv":  "video/jpg",
	".wasm":  "application/wasm",
	".mkv":   "video/x-matroska",
	".csv":   "text/csv",
	".ics":   "text/calendar",
	".sh":    "application/x-sh",
	".swf":   "application/x-shockwave-flash",
	".tar":   "application/x-tar",
	".xls":   "application/vnd.ms-excel",
	".xml":   "application/xml",
	".xul":   "application/vnd.mozilla.xul+xml",
	".zip":   "application/zip",
	".7z":    "application/x-7z-compressed",
	".apk":   "application/vnd.android.package-archive",
	".jar":   "application/java-archive",
	".vsd":   "application/vnd.visio",
	".xhtml": "application/xhtml+xml",
	".mpkg":  "application/vnd.apple.installer+xml",
	".ppt":   "application/vnd.ms-powerpoint",
	".rar":   "application/x-rar-compressed",
	".rtf":   "application/rtf",
	".3gp":   "video/3gpp",
	".wav":   "audio/x-wav",
	".weba":  "audio/webm",
	".mp3":   "audio/mpeg",
	".3g2":   "video/3gpp2",
	".aac":   "audio/aac",
	".midi":  "audio/midi",
	".mid":   "audio/midi",
	".oga":   "audio/og",
	".abw":   "application/x-abiword",
	".arc":   "application/octet-stream",
	".azw":   "application/vnd.amazon.ebook",
	".bin":   "application/octet-stream",
	".bz":    "application/x-bzip",
	".bz2":   "application/x-bzip2",
	".csh":   "application/x-csh",
	".doc":   "application/msword",
	".epub":  "application/epub+zip",
	".odp":   "application/vnd.oasis.opendocument.presentation",
	".ods":   "application/vnd.oasis.opendocument.spreadsheet",
	".odt":   "application/vnd.oasis.opendocument.text",
	".ogx":   "application/ogg",
}

func mime(fileName string) string {
	extensionName := filepath.Ext(fileName)
	mimeName, ok := mimes[extensionName]

	if ok {
		return mimeName
	}

	return "text/plain"
}
