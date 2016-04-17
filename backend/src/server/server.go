package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"server/controller"

	"golang.org/x/net/websocket"

	qr "github.com/RaymondChou/goqr/pkg"
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
	ctrl := controller.NewController()

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
	room := r.Group("/room/:room_id")
	{
		room.GET("/ws_chat", func(c *gin.Context) {
			room_id := c.Param("room_id")
			s := websocket.Server{
				Handler: websocket.Handler(ctrl.StartChatWS(room_id)),
				Config: websocket.Config{
					Header: http.Header{
						"Access-Control-Allow-Credentials": []string{"true"},
						"Access-Control-Allow-Headers":     []string{"x-websocket-protocol", "x-websocket-version", "x-websocket-extensions", "content-type", "authorization"},
						"Access-Control-Allow-Origin":      []string{"http://176.112.197.64"},
					},
				},
			}
			s.ServeHTTP(c.Writer, c.Request)
			// handler := websocket.Handler(ctrl.StartChatWS(room_id))
			// handler.ServeHTTP(c.Writer, c.Request)
		})
	}
	r.GET("/qr/:value", func(c *gin.Context) {
		value := c.Param("value")
		filename := path.Join("qr", value)
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			os.Mkdir("qr", 0777)
			co, err := qr.Encode(value, qr.L)

			if err != nil {
				fmt.Println(err)
			}

			pngdat := co.PNG()
			ioutil.WriteFile(filename, pngdat, 0666)
		}
		c.File(filename)
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
