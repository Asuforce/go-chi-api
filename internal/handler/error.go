package handler

import (
	"fmt"
	"log"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}

func RespondError(w http.ResponseWriter, code int, err error) {
	log.Printf("err: %v", err)
	if e, ok := err.(*HTTPError); ok {
		respondJSON(w, e.Code, e)
	} else if err != nil {
		he := HTTPError{
			Code:    code,
			Message: err.Error(),
		}
		respondJSON(w, code, he)
	}
}
