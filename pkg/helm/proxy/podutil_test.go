package proxy

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
)

func TestGetPodCondition(t *testing.T) {
	testCases := []struct {
		status        *apiv1.PodStatus
		conditionType apiv1.PodConditionType

		expectedCode      int
		expectedCondition *apiv1.PodCondition
	}{
		{
			expectedCode: -1,
		},
		{
			status: &apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type: apiv1.PodReady,
					},
				},
			},
			conditionType: apiv1.PodReady,
			expectedCode:  0,
			expectedCondition: &apiv1.PodCondition{
				Type: apiv1.PodReady,
			},
		},
		{
			status: &apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type: apiv1.PodReady,
					},
				},
			},
			conditionType:     apiv1.PodScheduled,
			expectedCode:      -1,
			expectedCondition: nil,
		},
	}

	for _, testCase := range testCases {
		code, condition := getPodCondition(testCase.status,
			testCase.conditionType)

		if code != testCase.expectedCode {
			t.Errorf("Wrong code expected %d actual %d",
				testCase.expectedCode, code)
		}

		if testCase.expectedCondition != nil && condition.Type != testCase.expectedCondition.Type {
			t.Errorf("Wrong conditions expected %v actual %v",
				testCase.expectedCondition, condition)
		}
	}
}

func TestGetPodReadyCondition(t *testing.T) {
	testCases := []struct {
		status            apiv1.PodStatus
		expectedCondition apiv1.PodCondition
	}{
		{
			status: apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type: apiv1.PodScheduled,
					},
				},
			},
		},
		{
			status: apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type: apiv1.PodScheduled,
					},
					{
						Type: apiv1.PodReady,
					},
					{
						Type: apiv1.PodInitialized,
					},
				},
			},
			expectedCondition: apiv1.PodCondition{
				Type: apiv1.PodReady,
			},
		},
	}

	for _, testCase := range testCases {
		actual := getPodReadyCondition(testCase.status)

		if actual != nil && actual.Type != testCase.expectedCondition.Type {
			t.Errorf("Wrong condition expected %v actual %v",
				testCase.expectedCondition, actual)
		}
	}
}

func TestIsPodReadyConditionTrue(t *testing.T) {
	testCases := []struct {
		status   apiv1.PodStatus
		expected bool
	}{
		{
			status: apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type: apiv1.PodScheduled,
					},
					{
						Type: apiv1.PodInitialized,
					},
				},
			},
			expected: false,
		},
		{
			status: apiv1.PodStatus{
				Conditions: []apiv1.PodCondition{
					{
						Type:   apiv1.PodReady,
						Status: apiv1.ConditionTrue,
					},
				},
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		actual := isPodReadyConditionTrue(testCase.status)

		if actual != testCase.expected {
			t.Errorf("Wrong result expected %v actual %v",
				testCase.expected, actual)
		}
	}
}

func TestIsPodReadyCondition(t *testing.T) {
	testCases := []struct {
		pod      apiv1.Pod
		expected bool
	}{
		{
			pod: apiv1.Pod{
				Status: apiv1.PodStatus{
					Conditions: []apiv1.PodCondition{
						{
							Type: apiv1.PodScheduled,
						},
						{
							Type: apiv1.PodInitialized,
						},
					},
				},
			},
			expected: false,
		},
		{
			pod: apiv1.Pod{
				Status: apiv1.PodStatus{
					Conditions: []apiv1.PodCondition{
						{
							Type:   apiv1.PodReady,
							Status: apiv1.ConditionTrue,
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		actual := isPodReady(&testCase.pod)

		if actual != testCase.expected {
			t.Errorf("Wrong result expected %v actual %v",
				testCase.expected, actual)
		}
	}
}
