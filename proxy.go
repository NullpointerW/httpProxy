package main

import (
	"fmt"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies/shadowsocks"
	"io"
	"log"
	"net/http"
)

const (
	sgp = 8964 + iota
	jp
	taiwan
	usa
	can
	gbr
	nld
	ind
	rus
)

func convertAddr(port int) string {
	return fmt.Sprintf(":%d", port)
}
func NewHttpClient(cfg string) (*http.Client, *shadowsocks.Dialer) {
	dialer, err := shadowsocks.NewDialerWithCfg(proxies.StringResolver, cfg)
	if err != nil {
		panic(err)
	}
	return proxies.NewHttpClient(new(http.Client), dialer), dialer
}

type ProxySrv struct {
	http.Server
	httpClient *http.Client
	dialer     *shadowsocks.Dialer
}

func NewProxyServer(cfg string, loc int) *ProxySrv {
	p := ProxySrv{}
	p.httpClient, p.dialer = NewHttpClient(cfg)
	p.Addr = convertAddr(loc)
	p.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			p.handleTunneling(w, r)
		} else {
			p.handleHTTP(w, r)
		}
	})
	return &p
}

func (p *ProxySrv) handleTunneling(w http.ResponseWriter, req *http.Request) {
	// 目标地址
	log.Println("get request: ", req.URL.Host)
	destConn, err := p.dialer.NewConn(nil, "", req.URL.Host)
	//destConn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// 响应客户端
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// 数据转发
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func (p *ProxySrv) handleHTTP(w http.ResponseWriter, req *http.Request) {
	// 发送请求到目标服务器
	//client := &http.Client{}
	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 复制请求头
	copyHeader(req.Header, request.Header)

	response, err := p.httpClient.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer response.Body.Close()

	// 复制响应头
	copyHeader(response.Header, w.Header())
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}
func copyHeader(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
