package receive

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
)

func TestUpdateStats(t *testing.T) {
	cs := NewCardinalityStats()
	lbls := labels.Labels{
		labels.Label{Name: "__name__", Value: "metricA"},
		labels.Label{Name: "cluster", Value: "us"},
		labels.Label{Name: "instance", Value: "A"},
	}

	cs.UpdateStats(lbls, "cluster")

	if cs.totalSeries != 1 {
		t.Fatalf("expected totalSeries to be 1, got %v", cs.totalSeries)
	}

	if cs.totalLabelValuePairs != 3 {
		t.Fatalf("expected totalLabelValuePairs to be 3, got %v", cs.totalLabelValuePairs)
	}

	if cs.seriesCountByMetricName["metricA"] != 1 {
		t.Fatalf("expected seriesCountByMetricName[\"metricA\"] to be 1, got %v", cs.seriesCountByMetricName["metricA"])
	}

	if cs.seriesCountByLabelValuePair["__name__:metricA"] != 1 {
		t.Fatalf("expected seriesCountByLabelValuePair[\"__name__:metricA\"] to be 1, got %v", cs.seriesCountByLabelValuePair["__name__:metricA"])
	}

	if cs.labelValueCountByLabelName["__name__"] != 1 {
		t.Fatalf("expected labelValueCountByLabelName[\"__name__\"] to be 1, got %v", cs.labelValueCountByLabelName["__name__"])
	}

	if cs.seriesCountByFocusLabelValue["us"] != 1 {
		t.Fatalf("expected seriesCountByFocusLabelValue[\"us\"] to be 1, got %v", cs.seriesCountByFocusLabelValue["us"])
	}
}

func TestCalculateCardinalityStats(t *testing.T) {
	//need help here
}

func TestTopKMap(t *testing.T) {
	// create some test data
	testData := map[string]int{
		"A": 5,
		"B": 3,
		"C": 2,
		"D": 1,
		"E": 4,
	}

	// run the function
	result := topKMap(testData, 3)

	// check the size
	if len(result) != 3 {
		t.Fatalf("expected size to be 3, got %v", len(result))
	}

	// check the values
	if result["A"] != 5 {
		t.Fatalf("expected A to be 5, got %v", result["A"])
	}

	if result["E"] != 4 {
		t.Fatalf("expected E to be 4, got %v", result["E"])
	}

	if result["B"] != 3 {
		t.Fatalf("expected B to be 3, got %v", result["B"])
	}
}
