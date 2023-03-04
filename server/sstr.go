package sstr

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ReqHeader struct { // string strings  strings.... Also more header fields with be supported
	HTTPver        string
	AcceptType     string
	AcceptCharset  string
	AcceptDatetime string
	AcceptEncoding string
	AcceptLanguage string
	Connection     string
	From           string
	Host           string
	Method         string
	RequestPath    string
	UserAgent      string
}
type RespHeader struct {
	HTTPver            string
	StatusCode         string
	StatusAsint        int
	Date               string
	Server             string
	LastModified       string
	ContentLength      int
	ContentType        string
	ContentDisposition string
	ConnectionType     string
}

// Feels kinda clunky sorry for that
var Mediasoundext = []string{".m4a", ".opus", ".flac", ".wav", ".mp3", ".m4b"}
var Vidext = []string{".mp4", ".mkv", ".ogg", ".avi", ".mpeg", ".svi", ".mov", ".flv", ".f4v", ".webm"}
var Imgex = []string{".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico", ".bmp"}

func HTMLDirList(pathto string, a string) string {
	filesListRaw, err := ioutil.ReadDir("./" + pathto + a)
	if err != nil {
		log.Printf("%v %v\n" + err.Error())
	}
	if len(filesListRaw) == 0 {
		return "<!DOCTYPE html><body style=\"background-color:black\"><p style=\"color: white;font-size:20px;\"><b>No files are found in " + pathto + " </b></p></body>"
	}
	filesList := "<!DOCTYPE html><body style=\"background-color:black\"><p style=\"color: white;font-size:20px;\"><b>Index of " + pathto + "</b></p>"
	for index, file := range filesListRaw {
		link := a
		if !strings.HasSuffix(link, "/") {
			link += "/"
		}
		link += file.Name()

		fmt.Println(a + "/" + file.Name())
		filesList += "<a href=\"" + link + "\"><u style=\"text-decoration-color: black;\"><p style=\"font-size: 0.6cm;color:white\">" + strconv.Itoa(index+1) + ". " + func(currfile fs.FileInfo) string {
			currfile.IsDir()
			if currfile.IsDir() {
				return "\U0001F4C1"
			} else {
				for _, ex := range Imgex {
					if filepath.Ext(currfile.Name()) == ex {
						return "\U0001f5bc"
					}
				}
				for _, ex := range Mediasoundext {
					if filepath.Ext(currfile.Name()) == ex {
						return "\U0001f3b5"
					}
				}
				for _, ex := range Vidext {
					if filepath.Ext(currfile.Name()) == ex {
						return "▶️"
					}
				}
				return "\U0001F4C4"
			}
		}(file) + "<b>" + file.Name() + " </b>" + "</p></u></a><br>"
	}
	filesList += "</body>"
	return filesList
}

func RetDefaultTime() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	time.Local = loc
	return time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
func (h *RespHeader) PrepRespHeader() string {
	compiled := fmt.Sprintf("%s %s\r\nDate: %s\r\nServer: %s\r\nLast-Modified: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nContent-Disposition: %s\r\nConnection: %s\r\n\r\n", h.HTTPver, h.StatusCode, h.Date, h.Server, h.LastModified, h.ContentLength, h.ContentType, h.ContentDisposition, h.ConnectionType)
	return compiled
}
func NewDefaultRespHeader(status int, size int, mimetype string, dispositiontype, conntype string) *RespHeader {
	h := &RespHeader{}
	h.HTTPver = "HTTP/1.1"
	h.StatusAsint = status
	h.StatusCode = strconv.Itoa(status)
	switch status {
	case 200:
		h.StatusCode = "200 OK"
	// Wont handle other codes yet
	case 400:
		h.StatusCode = "400 Bad Request"
	case 404:
		h.StatusCode = "404 Not Found"
	case 500:
		h.StatusCode = "500 Internal Server Error"
	}
	h.Date = RetDefaultTime()
	h.Server = "shitserver/0.0 (MicrosoftSucksCocks32)"
	h.LastModified = h.Date
	h.ContentLength = size
	h.ContentType = mimetype
	h.ContentDisposition = dispositiontype
	h.ConnectionType = conntype
	return h
}

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

func GetMimeByExt(ext string) string {
	if mime, ok := mimeTypes[ext]; ok {
		return mime
	} else {
		return "application/octet-stream"
	}
}
func ParseReqHeadersbyString(n net.Conn) (*ReqHeader, error) {
	headerbuf := make([]byte, 8190)
	length, err := n.Read(headerbuf)
	if err != nil {
		return nil, err
	}
	if length < 26 {
		return nil, errors.New(fmt.Sprintf("Header from %s / Length too short!", n.RemoteAddr().String()))
	}
	headerstring := string(headerbuf)
	lines := strings.Split(headerstring, "\r\n")
	h := &ReqHeader{
		HTTPver:        "Not Provided",
		AcceptType:     "Not Provided",
		AcceptCharset:  "Not Provided",
		AcceptDatetime: "Not Provided",
		AcceptEncoding: "Not Provided",
		AcceptLanguage: "Not Provided",
		Connection:     "Not Provided",
		From:           "Not Provided",
		Host:           "Not Provided",
		Method:         "Not Provided",
		RequestPath:    "Not Provided",
		UserAgent:      "Not Provided",
	}
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		fields := strings.Split(line, " ")
		if i == 0 && len(fields) >= 3 {
			h.Method = fields[0]
			h.RequestPath = fields[1]
			h.HTTPver = fields[2]
		} else if len(fields) >= 2 {
			switch strings.ToLower(fields[0]) {
			case "host:":
				h.Host = fields[1]
			case "user-agent:":
				h.UserAgent = strings.Join(fields[1:], " ")
			case "accept:":
				h.AcceptType = fields[1]
			case "accept-encoding:":
				h.AcceptEncoding = strings.Join(fields[1:], " ")
			case "connection:":
				h.Connection = fields[1]
			case "accept-charset:":
				h.AcceptCharset = fields[1]
			case "from:":
				h.From = fields[1]
			case "accept-language:":
				h.AcceptLanguage = fields[1]
			case "accept-datetime:":
				h.AcceptDatetime = strings.Join(fields[1:], " ")
			}
		}
	}
	return h, nil
}

func BadRequest400() string {
	return "<!DOCTYPE html><head><title>400 Bad Request</title></head><body style=\"background-color: black;\"></body><h1 style=\"text-align: center;color:white;font-size: 90px\">400</h1><br><p style=\"text-align: center;color:white;font-size: 30px\"> Bad Request " + "</p></body>"
}
func NotFound404(pathto string) string {
	return "<!DOCTYPE html><head><title>404 Not Found</title></head><body style=\"background-color: black;\"></body><h1 style=\"text-align: center;color:white;font-size: 90px\">404</h1><br><p style=\"text-align: center;color:white;font-size: 30px\"> The specified content is not found: " + path.Base(pathto) + "</p></body>"
}
func ServerErr500() string {
	return "<!DOCTYPE html><head><title>500 Server Error!></title></head><body style=\"background-color: black;\"></body><h1 style=\"text-align: center;color:white;font-size: 90px\">404</h1><br><p style=\"text-align: center;color:white;font-size: 30px\">Internal Server Error" + "</p></body>"
}
