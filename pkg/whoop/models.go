package whoop

// BodyMeasurement represents the user's body measurements.
type BodyMeasurement struct {
	UserID         int64   `json:"user_id"`
	HeightMeter    float64 `json:"height_meter"`
	WeightKilogram float64 `json:"weight_kilogram"`
	MaxHeartRate   int     `json:"max_heart_rate"`
}

// CycleScore contains scored metrics for a physiological cycle.
type CycleScore struct {
	Strain           float64 `json:"strain"`
	Kilojoule        float64 `json:"kilojoule"`
	AverageHeartRate int     `json:"average_heart_rate"`
	MaxHeartRate     int     `json:"max_heart_rate"`
}

// Cycle represents a single physiological cycle (typically one day).
type Cycle struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	Start          string      `json:"start"`
	End            string      `json:"end"`
	TimezoneOffset string      `json:"timezone_offset"`
	ScoreState     string      `json:"score_state"`
	Score          *CycleScore `json:"score"`
}

// CycleResponse is the paginated response for cycles.
type CycleResponse struct {
	Records   []Cycle `json:"records"`
	NextToken *string `json:"next_token"`
}

// RecoveryScore contains scored recovery metrics.
type RecoveryScore struct {
	UserCalibrating  bool    `json:"user_calibrating"`
	RecoveryScore    float64 `json:"recovery_score"`
	RestingHeartRate float64 `json:"resting_heart_rate"`
	HRVRmssdMilli    float64 `json:"hrv_rmssd_milli"`
	SpO2Percentage   float64 `json:"spo2_percentage"`
	SkinTempCelsius  float64 `json:"skin_temp_celsius"`
}

// Recovery represents a single recovery record.
type Recovery struct {
	CycleID    int64          `json:"cycle_id"`
	SleepID    string         `json:"sleep_id"`
	UserID     int64          `json:"user_id"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	ScoreState string         `json:"score_state"`
	Score      *RecoveryScore `json:"score"`
}

// RecoveryResponse is the paginated response for recovery data.
type RecoveryResponse struct {
	Records   []Recovery `json:"records"`
	NextToken *string    `json:"next_token"`
}

// StageSummary contains sleep stage durations in milliseconds.
type StageSummary struct {
	TotalInBedTimeMilli         int `json:"total_in_bed_time_milli"`
	TotalAwakeTimeMilli         int `json:"total_awake_time_milli"`
	TotalNoDataTimeMilli        int `json:"total_no_data_time_milli"`
	TotalLightSleepTimeMilli    int `json:"total_light_sleep_time_milli"`
	TotalSlowWaveSleepTimeMilli int `json:"total_slow_wave_sleep_time_milli"`
	TotalREMSleepTimeMilli      int `json:"total_rem_sleep_time_milli"`
	SleepCycleCount             int `json:"sleep_cycle_count"`
	DisturbanceCount            int `json:"disturbance_count"`
}

// SleepNeeded contains sleep need calculations in milliseconds.
type SleepNeeded struct {
	BaselineMilli             int `json:"baseline_milli"`
	NeedFromSleepDebtMilli    int `json:"need_from_sleep_debt_milli"`
	NeedFromRecentStrainMilli int `json:"need_from_recent_strain_milli"`
	NeedFromRecentNapMilli    int `json:"need_from_recent_nap_milli"`
}

// SleepScore contains scored sleep metrics.
type SleepScore struct {
	StageSummary               StageSummary `json:"stage_summary"`
	SleepNeeded                SleepNeeded  `json:"sleep_needed"`
	RespiratoryRate            float64      `json:"respiratory_rate"`
	SleepPerformancePercentage float64      `json:"sleep_performance_percentage"`
	SleepConsistencyPercentage float64      `json:"sleep_consistency_percentage"`
	SleepEfficiencyPercentage  float64      `json:"sleep_efficiency_percentage"`
}

// Sleep represents a single sleep record.
type Sleep struct {
	ID         string      `json:"id"`
	UserID     int64       `json:"user_id"`
	Start      string      `json:"start"`
	End        string      `json:"end"`
	Nap        bool        `json:"nap"`
	ScoreState string      `json:"score_state"`
	Score      *SleepScore `json:"score"`
}

// SleepResponse is the paginated response for sleep data.
type SleepResponse struct {
	Records   []Sleep `json:"records"`
	NextToken *string `json:"next_token"`
}

// ZoneDuration contains HR zone durations in milliseconds.
type ZoneDuration struct {
	ZoneZeroMilli  int `json:"zone_zero_milli"`
	ZoneOneMilli   int `json:"zone_one_milli"`
	ZoneTwoMilli   int `json:"zone_two_milli"`
	ZoneThreeMilli int `json:"zone_three_milli"`
	ZoneFourMilli  int `json:"zone_four_milli"`
	ZoneFiveMilli  int `json:"zone_five_milli"`
}

// WorkoutScore contains scored workout metrics.
type WorkoutScore struct {
	Strain              float64      `json:"strain"`
	AverageHeartRate    int          `json:"average_heart_rate"`
	MaxHeartRate        int          `json:"max_heart_rate"`
	Kilojoule           float64      `json:"kilojoule"`
	PercentRecorded     float64      `json:"percent_recorded"`
	DistanceMeter       float64      `json:"distance_meter"`
	AltitudeGainMeter   float64      `json:"altitude_gain_meter"`
	AltitudeChangeMeter float64      `json:"altitude_change_meter"`
	ZoneDuration        ZoneDuration `json:"zone_duration"`
}

// Workout represents a single workout record.
type Workout struct {
	ID         string        `json:"id"`
	UserID     int64         `json:"user_id"`
	Start      string        `json:"start"`
	End        string        `json:"end"`
	SportID    int           `json:"sport_id"`
	SportName  string        `json:"sport_name"`
	ScoreState string        `json:"score_state"`
	Score      *WorkoutScore `json:"score"`
}

// WorkoutResponse is the paginated response for workout data.
type WorkoutResponse struct {
	Records   []Workout `json:"records"`
	NextToken *string   `json:"next_token"`
}

// QueryParams holds common query parameters for paginated endpoints.
type QueryParams struct {
	Start     string
	End       string
	Limit     int
	NextToken string
}
