package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrFormatUnknown = errors.New("Format unknown")
)

type metricScanner struct {
	s *bufio.Scanner
}

// Returns a new metric scanner
func NewMetricScanner(r io.Reader) *metricScanner {
	return &metricScanner{
		s: bufio.NewScanner(r),
	}
}

// Scan advances the scanner to the metric
func (s *metricScanner) Scan() bool {
	return s.s.Scan()
}

// Metric returns the most recent key, value and error on parse or IO error
// If header is empty, it expects to find two elements: key and value.
// If header is not empty, it expects to find one element per header key + value
func (s *metricScanner) Metric(headers []string) (map[string]string, float64, error) {
	if err := s.s.Err(); err != nil {
		return nil, 0, err
	}

	//header: device, operation
	line := s.s.Text()
	fields := strings.Split(line, " ")
	if len(headers) == 0 {
		headers = []string{"type"}
	}

	if len(fields) != len(headers)+1 {
		return nil, 0, ErrFormatUnknown
	}

	labels := map[string]string{}
	for i, header := range headers {
		labels[header] = fields[i]
	}
	value, err := strconv.ParseFloat(fields[len(headers)], 64)
	if err != nil {
		return nil, 0, err
	}
	return labels, value, nil
}
