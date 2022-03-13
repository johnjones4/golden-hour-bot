package service

import (
	"fmt"
	"time"

	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/kelvins/sunrisesunset"
)

func MakePrediction(req shared.ParsedPredictionRequest) (string, error) {
	loc, err := time.LoadLocation(req.Timezone)
	if err != nil {
		return "", err
	}

	_, offset := time.Now().In(loc).Zone()

	sunrise, sunset, err := sunrisesunset.GetSunriseSunset(req.Location.Latitude, req.Location.Longitude, float64(offset)/3600.0, req.Date)
	if err != nil {
		return "", err
	}

	var event string
	var eventTime time.Time

	if req.PredictionType == shared.PredictionTypeSunrise {
		event = "Sunrise"
		eventTime = sunrise
	} else if req.PredictionType == shared.PredictionTypeSunset {
		event = "Sunset"
		eventTime = sunset
	} else {
		return "", fmt.Errorf("invalid type: %s", req.PredictionType)
	}

	text := fmt.Sprintf("%s for %s on %s will be at %s",
		event,
		req.LocationDetails.String(),
		eventTime.Format("Monday, January 2, 2006"),
		eventTime.Format("3:04 PM"),
	)

	return text, nil
}
