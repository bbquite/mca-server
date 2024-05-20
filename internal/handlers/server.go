package handlers

import (
	_ "embed"
	"errors"
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//go:embed html/index.gohtml
var htmlTemplateEmbed string

type Handler struct {
	services      *service.MetricService
	indexTemplate *template.Template
}

func NewHandler(services *service.MetricService) (*Handler, error) {
	tml, err := template.New("indexTemplate").Parse(htmlTemplateEmbed)
	if err != nil {
		return &Handler{}, err
	}
	return &Handler{
		services:      services,
		indexTemplate: tml,
	}, nil
}

//indexTemplate := template.Must(template.New("indexTemplate").Parse(htmlTemplateEmbed))

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
			http.Error(w, "", http.StatusInternalServerError)
			log.Printf("ERROR | caught the problem: %v", err)
			return
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
			http.Error(w, "", http.StatusInternalServerError)
			log.Printf("ERROR | caught the problem: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")

	gauge, err := h.services.GetAllGaugeItems()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Printf("ERROR | caught the problem: %v", err)
		return
	}

	counter, err := h.services.GetAllCounterItems()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Printf("ERROR | caught the problem: %v", err)
		return
	}

	data := map[string]map[string]map[string]string{
		"metrics": {
			"counter": counter,
			"gauge":   gauge,
		},
	}
	h.indexTemplate.Execute(w, data)
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

		body := strconv.FormatFloat(float64(value), 'f', -1, 64)

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

		body := strconv.Itoa(int(value))

		w.Write([]byte(body))
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusOK)

		return

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
