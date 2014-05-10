package metrics

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRegistryMarshallJSON(t *testing.T) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	r := NewRegistry()
	r.Register("counter", NewCounter())
	enc.Encode(r)
	if s := b.String(); "{\"counter\":{\"count\":0}}\n" != s {
		t.Fatalf(s)
	}
}
