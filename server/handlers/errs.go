package handlers

import "errors"

var (
	ErrBind          = errors.New("Invalid data")
	ErrGetSiteStruct = errors.New("Failed to get site struct")
	ErrGetDetails    = errors.New("Failed to get link details")
)

func newHTTPError(err error) map[string]any {
	return map[string]any{
		"error": err.Error(),
	}
}
