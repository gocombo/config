package config

import (
	"time"

	"github.com/gocombo/config"
)

type Server struct {
	Port              config.Value[int]
	IdleTimeout       config.Value[time.Duration]
	ReadHeaderTimeout config.Value[time.Duration]
	ReadTimeout       config.Value[time.Duration]
	WriteTimeout      config.Value[time.Duration]
}

type Hello struct {
	Message config.Value[string]
}

type HelloConfig struct {
	SayHelloTimes config.Value[int]
	Server        *Server
	Hello         *Hello
}
