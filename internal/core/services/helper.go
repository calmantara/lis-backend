package services

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func setHTTPClient() (client *resty.Client) {
	// Create a new Resty client with retry configuration
	client = resty.New()

	// Configure retry settings
	client.
		SetRetryCount(3).                      // Number of retries
		SetRetryWaitTime(5 * time.Second).     // Wait time between retries
		SetRetryMaxWaitTime(10 * time.Second). // Maximum wait time
		AddRetryCondition(                     // Custom retry condition
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == 429 || r.StatusCode() >= 500
			},
		)

	return
}
