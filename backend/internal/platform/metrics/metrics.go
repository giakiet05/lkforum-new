package metrics

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

var registry = &metricRegistry{
	counters:   map[string]float64{},
	gauges:     map[string]float64{},
	histograms: map[string]*histogram{},
}

type metricRegistry struct {
	mu         sync.RWMutex
	counters   map[string]float64
	gauges     map[string]float64
	histograms map[string]*histogram
}

type histogram struct {
	buckets map[float64]uint64
	count   uint64
	sum     float64
}

var defaultBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

func IncCounter(name string, labels map[string]string) {
	AddCounter(name, labels, 1)
}

func AddCounter(name string, labels map[string]string, value float64) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.counters[metricKey(name, labels)] += value
}

func SetGauge(name string, labels map[string]string, value float64) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.gauges[metricKey(name, labels)] = value
}

func AddGauge(name string, labels map[string]string, delta float64) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.gauges[metricKey(name, labels)] += delta
}

func ObserveDuration(name string, labels map[string]string, duration time.Duration) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	key := metricKey(name, labels)
	h := registry.histograms[key]
	if h == nil {
		h = &histogram{buckets: map[float64]uint64{}}
		registry.histograms[key] = h
	}

	seconds := duration.Seconds()
	for _, bucket := range defaultBuckets {
		if seconds <= bucket {
			h.buckets[bucket]++
		}
	}
	h.count++
	h.sum += seconds
}

func Render() string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	var out strings.Builder
	writeSamples(&out, registry.counters, "counter")
	writeSamples(&out, registry.gauges, "gauge")

	keys := make([]string, 0, len(registry.histograms))
	for key := range registry.histograms {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	declaredHistograms := map[string]bool{}
	for _, key := range keys {
		name, labelText := splitMetricKey(key)
		h := registry.histograms[key]
		if !declaredHistograms[name] {
			out.WriteString(fmt.Sprintf("# TYPE %s histogram\n", name))
			declaredHistograms[name] = true
		}
		for _, bucket := range defaultBuckets {
			out.WriteString(fmt.Sprintf("%s_bucket%s %d\n", name, addLabel(labelText, "le", fmt.Sprintf("%g", bucket)), h.buckets[bucket]))
		}
		out.WriteString(fmt.Sprintf("%s_bucket%s %d\n", name, addLabel(labelText, "le", "+Inf"), h.count))
		out.WriteString(fmt.Sprintf("%s_sum%s %g\n", name, labelText, h.sum))
		out.WriteString(fmt.Sprintf("%s_count%s %d\n", name, labelText, h.count))
	}

	return out.String()
}

func writeSamples(out *strings.Builder, samples map[string]float64, typ string) {
	keys := make([]string, 0, len(samples))
	for key := range samples {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	declared := map[string]bool{}
	for _, key := range keys {
		name, labelText := splitMetricKey(key)
		if !declared[name] {
			out.WriteString(fmt.Sprintf("# TYPE %s %s\n", name, typ))
			declared[name] = true
		}
		out.WriteString(fmt.Sprintf("%s%s %g\n", name, labelText, samples[key]))
	}
}

func metricKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%s=%q`, sanitizeLabelName(key), escapeLabelValue(labels[key])))
	}
	return name + "{" + strings.Join(parts, ",") + "}"
}

func splitMetricKey(key string) (string, string) {
	i := strings.IndexByte(key, '{')
	if i == -1 {
		return key, ""
	}
	return key[:i], key[i:]
}

func addLabel(labelText, name, value string) string {
	pair := fmt.Sprintf(`%s=%q`, name, escapeLabelValue(value))
	if labelText == "" {
		return "{" + pair + "}"
	}
	return strings.TrimSuffix(labelText, "}") + "," + pair + "}"
}

func sanitizeLabelName(name string) string {
	return strings.NewReplacer("-", "_", ".", "_").Replace(name)
}

func escapeLabelValue(value string) string {
	return strings.NewReplacer("\\", "\\\\", "\n", "\\n", `"`, `\"`).Replace(value)
}
