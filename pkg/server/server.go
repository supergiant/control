package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	"github.com/supergiant/supergiant/pkg/api"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/ui"
)

// type SecureInfoHandler struct {
// 	core *core.Core
// }
//
// func (s *SecureInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	msg := "<p>Supergiant is running securely at <a href='" + s.core.BaseURL() + "'>" + s.core.BaseURL() + "</a>.</p>"
// 	msg += "<p>Unless you have provided your own SSL certificate, this will be a self-signed certificate.</p>"
// 	msg += "<p>If using self-signed, your browser will most likely warn of an insecure connection. <strong>You must manually trust the certificate to utilize SSL.</strong></p>"
// 	w.Write([]byte(msg))
// }

//------------------------------------------------------------------------------

func New(c *core.Core) (server *Server, err error) {
	server = &Server{Core: c}

	router := api.NewRouter(c)
	c.Log.Info(c.APIURL())

	if c.UIEnabled {
		router = ui.NewRouter(c, router)
		c.Log.Info(c.UIURL())
	}

	server.primaryHandler = router

	if c.SSLEnabled() {
		server.primaryListener, err = newTLSListener(c.HTTPSPort, c.SSLCertFile, c.SSLKeyFile)
		// if err == nil {
		// 	server.secondaryHandler = &SecureInfoHandler{c}
		// 	server.secondaryListener, err = newListener(c.HTTPPort)
		// }
	} else {
		server.primaryListener, err = newListener(c.HTTPPort)
	}

	return
}

//------------------------------------------------------------------------------

type Server struct {
	Core *core.Core

	primaryHandler http.Handler
	// secondaryHandler http.Handler

	primaryListener net.Listener
	// secondaryListener net.Listener
}

func (s *Server) Start() error {
	// if s.secondaryListener != nil {
	// 	// NOTE we just kinda lose the error here, probably should do something
	// 	go http.Serve(s.secondaryListener, s.secondaryHandler)
	// }
	// CORS options added here
	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Request-Headers", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "PUT", "UPDATE", "POST", "DELETE"})
	return http.Serve(s.primaryListener, handlers.CORS(headersOk, methodsOk)(s.primaryHandler))
}

func (s *Server) Stop() error {
	// if s.secondaryListener != nil {
	// 	if err := s.secondaryListener.Close(); err != nil {
	// 		return err
	// 	}
	// }
	return s.primaryListener.Close()
}

//------------------------------------------------------------------------------

func newListener(port string) (ln net.Listener, err error) {
	ln, err = net.Listen("tcp", ":"+port)
	if err == nil {
		ln = tcpKeepAliveListener{ln.(*net.TCPListener)}
	}
	return
}

func newTLSListener(addr string, certFile string, keyFile string) (ln net.Listener, err error) {
	config := new(tls.Config)
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	ln, err = newListener(addr)
	if err == nil {
		ln = tls.NewListener(ln, config)
	}
	return
}

// Ripped straight from https://golang.org/src/net/http/server.go. Have to
// define our own Listener because we need to be able to Close() it manually.
// (not really sure why it's not exposed in the http package)
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
