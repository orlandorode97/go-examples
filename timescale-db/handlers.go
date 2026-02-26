package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Handler holds the database reference.
type Handler struct {
	DB *DB
}

// New creates a new Handler.
func NewHandlers(database *DB) *Handler {
	return &Handler{DB: database}
}

// PostMetric handles POST /metrics — inserts a new sensor reading.
func (h *Handler) PostMetric(w http.ResponseWriter, r *http.Request) {
	var m SensorMetric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if m.SensorID == "" || m.Location == "" {
		http.Error(w, "sensor_id and location are required", http.StatusBadRequest)
		return
	}
	if m.Time.IsZero() {
		m.Time = time.Now()
	}

	if err := h.DB.InsertMetric(r.Context(), m); err != nil {
		http.Error(w, "failed to insert metric: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

// GetMetrics handles GET /metrics — returns recent raw or bucketed readings.
//
// Query params:
//
//	hours  int  — how many hours back to look (default 24)
//	bucket int  — if set, aggregate into N-minute buckets using time_bucket()
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	hours := queryInt(r, "hours", 24)

	w.Header().Set("Content-Type", "application/json")

	if bucketMinutes := queryInt(r, "bucket", 0); bucketMinutes > 0 {
		// ⚡ TimescaleDB time_bucket() aggregation
		since := time.Now().Add(-time.Duration(hours) * time.Hour)
		data, err := h.DB.GetBucketedMetrics(r.Context(), bucketMinutes, since)
		if err != nil {
			http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if data == nil {
			data = []BucketedMetric{}
		}
		json.NewEncoder(w).Encode(map[string]any{
			"bucket_minutes": bucketMinutes,
			"hours":          hours,
			"count":          len(data),
			"data":           data,
		})
		return
	}

	// Raw readings
	data, err := h.DB.GetRecentMetrics(r.Context(), hours)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if data == nil {
		data = []SensorMetric{}
	}
	json.NewEncoder(w).Encode(map[string]any{
		"hours": hours,
		"count": len(data),
		"data":  data,
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Pool.Ping(r.Context()); err != nil {
		http.Error(w, `{"status":"unhealthy"}`, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func queryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
