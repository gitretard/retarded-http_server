package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"server/html"
	"server/server"
)

const (
	StdGETfmt = "HTTP/1.1 %s\r\nDate: %s\r\nServer: %s\r\nLast-Modified: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nContent-Disposition: %s\r\nConnection: %s\r\n\r\n"
	/*func (h *RespHeader) PrepRespHeader() string {
		compiled := fmt.Sprintf("%s %s\r\nDate: %s\r\nServer: %s\r\nLast-Modified: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nContent-Disposition: %s\r\nConnection: %s\r\n\r\n", h.HTTPver, h.StatusCode, h.Date, h.Server, h.LastModified, h.ContentLength, h.ContentType, h.ContentDisposition, h.ConnectionType)
		return compiled
	}*/
)

var cfg = std.LoadConfig("config")

func GetStat(p string) (os.FileInfo, error) {
	f, e := os.Stat(p)
	if e != nil {
		log.Println(e.Error())
		return nil, e
	}
	return f, nil
}
func sendFile(conn net.Conn, p string) error {
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the file to the network connection
	_, err = io.Copy(conn, file)
	if err != nil && err.Error() != io.EOF.Error() {
		return err
	}

	return nil
}
func main() {
	cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		log.Fatalf("\x1b[31mFailed to load TLS certificate: %s\x1b[m\n", err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("\x1b[31m[-]Failed to listen on port %s: %s\x1b[m\n", cfg.Port, err)
		return
	}
	defer ln.Close()
	tlsListener := tls.NewListener(ln, config)
	defer tlsListener.Close()
	if string(cfg.RootDirectory[len(cfg.RootDirectory)-1]) != "/" {
		cfg.RootDirectory += "/"
	}
	fmt.Println(cfg.RootDirectory)
	fmt.Printf("\x1b[32m[+]Sucessfully started a server on %s\x1b[m\n", cfg.Port)
	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("\x1b[31mFailed to accept connection: %s\x1b[31m\n", err)
			conn.Close()
			continue
		}
		defer conn.Close()

		go DefaultHandler(conn)
	}
}
func DefaultHandler(n net.Conn) {
	// I am still wondering how it works like you are in a for loop and how does it return to this function every time? damnit i need my brain checked // Wait is the goroutine still stuck?
	req, err := std.ParseRequest(n)
	if err != nil {
		n.Close()
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\x1b[32mFrom: %s\n\x1b[33m%s %s %s\n\x1b[32mUser-Agent: \x1b[33m%s\n\n\x1b[m", n.RemoteAddr(), req.Method, req.Path, req.HTTPver, req.UserAgent)
	if req.Method == "POST" {
		if cfg.TestStuff {
			for {
				//TODO find a way to make POST work properly
				donewithconn := TestPOST(n, req)
				if donewithconn {
					return
				}
			}
		}
	}
	for {
		donewithconn := DefaultGET(n, req)
		if donewithconn {
			return
		}
	}
}

// RespHeader struct for normal responses
func TestPOST(n net.Conn, req *std.Req) bool {
	// I don't know how the browser expects the connection to be . Also ill leave the fucntion here
	req.ParseFormData()
	n.Write([]byte("HTTP/1.1 200 OK \r\n\r\n"))
	fmt.Println(req.Data.FormData["text-input"])
	return true
}

// Will change from struct to fmt string but it ended up worse than i thought so later on
func DefaultGET(n net.Conn, req *std.Req) bool {
	stat, err := os.Stat(cfg.RootDirectory + req.Path)
	if err != nil && os.IsNotExist(err) {
		h := fmt.Sprintf(StdGETfmt, "404 Not found", std.RetDefaultTime(), cfg.ServerName, std.RetDefaultTime(), len(bakedhtml.NotFound404(req.Path)), "text/html", "inline;", "close")
		n.Write([]byte(h + bakedhtml.NotFound404(req.Path)))
		return false
	}
	if err != nil {
		h := fmt.Sprintf(StdGETfmt, "500 Internal server error", std.RetDefaultTime(), cfg.ServerName, std.RetDefaultTime(), len(bakedhtml.ServerErr500()), "text/html", "inline;", "close")
		n.Write([]byte(h + bakedhtml.ServerErr500()))
		return true
	}
	if cfg.IndexFirst {
		s, e := GetStat(cfg.RootDirectory + "index.html")
		if err != nil && os.IsNotExist(err) {
			log.Printf("\x1b[31m\nDude index.html doesnt exist\n\x1b[m")
			h := fmt.Sprintf(StdGETfmt, "200 OK", std.RetDefaultTime(), cfg.ServerName, std.RetDefaultTime(), len(bakedhtml.HTMLDirList(cfg.RootDirectory, "/")), "text/html;charset=utf-8", "inline", "close")
			n.Write([]byte(h + bakedhtml.HTMLDirList(cfg.RootDirectory, req.Path)))
			return false
		} else if err != nil {
			log.Printf("\x1b[Cannot open index.html: %s\n\x1b[m", e.Error())
			h := fmt.Sprintf(StdGETfmt, "500 Internal server error", std.RetDefaultTime(), cfg.ServerName, std.RetDefaultTime(), len(bakedhtml.ServerErr500()), "text/html", "inline;", "close")
			n.Write([]byte(h + bakedhtml.ServerErr500()))
			return true
		}
		h := fmt.Sprintf(StdGETfmt, "200 OK", std.RetDefaultTime(), cfg.ServerName, s.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"), s.Size(), std.GetMimeByExt(s.Name()), "inline", "close")
		n.Write([]byte(h))
		sendFile(n, cfg.RootDirectory+"index.html")
		return false
	}
	if stat.IsDir() {
		if cfg.AllowDirView {
			h := fmt.Sprintf(StdGETfmt, "200 OK", std.RetDefaultTime(), cfg.ServerName, std.RetDefaultTime(), len(bakedhtml.HTMLDirList(cfg.RootDirectory, req.Path)), "text/html;charset=utf-8", "inline", "close")
			n.Write([]byte(h + bakedhtml.HTMLDirList(cfg.RootDirectory, req.Path)))
			return false
		}
	} else {
		h := fmt.Sprintf(StdGETfmt, "200 OK", std.RetDefaultTime(), cfg.ServerName, stat.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"), stat.Size(), std.GetMimeByExt(cfg.RootDirectory+req.Path), "inline", "close")
		n.Write([]byte(h))
		sendFile(n, cfg.RootDirectory+req.Path)
		return false
	}
	h := fmt.Sprintf(StdGETfmt, "404 Not found", std.RetDefaultTime(), cfg.ServerName, std.RetDefaultTime(), len(bakedhtml.NotFound404(req.Path)), "text/html", "inline;", "close")
	n.Write([]byte(h + bakedhtml.NotFound404(req.Path)))
	return false
}
