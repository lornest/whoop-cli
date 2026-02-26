package screens

import "github.com/lornest/whoop-cli/pkg/whoop"

// Typed messages for data fetching results.

type BodyMsg struct {
	Data *whoop.BodyMeasurement
	Err  error
}

type CyclesMsg struct {
	Data *whoop.CycleResponse
	Err  error
}

type RecoveryMsg struct {
	Data *whoop.RecoveryResponse
	Err  error
}

type SleepMsg struct {
	Data *whoop.SleepResponse
	Err  error
}

type WorkoutsMsg struct {
	Data *whoop.WorkoutResponse
	Err  error
}
