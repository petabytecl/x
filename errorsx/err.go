package errorsx

import (
	stderr "errors"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Err struct {
	CodeField    int                    `json:"code,omitempty"`
	StatusField  string                 `json:"status,omitempty"`
	RIDField     string                 `json:"request,omitempty"`
	ReasonField  string                 `json:"reason,omitempty"`
	DebugField   string                 `json:"debug,omitempty"`
	DetailsField map[string]interface{} `json:"details,omitempty"`
	ErrorField   string                 `json:"message"`

	err error
}

func (e Err) Error() string {
	return e.ErrorField
}

func (e Err) Unwrap() error {
	return e.err
}

func (e *Err) Wrap(err error) {
	e.err = err
}

func (e Err) WithWrap(err error) *Err {
	e.err = err
	return &e
}

func (e *Err) WithTrace(err error) *Err {
	if st := StackTracer(nil); !stderr.As(e.err, &st) {
		e.Wrap(errors.WithStack(err))
	} else {
		e.Wrap(err)
	}
	return e
}

func (e Err) Is(err error) bool {
	switch te := err.(type) {
	case Err:
		return e.ErrorField == te.ErrorField &&
			e.StatusField == te.StatusField &&
			e.CodeField == te.CodeField
	case *Err:
		return e.ErrorField == te.ErrorField &&
			e.StatusField == te.StatusField &&
			e.CodeField == te.CodeField
	default:
		return false
	}
}

func (e Err) StatusCode() int {
	return e.CodeField
}

func (e Err) RequestID() string {
	return e.RIDField
}

func (e Err) Reason() string {
	return e.ReasonField
}

func (e Err) Debug() string {
	return e.DebugField
}

func (e Err) Status() string {
	return e.StatusField
}

func (e Err) Details() map[string]interface{} {
	return e.DetailsField
}

// StackTrace returns the error's stack trace.
func (e *Err) StackTrace() (trace errors.StackTrace) {
	if e.err == e {
		return
	}

	if st := StackTracer(nil); stderr.As(e.err, &st) {
		trace = st.StackTrace()
	}

	return
}

func (e Err) WithRequestID(id string) *Err {
	e.RIDField = id
	return &e
}

func (e Err) WithReason(reason string) *Err {
	e.ReasonField = reason
	return &e
}

func (e Err) WithReasonf(reason string, args ...interface{}) *Err {
	return e.WithReason(fmt.Sprintf(reason, args...))
}

func (e Err) WithError(message string) *Err {
	e.ErrorField = message
	return &e
}

func (e Err) WithErrorf(message string, args ...interface{}) *Err {
	return e.WithError(fmt.Sprintf(message, args...))
}

func (e Err) WithDebugf(debug string, args ...interface{}) *Err {
	return e.WithDebug(fmt.Sprintf(debug, args...))
}

func (e Err) WithDebug(debug string) *Err {
	e.DebugField = debug
	return &e
}

func (e Err) WithDetail(key string, detail interface{}) *Err {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = detail
	return &e
}

func (e Err) WithDetailf(key string, message string, args ...interface{}) *Err {
	if e.DetailsField == nil {
		e.DetailsField = map[string]interface{}{}
	}
	e.DetailsField[key] = fmt.Sprintf(message, args...)
	return &e
}

func ToErr(err error, id string) *Err {
	de := &Err{
		RIDField:     id,
		CodeField:    http.StatusInternalServerError,
		DetailsField: map[string]interface{}{},
		ErrorField:   err.Error(),
	}
	de.Wrap(err)

	if c := ReasonCarrier(nil); stderr.As(err, &c) {
		de.ReasonField = c.Reason()
	}
	if c := RequestIDCarrier(nil); stderr.As(err, &c) && c.RequestID() != "" {
		de.RIDField = c.RequestID()
	}
	if c := DetailsCarrier(nil); stderr.As(err, &c) && c.Details() != nil {
		de.DetailsField = c.Details()
	}
	if c := StatusCarrier(nil); stderr.As(err, &c) && c.Status() != "" {
		de.StatusField = c.Status()
	}
	if c := StatusCodeCarrier(nil); stderr.As(err, &c) && c.StatusCode() != 0 {
		de.CodeField = c.StatusCode()
	}
	if c := DebugCarrier(nil); stderr.As(err, &c) {
		de.DebugField = c.Debug()
	}

	if de.StatusField == "" {
		de.StatusField = http.StatusText(de.StatusCode())
	}

	return de
}
