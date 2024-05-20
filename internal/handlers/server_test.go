package handlers

import (
	_ "embed"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed html/index.gohtml
var htmlTemplateEmbed string

func Test_apiHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "positive test #1",
			url:  "/update/gauge/MSpanInuse/34400",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #2",
			url:  "/update/counter/PollCount/10",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #1",
			url:  "/updateerror/gauge/MSpanInuse/34400",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2",
			url:  "/update/counter/PollCount/10.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}

	db := storage.NewMemStorage()
	serv := service.NewMetricService(db)

	handler, err := NewHandler(serv)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	mux := handler.InitChiRoutes()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, test.url, nil)

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()
		})
	}
}
