package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var ctx = context.Background()

var (
	hostTarget = map[string]string{
		"baidu.sun.com": "http://www.baidu.com",
	}
	hostProxy = make(map[string]*httputil.ReverseProxy)
)

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if fn, ok := hostProxy[host]; ok {
		fn.ServeHTTP(w, r)
		return
	}

	if target, ok := hostTarget[host]; ok {
		remoteUrl, err := url.Parse(target)
		if err != nil {
			log.Println("target parse fail:", err)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
		hostProxy[host] = proxy
		r.Host = remoteUrl.Host
		proxy.ServeHTTP(w, r)
		return
	}
	w.Write([]byte("403: Host forbidden " + host))
}

func healthCheckhandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	h := &baseHandle{}
	http.Handle("/", h)
	server := &http.Server{
		Addr:    ":80",
		Handler: h,
	}
	server.ListenAndServe()
}
