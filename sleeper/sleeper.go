package sleeper

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Sleeper struct {
	offset      int
	until       *int64
	triggerFunc func()
	run         *bool
	triggered   *bool
}

func NewSleeper(offset int, triggerFunc func()) Sleeper {
	until := time.Now().Add(time.Millisecond * time.Duration(offset)).UnixMilli()
	triggered := false
	run := true
	s := Sleeper{
		offset,
		&until,
		triggerFunc,
		&run,
		&triggered,
	}
	defer func() { go s.start() }()
	return s
}

func (s *Sleeper) start() {
	go s.sleepRoutine()
}

func (s *Sleeper) sleepRoutine() {
out:
	for {
		log.Debug().Bool("triggered", *s.triggered).Str("triggerAt", time.Now().String()).Msg("sleeper woke up")
		now := time.Now().UnixMilli()
		if now >= *s.until && !*s.triggered {
			log.Info().Msg("sleeper triggering func")
			s.triggerFunc()
			*s.triggered = true
		}

		if !*s.run {
			break out
		}
		log.Debug().Msg("sleeper sleeping")
		time.Sleep(time.Second)
	}
}

func (s *Sleeper) Stop() {
	*s.run = false
}

func (s *Sleeper) Notify() {
	log.Info().Msg("sleeper notified")
	foo := time.Now().Add(time.Millisecond * time.Duration(s.offset)).UnixMilli()
	s.until = &foo
	fl := false
	s.triggered = &fl
}
