package main

import (
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
const(
	rootdir = "rootdir"
	allowdirectoryview = true
)

func Checkerr(err error) {
	if err != nil {
		log.Printf("\n" + err.Error())
	}
}
func main() {
	fmt.Println("Starting!")
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Println(err.Error())
	}
	for {
		acp, err := ln.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		DefaultHandler(acp)
	}
}
func DefaultHandler(n net.Conn) {
	req, err := sstr.ParseReqHeadersbyString(n)
	Checkerr(err)
	if req.Method == "Not Provided" || req.RequestPath == "Not Provided" {
		n.Write([]byte("Kill youself"))
	}
	if req.Method == "GET" {
		GET(req, n)
	}
}
// Very ugly code sorry for that
func GET(req *sstr.ReqHeader, n net.Conn) {
	header := &sstr.RespHeader{}
	localpath, err := url.QueryUnescape(req.RequestPath)
	Checkerr(err)
	fmt.Printf("\x1b[32mFrom: \x1b[33m%s\n%s %s %s\x1b[m\n\x1b[32mUser-Agent: \x1b[33m%s\x1b[m\n\x1b[32mAccepted types: \x1b[33m%s\x1b[m\n\n", n.RemoteAddr(), req.Method, req.RequestPath, req.HTTPver, req.UserAgent, req.AcceptType)
	if req.RequestPath == "Not Provided" {
		header = sstr.NewDefaultRespHeader(400, len(sstr.BadRequest400()), "text/html", "inline;", "close")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts + sstr.BadRequest400()))
		return
	}
	stat, err := os.Stat(rootdir + localpath)
	if err != nil {
		header = sstr.NewDefaultRespHeader(404, len(sstr.NotFound404(localpath)), "text/html; charset=utf-8", "inline;", "close")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts + sstr.NotFound404(localpath)))
		return
	}
	if stat.IsDir() {
		if allowdirectoryview {
			header = sstr.NewDefaultRespHeader(200, len(sstr.HTMLDirList(rootdir,localpath)), "text/html; charset=utf-8", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\x1b[34m\n%s\x1b[m", headerts)
			n.Write([]byte(headerts + sstr.HTMLDirList(rootdir,localpath)))
			return
		} else {
			header = sstr.NewDefaultRespHeader(404, len(sstr.NotFound404(req.RequestPath)), "text/html", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
			n.Write([]byte(headerts + sstr.NotFound404(req.RequestPath)))
			return
		}
	} else {
		ftype := sstr.GetMimeByExt(filepath.Ext(rootdir + req.RequestPath))
		header = sstr.NewDefaultRespHeader(200, int(stat.Size()), ftype, "inline", "keep-alive")
		headerts := header.PrepRespHeader()
		fmt.Printf("Sent Header:\n\x1b[34m%s\x1b[m", headerts)
		n.Write([]byte(headerts))
		err = sendFile(n, rootdir+req.RequestPath)
		if err != nil {
			log.Printf("%v\n" + err.Error())
			header = sstr.NewDefaultRespHeader(500, len(sstr.ServerErr500()), "text/html; charset=utf-8", "inline;", "close")
			headerts := header.PrepRespHeader()
			fmt.Printf("Sent Header:\n\n\x1b[34m%s\x1b[m", headerts)
			n.Write([]byte(headerts + sstr.ServerErr500()))
		}
	}

}
func sendFile(conn net.Conn, filename string) error {
	// Open the file
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
