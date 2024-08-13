package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	RepositoryModelsName   = "models"
	RepositoryModelsDesc   = "Models located in the triton model repository"
	RepositoryModelsLabels = []string{"nodename", "modelname", "state", "version"}
)

type RepositoryCollector struct {
	hostname     string
	tritonAddr   string
	modelsMetric *prometheus.Desc
}

func (c *RepositoryCollector) Describe(ch chan<- *prometheus.Desc) {
	c.modelsMetric = prometheus.NewDesc(
		prometheus.BuildFQName("triton_repo_exporter", "", RepositoryModelsName),
		RepositoryModelsDesc, RepositoryModelsLabels, nil,
	)
}

func (c *RepositoryCollector) Collect(ch chan<- prometheus.Metric) {
	url := fmt.Sprintf("http://%s/metrics", c.tritonAddr)
	r, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Errorf("Failed to make http request to the triton server: %s", err)
	}
	defer r.Body.Close()

	var result []tritonRepositoryModelResp
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		log.Errorf("Failed to parse response from the triton server: %s", err)
	}

	for _, it := range result {
		if it.State == "" {
			continue
		}

		labelVals := []string{c.hostname, it.Name, it.State, it.Version}
		ch <- prometheus.MustNewConstMetric(c.modelsMetric, prometheus.GaugeValue, 1, labelVals...)
	}

}

type tritonRepositoryModelResp struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	State   string `json:"state,omitempty"`
	Reason  string `json:"reason,omitempty"`
}
