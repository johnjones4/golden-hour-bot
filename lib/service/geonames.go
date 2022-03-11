package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/johnjones4/golden-hour-bot/lib/shared"
)

type geonamesResponse struct {
	GmtOffset int `json:"gmtOffset"`
}

type GeoNames struct {
	Username string
}

func (g GeoNames) GetUTCOffset(coords shared.Coordinates) (float64, error) {
	params := make(url.Values)
	params.Set("lat", fmt.Sprint(coords.Latitude))
	params.Set("lng", fmt.Sprint(coords.Longitude))
	params.Set("username", g.Username)
	res, err := http.Get("http://api.geonames.org/timezoneJSON?" + params.Encode())
	if err != nil {
		return 0, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var body geonamesResponse
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return 0, err
	}

	return float64(body.GmtOffset), nil
}
