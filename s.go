package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"server/server"
)

// Sorry for copy pasting
const (
	// Set root directory here
	rootdir = "./rootdir/"
	// Allow HTMl directory listing
	allowdirectoryview = true
	// Will always serve index.html from rootdir when asking for /
	indexfirst = false
	// Set port ofc
	port = ":443"
	// Certificate and key path i self signed one and it works okay ig
	certFile = "./cert.pem"
	keyFile  = "./key.pem"
	// test stuff
	ts = true
)

// func
func GetSize(p string) int {
	f, e := os.Stat(p)
	if e != nil {
		log.Println(e.Error())
		return 0
	}
	return int(f.Size())
}
func Checkerr(err error) {
	if err != nil {
		log.Printf("\n" + err.Error())
	}
}
func sendFile(conn net.Conn,p string) error {
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
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("\x1b[31mFailed to load TLS certificate: %s\x1b[m\n", err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("\x1b[31m[-]Failed to listen on port %s: %s\x1b[m\n", port, err)
		return
	}

	tlsListener := tls.NewListener(ln, config)
	defer tlsListener.Close()
	fmt.Printf("\x1b[32m[+]Sucessfully started a server on %s\x1b[m\n", port)
	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("\x1b[31mFailed to accept connection: %s\x1b[31m\n", err)
			continue
		}
		defer conn.Close()

		go DefaultHandler(conn)
	}
}
func DefaultHandler(n net.Conn) {
	for {
		req, err := std.ParseReqHeadersbyString(n)
		// Funni
		if err != nil {
			/*log.Printf("\x1b[31m%s\x1b[m", err.Error()) The only place where a red error log appears as of editing rn i will add soon this may be the culprit of the EOF error idk why*/
			n.Write([]byte("Kill yourself"))
			return
		} else if req.Method == "Not Provided" || req.Path == "Not Provided" {
			n.Write([]byte("Kill yourself"))
			return
		}
		if req.Connection == "close" {
			n.Close()
			return
		}
		if ts {
			if req.Method == "POST" {
				FormTest(req, n)
				return
			}
		}
		if req.Method == "GET" {
			GET(req, n)
		}
	}
}

// Very ugly code sorry for that
func GET(req *std.Req, n net.Conn) {
	localpath, err := url.QueryUnescape(req.Path)
	Checkerr(err)
	fmt.Printf("\x1b[32mFrom: \x1b[33m%s\n%s %s %s\x1b[m\n\x1b[32mUser-Agent: \x1b[33m%s\x1b[m\n\x1b[32mAccepted types: \x1b[33m%s\x1b[m\n\n", n.RemoteAddr(), req.Method, req.Path, req.HTTPver, req.UserAgent, req.AcceptType)
	if indexfirst {
		header := std.NewDefaultRespHeader(200, GetSize(rootdir+"/"+"index.html"), "text/html; charset=utf-8", "inline;", "close")
		n.Write([]byte(header.PrepRespHeader()))
		err := sendFile(n, rootdir+"/"+"index.html")
		if err != nil {
			err = sendFile(n, rootdir+"/"+"index.html")
			if err != nil {
				log.Println(err.Error())
			}
		}
		return
	}
	if req.Path == "Not Provided" {
		header := std.NewDefaultRespHeader(400, len(std.BadRequest400()), "text/html", "inline;", "close")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts + std.BadRequest400()))
		return
	}
	stat, err := os.Stat(rootdir + localpath)
	if err != nil {
		header := std.NewDefaultRespHeader(404, len(std.NotFound404(localpath)), "text/html; charset=utf-8", "inline;", "close")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts + std.NotFound404(localpath)))
		return
	}
	if stat.IsDir() {
		if allowdirectoryview {
			header := std.NewDefaultRespHeader(200, len(std.HTMLDirList(rootdir, req.Path)), "text/html; charset=utf-8", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\x1b[34m\n%s\x1b[m", headerts)
			n.Write([]byte(headerts + std.HTMLDirList(rootdir, localpath)))
			return
		} else {
			header := std.NewDefaultRespHeader(404, len(std.NotFound404(req.Path)), "text/html", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
			n.Write([]byte(headerts + std.NotFound404(req.Path)))
			return
		}
	} else {
		ftype := std.GetMimeByExt(filepath.Ext(rootdir + req.Path))
		header := std.NewDefaultRespHeader(200, int(stat.Size()), ftype, "inline", "keep-alive")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts))
		err = sendFile(n, rootdir+req.Path)
		if err != nil {
			log.Printf("%v\n" + err.Error())
			header = std.NewDefaultRespHeader(500, len(std.ServerErr500()), "text/html; charset=utf-8", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\n\n\x1b[34m%s\x1b[m", headerts)
			n.Write([]byte(headerts + std.ServerErr500()))
		}
		return
	}

}
func FormTest(req *std.Req, n net.Conn) {
	if req.Method == "POST" {
		req.ParseFormData()
		fmt.Printf("\nBody: %s\n", req.Data.FormData["text-input"])
	} else {
		header := std.NewDefaultRespHeader(200, GetSize(rootdir+"index.html"), "text/html; charset=utf-8", "inline;", "keep-alive;")
		n.Write([]byte(header.PrepRespHeader()))
		sendFile(n, rootdir+"index.html")
	}
}
