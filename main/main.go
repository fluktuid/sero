package main

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	cfg "github.com/fluktuid/sero/config"
	t "github.com/fluktuid/sero/target"
)

var target t.Target

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	config := cfg.LoadConfig()

	target = *t.Init(config.Target.Deployment)
	ln, err := net.Listen("tcp", config.Host)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(config.Target.Protocol, config.Target.Host, config.Target.Timeout.Ping, config.Target.Timeout.Forward, conn)
	}
}

func handleRequest(protocol string, targetHost string, timeout, scaleUpTimeout int, conn net.Conn) {
	log.Info().Msg("new client")

	proxy, err := net.DialTimeout(protocol, targetHost, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Info().Msg("notify failed request")
		readyChan := target.NotifyFailedRequest(scaleUpTimeout)

		for range readyChan {
		}

		proxy, err = net.DialTimeout(protocol, targetHost, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			log.Panic().
				Err(err).
				Str("target", targetHost).
				Msg("Failed dialing Target")
		}
	}

	log.Info().
		Str("target", targetHost).
		Msg("Proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
