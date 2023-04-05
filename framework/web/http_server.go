package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpServer struct {
	router *gin.Engine
	ln     net.Listener
	addr   string
}

func NewHttpServer() *HttpServer {

	gin.DisableConsoleColor()
	gin.DefaultWriter = ioutil.Discard
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Debug("route:", httpMethod, ",", absolutePath, ",", handlerName, ",", nuHandlers)
	}
	server := &HttpServer{router: gin.Default()}
	return server
}

func (s *HttpServer) Init(addr string) error {
	s.addr = addr
	var err error
	s.ln, err = net.Listen("tcp", s.addr)
	if err != nil {
		logger.Error("Http listen failed:", s.addr, err)
		return err
	}
	s.router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "I'm ok, "+time.Now().String())
	})
	return nil
}

func (s *HttpServer) GlobalUse(middleware ...gin.HandlerFunc) {
	s.router.Use(middleware...)
}

/*
 method for GET, POST
*/
func (s *HttpServer) RegisterHandler(method, relativePath string, fn gin.HandlerFunc) error {
	s.router.Handle(method, relativePath, fn)
	return nil
}

func (s *HttpServer) Start() error {

	if err := s.router.RunListener(s.ln); err != nil {
		logger.Error("http listen failed: ", s.addr, ", err= ", err)
		return err
	}
	return nil
}
