package target

import (
	"sync"
	"time"

	"github.com/fluktuid/sero/util"
)

type Scaler interface {
	ScaleUP()
	ScaleDown()
	Status() util.Status
	StatusReady(time.Duration) *sync.WaitGroup
}
