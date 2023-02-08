package target

import (
	"github.com/rs/zerolog/log"

	"github.com/fluktuid/sero/cluster"
	"github.com/fluktuid/sero/util"
)

type Target struct {
	deployment string
	scaler     Scaler
}

func Init(deployment string) *Target {
	scaler := cluster.InitDeploymentScaler(deployment)
	return &Target{
		deployment: deployment,
		scaler:     scaler,
	}
}

func (t *Target) Status() util.Status {
	return t.scaler.Status()
}

func (t *Target) Deployment() string {
	return t.deployment
}

func (t *Target) NotifyFailedRequest(chanTimeout int) <-chan util.Void {
	if t.scaler.Status() != util.StatusUpscaling {
		log.Info().Msg("scaling up")
		t.scaler.ScaleUP()
	}

	// returns 'continue' chan
	return t.scaler.StatusReadyChan(chanTimeout)
}

func (t *Target) NotifyScaleDown() {
	if t.scaler.Status() != util.StatusDownscaling {
		log.Info().Msg("scaling down")
		t.scaler.ScaleDown()
	}
}
