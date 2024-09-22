package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bbquite/mca-server/internal/service"
	"go.uber.org/zap"
)

func MetricsURIRequest(services *service.MetricService, host string, logger *zap.SugaredLogger) error {

	var url string
	var value any
	client := http.Client{}

	metricsPack, err := services.GetAllMetrics()
	if err != nil {
		return err
	}

	logger.Infof("Sending metrics to %s", host)

	err = services.ResetCounterItem("PollCount")
	if err != nil {
		logger.Errorf("PollCount reset error")
	}

	for _, el := range metricsPack.Metrics {
		switch el.MType {
		case "gauge":
			value = fmt.Sprintf("%.2f", *el.Value)
		case "counter":
			value = fmt.Sprintf("%v", *el.Delta)
		}

		url = fmt.Sprintf("http://%s/update/%s/%s/%s", host, el.MType, el.ID, value)

		logger.Debugf("SEND %s", url)

		request, err := http.NewRequest(http.MethodPost, url, http.NoBody)
		if err != nil {
			logger.Error(err)
			return nil
		}

		request.Header.Set("Content-Type", "Content-Type: text/plain")

		response, err := client.Do(request)
		if err != nil {
			logger.Error(err)
			return nil
		}

		err = response.Body.Close()
		if err != nil {
			logger.Error(err)
			return nil
		}
	}

	return nil
}

func MetricsJSONRequest(services *service.MetricService, host string, logger *zap.SugaredLogger) error {

	url := fmt.Sprintf("http://%s/update/", host)
	client := http.Client{}

	metricsPack, err := services.GetAllMetrics()
	if err != nil {
		return err
	}

	logger.Infof("Sending metrics to %s", host)

	err = services.ResetCounterItem("PollCount")
	if err != nil {
		logger.Errorf("PollCount reset error")
	}

	for _, el := range metricsPack.Metrics {

		body, err := json.Marshal(el)
		if err != nil {
			logger.Error(err)
			return err
		}

		bodySend := bytes.NewBuffer(body)

		request, err := http.NewRequest(http.MethodPost, url, bodySend)
		if err != nil {
			logger.Error(err)
			return err
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept-Encoding", "gzip")

		logger.Debugf("SEND %s %s", url, request.Body)

		response, err := client.Do(request)
		if err != nil {
			logger.Error(err)
			return nil
		}

		defer response.Body.Close()
		logger.Debugf("RESP %s %s", url, response.Status)
	}

	return nil
}

func MetricsPackJSONRequest(services *service.MetricService, host string, logger *zap.SugaredLogger) error {
	url := fmt.Sprintf("http://%s/updates/", host)
	client := http.Client{}

	metricsJSON, err := services.ExportToJSON()
	if err != nil {
		logger.Error(err)
		return err
	}
	bodySend := bytes.NewBuffer(metricsJSON)

	request, err := http.NewRequest(http.MethodPost, url, bodySend)
	if err != nil {
		logger.Error(err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept-Encoding", "gzip")

	logger.Debugf("SEND %s %s", url, request.Body)

	response, err := client.Do(request)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer response.Body.Close()
	logger.Debugf("RESP %s %s", url, response.Status)

	return nil
}
