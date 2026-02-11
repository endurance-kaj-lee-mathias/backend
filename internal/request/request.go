package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxBodySize = 1_048_576

func Decode(r *http.Request, data any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))

		if mediaType != "application/json" {
			return ContentType
		}
	}

	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(data); err != nil {
		var syntax *json.SyntaxError
		var unmarshalType *json.UnmarshalTypeError
		var maxBytes *http.MaxBytesError

		switch {
		case errors.As(err, &syntax):
			return InvalidJSON
		case errors.Is(err, io.ErrUnexpectedEOF):
			return InvalidJSON
		case errors.As(err, &unmarshalType):
			return fmt.Errorf("invalid value for field %q", unmarshalType.Field)
		case errors.Is(err, io.EOF):
			return EmptyBody
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s", fieldName)
		case errors.As(err, &maxBytes):
			return BodyTooLarge
		default:
			return err
		}
	}

	if dec.More() {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
