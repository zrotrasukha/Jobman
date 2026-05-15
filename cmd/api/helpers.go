package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
)

type envelop map[string]any

func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	MAX_BYTESIZE := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(MAX_BYTESIZE))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("malformed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("invalid value for field %q (at character %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			}
			return fmt.Errorf("invalid JSON (at character %d)", unmarshalTypeError.Offset)
		case strings.Contains(err.Error(), "http: request body too large"):
			return fmt.Errorf("body must not be larger than %d bytes", MAX_BYTESIZE)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
