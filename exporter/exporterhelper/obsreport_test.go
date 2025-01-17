// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exporterhelper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/metric"
	"go.opencensus.io/tag"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/obsreport/obsreporttest"
)

func TestExportEnqueueFailure(t *testing.T) {
	exporterID := component.NewID("fakeExporter")
	tt, err := obsreporttest.SetupTelemetry(exporterID)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, tt.Shutdown(context.Background())) })

	insts := newInstruments(metric.NewRegistry())
	obsrep, err := newObsExporter(ObsReportSettings{
		ExporterID:             exporterID,
		ExporterCreateSettings: exporter.CreateSettings{ID: exporterID, TelemetrySettings: tt.TelemetrySettings, BuildInfo: component.NewDefaultBuildInfo()},
	}, insts)
	require.NoError(t, err)

	logRecords := int64(7)
	obsrep.recordLogsEnqueueFailure(context.Background(), logRecords)
	checkExporterEnqueueFailedLogsStats(t, insts, exporterID, logRecords)

	spans := int64(12)
	obsrep.recordTracesEnqueueFailure(context.Background(), spans)
	checkExporterEnqueueFailedTracesStats(t, insts, exporterID, spans)

	metricPoints := int64(21)
	obsrep.recordMetricsEnqueueFailure(context.Background(), metricPoints)
	checkExporterEnqueueFailedMetricsStats(t, insts, exporterID, metricPoints)
}

// checkExporterEnqueueFailedTracesStats checks that reported number of spans failed to enqueue match given values.
// When this function is called it is required to also call SetupTelemetry as first thing.
func checkExporterEnqueueFailedTracesStats(t *testing.T, insts *instruments, exporter component.ID, spans int64) {
	checkValueForProducer(t, insts.registry, tagsForExporterView(exporter), spans, "exporter/enqueue_failed_spans")
}

// checkExporterEnqueueFailedMetricsStats checks that reported number of metric points failed to enqueue match given values.
// When this function is called it is required to also call SetupTelemetry as first thing.
func checkExporterEnqueueFailedMetricsStats(t *testing.T, insts *instruments, exporter component.ID, metricPoints int64) {
	checkValueForProducer(t, insts.registry, tagsForExporterView(exporter), metricPoints, "exporter/enqueue_failed_metric_points")
}

// checkExporterEnqueueFailedLogsStats checks that reported number of log records failed to enqueue match given values.
// When this function is called it is required to also call SetupTelemetry as first thing.
func checkExporterEnqueueFailedLogsStats(t *testing.T, insts *instruments, exporter component.ID, logRecords int64) {
	checkValueForProducer(t, insts.registry, tagsForExporterView(exporter), logRecords, "exporter/enqueue_failed_log_records")
}

// tagsForExporterView returns the tags that are needed for the exporter views.
func tagsForExporterView(exporter component.ID) []tag.Tag {
	return []tag.Tag{
		{Key: exporterTag, Value: exporter.String()},
	}
}
