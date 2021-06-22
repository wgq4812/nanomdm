package service

import (
	"fmt"
	"net/http"

	"github.com/micromdm/nanomdm/mdm"
)

type HTTPStatusError struct {
	Status int
	Err    error
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("HTTP status %d (%s): %v", e.Status, http.StatusText(e.Status), e.Err)
}

func (e *HTTPStatusError) Unwrap() error {
	return e.Err
}

func NewHTTPStatusError(status int, err error) *HTTPStatusError {
	return &HTTPStatusError{Status: status, Err: err}
}

// CheckinRequest is a simple adapter that takes the raw check-in bodyBytes
// and dispatches to the respective check-in method on svc.
func CheckinRequest(svc Checkin, r *mdm.Request, bodyBytes []byte) ([]byte, error) {
	msg, err := mdm.DecodeCheckin(bodyBytes)
	if err != nil {
		return nil, NewHTTPStatusError(http.StatusBadRequest, fmt.Errorf("decoding check-in: %w", err))
	}
	switch m := msg.(type) {
	case *mdm.Authenticate:
		err = svc.Authenticate(r, m)
		if err != nil {
			err = fmt.Errorf("authenticate service: %w", err)
		}
	case *mdm.TokenUpdate:
		err = svc.TokenUpdate(r, m)
		if err != nil {
			err = fmt.Errorf("tokenupdate service: %w", err)
		}
	case *mdm.CheckOut:
		err = svc.CheckOut(r, m)
		if err != nil {
			err = fmt.Errorf("checkout service: %w", err)
		}
	default:
		return nil, NewHTTPStatusError(http.StatusBadRequest, mdm.ErrUnrecognizedMessageType)
	}
	return nil, err
}

// CommandAndReportResultsRequest is a simple adapter that takes the raw
// command result report bodyBytes, dispatches to svc, and returns the
// response.
func CommandAndReportResultsRequest(svc CommandAndReportResults, r *mdm.Request, bodyBytes []byte) ([]byte, error) {
	report, err := mdm.DecodeCommandResults(bodyBytes)
	if err != nil {
		return nil, NewHTTPStatusError(http.StatusBadRequest, fmt.Errorf("decoding command results: %w", err))
	}
	cmd, err := svc.CommandAndReportResults(r, report)
	if err != nil {
		return nil, fmt.Errorf("command and report results service: %w", err)
	}
	if cmd != nil {
		return cmd.Raw, nil
	}
	return nil, nil
}
