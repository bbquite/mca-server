package main

import (
	"net/http"
	"slices"
	"strconv"

	gorilla_mux "github.com/gorilla/mux"
)

func reqPostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
		}
		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	// type gauge float64
	// type counter int64

	pathVars := gorilla_mux.Vars(r)

	metricType := pathVars["metric_type"]
	metricTypeSlice := []string{"gauge", "counter"}
	if !slices.Contains(metricTypeSlice, metricType) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metricValue := pathVars["metric_value"]
	metricValueInt, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricValueInt++

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Это страница test"))

	// log.Println("")

}

func main() {

	mux := gorilla_mux.NewRouter()
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", apiHandler)
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

	err := http.ListenAndServe(":8080", reqPostMiddleware(mux))
	if err != nil {
		panic(err)
	}

}
