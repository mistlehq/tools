package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type proxyAuthMode string

const (
	ProxyAuthModeBasic  proxyAuthMode = "basic"
	ProxyAuthModeBearer proxyAuthMode = "bearer"
)

type proxyConfig struct {
	UpstreamBaseURL string
	AuthMode        proxyAuthMode
	Username        string
	Password        string
	Token           string
}

type proxyServer struct {
	BaseURL string
	Close   func() error
}

func basicAuthorizationHeader(username string, password string) string {
	credentials := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}

func startProxyServer(config proxyConfig) (*proxyServer, error) {
	upstreamURL, err := url.Parse(config.UpstreamBaseURL)
	if err != nil {
		return nil, err
	}

	switch config.AuthMode {
	case ProxyAuthModeBasic:
		if config.Username == "" {
			return nil, fmt.Errorf("basic auth requires username")
		}

		if config.Password == "" {
			return nil, fmt.Errorf("basic auth requires password")
		}
	case ProxyAuthModeBearer:
		if config.Token == "" {
			return nil, fmt.Errorf("bearer auth requires token")
		}

	default:
		return nil, fmt.Errorf("unsupported auth mode: %s", config.AuthMode)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(upstreamURL)
			r.Out.Host = upstreamURL.Host
			r.Out.Header.Del("Authorization")

			switch config.AuthMode {
			case ProxyAuthModeBasic:
				r.Out.Header.Set("Authorization", basicAuthorizationHeader(config.Username, config.Password))
			case ProxyAuthModeBearer:
				r.Out.Header.Set("Authorization", "Bearer "+config.Token)
			}
		},
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Handler: proxy,
	}

	go func() {
		_ = server.Serve(listener)
	}()

	return &proxyServer{
		BaseURL: "http://" + listener.Addr().String(),
		Close:   server.Close,
	}, nil
}
