package handlers

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/bbquite/mca-server/internal/middleware"
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//go:embed html/index.gohtml
var htmlTemplateEmbed string

type Handler struct {
	services      *service.MetricService
	indexTemplate *template.Template
	logger        *zap.SugaredLogger
}

func NewHandler(services *service.MetricService) (*Handler, error) {
	tml, err := template.New("indexTemplate").Parse(htmlTemplateEmbed)
	if err != nil {
		return &Handler{}, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return &Handler{}, err
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	return &Handler{
		services:      services,
		indexTemplate: tml,
		logger:        sugar,
	}, nil
}

func (h *Handler) InitChiRoutes() *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.RequestsLoggingMiddleware(h.logger))
	chiRouter.Route("/", func(r chi.Router) {
		r.Get("/", h.renderMetricsPage)
		r.Route("/value/", func(r chi.Router) {
			r.Post("/", h.valueMetricJSON)
			r.Get("/{m_type}/{m_name}", h.valueMetricURI)
		})
		r.Route("/update/", func(r chi.Router) {
			r.Post("/", h.updateMetricJSON)
			r.Post("/{m_type}/{m_name}/{m_value}", h.updateMetricURI)
		})
	})

	return chiRouter
}

func (h *Handler) updateMetricJSON(w http.ResponseWriter, r *http.Request) {

	var metric model.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	switch metric.MType {
	case "gauge":
		_, err = h.services.AddGaugeItem(metric.ID, model.Gauge(metric.Value))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			h.logger.Error(err)
			return
		}

	case "counter":
		_, err = h.services.AddCounterItem(metric.ID, model.Counter(metric.Delta))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			h.logger.Error(err)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) updateMetricURI(w http.ResponseWriter, r *http.Request) {
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
			h.logger.Error(err)
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
			h.logger.Error(err)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) renderMetricsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")

	gauge, err := h.services.GetAllGaugeItems()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	counter, err := h.services.GetAllCounterItems()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		h.logger.Error(err)
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

func (h *Handler) valueMetricJSON(w http.ResponseWriter, r *http.Request) {

	var metric model.Metric
	var metricResponse model.Metric

	var metricGaugeValue model.Gauge
	var metricCounterValue model.Counter

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		metricGaugeValue, err = h.services.GetGaugeItem(metric.ID)
		if err != nil {
			if !errors.Is(err, service.ErrorGaugeNotFound) {
				http.Error(w, "", http.StatusInternalServerError)
				h.logger.Error(err)
				return
			}
		}

		metricResponse = model.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Value: float64(metricGaugeValue),
		}

	case "counter":
		metricCounterValue, err = h.services.GetCounterItem(metric.ID)
		if err != nil {
			if !errors.Is(err, service.ErrorCounterNotFound) {
				http.Error(w, "", http.StatusInternalServerError)
				h.logger.Error(err)
				return
			}
		}

		metricResponse = model.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: int64(metricCounterValue),
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metricResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) valueMetricURI(w http.ResponseWriter, r *http.Request) {
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
