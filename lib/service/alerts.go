package service

import (
	"fmt"
	"log"
	"time"

	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
	"github.com/kelvins/sunrisesunset"
)

func RunAlertCycle(rs shared.ReminderStorage, aq shared.AlertQueue) error {
	regions, err := rs.GetRegions()
	if err != nil {
		return err
	}

	for i, region := range regions {
		madeChanges, err := scheduleNextAlerts(&region)
		if err != nil {
			return err
		}
		if madeChanges {
			err = rs.UpdateRegion(region)
			if err != nil {
				return err
			}
			regions[i] = region
		}
	}

	now := time.Now().UTC()
	for _, region := range regions {
		sentAlerts := false
		if !region.NextSunriseAlert.IsZero() && region.NextSunriseAlert.Before(now) {
			err = enqueueAlerts(rs, aq, shared.PredictionTypeSunrise, region)
			if err != nil {
				return err
			}
			region.NextSunriseAlert = time.Time{}
			sentAlerts = true
		}
		if !region.NextSunsetAlert.IsZero() && region.NextSunsetAlert.Before(now) {
			err = enqueueAlerts(rs, aq, shared.PredictionTypeSunset, region)
			if err != nil {
				return err
			}
			region.NextSunsetAlert = time.Time{}
			sentAlerts = true
		}
		if sentAlerts {
			err = rs.UpdateRegion(region)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SendAlert(client telegram.Telegram, predType string, reminder shared.Reminder) error {
	log.Printf("Sending alert to %d for %f,%f", reminder.ChatId, reminder.Location.Latitude, reminder.Location.Longitude)
	messageText, err := MakePrediction(shared.ParsedPredictionRequest{
		PredictionType:  predType,
		Location:        reminder.Location,
		LocationDetails: reminder.LocationDetails,
		Offset:          reminder.UTCOffset,
		Date:            time.Now(),
	})
	if err != nil {
		return err
	}

	_, err = client.SendMessage(telegram.OutgoingMessage{
		ChatId: reminder.ChatId,
		Message: telegram.Message{
			Text: messageText,
		},
	})

	return err
}

func enqueueAlerts(rs shared.ReminderStorage, aq shared.AlertQueue, alertType string, region shared.Region) error {
	log.Printf("Enqueueing %s alerts for %s", alertType, region.Region)
	reminders, err := rs.GetRemindersInRegion(region.Region)
	if err != nil {
		return err
	}
	log.Printf("Will enqueue %d alerts", len(reminders))

	err = aq.EnqueueAlerts(alertType, reminders)
	if err != nil {
		return err
	}

	return nil
}

func scheduleNextAlerts(region *shared.Region) (bool, error) {
	var err error
	now := time.Now().In(time.FixedZone("myzone", int(region.UTCOffset*3600)))
	madeChanges := false

	if region.NextSunriseAlert.IsZero() {
		log.Printf("Scheduling sunrise alerts for %s", region.Region)
		nextSunrise := now.Add(-time.Hour)
		sunriseDay := now
		for nextSunrise.Before(now) {
			fmt.Println(nextSunrise, sunriseDay, now)
			nextSunrise, _, err = sunrisesunset.GetSunriseSunset(region.Location.Latitude, region.Location.Longitude, region.UTCOffset, sunriseDay)
			if err != nil {
				return false, err
			}
			sunriseDay = sunriseDay.Add(time.Hour)
		}
		madeChanges = true
		region.NextSunriseAlert = nextSunrise.Add(-time.Hour).UTC()
		log.Printf("Sunrise alert will go out at %s", region.NextSunriseAlert.String())
	}

	if region.NextSunsetAlert.IsZero() {
		log.Printf("Scheduling sunset alerts for %s", region.Region)
		nextSunset := now.Add(-time.Hour)
		sunsetDay := now
		for nextSunset.Before(now) {
			fmt.Println(nextSunset, sunsetDay, now)
			_, nextSunset, err = sunrisesunset.GetSunriseSunset(region.Location.Latitude, region.Location.Longitude, region.UTCOffset, sunsetDay)
			if err != nil {
				return false, err
			}
			sunsetDay = sunsetDay.Add(time.Hour)
		}
		madeChanges = true
		region.NextSunsetAlert = nextSunset.Add(-time.Hour).UTC()
		log.Printf("Sunset alert will go out at %s", region.NextSunsetAlert.String())
	}

	return madeChanges, nil
}