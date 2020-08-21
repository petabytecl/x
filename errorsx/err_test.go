package errorsx

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestToErr(t *testing.T) {
	t.Run("case=stack", func(t *testing.T) {
		e := errors.New("hi")
		assert.EqualValues(t, e.(StackTracer).StackTrace(), ToErr(e, "").StackTrace())
	})

	t.Run("case=wrap", func(t *testing.T) {
		orig := errors.New("hi")
		wrap := new(Err)
		wrap.Wrap(orig)

		assert.EqualValues(t, orig.(StackTracer).StackTrace(), wrap.StackTrace())
	})

	t.Run("case=wrap_self", func(t *testing.T) {
		wrap := new(Err)
		wrap.Wrap(wrap)

		assert.Empty(t, wrap.StackTrace())
	})

	t.Run("case=status", func(t *testing.T) {
		e := &Err{
			StatusField: "foo-status",
		}
		assert.EqualValues(t, "foo-status", ToErr(e, "").Status())
	})

	t.Run("case=reason", func(t *testing.T) {
		e := &Err{
			ReasonField: "foo-reason",
		}
		assert.EqualValues(t, "foo-reason", ToErr(e, "").Reason())
	})

	t.Run("case=debug", func(t *testing.T) {
		e := &Err{
			DebugField: "foo-debug",
		}
		assert.EqualValues(t, "foo-debug", ToErr(e, "").Debug())
	})

	t.Run("case=details", func(t *testing.T) {
		e := &Err{
			DetailsField: map[string]interface{}{"foo-debug": "bar"},
		}
		assert.EqualValues(t, map[string]interface{}{"foo-debug": "bar"}, ToErr(e, "").Details())
	})

	t.Run("case=rid", func(t *testing.T) {
		e := &Err{
			RIDField: "foo-rid",
		}
		assert.EqualValues(t, "foo-rid", ToErr(e, "").RequestID())
		assert.EqualValues(t, "fallback-rid", ToErr(new(Err), "fallback-rid").RequestID())
	})

	t.Run("case=code", func(t *testing.T) {
		e := &Err{CodeField: 501}
		assert.EqualValues(t, 501, ToErr(e, "").StatusCode())
		assert.EqualValues(t, http.StatusText(501), ToErr(e, "").Status())

		e = &Err{CodeField: 501, StatusField: "foobar"}
		assert.EqualValues(t, 501, ToErr(e, "").StatusCode())
		assert.EqualValues(t, "foobar", ToErr(e, "").Status())

		assert.EqualValues(t, 500, ToErr(errors.New(""), "").StatusCode())
	})
}
