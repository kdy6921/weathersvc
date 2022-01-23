package utils

import (
	"math/rand"
	"time"
)

type GeoPosition struct {
	Latitude  float64
	Longitude float64
}

func getRandom(min, max float64) float64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return r1.Float64()*(max-min) + min
}

// GetCurrentGeoPosition returns the current geo position (actually it returns random coordinate)
func GetCurrentGeoPosition() GeoPosition {
	return GeoPosition{
		Latitude:  getRandom(-90, 90),
		Longitude: getRandom(-180, 180),
	}
}
