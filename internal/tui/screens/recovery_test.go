package screens

import (
	"testing"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryModel_Loading(t *testing.T) {
	m := RecoveryModel{loading: true, width: 80, height: 24}
	view := m.View()
	assert.Contains(t, view, "Loading")
}

func TestRecoveryModel_WithData(t *testing.T) {
	m := RecoveryModel{
		width:  100,
		height: 30,
		recovery: &whoop.RecoveryResponse{
			Records: []whoop.Recovery{
				{ScoreState: "SCORED", Score: &whoop.RecoveryScore{
					RecoveryScore: 78, RestingHeartRate: 52, HRVRmssdMilli: 65.3,
				}},
				{ScoreState: "SCORED", Score: &whoop.RecoveryScore{
					RecoveryScore: 45, RestingHeartRate: 58, HRVRmssdMilli: 42.1,
				}},
			},
		},
	}
	view := m.View()
	assert.Contains(t, view, "Recovery")
	assert.Contains(t, view, "78%")
	assert.Contains(t, view, "52 bpm")
}

func TestRecoveryModel_Empty(t *testing.T) {
	m := RecoveryModel{
		width:    80,
		height:   24,
		recovery: &whoop.RecoveryResponse{Records: []whoop.Recovery{}},
	}
	view := m.View()
	assert.Contains(t, view, "No recovery data")
}

func TestRecoveryModel_Error(t *testing.T) {
	m := RecoveryModel{width: 80, height: 24, err: assert.AnError}
	view := m.View()
	assert.Contains(t, view, "Error")
}
