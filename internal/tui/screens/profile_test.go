package screens

import (
	"testing"

	"github.com/lornest/whoop-cli/pkg/whoop"
	"github.com/stretchr/testify/assert"
)

func TestProfileModel_Loading(t *testing.T) {
	m := ProfileModel{loading: true, width: 80, height: 24}
	view := m.View()
	assert.Contains(t, view, "Loading")
}

func TestProfileModel_WithData(t *testing.T) {
	m := ProfileModel{
		width:  120,
		height: 30,
		body: &whoop.BodyMeasurement{
			HeightMeter:    1.8288,
			WeightKilogram: 81.6466,
			MaxHeartRate:   195,
		},
	}
	view := m.View()
	assert.Contains(t, view, "Profile")
	assert.Contains(t, view, "Height")
	assert.Contains(t, view, "Weight")
	assert.Contains(t, view, "Max HR")
	assert.Contains(t, view, "195 bpm")
}

func TestProfileModel_Error(t *testing.T) {
	m := ProfileModel{width: 80, height: 24, err: assert.AnError}
	view := m.View()
	assert.Contains(t, view, "Error")
}

func TestProfileModel_NilBody(t *testing.T) {
	m := ProfileModel{width: 80, height: 24}
	view := m.View()
	assert.Contains(t, view, "No profile data")
}
