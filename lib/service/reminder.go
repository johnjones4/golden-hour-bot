package service

import (
	"github.com/bradfitz/latlong"
	"github.com/codingsince1985/geo-golang"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
)

func NewReminder(geocoder geo.Geocoder, req shared.RemindRequest) (shared.Reminder, error) {
	reminder := shared.Reminder{
		ChatId: req.ChatId,
	}

	if !req.Location.IsZero() {
		reminder.Location = req.Location
	} else {
		if req.LocationString == "" {
			return shared.Reminder{}, shared.ErrorNoLocation()
		}
		location, err := geocoder.Geocode(req.LocationString)
		if err != nil {
			return shared.Reminder{}, err
		}
		reminder.Location.Latitude = location.Lat
		reminder.Location.Longitude = location.Lng
	}

	address, err := geocoder.ReverseGeocode(reminder.Location.Latitude, reminder.Location.Longitude)
	if err != nil {
		return shared.Reminder{}, err
	}
	if address != nil {
		reminder.LocationDetails.City = address.City
		reminder.LocationDetails.State = address.State
	}

	reminder.Timezone = latlong.LookupZoneName(reminder.Location.Latitude, reminder.Location.Longitude)

	return reminder, nil
}
