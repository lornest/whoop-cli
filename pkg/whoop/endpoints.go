package whoop

const (
	BaseURL = "https://api.prod.whoop.com"

	PathBodyMeasurement = "/developer/v2/user/measurement/body"
	PathCycles          = "/developer/v2/cycle"
	PathRecovery        = "/developer/v2/recovery"
	PathSleep           = "/developer/v2/activity/sleep"
	PathWorkouts        = "/developer/v2/activity/workout"

	AuthURL  = "https://api.prod.whoop.com/oauth/oauth2/auth"
	TokenURL = "https://api.prod.whoop.com/oauth/oauth2/token"

	RedirectURI = "http://localhost:8080/callback"

	Scopes = "read:body_measurement read:cycles read:recovery read:sleep read:workout offline"
)
