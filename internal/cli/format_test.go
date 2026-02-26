package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestFormatRecoveryTable(t *testing.T) {
	records := []whoop.Recovery{
		{
			CreatedAt:  "2026-02-26T08:00:00.000Z",
			ScoreState: "SCORED",
			Score: &whoop.RecoveryScore{
				RecoveryScore:    78,
				RestingHeartRate: 52,
				HRVRmssdMilli:    65.3,
				SpO2Percentage:   97,
			},
		},
		{
			CreatedAt:  "2026-02-25T08:00:00.000Z",
			ScoreState: "SCORED",
			Score: &whoop.RecoveryScore{
				RecoveryScore:    45,
				RestingHeartRate: 58,
				HRVRmssdMilli:    42.1,
				SpO2Percentage:   96,
			},
		},
	}

	output := captureOutput(func() {
		err := formatOutput("table", "recovery", records)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "Date")
	assert.Contains(t, output, "Score")
	assert.Contains(t, output, "78%")
	assert.Contains(t, output, "45%")
	assert.Contains(t, output, "52 bpm")
	assert.Contains(t, output, "65.3 ms")
}

func TestFormatRecoveryText(t *testing.T) {
	records := []whoop.Recovery{
		{
			CreatedAt: "2026-02-26T08:00:00.000Z",
			Score: &whoop.RecoveryScore{
				RecoveryScore:    78,
				RestingHeartRate: 52,
				HRVRmssdMilli:    65.3,
				SpO2Percentage:   97,
				SkinTempCelsius:  33.2,
			},
		},
	}

	output := captureOutput(func() {
		err := formatOutput("text", "recovery", records)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "recovery: 78%")
	assert.Contains(t, output, "resting_heart_rate: 52 bpm")
	assert.Contains(t, output, "hrv: 65.3 ms")
}

func TestFormatRecoveryJSON(t *testing.T) {
	records := []whoop.Recovery{
		{
			CycleID:   101,
			CreatedAt: "2026-02-26T08:00:00.000Z",
			Score: &whoop.RecoveryScore{
				RecoveryScore: 78,
			},
		},
	}

	output := captureOutput(func() {
		err := formatOutput("json", "recovery", records)
		require.NoError(t, err)
	})

	var parsed []whoop.Recovery
	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
	assert.Len(t, parsed, 1)
	assert.Equal(t, int64(101), parsed[0].CycleID)
}

func TestFormatSleepTable(t *testing.T) {
	records := []whoop.Sleep{
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
				SleepEfficiencyPercentage:  91,
			},
		},
	}

	output := captureOutput(func() {
		err := formatOutput("table", "sleep", records)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "Date")
	assert.Contains(t, output, "In Bed")
	assert.Contains(t, output, "Awake")
	assert.Contains(t, output, "Light")
	assert.Contains(t, output, "Deep")
	assert.Contains(t, output, "REM")
	assert.Contains(t, output, "Perf %")
	assert.Contains(t, output, "Eff %")
	assert.Contains(t, output, "8h 15m")
	assert.Contains(t, output, "45m")
	assert.Contains(t, output, "3h 30m")
	assert.Contains(t, output, "2h 0m")
	assert.Contains(t, output, "92%")
	assert.Contains(t, output, "91%")
}

func TestFormatWorkoutsTable(t *testing.T) {
	records := []whoop.Workout{
		{
			Start:      "2026-02-25T17:00:00.000Z",
			End:        "2026-02-25T18:05:00.000Z",
			SportID:    1,
			ScoreState: "SCORED",
			Score: &whoop.WorkoutScore{
				Strain:           14.2,
				AverageHeartRate: 145,
				MaxHeartRate:     178,
			},
		},
	}

	output := captureOutput(func() {
		err := formatOutput("table", "workouts", records)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "Running")
	assert.Contains(t, output, "14.2")
	assert.Contains(t, output, "145")
	assert.Contains(t, output, "1h 5m")
}

func TestFormatCyclesTable(t *testing.T) {
	records := []whoop.Cycle{
		{
			Start:      "2026-02-25T06:00:00.000Z",
			ScoreState: "SCORED",
			Score: &whoop.CycleScore{
				Strain:           12.5,
				Kilojoule:        8500,
				AverageHeartRate: 72,
				MaxHeartRate:     175,
			},
		},
	}

	output := captureOutput(func() {
		err := formatOutput("table", "cycles", records)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "12.5")
	assert.Contains(t, output, "8500")
	assert.Contains(t, output, "72")
	assert.Contains(t, output, "175")
}

func TestFormatProfileTable(t *testing.T) {
	body := &whoop.BodyMeasurement{
		HeightMeter:    1.8288,
		WeightKilogram: 81.6466,
		MaxHeartRate:   195,
	}

	output := captureOutput(func() {
		err := formatOutput("table", "profile", body)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "Height")
	assert.Contains(t, output, "Weight")
	assert.Contains(t, output, "195 bpm")
}

func TestFormatEmptyRecovery(t *testing.T) {
	output := captureOutput(func() {
		err := formatOutput("table", "recovery", []whoop.Recovery{})
		require.NoError(t, err)
	})

	assert.Contains(t, output, "No recovery data")
}

func TestFormatUnknownFormat(t *testing.T) {
	err := formatOutput("yaml", "recovery", []whoop.Recovery{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format")
}

func TestParseDate(t *testing.T) {
	assert.Equal(t, "Feb 26", parseDate("2026-02-26T08:00:00.000Z"))
	assert.Equal(t, "bad-date", parseDate("bad-date"))
}

func TestParseDatetime(t *testing.T) {
	assert.Equal(t, "Feb 25 17:00", parseDatetime("2026-02-25T17:00:00.000Z"))
	assert.Equal(t, "bad-date", parseDatetime("bad-date"))
}
