package xrequestid

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

type Option interface {
	apply(*options)
}

type optionApplyer func(*options)

func (a optionApplyer) apply(opt *options) {
	a(opt)
}

type options struct {
	chainRequestID   bool
	persistRequestID bool
	logRequest       bool
	validator        requestIDValidator
}

func ChainRequestID() Option {
	return optionApplyer(func(opt *options) {
		opt.chainRequestID = true
	})
}

// Attach the request id to the outgoing context
func PersistRequestID() Option {
	return optionApplyer(func(opt *options) {
		opt.persistRequestID = true
	})
}

// Logs the incoming request with the request id and the method destination
func LogRequest() Option {
	return optionApplyer(func(opt *options) {
		opt.logRequest = true
	})
}

type requestIDValidator func(string) bool

// RequestIDValidator is validator function that returns true if
// request id is valid, or false if invalid.
func RequestIDValidator(validator requestIDValidator) Option {
	return optionApplyer(func(opt *options) {
		opt.validator = validator
	})
}

func defaultReqeustIDValidator(requestID string) bool {
	return true
}

// Logs the incoming request logrus in context
func addRequestToLogger(ctx context.Context, requestID string, requestData interface{}) context.Context {
	log := ctxlogrus.Extract(ctx).WithFields(
		logrus.Fields{
			"Request ID":   requestID,
			"Request Data": fmt.Sprintf("%+v", requestData),
		})

	log.Info("request made")
	return ctxlogrus.ToContext(ctx, log)
}
