package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/thanos-io/thanos/pkg/receive"
)

type API struct {
	db       *tsdb.DB
	logger   log.Logger
	receiver *receive.CardinalityStats
}

func New(db *tsdb.DB, logger log.Logger) *API {
	return &API{
		db:       db,
		logger:   logger,
		receiver: receive.NewCardinalityStats(),
	}
}

func (api *API) Register(r *http.ServeMux) {
	r.HandleFunc("/api/v1/cardinality", api.cardinality)
}

//need help in creatigng api instance and creating API response, below is rough overview of what I am trying to do
func (api *API) cardinality(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters.
	topK := r.URL.Query().Get("topK")
	focusLabel := r.URL.Query().Get("focusLabel")
	matcherStr := r.URL.Query().Get("matcher")

	var topKVal int
	if topK != "" {
		var err error
		topKVal, err = strconv.Atoi(topK)
		if err != nil {
			http.Error(w, "Invalid topK value", http.StatusBadRequest)
			return
		}
	} else {
		topKVal = 10 // default value
	}

	var matcher *labels.Matcher
	if matcherStr != "" {
		var err error
		matcher, err = //How to parse matcherStr?
		if err != nil {
			http.Error(w, "Invalid matcher", http.StatusBadRequest)
			return
		}
	}

	// Calculate cardinality.
	err := api.receiver.CalculateCardinalityStats(api.db, focusLabel, []*labels.Matcher{matcher}, topKVal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Construct the API response.
	resp := map[string]interface{}{
		"status": "success",
		"data":   api.receiver,
	}

	// Encode the response as JSON.
	json.NewEncoder(w).Encode(resp)
}