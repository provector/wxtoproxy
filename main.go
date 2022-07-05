// Credit for proxy server:  yowu/HttpProxy.go
// https://gist.github.com/yowu/f7dc34bd4736a65ff28d
// MIT License

package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

const internalVer = "1.0.0"

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, " ", req.Method, " ", req.URL)

	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		msg := "unsupported protocal scheme " + req.URL.Scheme
		http.Error(wr, msg, http.StatusBadRequest)
		log.Println(msg)
		return
	}

	client := &http.Client{}

	//http: Request.RequestURI can't be set in client requests.
	//http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	delHopHeaders(req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		appendHostToXForwardHeader(req.Header, clientIP)
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
		time.Sleep(time.Second * 5)
	}

	log.Println(req.RemoteAddr, " ", resp.Status)

	delHopHeaders(resp.Header)

	copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	time.Sleep(10 * time.Millisecond)
	resp.Body.Close()
}

func main() {
	log.Printf("Starting WxToProxy wrapper ver %s, by provector '22 MIT License...", internalVer)
	go func() {
		var addr = flag.String("addr", "127.0.0.1:8080", "The addr of the application.")
		flag.Parse()
		handler := &proxy{}
		log.Println("Starting proxy server on", *addr)
		if err := http.ListenAndServe(*addr, handler); err != nil {
			log.Fatal("ListenAndServe:", err)
			time.Sleep(time.Second * 5)
		}
	}()
	cmd := exec.Command("xwxtoimg.exe")
	runErr := cmd.Start()
	if runErr != nil {
		log.Fatal("RunError: ", runErr)
		time.Sleep(time.Second * 5)
	}
	cmd.Wait()
	log.Println("Finished")
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type proxy struct{}
