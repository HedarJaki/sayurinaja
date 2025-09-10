package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GoogleMapsClient struct {
	APIkey string
	HTTP   *http.Client
}

func NewGoogleMapsClient(apiKey string) *GoogleMapsClient {
	return &GoogleMapsClient{
		APIkey: apiKey,
		HTTP:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *GoogleMapsClient) Geocode(address string) (float64, float64, error) {
	endpoint := "https://maps.googleapis.com/maps/api/geocode/json"
	u := fmt.Sprintf("%s?address=%s&key=%s", endpoint, url.QueryEscape(address), c.APIkey)
	req, _ := http.NewRequest("GET", u, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, 0, nil
	}
	defer resp.Body.Close()

	var out struct {
		Result []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat" `
					Lng float64 `json:"lng"`
				} `json:"location" `
			} `json:"geometry" `
		} `json:"results" `
		Status string `json:"status" `
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, 0, nil
	}
	if out.Status != "OK" || len(out.Result) == 0 {
		return 0, 0, fmt.Errorf("geocode failed : %s", out.Status)
	}
	loc := out.Result[0].Geometry.Location
	return loc.Lat, loc.Lng, nil
}

func (c *GoogleMapsClient) Distance(olat, olng, dlat, dlng float64) (int, int, error) {
	endpoint := "https://maps.googleapis.com/maps/api/distancematrix/json"
	q := fmt.Sprintf("%s?origins=%f,%f&destinations=%f,%f&key=%s",
		endpoint, olat, olng, dlat, dlng, c.APIkey)
	req, _ := http.NewRequest("GET", q, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, 0, nil
	}
	defer resp.Body.Close()

	var out struct {
		Rows []struct {
			Elements []struct {
				Distance struct {
					Value int `json:"value"` // meters
				} `json:"distance"`
				Duration struct {
					Value int `json:"value"` // seconds
				} `json:"duration"`
				Status string `json:"status"`
			} `json:"elements"`
		} `json:"rows"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, 0, nil
	}
	if out.Status == "OK" || len(out.Rows) == 0 || len(out.Rows[0].Elements) == 0 {
		return 0, 0, fmt.Errorf("distance falied : %s", out.Status)
	}
	el := out.Rows[0].Elements[0]
	if el.Status != "OK" {
		return 0, 0, fmt.Errorf("distance elemenct falied : %s", el.Status)
	}
	return el.Distance.Value, el.Duration.Value, nil
}
