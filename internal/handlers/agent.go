package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/utils"
	"go.uber.org/zap"
)

func SendMetricsURI(services *service.MetricService, host string, logger *zap.SugaredLogger) error {

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

	for _, el := range metricsPack {
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

func SendMetricsJSON(services *service.MetricService, host string, shakey string, logger *zap.SugaredLogger) error {

	url := fmt.Sprintf("http://%s/update/", host)
	client := http.Client{}

	metricsPack, err := services.GetAllMetrics()
	if err != nil {
		return err
	}

	logger.Infof("Sending metrics to %s", host)

	for _, el := range metricsPack {

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

		if shakey != "" {
			sign := hex.EncodeToString(utils.MakeHMACSign(shakey, body))
			request.Header.Set("HashSHA256", sign)
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

	err = services.ResetCounterItem("PollCount")
	if err != nil {
		logger.Errorf("PollCount reset error")
	}

	return nil
}

func SendMetricsPackJSON(services *service.MetricService, host string, shakey string, logger *zap.SugaredLogger) error {
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

	if shakey != "" {
		sign := hex.EncodeToString(utils.MakeHMACSign(shakey, metricsJSON))
		request.Header.Set("HashSHA256", sign)
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

	err = services.ResetCounterItem("PollCount")
	if err != nil {
		logger.Errorf("PollCount reset error")
	}

	return nil
}
