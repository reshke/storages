package s3

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"time"
)

var defaultPermanentErrorCodes = map[int]struct{} {
	404 : {},
	418 : {}, // always retry if aws is teapot
}

type ThrottlerRetries struct {
	request.Retryer
}

// RetryRules return the retry delay that should be used by the SDK before
// making another request attempt for the failed request.

func (t ThrottlerRetries) RetryRules(r *request.Request) time.Duration {
	return t.Retryer.RetryRules(r)
}

// ShouldRetry returns if the failed request is retryable.
//
// Implementations may consider request attempt count when determining if a
// request is retryable, but the SDK will use MaxRetries to limit the
// number of attempts a request are made.

func (t ThrottlerRetries) ShouldRetry(r *request.Request) bool {
	statusCode := r.HTTPRequest.Response.StatusCode
	if _, found := defaultPermanentErrorCodes[statusCode]; found {
		return true
	}
	return t.Retryer.ShouldRetry(r)
}

// MaxRetries is the number of times a request may be retried before
// failing.

func (t ThrottlerRetries) MaxRetries() int {
	return t.Retryer.MaxRetries()
}

func NewThrottlerRetries(maxRetries int) ThrottlerRetries {
	return ThrottlerRetries{client.DefaultRetryer{NumMaxRetries: maxRetries}}
}
