package receive

import (
	"sort"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunks"
)

// CardinalityStats holds the cardinality statistics
type CardinalityStats struct {
	totalSeries                  int            `json:"totalSeries"`
	totalLabelValuePairs         int            `json:"totalLabelValuePairs"`
	seriesCountByMetricName      map[string]int `json:"seriesCountByMetricName"`
	seriesCountByLabelName       map[string]int `json:"seriesCountByLabelName"`
	seriesCountByLabelValuePair  map[string]int `json:"seriesCountByLabelValuePair"`
	labelValueCountByLabelName   map[string]int `json:"labelValueCountByLabelName"`
	seriesCountByFocusLabelValue map[string]int `json:"seriesCountByFocusLabelValue"`
}

type Pair struct {
	Key   string
	Value int
}

func NewCardinalityStats() *CardinalityStats {
	return &CardinalityStats{
		seriesCountByMetricName:      make(map[string]int),
		seriesCountByLabelName:       make(map[string]int),
		seriesCountByLabelValuePair:  make(map[string]int),
		labelValueCountByLabelName:   make(map[string]int),
		seriesCountByFocusLabelValue: make(map[string]int),
	}
}

func (cs *CardinalityStats) UpdateStats(lbls labels.Labels, focusLabel string) {
	cs.totalSeries++

	for _, l := range lbls {
		cs.seriesCountByLabelName[l.Name]++
		cs.seriesCountByLabelValuePair[l.Name+":"+l.Value]++

		if l.Name == focusLabel {
			cs.seriesCountByFocusLabelValue[l.Value]++
		}

		if l.Name == "__name__" {
			cs.seriesCountByMetricName[l.Value]++
		}

		cs.totalLabelValuePairs++
		cs.labelValueCountByLabelName[l.Name]++
	}
}

func (cs *CardinalityStats) CalculateCardinalityStats(db *tsdb.DB, focusLabel string, matchers []*labels.Matcher, topK int) error {
	head := db.Head()
	indexReader, err := head.Index()
	if err != nil {
		return err
	}

	postings, err := indexReader.Postings("", "")
	if err != nil {
		return err
	}

	var chks []chunks.Meta
	var builder labels.ScratchBuilder
	for postings.Next() {
		p := postings.At()
		if err := indexReader.Series(p, &builder, &chks); err != nil {
			return err
		}
		lbls := builder.Labels()

		for _, matcher := range matchers {
			value := lbls.Get(matcher.Name)
			if value != "" {
				if !matcher.Matches(value) {
					continue
				}
			}
		}
		cs.UpdateStats(lbls, focusLabel)
	}

	cs.seriesCountByMetricName = topKMap(cs.seriesCountByMetricName, topK)
	cs.seriesCountByLabelName = topKMap(cs.seriesCountByLabelName, topK)
	cs.seriesCountByLabelValuePair = topKMap(cs.seriesCountByLabelValuePair, topK)
	cs.labelValueCountByLabelName = topKMap(cs.labelValueCountByLabelName, topK)
	cs.seriesCountByFocusLabelValue = topKMap(cs.seriesCountByFocusLabelValue, topK)

	return nil
}

func topKMap(m map[string]int, topK int) map[string]int {
	pairs := make([]Pair, len(m))
	i := 0
	for k, v := range m {
		pairs[i] = Pair{k, v}
		i++
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	if len(pairs) > topK {
		pairs = pairs[:topK]
	}

	topKMap := make(map[string]int, len(pairs))
	for _, p := range pairs {
		topKMap[p.Key] = p.Value
	}

	return topKMap
}
