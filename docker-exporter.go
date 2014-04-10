package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/exp"
)

// We expect this cgroup to contains all tasks/containers we are interested in
const (
	defaultCgroup = "memory"
)

var (
	listen   = flag.String("listen", ":8080", "Address to listen on")
	addr     = flag.String("addr", "unix:///var/run/docker.sock", "Docker address to connect to")
	root     = flag.String("root", "/sys/fs/cgroup", "cgroup root")
	interval = flag.Duration("interval", 15*time.Second, "refresh interval")
)

func main() {
	flag.Parse()
	registry := prometheus.NewRegistry()
	scrapeDurations := prometheus.NewDefaultHistogram()
	registry.Register("docker_scrape_duration_seconds", "node_exporter: Duration of a scrape job.", prometheus.NilLabels, scrapeDurations)

	go func() {
		exp.Handle(prometheus.ExpositionResource, registry.Handler())
		log.Fatal(http.ListenAndServe(*listen, exp.DefaultCoarseMux))
	}()

	tick := time.Tick(*interval)
	for {
		log.Print("Updating metrics")
		de, err := newExporter(registry, *addr)
		if err != nil {
			log.Fatal(err)
		}
		begin := time.Now()
		if err := filepath.Walk(*root, de.walkFile); err != nil {
			log.Fatal(err)
		}
		scrapeDurations.Add(prometheus.NilLabels, time.Since(begin).Seconds())
		<-tick
	}
}
