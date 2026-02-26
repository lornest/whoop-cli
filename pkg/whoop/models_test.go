package whoop

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + name)
	require.NoError(t, err)
	return data
}

func TestBodyMeasurement_Unmarshal(t *testing.T) {
	data := loadFixture(t, "body.json")
	var body BodyMeasurement
	require.NoError(t, json.Unmarshal(data, &body))

	assert.Equal(t, int64(12345), body.UserID)
	assert.InDelta(t, 1.8288, body.HeightMeter, 0.0001)
	assert.InDelta(t, 81.6466, body.WeightKilogram, 0.0001)
	assert.Equal(t, 195, body.MaxHeartRate)
}

func TestCycleResponse_Unmarshal(t *testing.T) {
	data := loadFixture(t, "cycles.json")
	var resp CycleResponse
	require.NoError(t, json.Unmarshal(data, &resp))

	assert.Len(t, resp.Records, 3)
	assert.Nil(t, resp.NextToken)

	// First cycle: scored
	c := resp.Records[0]
	assert.Equal(t, int64(101), c.ID)
	assert.Equal(t, "SCORED", c.ScoreState)
	require.NotNil(t, c.Score)
	assert.InDelta(t, 12.5, c.Score.Strain, 0.01)
	assert.Equal(t, 72, c.Score.AverageHeartRate)
	assert.Equal(t, 175, c.Score.MaxHeartRate)

	// Third cycle: pending, score is nil
	pending := resp.Records[2]
	assert.Equal(t, "PENDING", pending.ScoreState)
	assert.Nil(t, pending.Score)
}

func TestRecoveryResponse_Unmarshal(t *testing.T) {
	data := loadFixture(t, "recovery.json")
	var resp RecoveryResponse
	require.NoError(t, json.Unmarshal(data, &resp))

	assert.Len(t, resp.Records, 3)

	r := resp.Records[0]
	assert.Equal(t, int64(101), r.CycleID)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", r.SleepID)
	assert.Equal(t, "SCORED", r.ScoreState)
	require.NotNil(t, r.Score)
	assert.InDelta(t, 78.0, r.Score.RecoveryScore, 0.01)
	assert.InDelta(t, 52.0, r.Score.RestingHeartRate, 0.01)
	assert.InDelta(t, 65.3, r.Score.HRVRmssdMilli, 0.01)
	assert.InDelta(t, 97.5, r.Score.SpO2Percentage, 0.01)
	assert.InDelta(t, 33.2, r.Score.SkinTempCelsius, 0.01)
}

func TestSleepResponse_Unmarshal(t *testing.T) {
	data := loadFixture(t, "sleep.json")
	var resp SleepResponse
	require.NoError(t, json.Unmarshal(data, &resp))

	assert.Len(t, resp.Records, 2)

	s := resp.Records[0]
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", s.ID)
	assert.False(t, s.Nap)
	assert.Equal(t, "SCORED", s.ScoreState)
	require.NotNil(t, s.Score)

	assert.InDelta(t, 92.0, s.Score.SleepPerformancePercentage, 0.01)
	assert.InDelta(t, 90.9, s.Score.SleepEfficiencyPercentage, 0.01)
	assert.InDelta(t, 15.2, s.Score.RespiratoryRate, 0.01)

	stages := s.Score.StageSummary
	assert.Equal(t, 29700000, stages.TotalInBedTimeMilli)
	assert.Equal(t, 2700000, stages.TotalAwakeTimeMilli)
	assert.Equal(t, 12600000, stages.TotalLightSleepTimeMilli)
	assert.Equal(t, 7200000, stages.TotalSlowWaveSleepTimeMilli)
	assert.Equal(t, 7200000, stages.TotalREMSleepTimeMilli)
	assert.Equal(t, 4, stages.SleepCycleCount)
}

func TestWorkoutResponse_Unmarshal(t *testing.T) {
	data := loadFixture(t, "workouts.json")
	var resp WorkoutResponse
	require.NoError(t, json.Unmarshal(data, &resp))

	assert.Len(t, resp.Records, 2)
	require.NotNil(t, resp.NextToken)
	assert.Equal(t, "abc123", *resp.NextToken)

	w := resp.Records[0]
	assert.Equal(t, "d4e5f6a7-b8c9-0123-defg-456789012345", w.ID)
	assert.Equal(t, 1, w.SportID)
	assert.Equal(t, "Running", w.SportName)
	assert.Equal(t, "SCORED", w.ScoreState)
	require.NotNil(t, w.Score)
	assert.InDelta(t, 14.2, w.Score.Strain, 0.01)
	assert.Equal(t, 145, w.Score.AverageHeartRate)
	assert.Equal(t, 178, w.Score.MaxHeartRate)
	assert.InDelta(t, 8045.0, w.Score.DistanceMeter, 0.01)
	assert.InDelta(t, 52.3, w.Score.AltitudeGainMeter, 0.01)

	// Zone durations
	assert.Equal(t, 600000, w.Score.ZoneDuration.ZoneOneMilli)
	assert.Equal(t, 1500000, w.Score.ZoneDuration.ZoneThreeMilli)
}
