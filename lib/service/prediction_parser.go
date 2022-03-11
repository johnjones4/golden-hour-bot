package service

import (
	"time"

	"github.com/codingsince1985/geo-golang"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/olebedev/when"
)

type PredictionRequestParser struct {
	GeoNames   GeoNames
	DateParser *when.Parser
	Geocoder   geo.Geocoder
}

func (p *PredictionRequestParser) geocode(desc string) (shared.Coordinates, error) {
	location, err := p.Geocoder.Geocode(desc)
	if err != nil {
		return shared.Coordinates{}, err
	}
	return shared.Coordinates{Latitude: location.Lat, Longitude: location.Lng}, nil
}

func (p *PredictionRequestParser) parseDate(dateStr string, offset float64) (time.Time, error) {
	base := time.Now().In(time.FixedZone("myzone", int(offset*3600)))
	res, err := p.DateParser.Parse(dateStr, base)
	if err != nil {
		return time.Time{}, err
	}
	if res == nil {
		return time.Time{}, shared.ErrorTimeNotParsable(dateStr)
	}
	return res.Time, nil
}

func (p *PredictionRequestParser) NewParsedPredictionRequest(req shared.PredictionRequest) (shared.ParsedPredictionRequest, error) {
	var c shared.Coordinates
	var ld shared.LocationDetails
	var err error

	if !req.Location.IsZero() {
		c = req.Location
	} else if req.LocationString != "" {
		c, err = p.geocode(req.LocationString)
		if err != nil {
			return shared.ParsedPredictionRequest{}, err
		}
	} else {
		return shared.ParsedPredictionRequest{}, shared.ErrorNoLocation()
	}

	address, err := p.Geocoder.ReverseGeocode(c.Latitude, c.Longitude)
	if err != nil {
		return shared.ParsedPredictionRequest{}, err
	}
	if address != nil {
		ld.City = address.City
		ld.State = address.State
	}

	offset, err := p.GeoNames.GetUTCOffset(c)
	if err != nil {
		return shared.ParsedPredictionRequest{}, err
	}

	date, err := p.parseDate(req.When, offset)
	if err != nil {
		return shared.ParsedPredictionRequest{}, err
	}

	return shared.ParsedPredictionRequest{
		PredictionType:  req.PredictionType,
		Location:        c,
		LocationDetails: ld,
		Offset:          offset,
		Date:            date,
	}, nil
}
