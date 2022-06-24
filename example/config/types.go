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

func newConfig(p val.Provider) *HelloConfig {
	return &HelloConfig{
		SayHelloTimes: val.Get[int](p, "sayHelloTimes"),
		Server: &Server{
			Port: val.Get[int](p, "server/port"),
		},
		Hello: &Hello{
			Message: val.Get[string](p, "hello/message"),
		},
		String: val.Get[string](p, "string"),
	}
}
