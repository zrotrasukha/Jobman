package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// envelop is a simple wrapper type that allows us to easily create JSON responses with a consistent structure. It is defined as a map with string keys and values of any type, which can be used in helpers like writeJSON.
type envelop map[string]any

// writeJSON writes the given data as JSON to the response writer, along with any additional headers. It also sets the appropriate content type and status code for the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {

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

// readJSON reads the JSON from the request body and decodes it into the destination struct. It also performs validation checks for syntax errors, unknown fields, and body size limits, and returns appropriate error messages for each case.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
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

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldname := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s", fieldname)
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

// readParamID extracts the "id" parameter from the URL path, converts it to an int64, and returns it. If the parameter is missing, not a valid integer, or less than 1, it returns an error indicating that the ID parameter is invalid.
func (app *application) readParamID(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// readString retrieves the value of the specified key from the query string parameters. If the key is not present or has an empty value, it returns the provided default value instead.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

// readCSV retrieves the value of the specified key from the query string parameters, splits it by commas, and returns a slice of strings. If the key is not present or has an empty value, it returns the provided default slice instead.
func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
