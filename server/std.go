package std

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// structs for yknow
type Req struct { // string strings  strings.... Also more header fields with be supported
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
	Path    string
	UserAgent      string
	CurrentConnection  net.Conn // haha very secure
	Data struct{
		ContentLength  int
		FormData map[string]string
		Body string	
	}
	ContentType string
}
type RespHeader struct {
	HTTPver            string
	StatusCode         string
	Date               string
	Server             string
	LastModified       string
	ContentLength      int
	ContentType        string
	ContentDisposition string
	ConnectionType     string
}

// Returns the time in HTTP format
func RetDefaultTime() string {
	loc, _ := time.LoadLocation("Asia/Bangkok") // Set your locale here?
	time.Local = loc
	return time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

// Compiles the header
func (h *RespHeader) PrepRespHeader() string {
	compiled := fmt.Sprintf("%s %s\r\nDate: %s\r\nServer: %s\r\nLast-Modified: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nContent-Disposition: %s\r\nConnection: %s\r\n\r\n", h.HTTPver, h.StatusCode, h.Date, h.Server, h.LastModified, h.ContentLength, h.ContentType, h.ContentDisposition, h.ConnectionType)
	return compiled
}

// Parse request headers
func ParseRequest(n net.Conn) (*Req, error) {
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
	h := &Req{
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
		Path:    "Not Provided",
		UserAgent:      "Not Provided",
	}
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		fields := strings.Split(line, " ")
		if i == 0 && len(fields) >= 3 {
			h.Method = fields[0]
			h.Path = fields[1]
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
			case "content-length":
				h.Data.ContentLength,err = strconv.Atoi(fields[1])
				if err!=nil{
					return nil,err
				}
			}

		}
	}
	tmp := strings.Split(headerstring,"\r\n\r\n")
	h.Data.Body = tmp[1]
	h.Data.FormData = make(map[string]string)
	return h, nil
}
func AckHeader(lent int) *RespHeader{
	h := &RespHeader{}
	h.HTTPver = "HTTP/1.1"
	h.StatusCode = "200 OK"
	h.ContentDisposition = "inline;"
	h.Server = "shitserver/0.0"
	h.LastModified = RetDefaultTime()
	h.Date = RetDefaultTime()
	h.ContentLength = lent
	return h
}
func URLunescape(s string) (string, error) {
    var buf bytes.Buffer
    for i := 0; i < len(s); i++ {
        switch s[i] {
        case '%':
            if i+2 >= len(s) {
                return "", errors.New("Malformed URL-encoded string")
            }
            b1 := s[i+1]
            b2 := s[i+2]
            hexStr := string([]byte{b1, b2})
            hexValue, err := strconv.ParseUint(hexStr, 16, 8)
            if err != nil {
                return "", err
            }
            buf.WriteByte(byte(hexValue))
            i += 2
        case '+':
            buf.WriteByte(' ')
        default:
            buf.WriteByte(s[i])
        }
    }
    return buf.String(), nil
}
type Config struct{
	Port string
	CertPath string
	KeyPath string
	TestStuff bool
	AllowDirView bool
	RootDirectory string
	IndexFirst bool
}
func LoadConfig(p string) (*Config){
	cnfgf,err := os.Open(p)
	if err != nil {
		log.Fatalf("\x1b[31mUnable to open config file: %s",err.Error())
	}
	sc := bufio.NewScanner(cnfgf)
	cn := &Config{}
	for sc.Scan(){
		c := strings.Split(strings.ReplaceAll(sc.Text()," ",""),":")
		switch strings.ToLower(c[0]){
			// use default :port strat
		case "port":
			cn.Port = ":"+c[1]
		case "cert_path":
			cn.CertPath = c[1]
		case "key_path":
			cn.KeyPath = c[1]
		case "allow_directory_view":
			if c[1] == "true"{
				cn.AllowDirView = true
			} else {
				cn.AllowDirView = false
			}
		case "teststuff":
			if c[1] == "true"{
				cn.TestStuff= true
			} else {
				cn.TestStuff = false
			}
		case "root_directory":
			cn.RootDirectory = c[1]
		case "index_first":
			if c[1] == "true"{
				cn.IndexFirst = true
			} else {
				cn.IndexFirst = false
			}
		}

	}
	return cn
}