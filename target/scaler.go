package target

import (
	"github.com/fluktuid/sero/util"
)

type Scaler interface {
	ScaleUP()
	Status() util.Status
	StatusReadyChan(int) <-chan util.Void
}
