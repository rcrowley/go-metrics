package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRegistryMarshallJSON(t *testing.T) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	r := NewRegistry()
	r.Register("counter", NewCounter())
	enc.Encode(r)
	if s := b.String(); strings.Contains(s, "{\"counter\":{\"count\":0}}") == false {
		t.Fatalf(s)
	}
}

func TestRegistryWriteJSONOnce(t *testing.T) {
	r := NewRegistry()
	r.Register("counter", NewCounter())
	b := &bytes.Buffer{}
	WriteJSONOnce(r, b)
	if s := b.String(); strings.Contains(s, "{\"counter\":{\"count\":0}}") == false {
		t.Fail()
	}
}
