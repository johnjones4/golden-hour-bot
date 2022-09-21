package shared

const (
	StateIdle                  = "idle"
	StateImmediateNeedLocation = "immediate-needlocation"
	StateImmediateNeedDate     = "immediate-needdate"
	StateRemindNeedLocation    = "remind-needlocation"
)

const (
	DefaultState = StateIdle
)

const (
	MessageWelcome           = "Hello! I can tell you the sunrise or sunset time for any location and I can send you alerts when it is sunset or sunrise for a given area. Use the commands to get started."
	MessageConfused          = "I'm sorry I don't understand."
	MessageShareLocation     = "Where are you?"
	MessageShareTime         = "What day would you like to know this for?"
	MessageLocationThanks    = "Thanks for sharing your location."
	MessageReminderSet       = "Ok. I'll send you alerts when it's sunrise and sunset for that location."
	MessageDuplicateReminder = "You already have a similar reminder established."
	MessageBadDate           = "I don't understand that date/time string."
	MessageNoLocation        = "You did not share a location with me."
)

const (
	ButtonShareLocation = "Share my location"
)

const (
	PredictionTypeSunrise = "sunrise"
	PredictionTypeSunset  = "sunset"
)
