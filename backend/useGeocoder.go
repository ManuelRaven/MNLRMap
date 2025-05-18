package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/pocketbase/pocketbase/core"
)

type Address struct {
	ID          int64   `json:"id"`
	Street      string  `json:"street"`
	HouseNumber string  `json:"house_number"`
	City        string  `json:"city"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

type MNLRAddressServerResponse struct {
	Schema    string    `json:"$schema"`
	Addresses []Address `json:"addresses"`
}

type BackendInfoResponse struct {
	Reachable bool   `json:"reachable"`
	Name      string `json:"name"`
}

// fs to maps folder

// mountFs configures the embedded file system for the application's
// front-end assets when building for production.
func (app *application) useGeocoder() {

	// register "GET /hello/{name}" route (allowed for everyone)
	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		addresserver := getEnvOrDefault("MNLRADDRESSSERVER", "")
		// register "GET /hello/{name}" route (allowed for everyone)
		se.Router.GET("/geoapi/geocode", func(e *core.RequestEvent) error {

			addresserver := getEnvOrDefault("MNLRADDRESSSERVER", "")

			var resultAddress []Address
			search := e.Request.URL.Query().Get("q")

			if addresserver == "" {
				resultAddress, _ = GetGeoFromOSM(search)
			}
			if addresserver != "" {
				resultAddress, _ = GetGeoFromMNLR(search, addresserver)
			}

			return e.JSON(http.StatusOK, resultAddress)

		})

		// register "GET /hello/{name}" route (allowed for everyone)
		se.Router.GET("/geoapi/backend", func(e *core.RequestEvent) error {

			var backendInfo BackendInfoResponse = BackendInfoResponse{
				Reachable: false,
				Name:      "OSM",
			}

			// Check if the address server is reachable by performing a demo request
			if addresserver != "" {
				_, err := GetGeoFromMNLR("test", addresserver)
				if err != nil {
					backendInfo.Reachable = false
					backendInfo.Name = addresserver + " (MNLRAddressServer)"
				} else {
					backendInfo.Reachable = true
					backendInfo.Name = addresserver + " (MNLRAddressServer)"
				}
			}

			if addresserver == "" {
				_, err := GetGeoFromOSM("test")
				if err != nil {
					backendInfo.Reachable = false
					backendInfo.Name = "OSM"
				} else {
					backendInfo.Reachable = true
					backendInfo.Name = "OSM"
				}
			}

			return e.JSON(http.StatusOK, backendInfo)

		})

		return se.Next()
	})

}

func GetGeoFromMNLR(address string, url string) ([]Address, error) {
	// Get Request to url/geocode?q=address

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+"/api/search", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("q", address)
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	// Parse the response
	var mnlrrespnse MNLRAddressServerResponse

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &mnlrrespnse)
	if err != nil {
		return nil, err
	}

	return mnlrrespnse.Addresses, nil
}

func GetGeoFromOSM(address string) ([]Address, error) {
	// Use OSM as default
	geocoder := openstreetmap.Geocoder()

	// Geocode the address
	resultLatLong, err := geocoder.Geocode(address)
	if err != nil {
		return nil, err
	}

	resulltReverse, err := geocoder.ReverseGeocode(resultLatLong.Lat, resultLatLong.Lng)
	if err != nil {
		return nil, err
	}

	// Return results
	addresses := make([]Address, 1)
	addresses[0].ID = 1
	addresses[0].Street = resulltReverse.Street
	addresses[0].HouseNumber = resulltReverse.HouseNumber
	addresses[0].City = resulltReverse.City
	addresses[0].Longitude = resultLatLong.Lng
	addresses[0].Latitude = resultLatLong.Lat

	return addresses, nil
}
