package server

import (
	"net/http"
	"server/controller"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

func (s *Server) Run() error {
	go s.startServer()
	return nil
}

func (s *Server) startServer() {
	c := &controller.Controller{}

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

	api(&r.RouterGroup, "/login", c.Login)
	api(&r.RouterGroup, "/events", c.GetEvents)

	r.Run(":19888")
}

func (s *Server) Stop() {

}

func api(g *gin.RouterGroup, path string, f func(b controller.Binder) (int, *controller.R)) {
	g.POST(path, func(c *gin.Context) {
		i, r := f(c)
		c.JSON(i, r)
	})
}
