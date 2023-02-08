package sleeper

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Sleeper struct {
	offset      int
	until       *time.Time
	triggerFunc func()
	run         *bool
	triggered   *bool
}

func NewSleeper(offset int, triggerFunc func()) Sleeper {
	until := time.Now().Add(time.Millisecond * time.Duration(offset))
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
		log.Debug().Bool("triggered", *s.triggered).Str("triggerAt", s.until.Local().String()).Msg("sleeper woke up")
		now := time.Now()
		if !now.Before(*s.until) && !*s.triggered {
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
	*s.until = time.Now().Add(time.Millisecond * time.Duration(s.offset))
	*s.triggered = false
}
