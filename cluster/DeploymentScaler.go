package cluster

import (
	"context"
	"time"

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
		clientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, v1.UpdateOptions{})
	}
}

func (d DeploymentScaler) Status() util.Status {
	deploy, _ := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), d.targetDeploymentName, v1.GetOptions{})
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

func (d DeploymentScaler) StatusReadyChan(timeoutMillis int) <-chan util.Void {
	chn := make(chan util.Void)
	go func() {
		replicas := d.readyReplicas(d.targetDeploymentName)
		limit := timeoutMillis / 4
		for replicas < 1 && limit > 0 {
			time.Sleep(time.Duration(250) * time.Millisecond)
			replicas = d.readyReplicas(d.targetDeploymentName)
			limit--
		}
		close(chn)
	}()
	return chn
}

func (d DeploymentScaler) readyReplicas(deployment string) int {
	deploy, _ := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, v1.GetOptions{})
	return int(deploy.Status.ReadyReplicas)
}
