package analytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/atharvamhaske/supaspy/internal/models"
	"google.golang.org/grpc/benchmark/latency"
)

const (
	slowQueryTH        = 300
	ineffQueryTH       = 700
	metricsLogInterval = 10 * time.Second

	// N+1 detection: flag a query excecuted more than maxN1Freq times within n1Window.
	// here we are hardcoding all values as this is just tiny workflow, this will be changed in production systems.
	maxN1Freq = 6
	n1Window  = 3 * time.Second
)

// analyzer inspects a single QueryEvent and optionally produces an Alert.
// implemntations must be safe for concurrent use.
type Analyzer interface {
	Name() string
	Analyzer(event models.QueryEvent) (models.Alert, bool)
}

type SlowQueryAnalyzer struct{}

// note: SlowQueryAnalyzer struct implements Analyzer interface

// SlowQueryAnlayzer fires when a query exceeds the slow query threshold.

func (s *SlowQueryAnalyzer) Name() string {
	return "Slow Query"
}

func (s *SlowQueryAnalyzer) Analyzer(e models.QueryEvent) (models.Alert, bool) {
	if e.Duration < slowQueryTH {
		return models.Alert{}, false
		// here returning false as this is not depicted as slow query
	}
	return models.Alert{
		Title:     "Slow Query Detected",
		Message:   fmt.Sprintf("Your query on table %q took %dms (threshold %dms). Please fix this.",
	e.Table, e.Duration, slowQueryTH),
		Severity:  models.SeverityWarning,
		Event:     e,
		Timestamp: time.Now(),
	}, true
}

// FailedQueryAnlayzer fires whenever a query return a non empty error string and not an expected result.
type FailedQueryAnalyzer struct{}

func (f *FailedQueryAnalyzer) Name() string {
	return "Failed Query"
}

func (f *FailedQueryAnalyzer) Analyzer(e models.QueryEvent) (models.Alert, bool) {
	if e.Error == "" {
		return models.Alert{}, false
	}

	return models.Alert{
		Title: "Query Failure",
		Message: fmt.Sprintf("Query on table %q failed: %s", e.Table, e.Error),
		Severity: models.SeverityCritical,
		Event: e,
		Timestamp: time.Now(),
	}, true
}

type IneffecientQueryAnalyzer struct{}

func(i *IneffecientQueryAnalyzer) Name() string {
	return "Ineffeicent Query"
}

func (i *IneffecientQueryAnalyzer) Analyze(e models.QueryEvent) (models.Alert, bool) {
	upper := strings.ToUpper(e.Query)

	hasSelectStart := strings.Contains(upper, "SELECT *")
	hasWhere := strings.Contains(upper, "WHERE")
	hasLimit := strings.Contains(upper, "LIMIT")

	if hasSelectStart && !hasWhere && !hasLimit &&e.Duration >= ineffQueryTH {
		return models.Alert{
			Title: "Ineffecient Query, Likely a full Table Scan attempt",
			Message: fmt.Sprintf(
				"SELECT * on %q with no WHERE/LIMIT ran for %dms. Add a filter a pagination.",
				e.Table, e.Duration,
			),
			Severity: models.SeverityWarning,
			Event: e,
			Timestamp: time.Now(),
		}, true
	}
	return models.Alert{}, false
}

// Think about how can we predict N+1 query other than normal window approach.


// RunAnalytics is the analytics branch of the TEE pipline
// it runs every incoming queryevent through all registered list of analyzers
// feeds the latency/error metrics trackers, and emits Alert on the returned channel.
// The returned channel is closed when it is drained or ctx is cancelled

func RunAnalytics (ctx context.Context, in <- chan models.QueryEvent) <- chan models.Alert {
	out := make(chan models.Alert, 64)

}