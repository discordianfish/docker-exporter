package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/prometheus/client_golang/prometheus"
)

var blkioKeys = []string{"device", "operation"}

type exporter struct {
	memoryStats  prometheus.Gauge
	cpuStats     prometheus.Gauge
	blockIOps    prometheus.Counter
	blockIOBytes prometheus.Counter

	containers map[string]docker.APIContainers
	registry   prometheus.Registry
}

func newExporter(registry prometheus.Registry, dockerUrl string) (*exporter, error) {
	dc, err := docker.NewClient(dockerUrl)
	if err != nil {
		return nil, err
	}

	cs, err := dc.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		return nil, err
	}
	containers := map[string]docker.APIContainers{}
	for _, container := range cs {
		containers[container.ID] = container
	}

	memoryStats := prometheus.NewGauge()
	cpuStats := prometheus.NewGauge()
	blockIOps := prometheus.NewCounter()
	blockIOBytes := prometheus.NewCounter()

	registry.Register("container_memory", "docker_exporter: container memory metrics", prometheus.NilLabels, memoryStats)
	registry.Register("container_cpu", "docker_exporter: container cpu metrics", prometheus.NilLabels, cpuStats)
	registry.Register("container_blockiops", "docker_exporter: IO operations to/from container", prometheus.NilLabels, blockIOps)
	registry.Register("container_blockio_bytes", "docker_exporter: IO bytes to/from container", prometheus.NilLabels, blockIOBytes)

	return &exporter{
		memoryStats:  memoryStats,
		cpuStats:     cpuStats,
		blockIOps:    blockIOps,
		blockIOBytes: blockIOBytes,

		containers: containers,
		registry:   registry,
	}, nil
}

// update metric if Base(Dir(path)) is in list of container and we can handle Name
func (e *exporter) walkFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	parentDir := filepath.Base(filepath.Dir(path))
	for id, container := range e.containers {
		if parentDir == id {
			var (
				gauge prometheus.Gauge
				keys  []string
			)
			switch info.Name() {
			case "memory.stat":
				gauge = e.memoryStats

			case "cpuacct.stat":
				gauge = e.cpuStats

			case "blkio.throttle.io_serviced":
				gauge = e.blockIOps
				keys = blkioKeys

			case "blkio.throttle.io_service_bytes":
				gauge = e.blockIOBytes
				keys = blkioKeys
			default:
				continue
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			ms := NewMetricScanner(file)
			for ms.Scan() {
				labels, value, err := ms.Metric(keys)
				if err == ErrFormatUnknown { // we are not interested in 'total' entry, so ignore unknown
					continue
				}
				if err != nil {
					return err
				}

				labels["container"] = strings.TrimPrefix(container.Names[0], "/")
				labels["container_id"] = container.ID
				gauge.Set(labels, value)
			}
		}
	}
	return nil
}
