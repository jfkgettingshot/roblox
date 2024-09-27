package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Always check if PreviousPageCursor or NextPageCursor is not blank before using
type CursorResponse[T any] struct {
	PreviousPageCursor string `json:"previousPageCursor"`
	NextPageCursor     string `json:"nextPageCursor"`
	Data               []T    `json:"data"`
}

func ReadCursorResponse[T any](request *http.Response) (*CursorResponse[T], error) {
	defer request.Body.Close()

	if request.StatusCode != http.StatusOK {
		return nil, readErrorResponse(request.Body, request.StatusCode)
	}

	cursorResponse := new(CursorResponse[T])
	jsonDecoder := json.NewDecoder(request.Body)

	if err := jsonDecoder.Decode(cursorResponse); err != nil {
		return nil, fmt.Errorf("internal/responses.go/ReadCursorResponse failed to decode response: %w", err)
	}

	return cursorResponse, nil
}

func readErrorResponse(body io.Reader, status int) error {
	var response ErrorResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return fmt.Errorf("internal/responses.go/readErrorResponse failed to decode error response: %w", err)
	}

	if len(response.Errors) == 0 {
		return fmt.Errorf("internal/responses.go/readErrorResponse no error data found in response")
	}

	finalError := fmt.Errorf("recieved status code: %d expected 200", status)
	for _, err := range response.Errors {
		finalError = errors.Join(finalError, err.Error())
	}

	return finalError
}

type ErrorResponse struct {
	Errors []ErrorData `json:"errors"`
}

type ErrorData struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	UserFacingMessage string `json:"userFacingMessage"`
}

// Fun fact
// if this function dosent return an error
// there was an error
func (e *ErrorData) Error() error {
	// Exclude userFacingMessage from the error message
	// as the value is nearly always the same as message
	return fmt.Errorf("code: %d, message: %s", e.Code, e.Message)
}

// CursorHandler is a generic function to handle cursor responses
// still requires some repeated code but helps reduce it
func CursorHandler[T any](output *[]T, response *CursorResponse[T], err error) (string, error) {
	if err != nil {
		return "", err
	}

	*output = append(*output, response.Data...)
	return response.NextPageCursor, nil
}
