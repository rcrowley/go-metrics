package metrics

import (
	"errors"
	"testing"
)

func TestGetOrRegisterHealthcheck(t *testing.T) {
	r := NewRegistry()
	check := func(h Healthcheck) {
		h.Unhealthy(errors.New("foo"))
	}
	NewRegisteredHealthcheck("foo", r, check).Check()
	if h := GetOrRegisterHealthcheck("foo", r, check); h.Error().Error() != "foo" {
		t.Fatal(h)
	}
}
