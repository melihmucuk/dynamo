package dynamo

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/cenkalti/backoff"
)

const retryTimeout = 1 * time.Minute // TODO: make this configurable

func retry(f func() error) error {
	var err error
	var next time.Duration
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = retryTimeout

	for {
		if err = f(); err == nil {
			return nil
		}

		if !canRetry(err) {
			return err
		}

		if next = b.NextBackOff(); next == backoff.Stop {
			return err
		}

		time.Sleep(next)
	}
}

func canRetry(err error) bool {
	if ae, ok := err.(awserr.RequestFailure); ok {
		switch ae.StatusCode() {
		case 500, 503:
			return true
		case 400:
			switch ae.Code() {
			case "ProvisionedThroughputExceededException",
				"ThrottlingException":
				return true
			}
		}
	}
	return false
}
