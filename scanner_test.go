package main

import (
	"bytes"
	"testing"
)

const (
	testData        = "cache 20480\nrss 12288\nmapped_file 0\n"
	testDataWHeader = "252:0 Async 357\n252:0 Total 357\nTotal 714\n"
)

var (
	testHeader = []string{"device", "operation"}
)

func TestScanMetric(t *testing.T) {
	buf := bytes.NewBufferString(testData)
	s := NewMetricScanner(buf)

	if !s.Scan() {
		t.Error("Scan ended too early")
	}
	labels, value, err := s.Metric(nil)
	if err != nil {
		t.Error(t)
	}
	if value != 20480 {
		t.Errorf("%f != %f", value, 20480)
	}

	if !s.Scan() {
		t.Error("Scan ended too early")
	}

	if !s.Scan() {
		t.Error("Scan ended too early")
	}
	labels, value, err = s.Metric(nil)
	if err != nil {
		t.Error(t)
	}
	if labels["type"] != "mapped_file" {
		t.Errorf("Expected 'rss' but got: %s", labels["type"])
	}

	if s.Scan() {
		t.Error("Scan ran too long")
	}
}

func TestScanMetricHeader(t *testing.T) {
	buf := bytes.NewBufferString(testDataWHeader)
	s := NewMetricScanner(buf)

	if !s.Scan() {
		t.Error("Scan ended too early")
	}
	labels, value, err := s.Metric(testHeader)
	if err != nil {
		t.Error(t)
	}
	if value != 357 {
		t.Errorf("%f != %f", value, 357)
	}
	if labels["device"] != "252:0" || labels["operation"] != "Async" {
		t.Errorf("Expected device: 252:0 and operation: Total but got: %s/%s", labels["device"], labels["operation"])
	}

	if !s.Scan() {
		t.Error("Scan ended too early")
	}
	labels, value, err = s.Metric(nil)
	if err != ErrFormatUnknown {
		t.Errorf("Expected to get error %s but got: %s", ErrFormatUnknown, err)
	}

	if !s.Scan() {
		t.Error("Scan ended too early")
	}

	labels, value, err = s.Metric(testHeader)
	if err != ErrFormatUnknown {
		t.Errorf("Expected to get error %s but got: %s", ErrFormatUnknown, err)
	}

	if s.Scan() {
		t.Error("Scan ran too long")
	}
}
