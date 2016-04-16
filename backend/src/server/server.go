package server

import (
	"net/http"
	"server/controller"

	"golang.org/x/net/websocket"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type Server struct {
}

func (s *Server) Run() error {
	go s.startServer()
	return nil
}

func (s *Server) startServer() {
	ctrl := &controller.Controller{}

	r := gin.Default()

	r.POST("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pong": "ok",
		})
	})

	r.POST("/registration", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"token": "1234",
		})
	})

	api(&r.RouterGroup, "/login", ctrl.Login)
	api(&r.RouterGroup, "/events", ctrl.GetEvents)
	r.GET("/ws_chat", func(c *gin.Context) {
		handler := websocket.Handler(ctrl.ChatWS)
		handler.ServeHTTP(c.Writer, c.Request)
	})
	endless.ListenAndServe(":19888", r)
}

func (s *Server) Stop() {

}

func api(g *gin.RouterGroup, path string, f func(b controller.Binder) (int, *controller.R)) {
	g.POST(path, func(c *gin.Context) {
		i, r := f(c)
		c.JSON(i, r)
	})
}
