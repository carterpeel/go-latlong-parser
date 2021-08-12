package glp

import (
	"errors"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"strconv"
	"strings"
)

type latLong struct {
	Latitude     float64
	Longitude    float64
	mapsClient   *maps.Client
	bCtx         context.Context
}

func NewLatLong(Latitude, Longitude string, GoogleApiKey string) (*latLong, error) {
	ctx := context.Background()
	lat, err := strconv.ParseFloat(Latitude, 64)
	if err != nil {
		return nil, err
	}
	long, err := strconv.ParseFloat(Longitude, 64)
	if err != nil {
		return nil, err
	}
	switch {
	case lat > 90:
		return nil, ErrLatitudeTooLarge
	case lat < -90:
		return nil, ErrLatitudeTooSmall
	case long > 180:
		return nil, ErrLongitudeTooLarge
	case long < -180:
		return nil, ErrLongitudeTooSmall
	}
	MapsClient, err := newGoogleClientWithTokenCheck(GoogleApiKey, ctx)
	if err != nil {
		return nil, err
	}
	return &latLong{Latitude:lat, Longitude:long, mapsClient: MapsClient, bCtx: ctx}, nil
}

func newGoogleClientWithTokenCheck(token string, ctx context.Context) (*maps.Client, error) {
	c, err := maps.NewClient(maps.WithAPIKey(token))
	if err != nil {
		return nil, err
	}
	results, err := c.Geocode(ctx, &maps.GeocodingRequest{
		Address:    "2Q93+R28, Pyongyang, North Korea",
	})
	switch {
	case err != nil && strings.Contains(err.Error(), "REQUEST_DENIED"):
		return nil, ErrBadGoogleAPIKey
	case err != nil:
		return nil, err
	}
	switch {
	case len(results) < 0:
		return nil, ErrBadGoogleAPIKey
	default:
		return c, nil
	}

}

func (l *latLong) GetAddress() (string, error) {
	results, err := l.mapsClient.ReverseGeocode(l.bCtx, &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: l.Latitude,
			Lng: l.Longitude,
		},
	})
	if err != nil {
		return "", err
	}
	switch {
	case len(results) >= 1:
		return results[0].FormattedAddress, nil
	default:
		return "", ErrNoAddressFound
	}
}

var (
	ErrBadGoogleAPIKey   = errors.New("the key provided is not a valid Google API token")
	ErrNoAddressFound    = errors.New("no address was found at the provided coordinates")
	ErrLatitudeTooLarge  = errors.New("latitude value must be less than 90")
	ErrLatitudeTooSmall  = errors.New("latitude value must be more than -90")
	ErrLongitudeTooLarge = errors.New("longitude value must be less than 180")
	ErrLongitudeTooSmall = errors.New("longitude value must be more than -180")
)