package service

import (
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
	loc, err := time.LoadLocation(reminder.Timezone)
	if err != nil {
		return err
	}

	log.Printf("Sending alert to %d for %f,%f", reminder.ChatId, reminder.Location.Latitude, reminder.Location.Longitude)
	messageText, err := MakePrediction(shared.ParsedPredictionRequest{
		PredictionType:  predType,
		Location:        reminder.Location,
		LocationDetails: reminder.LocationDetails,
		Timezone:        reminder.Timezone,
		Date:            time.Now().In(loc),
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

	loc, err := time.LoadLocation(region.Timezone)
	if err != nil {
		return false, err
	}

	now := time.Now().In(loc)
	madeChanges := false
	_, offset := time.Now().In(loc).Zone()
	offsetHours := float64(offset) / 3600.0

	if region.NextSunriseAlert.IsZero() {
		log.Printf("Scheduling sunrise alerts for %s", region.Region)
		nextSunrise := now.Add(-time.Hour)
		sunriseDay := now
		for nextSunrise.Before(now) {
			nextSunrise, _, err = sunrisesunset.GetSunriseSunset(region.Location.Latitude, region.Location.Longitude, offsetHours, sunriseDay)
			if err != nil {
				return false, err
			}
			nextSunrise = nextSunrise.Add(-time.Hour)
			sunriseDay = sunriseDay.Add(time.Hour)
		}
		madeChanges = true
		region.NextSunriseAlert = nextSunrise.UTC()
		log.Printf("Sunrise alert will go out at %s", region.NextSunriseAlert.String())
	}

	if region.NextSunsetAlert.IsZero() {
		log.Printf("Scheduling sunset alerts for %s", region.Region)
		nextSunset := now.Add(-time.Hour)
		sunsetDay := now
		for nextSunset.Before(now) {
			_, nextSunset, err = sunrisesunset.GetSunriseSunset(region.Location.Latitude, region.Location.Longitude, offsetHours, sunsetDay)
			if err != nil {
				return false, err
			}
			nextSunset = nextSunset.Add(-time.Hour)
			sunsetDay = sunsetDay.Add(time.Hour)
		}
		madeChanges = true
		region.NextSunsetAlert = nextSunset.UTC()
		log.Printf("Sunset alert will go out at %s", region.NextSunsetAlert.String())
	}

	return madeChanges, nil
}
