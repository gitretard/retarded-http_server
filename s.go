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

// set config path
var cfg = std.LoadConfig("config")

// Sorry for copy pasting
const (
// Set root directory here
// rootdir = "./rootdir/"
// Allow HTMl directory listing
// allowdirectoryview = true
// Will always serve index.html from rootdir when asking for /
// indexfirst = false
// Set port ofc
// port = ":443"
// Certificate and key path i self signed one and it works okay ig
// certFile = "./cert.pem"
// keyFile  = "./key.pem"
// test stuff
// ts = true
)

// func
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
func DefaultGET(n net.Conn, req *std.Req) bool {
	stat, err := os.Stat(cfg.RootDirectory + req.Path)
	if err != nil && os.IsNotExist(err) {
		h := &std.RespHeader{
			StatusCode:         "404 Not found",
			HTTPver:            "HTTP/1.1",
			Date:               std.RetDefaultTime(),
			Server:             "retarded_server/1.1",
			LastModified:       std.RetDefaultTime(),
			ContentType:        "text/html; charset=utf-8",
			ContentLength:      len(bakedhtml.NotFound404(req.Path)),
			ContentDisposition: "inline",
			ConnectionType:     "close",
		}
		n.Write([]byte(h.PrepRespHeader() + bakedhtml.NotFound404(req.Path)))
		return false
	}
	if err != nil {
		h := &std.RespHeader{
			StatusCode:         "500 Server error. Closing connection",
			HTTPver:            "HTTP/1.1",
			Date:               std.RetDefaultTime(),
			Server:             "retarded_server/1.1",
			LastModified:       std.RetDefaultTime(),
			ContentType:        "text/html; charset=utf-8",
			ContentLength:      len(bakedhtml.ServerErr500()),
			ContentDisposition: "inline",
			ConnectionType:     "close",
		}
		n.Write([]byte(h.PrepRespHeader() + bakedhtml.ServerErr500()))
		return true
	}
	if cfg.IndexFirst {
		s, e := GetStat(cfg.RootDirectory + "index.html")
		if err != nil && os.IsNotExist(err) {
			log.Printf("\x1b[31m\nDude index.html doesnt exist\n\x1b[m")
			h := &std.RespHeader{
				StatusCode:         "200",
				HTTPver:            "HTTP/1.1",
				Date:               std.RetDefaultTime(),
				Server:             "retarded_server/1.1",
				LastModified:       std.RetDefaultTime(),
				ContentType:        "text/html; charset=utf-8",
				ContentLength:      len(bakedhtml.HTMLDirList(cfg.RootDirectory, "/")),
				ContentDisposition: "inline",
				ConnectionType:     "keep-alive",
			}
			n.Write([]byte(h.PrepRespHeader() + bakedhtml.HTMLDirList(cfg.RootDirectory, req.Path)))
			return false
		} else if err != nil {
			log.Printf("\x1b[Cannot open index.html: %s\n\x1b[m", e.Error())
			h := &std.RespHeader{
				StatusCode:         "500 Server error. Closing connection",
				HTTPver:            "HTTP/1.1",
				Date:               std.RetDefaultTime(),
				Server:             "retarded_server/1.1",
				LastModified:       std.RetDefaultTime(),
				ContentType:        "text/html; charset=utf-8",
				ContentLength:      len(bakedhtml.ServerErr500()),
				ContentDisposition: "inline",
				ConnectionType:     "close",
			}
			n.Write([]byte(h.PrepRespHeader() + bakedhtml.ServerErr500()))
			return true
		}
		h := &std.RespHeader{
			StatusCode:         "200 OK",
			HTTPver:            "HTTP/1.1",
			Date:               std.RetDefaultTime(),
			Server:             "retarded_server/1.1",
			LastModified:       stat.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
			ContentType:        "text/html; charset=utf-8",
			ContentLength:      int(s.Size()),
			ContentDisposition: "inline",
			ConnectionType:     "keep-alive",
		}
		n.Write([]byte(h.PrepRespHeader()))
		sendFile(n, cfg.RootDirectory+"/index.html")
		return false
	}
	if stat.IsDir() {
		if cfg.AllowDirView {
			h := &std.RespHeader{
				StatusCode:         "200 OK",
				HTTPver:            "HTTP/1.1",
				Date:               std.RetDefaultTime(),
				Server:             "retarded_server/1.1",
				LastModified:       std.RetDefaultTime(),
				ContentType:        "text/html; charset=utf-8",
				ContentLength:      len(bakedhtml.HTMLDirList(cfg.RootDirectory, req.Path)),
				ContentDisposition: "inline",
				ConnectionType:     "keep-alive",
			}
			n.Write([]byte(h.PrepRespHeader() + bakedhtml.HTMLDirList(cfg.RootDirectory, req.Path)))
			return false
		}
	} else {
		h := &std.RespHeader{
			StatusCode:         "200 OK",
			HTTPver:            "HTTP/1.1",
			Date:               std.RetDefaultTime(),
			Server:             "retarded_server/1.1",
			LastModified:       stat.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
			ContentType:        "text/html; charset=utf-8",
			ContentLength:      int(stat.Size()),
			ContentDisposition: "inline",
			ConnectionType:     "keep-alive",
		}
		n.Write([]byte(h.PrepRespHeader()))
		sendFile(n, cfg.RootDirectory+"/"+req.Path)
		return false
	}
	h := &std.RespHeader{
		StatusCode:         "404 Not found",
		HTTPver:            "HTTP/1.1",
		Date:               std.RetDefaultTime(),
		Server:             "retarded_server/1.1",
		LastModified:       std.RetDefaultTime(),
		ContentType:        "text/html; charset=utf-8",
		ContentLength:      len(bakedhtml.NotFound404(req.Path)),
		ContentDisposition: "inline",
		ConnectionType:     "close",
	}
	n.Write([]byte(h.PrepRespHeader() + bakedhtml.NotFound404(req.Path)))
	return false

}
