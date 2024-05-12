package handlers

import (
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	services *service.MetricService
}

func NewHandler(services *service.MetricService) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{m_type}/{m_name}/{m_value}", h.apiHandler)
	return mux
}

func (h *Handler) apiHandler(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("m_type")
	mName := r.PathValue("m_name")
	mValue := r.PathValue("m_value")

	log.Print(r.URL)
	log.Print(mType, mName, mValue)

	w.Header().Set("Content-type", "text/plain")

	switch mType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("gauge не флоат")
			return
		}

		_, err = h.services.AddGaugeItem(mName, model.Gauge(metricValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("caught the problem: %v", err)
			return
		}

	case "counter":
		metricValue, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("counter не инт")
			return
		}

		_, err = h.services.AddCounterItem(mName, model.Counter(metricValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("caught the problem: %v", err)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Print("просто бэд")
		return
	}

	defer logStorage(h)
}

func logStorage(h *Handler) {
	log.Print(h.services.GetAllCounterItems())
	log.Print(h.services.GetAllGaugeItems())
}

//metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server
//metricstest -test.v -test.run=^TestIteration2[AB]*$ -source-path=. -agent-binary-path=cmd/agent/agent
