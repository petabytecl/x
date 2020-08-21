package errorsx

import "net/http"

var ErrInternalServerError = Err{
	CodeField:   http.StatusInternalServerError,
	StatusField: http.StatusText(http.StatusInternalServerError),
	ErrorField:  "internal server error occurred",
}

var ErrNotFound = Err{
	CodeField:   http.StatusNotFound,
	StatusField: http.StatusText(http.StatusNotFound),
	ErrorField:  "resource not found",
}

var ErrUnauthorized = Err{
	CodeField:   http.StatusUnauthorized,
	StatusField: http.StatusText(http.StatusUnauthorized),
	ErrorField:  "authorized request",
}

var ErrForbidden = Err{
	CodeField:   http.StatusForbidden,
	StatusField: http.StatusText(http.StatusForbidden),
	ErrorField:  "forbidden requested action",
}

var ErrBadRequest = Err{
	CodeField:   http.StatusBadRequest,
	StatusField: http.StatusText(http.StatusBadRequest),
	ErrorField:  " malformed request or invalid parameters",
}

var ErrUnsupportedMediaType = Err{
	CodeField:   http.StatusUnsupportedMediaType,
	StatusField: http.StatusText(http.StatusUnsupportedMediaType),
	ErrorField:  "unsupported content type",
}
