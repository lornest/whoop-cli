package screens

import (
	"testing"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/stretchr/testify/assert"
)

func TestSleepModel_Loading(t *testing.T) {
	m := SleepModel{loading: true, width: 80, height: 24}
	view := m.View()
	assert.Contains(t, view, "Loading")
}

func TestSleepModel_WithData(t *testing.T) {
	m := SleepModel{
		width:  120,
		height: 40,
		sleep: &whoop.SleepResponse{
			Records: []whoop.Sleep{
				{
					Start:      "2026-02-25T22:30:00.000Z",
					ScoreState: "SCORED",
					Score: &whoop.SleepScore{
						StageSummary: whoop.StageSummary{
							TotalInBedTimeMilli:         29700000,
							TotalAwakeTimeMilli:         2700000,
							TotalLightSleepTimeMilli:    12600000,
							TotalSlowWaveSleepTimeMilli: 7200000,
							TotalREMSleepTimeMilli:      7200000,
						},
						SleepPerformancePercentage: 92,
						SleepEfficiencyPercentage:  90.9,
						RespiratoryRate:            15.2,
					},
				},
			},
		},
	}
	view := m.View()
	assert.Contains(t, view, "Sleep")
	assert.Contains(t, view, "92%")
	assert.Contains(t, view, "Performance")
	assert.Contains(t, view, "Efficiency")
}

func TestSleepModel_Empty(t *testing.T) {
	m := SleepModel{
		width:  80,
		height: 24,
		sleep:  &whoop.SleepResponse{Records: []whoop.Sleep{}},
	}
	view := m.View()
	assert.Contains(t, view, "No sleep data")
}
