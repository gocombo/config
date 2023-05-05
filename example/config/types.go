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
		SayHelloTimes: val.Define[int](p, "sayHelloTimes"),
		Server: &Server{
			Port: val.Define[int](p, "server/port"),
		},
		Hello: &Hello{
			Message: val.Define[string](p, "hello/message"),
		},
		String: val.Define[string](p, "string"),
	}
}
