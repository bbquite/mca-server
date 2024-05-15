package handlers

import (
	"errors"
	"fmt"
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type Handler struct {
	services *service.MetricService
}

func NewHandler(services *service.MetricService) *Handler {
	return &Handler{services: services}
}

// InitRoutes Оригинальньный роутер
func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{m_type}/{m_name}/{m_value}", h.addMetricByName)
	return mux
}

// InitChiRoutes Роутер с chi
func (h *Handler) InitChiRoutes() *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)
	chiRouter.Route("/", func(r chi.Router) {
		r.Get("/", h.getAllMetrics)
		r.Get("/value/{m_type}/{m_name}", h.getMetricByName)
		r.Post("/update/{m_type}/{m_name}/{m_value}", h.addMetricByName)
	})

	return chiRouter
}

func (h *Handler) addMetricByName(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "m_type")
	mName := chi.URLParam(r, "m_name")
	mValue := chi.URLParam(r, "m_value")

	w.Header().Set("Content-type", "text/plain")

	switch mType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		_, err = h.services.AddGaugeItem(mName, model.Gauge(metricValue))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("caught the problem: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		return

	case "counter":
		metricValue, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		_, err = h.services.AddCounterItem(mName, model.Counter(metricValue))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("caught the problem: %v", err)
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")

	if t, err := template.ParseFiles(filepath.Join("web", "index.gohtml")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR %s", err.Error())
	} else {
		gauge, _ := h.services.GetAllGaugeItems()
		counter, _ := h.services.GetAllCounterItems()
		data := map[string]map[string]map[string]string{
			"metrics": {
				"counter": counter,
				"gauge":   gauge,
			},
		}
		t.Execute(w, data)
	}
}

func (h *Handler) getMetricByName(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "m_type")
	mName := chi.URLParam(r, "m_name")

	switch mType {
	case "gauge":

		value, err := h.services.GetGaugeItem(mName)
		if errors.Is(err, service.ErrorGaugeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body := fmt.Sprintf("%v", value)

		w.Write([]byte(body))
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusOK)

		return

	case "counter":

		value, err := h.services.GetCounterItem(mName)
		if errors.Is(err, service.ErrorCounterNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body := fmt.Sprintf("%v", value)
		w.Write([]byte(body))
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusOK)

		return

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
