package main

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	cfg "github.com/fluktuid/sero/config"
	"github.com/fluktuid/sero/metrics"
	s "github.com/fluktuid/sero/sleeper"
	t "github.com/fluktuid/sero/target"
)

var target t.Target

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	config := cfg.LoadConfig()

	metrics.InitAsync()

	target = *t.Init(config.Target.Deployment)
	ln, err := net.Listen("tcp", config.Host)
	if err != nil {
		panic(err)
	}

	metrics.Ready(true)
	metrics.Healthy(true)

	var sleeper s.Sleeper
	if config.Target.Timeout.ScaleDown > 0 {
		notify := func() {
			target.NotifyScaleDown()
		}
		sleeper = s.NewSleeper(config.Target.Timeout.ScaleDown, notify)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		sleeper.Notify()

		timeout := time.Duration(config.Target.Timeout.Forward * int(time.Millisecond))
		scaleUp := time.Duration(config.Target.Timeout.ScaleUP * int(time.Millisecond))

		go func() {
			metrics.RecordRequest()
			success := handleRequest(config.Target.Protocol, config.Target.Host, timeout, scaleUp, conn)
			metrics.RecordRequestFinish(success)
			sleeper.Notify()
		}()
	}
}

func handleRequest(protocol string, targetHost string, timeout, scaleUpTimeout time.Duration, conn net.Conn) bool {
	log.Debug().Msg("new client")

	proxy, err := net.DialTimeout(protocol, targetHost, timeout)
	if err != nil {
		log.Info().Msg("notify failed request")
		target.
			NotifyFailedRequest(scaleUpTimeout).
			Wait()

		proxy, err = net.DialTimeout(protocol, targetHost, timeout)
		if err != nil {
			log.Warn().
				Err(err).
				Str("target", targetHost).
				Msg("Failed dialing Target")
			defer conn.Close()
			return false
		}
	}

	log.Debug().
		Str("target", targetHost).
		Msg("Proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
	return true
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
