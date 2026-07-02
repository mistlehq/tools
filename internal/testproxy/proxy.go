package testproxy

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type AuthMode string

const (
	AuthModeBasic             AuthMode = "basic"
	AuthModeBearer            AuthMode = "bearer"
	AuthModeHeader            AuthMode = "header"
	AuthModePathSegmentPrefix AuthMode = "path_segment_prefix"
)

type Config struct {
	UpstreamBaseURL   string
	AuthMode          AuthMode
	Username          string
	Password          string
	Token             string
	HeaderName        string
	HeaderValue       string
	SegmentPrefix     string
	SegmentCredential string
	Headers           map[string]string
}

type Server struct {
	BaseURL string
	Close   func() error
}

func basicAuthorizationHeader(username string, password string) string {
	credentials := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}

func prependPathSegment(outURL *url.URL, segment string) {
	outURL.Path = "/" + segment + outURL.Path
	outURL.RawPath = ""
}

func Start(config Config) (*Server, error) {
	upstreamURL, err := url.Parse(config.UpstreamBaseURL)
	if err != nil {
		return nil, err
	}

	switch config.AuthMode {
	case AuthModeBasic:
		if config.Username == "" {
			return nil, fmt.Errorf("basic auth requires username")
		}

		if config.Password == "" {
			return nil, fmt.Errorf("basic auth requires password")
		}
	case AuthModeBearer:
		if config.Token == "" {
			return nil, fmt.Errorf("bearer auth requires token")
		}
	case AuthModeHeader:
		if config.HeaderName == "" {
			return nil, fmt.Errorf("header auth requires header name")
		}
		if config.HeaderValue == "" {
			return nil, fmt.Errorf("header auth requires header value")
		}
	case AuthModePathSegmentPrefix:
		if config.SegmentPrefix == "" {
			return nil, fmt.Errorf("path segment prefix auth requires segment prefix")
		}
		if config.SegmentCredential == "" {
			return nil, fmt.Errorf("path segment prefix auth requires segment credential")
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
			case AuthModeBasic:
				r.Out.Header.Set("Authorization", basicAuthorizationHeader(config.Username, config.Password))
			case AuthModeBearer:
				r.Out.Header.Set("Authorization", "Bearer "+config.Token)
			case AuthModeHeader:
				r.Out.Header.Set(config.HeaderName, config.HeaderValue)
			case AuthModePathSegmentPrefix:
				prependPathSegment(r.Out.URL, config.SegmentPrefix+config.SegmentCredential)
			}
			for name, value := range config.Headers {
				r.Out.Header.Set(name, value)
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

	return &Server{
		BaseURL: "http://" + listener.Addr().String(),
		Close:   server.Close,
	}, nil
}
