package helm

import (
	apiv1 "k8s.io/api/core/v1"
)

// These functions are adapted from the "kubernetes" repository's file
//
//    kubernetes/pkg/api/v1/pod/util.go
//
// where they rely upon the API types specific to that repository. Here we recast them to operate
// upon the type from the "client-go" repository instead.

// isPodReady returns true if a pod is ready; false otherwise.
func isPodReady(pod *apiv1.Pod) bool {
	return isPodReadyConditionTrue(pod.Status)
}

// isPodReady retruns true if a pod is ready; false otherwise.
func isPodReadyConditionTrue(status apiv1.PodStatus) bool {
	condition := getPodReadyCondition(status)
	return condition != nil && condition.Status == apiv1.ConditionTrue
}

// getPodReadyCondition extracts the pod ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func getPodReadyCondition(status apiv1.PodStatus) *apiv1.PodCondition {
	_, condition := getPodCondition(&status, apiv1.PodReady)
	return condition
}

// getPodCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func getPodCondition(status *apiv1.PodStatus, conditionType apiv1.PodConditionType) (int, *apiv1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}
