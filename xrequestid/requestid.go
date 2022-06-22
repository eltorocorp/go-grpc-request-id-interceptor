package xrequestid

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "x-request-id"

func HandleRequestID(ctx context.Context, validator requestIDValidator) string {
	requestID := getStringFromContext(ctx, DefaultXRequestIDKey)
	if requestID == "" {
		return newRequestID()
	}

	if !validator(requestID) {
		return newRequestID()
	}

	return requestID
}

func HandleRequestIDChain(ctx context.Context, validator requestIDValidator) string {
	requestID := getStringFromContext(ctx, DefaultXRequestIDKey)
	if requestID == "" {
		return newRequestID()
	}

	if !validator(requestID) {
		return newRequestID()
	}

	return fmt.Sprintf("%s,%s", requestID, newRequestID())
}

func newRequestID() string {
	return uuid.NewString()
}
