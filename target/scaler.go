package target

import (
	"sync"

	"github.com/fluktuid/sero/util"
)

type Scaler interface {
	ScaleUP()
	ScaleDown()
	Status() util.Status
	StatusReady(int) *sync.WaitGroup
}
