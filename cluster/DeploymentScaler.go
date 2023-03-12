package cluster

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluktuid/sero/util"
)

type DeploymentScaler struct {
	targetDeploymentName string
}

func InitDeploymentScaler(deployment string) *DeploymentScaler {
	return &DeploymentScaler{
		targetDeploymentName: deployment,
	}
}

func (d DeploymentScaler) ScaleUP() {
	deploy, _ := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), d.targetDeploymentName, v1.GetOptions{})
	if *deploy.Spec.Replicas < 1 {
		one := int32(1)
		deploy.Spec.Replicas = &one
		_, err := clientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, v1.UpdateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("error getting kubernetes stuff")
		}
	}
}

func (d DeploymentScaler) ScaleDown() {
	deploy, _ := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), d.targetDeploymentName, v1.GetOptions{})
	if *deploy.Spec.Replicas > 0 {
		zero := int32(0)
		deploy.Spec.Replicas = &zero
		_, err := clientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, v1.UpdateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("error getting kubernetes stuff")
		}
	}
}

func (d DeploymentScaler) Status() util.Status {
	deploy, err := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), d.targetDeploymentName, v1.GetOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("error getting kubernetes stuff")
	}
	if deploy.Status.ReadyReplicas == 0 && *deploy.Spec.Replicas == 0 {
		return util.StatusDown
	} else if deploy.Status.ReadyReplicas > 0 && *deploy.Spec.Replicas > 0 {
		return util.StatusUp
	} else if deploy.Status.ReadyReplicas > 0 && *deploy.Spec.Replicas == 0 {
		return util.StatusDownscaling
	} else if deploy.Status.ReadyReplicas == 0 && *deploy.Spec.Replicas < 0 {
		return util.StatusUpscaling
	}
	return util.StatusDown
}

func (d DeploymentScaler) StatusReady(timeoutMillis int) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		desired, ready := d.replicas(d.targetDeploymentName)
		limit := timeoutMillis / 4
		for desired > ready {
			log.Info().Msg("unready replicas")
			time.Sleep(time.Duration(250) * time.Millisecond)
			desired, ready = d.replicas(d.targetDeploymentName)
			limit--
		}
		log.Info().Msg("ready replicas or waitlimit")
		wg.Done()
	}()
	return &wg
}

func (d DeploymentScaler) replicas(deployment string) (int, int) {
	deploy, _ := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, v1.GetOptions{})
	return int(deploy.Status.Replicas), int(deploy.Status.ReadyReplicas)
}
