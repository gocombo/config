package config

import (
	"time"

	"github.com/gocombo/config/val"
)

type Server struct {
	Port              int
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

type Hello struct {
	Message string
}

type HelloConfig struct {
	SayHelloTimes int
	Server        *Server
	Hello         *Hello
	String        string
}

func newConfig(loader val.Loader) *HelloConfig {
	return &HelloConfig{
		SayHelloTimes: val.Load[int](loader, "sayHelloTimes"),
		Server: &Server{
			Port: val.Load[int](loader, "server/port"),
		},
		Hello: &Hello{
			Message: val.Load[string](loader, "hello/message"),
		},
		String: val.Load[string](loader, "string"),
	}
}
