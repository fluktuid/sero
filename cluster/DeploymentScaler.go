package cluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

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

func (d DeploymentScaler) StatusReady(timeout time.Duration) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := waitForDeployRunning(clientSet, namespace, d.targetDeploymentName, timeout)
		wg.Done()
		log.Info().Time("now", time.Now()).Msg("pod ready")
		if err != nil {
			log.Warn().Err(err).Msg("Error while waiting for ready deployment")
		}
	}()
	return &wg
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func waitForDeployRunning(c kubernetes.Interface, namespace, deployName string, timeout time.Duration) error {
	return wait.PollImmediate(100*time.Millisecond, timeout, hasDeployReadyPod(c, deployName, namespace))
}

// return a condition function that indicates whether the given pod is
// currently running
func hasDeployReadyPod(c kubernetes.Interface, deployName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		fmt.Printf(".") // progress bar!
		deploy, err := c.AppsV1().Deployments(namespace).Get(context.TODO(), deployName, v1.GetOptions{})
		if err != nil {
			return false, err
		}

		if deploy.Status.ReadyReplicas > 0 {
			return true, nil
		} else if deploy.Status.ReadyReplicas == 0 && *deploy.Spec.Replicas == 0 {
			return false, nil // todo: add condition
		}
		return false, nil
	}
}
