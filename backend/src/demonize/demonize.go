package demonize

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

type Demonizer interface {
	Run() error
	Stop()
}

func Demonize(app Demonizer, name string) {
	os.Mkdir("pids", 0777)
	os.Mkdir("logs", 0777)
	ctx := &daemon.Context{
		PidFileName: fmt.Sprintf("pids/%s.pid", name),
		LogFileName: fmt.Sprintf("logs/%s.log", name),
		LogFilePerm: 0644,
		PidFilePerm: 0644,

		Umask: 027,
		Args:  []string{name},
	}
	dmn, err := ctx.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if dmn != nil {
		return
	}
	defer ctx.Release()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)
	go func() {
		_ = <-sigs
		done <- true
	}()

	err = app.Run()
	if err != nil {
		log.Printf("Error on run() %s: %v", name, err)
		return
	}
	<-done

	app.Stop()
}
