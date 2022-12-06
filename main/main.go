package main

import (
	"io"
	"net"

	"github.com/rs/zerolog/log"

	cfg "github.com/fluktuid/sero/config"
)

func main() {
	config := cfg.LoadConfig()
	ln, err := net.Listen(config.Protocol, config.Host)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(config.Protocol, config.TargetHost, conn)
	}
}

func handleRequest(protocol string, target string, conn net.Conn) {
	log.Info().Msg("new client")

	proxy, err := net.Dial(protocol, target)
	if err != nil {
		log.Panic().
			Err(err).
			Str("target", target).
			Msg("Failed dialing Target")
	}

	log.Info().
		Str("target", target).
		Msg("Proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
