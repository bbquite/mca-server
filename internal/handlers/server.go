package handlers

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/bbquite/mca-server/internal/middleware"
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/bbquite/mca-server/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//go:embed html/index.gohtml
var htmlTemplateEmbed string

type Handler struct {
	services      *service.MetricService
	indexTemplate *template.Template
	logger        *zap.SugaredLogger
	shaKey        string
}

func NewHandler(services *service.MetricService, shaKey string, logger *zap.SugaredLogger) (*Handler, error) {
	tml, err := template.New("indexTemplate").Parse(htmlTemplateEmbed)
	if err != nil {
		return &Handler{}, err
	}

	return &Handler{
		services:      services,
		indexTemplate: tml,
		logger:        logger,
		shaKey:        shaKey,
	}, nil
}

func (h *Handler) InitChiRoutes() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RequestsLoggingMiddleware(h.logger))
	// chiRouter.Use(chiMiddleware.Logger)
	chiRouter.Use(middleware.GzipMiddleware)

	chiRouter.Route("/", func(r chi.Router) {
		r.Get("/", h.renderMetricsPage)
		r.Get("/ping", h.databasePing)
		r.Route("/value/", func(r chi.Router) {
			r.Post("/", h.valueMetricJSON)
			r.Get("/{m_type}/{m_name}", h.valueMetricURI)
		})
		r.Route("/update/", func(r chi.Router) {
			r.Post("/", h.updateMetricJSON)
			r.Post("/{m_type}/{m_name}/{m_value}", h.updateMetricURI)
		})
		r.Post("/updates/", h.updatePackMetricsJSON)
	})

	return chiRouter
}

func (h *Handler) databasePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain")
	err := h.services.PingDatabase()
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updatePackMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.shaKey != "" {
		shaHeaderSign, err := hex.DecodeString(r.Header.Get("HashSHA256"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(err)
		}
		if utils.CheckHMACEqual(h.shaKey, shaHeaderSign, buf.Bytes()) {
			h.logger.Info("Норм подпись")
		} else {
			h.logger.Info("Подпись не оч")
		}
	}

	h.logger.Debugf("| req %s", buf.Bytes())

	err = h.services.ImportFromJSON(buf.Bytes())
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateMetricJSON(w http.ResponseWriter, r *http.Request) {

	var metric model.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.shaKey != "" {
		shaHeaderSign, err := hex.DecodeString(r.Header.Get("HashSHA256"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(err)
		}
		if utils.CheckHMACEqual(h.shaKey, shaHeaderSign, buf.Bytes()) {
			h.logger.Info("Норм подпись")
		} else {
			h.logger.Info("Подпись не оч")
		}
	}

	h.logger.Debugf("| req %s", buf.Bytes())

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	switch metric.MType {
	case "gauge":
		_, err = h.services.AddGaugeItem(metric.ID, model.Gauge(*metric.Value))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(err)
			return
		}

	case "counter":
		_, err = h.services.AddCounterItem(metric.ID, model.Counter(*metric.Delta))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	h.logger.Debugf("| resp %s", resp)
}

func (h *Handler) updateMetricURI(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "m_type")
	mName := chi.URLParam(r, "m_name")
	mValue := chi.URLParam(r, "m_value")

	w.Header().Set("Content-type", "text/plain")
	w.Header().Set("Content-Encoding", "gzip")

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
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")
	data, err := h.services.ExportToJSON()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		h.logger.Error(err)
	}

	var tmplContext map[string]interface{}
	if err := json.Unmarshal(data, &tmplContext); err != nil {
		h.logger.Error(err)
	}

	h.indexTemplate.Execute(w, tmplContext)
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

	if h.shaKey != "" {
		shaHeaderSign, err := hex.DecodeString(r.Header.Get("HashSHA256"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(err)
		}
		if utils.CheckHMACEqual(h.shaKey, shaHeaderSign, buf.Bytes()) {
			h.logger.Info("Норм подпись")
		} else {
			h.logger.Info("Подпись не оч")
		}
	}

	h.logger.Debugf("| req %s", buf.Bytes())

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		metricGaugeValue, err = h.services.GetGaugeItem(metric.ID)
		if err != nil {
			h.logger.Debug(err)
			if !errors.Is(err, storage.ErrorGaugeNotFound) {
				http.Error(w, "", http.StatusInternalServerError)
				h.logger.Error(err)
				return
			} else {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		val := float64(metricGaugeValue)

		metricResponse = model.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Value: &val,
		}

	case "counter":
		metricCounterValue, err = h.services.GetCounterItem(metric.ID)
		if err != nil {
			if !errors.Is(err, storage.ErrorCounterNotFound) {
				http.Error(w, "", http.StatusInternalServerError)
				h.logger.Error(err)
				return
			} else {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		val := int64(metricCounterValue)

		metricResponse = model.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: &val,
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
	h.logger.Debugf("| resp %s", resp)
}

func (h *Handler) valueMetricURI(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "m_type")
	mName := chi.URLParam(r, "m_name")

	switch mType {
	case "gauge":

		value, err := h.services.GetGaugeItem(mName)
		if errors.Is(err, storage.ErrorGaugeNotFound) {
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
		if errors.Is(err, storage.ErrorCounterNotFound) {
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
