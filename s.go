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
func main() {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %s", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

	tlsListener := tls.NewListener(ln, config)
	fmt.Printf("\x1b[32mSucessfully started a server on %s\x1b[m\n", port)
	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		defer conn.Close()

		go DefaultHandler(conn)
	}
}
func DefaultHandler(n net.Conn) {
	req, err := sstr.ParseReqHeadersbyString(n)
	// Funni
	if err != nil {
		log.Printf("\x1b[31m%s\x1b[m", err.Error())
		n.Write([]byte("Kill yourself"))
		return
	} else if req.Method == "Not Provided" || req.Path == "Not Provided" {
		n.Write([]byte("Kill yourself"))
		return
	}
	/*
	if req.Path == "/test" {
		FormTest(req, n)
		return
	}*/

	if req.Method == "GET" {
		GET(req, n)
	}
}

// Very ugly code sorry for that
func GET(req *sstr.Req, n net.Conn) {
	localpath, err := url.QueryUnescape(req.Path)
	Checkerr(err)
	fmt.Printf("\x1b[32mFrom: \x1b[33m%s\n%s %s %s\x1b[m\n\x1b[32mUser-Agent: \x1b[33m%s\x1b[m\n\x1b[32mAccepted types: \x1b[33m%s\x1b[m\n\n", n.RemoteAddr(), req.Method, req.Path, req.HTTPver, req.UserAgent, req.AcceptType)
	if indexfirst {
		header := sstr.NewDefaultRespHeader(200, GetSize(rootdir+"/"+"index.html"), "text/html; charset=utf-8", "inline;", "close")
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
		header := sstr.NewDefaultRespHeader(400, len(sstr.BadRequest400()), "text/html", "inline;", "close")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts + sstr.BadRequest400()))
		return
	}
	stat, err := os.Stat(rootdir + localpath)
	if err != nil {
		header := sstr.NewDefaultRespHeader(404, len(sstr.NotFound404(localpath)), "text/html; charset=utf-8", "inline;", "close")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts + sstr.NotFound404(localpath)))
		return
	}
	if stat.IsDir() {
		if allowdirectoryview {
			header := sstr.NewDefaultRespHeader(200, len(sstr.HTMLDirList(rootdir, req.Path)), "text/html; charset=utf-8", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\x1b[34m\n%s\x1b[m", headerts)
			n.Write([]byte(headerts + sstr.HTMLDirList(rootdir, localpath)))
			return
		} else {
			header := sstr.NewDefaultRespHeader(404, len(sstr.NotFound404(req.Path)), "text/html", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
			n.Write([]byte(headerts + sstr.NotFound404(req.Path)))
			return
		}
	} else {
		ftype := sstr.GetMimeByExt(filepath.Ext(rootdir + req.Path))
		header := sstr.NewDefaultRespHeader(200, int(stat.Size()), ftype, "inline", "keep-alive")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts))
		err = sendFile(n, rootdir+req.Path)
		if err != nil {
			log.Printf("%v\n" + err.Error())
			header = sstr.NewDefaultRespHeader(500, len(sstr.ServerErr500()), "text/html; charset=utf-8", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\n\n\x1b[34m%s\x1b[m", headerts)
			n.Write([]byte(headerts + sstr.ServerErr500()))
		}
	}

}
/*
func FormTest(req *sstr.Req, n net.Conn) {
    if req.Method == "POST" {
        req.ParseFormData()
		header := sstr.AckHeader(len("tysm"))
		fmt.Printf("\nBody: %s", req.Data.FormData["text-input"])
        n.Write([]byte(header.PrepRespHeader()))
		n.Write([]byte("tysm"))
    } else {
        header := sstr.NewDefaultRespHeader(200, GetSize(rootdir+"index.html"), "text/html; charset=utf-8", "inline;", "close;")
        n.Write([]byte(header.PrepRespHeader()))
        sendFile(n, rootdir+"index.html")
    }
}*/
func sendFile(conn net.Conn, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the file to the network connection
	_, err = io.Copy(conn, file)
	if err != nil {
		return err
	}

	return nil
}
