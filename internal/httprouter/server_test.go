package httprouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	tt := []struct {
		name string
		code int
		r    *http.Request
	}{
		{
			name: "redirect",
			code: http.StatusMovedPermanently,
			r: httptest.NewRequest(
				http.MethodGet,
				"http://romanyx.ru/",
				nil,
			),
		},
		{
			name: "recover",
			code: http.StatusInternalServerError,
			r: httptest.NewRequest(
				http.MethodGet,
				"http://romanyx.info/",
				nil,
			),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			prepareHandler().ServeHTTP(w, tc.r)
			if w.Code != tc.code {
				t.Errorf("expected code %d got %d", tc.code, w.Code)
			}
		})
	}
}

func prepareHandler() http.Handler {
	h := NewHandler(
		nil,
		func(err error) {},
		nil,
		nil,
	)

	s := NewServer("", h)

	return s.server.Handler
}
