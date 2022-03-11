package shared

import (
	"fmt"
	"math"
	"strings"
)

func (r Reminder) GetRegionKey() string {
	return fmt.Sprintf("%d,%d",
		int(math.Round(r.Location.Latitude)),
		int(math.Round(r.Location.Longitude)),
	)
}

func (r Reminder) GetRegion() Coordinates {
	return Coordinates{Latitude: math.Round(r.Location.Latitude), Longitude: math.Round(r.Location.Longitude)}
}

func (l Coordinates) IsZero() bool {
	return l.Latitude == 0 && l.Longitude == 0
}

func (ld LocationDetails) String() string {
	sb := strings.Builder{}

	if ld.City != "" {
		sb.WriteString(ld.City)
	}

	if ld.State != "" {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ld.State)
	}

	return sb.String()
}
