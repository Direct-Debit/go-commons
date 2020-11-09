package webutil

import (
	"bytes"
	"encoding/json"
	"github.com/Direct-Debit/go-commons/errlib"
	"io"
	"io/ioutil"
	"net/http"
)

// Write err.Error() in JSON object to w and return true if err != nil.
// Return from your handler function if this method returns true.
// The JSON object will be the same as one generated by ErrorResponder
func ClientError(w http.ResponseWriter, err error, code int) bool {
	if err == nil {
		return false
	}
	ErrorResponder(err.Error(), code)(w, nil)
	return true
}

func ErrorResponder(msg string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		writer := json.NewEncoder(w)
		errlib.PanicError(writer.Encode(struct {
			Error string `json:"error"`
		}{Error: msg}), "Couldn't encode JSON error")
	}
}

func RequestBody(r *http.Request, consume bool) string {
	buff := new(bytes.Buffer)
	_, err := io.Copy(buff, r.Body)
	errlib.PanicError(err, "Couldn't copy request body")
	if !consume {
		r.Body = ioutil.NopCloser(buff)
	}
	return buff.String()
}

func ParseBodyJSON(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}
