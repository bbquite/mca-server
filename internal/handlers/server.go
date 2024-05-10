package handlers

import (
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/server"
	"github.com/bbquite/mca-server/internal/service"
	"log"
	"net/http"
	"strconv"
)

func InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{m_type}/{m_name}/{m_value}", apiHandler)
	return mux
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("m_type")
	mName := r.PathValue("m_name")
	mValue := r.PathValue("m_value")

	log.Print(r)

	metricService, ok := r.Context().Value(server.ServiceCtx{}).(service.MetricService)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("bad conext")
		return
	}

	switch mType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = metricService.AddGaugeItem(mName, model.Gauge(metricValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("caught the problem: %v", err)
			return
		}

	case "counter":
		metricValue, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = metricService.AddCounterItem(mName, model.Counter(metricValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("caught the problem: %v", err)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return
}

//metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server
//metricstest -test.v -test.run=^TestIteration2[AB]*$ -source-path=. -agent-binary-path=cmd/agent/agent
