package service

import (
	"fmt"
	"time"

	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/kelvins/sunrisesunset"
)

func MakePrediction(req shared.ParsedPredictionRequest) (string, error) {
	sunrise, sunset, err := sunrisesunset.GetSunriseSunset(req.Location.Latitude, req.Location.Longitude, req.Offset, req.Date)
	if err != nil {
		return "", err
	}

	var event string
	var eventTime time.Time

	if req.PredictionType == shared.PredictionTypeSunrise {
		event = "Sunrise"
		eventTime = sunrise.In(time.FixedZone("local", int(req.Offset*3600)))
	} else if req.PredictionType == shared.PredictionTypeSunset {
		event = "Sunset"
		eventTime = sunset.In(time.FixedZone("local", int(req.Offset*3600)))
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
