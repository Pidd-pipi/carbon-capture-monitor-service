package main

import (
	"context"
	"log"

	"example.com/carbon-capture-monitor-service/config"
	"example.com/carbon-capture-monitor-service/httpapi"
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
	"example.com/carbon-capture-monitor-service/web"
)

func main() {
	address := ":" + config.Port()
	log.Printf("carbon-capture-monitor-service listening on %s", address)

	unitStore := store.New()
	alertService := ops.NewService(nil)
	configureDefaultLabels(config.NewLabelsProvider())
	readingStore := readings.NewStore(config.ReadingRetention(), config.MaxReadingsPerUnit())

	handler := httpapi.NewHandler(unitStore, alertService, readingStore, web.FS)

	monitor := newMonitor(alertService, readingStore, unitStore, config.MonitorInterval())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	monitor.Start(ctx)

	if err := serveAddress(address, handler); err != nil {
		cancel()
		log.Fatal(err)
	}
}

// configureDefaultLabels wires the default alert labels into the ops layer.
// When the provider has no labels configured (e.g. ALERT_DEFAULT_LABELS unset),
// a non-nil empty map is installed so the service runs with zero configuration;
// NormalizeRecord still seeds the "source" label on every created alert.
func configureDefaultLabels(provider config.LabelsProvider) {
	if provider != nil {
		if labels := provider.Labels(); labels != nil {
			ops.SetDefaultLabels(labels)
			return
		}
	}
	ops.SetDefaultLabels(map[string]string{})
}
