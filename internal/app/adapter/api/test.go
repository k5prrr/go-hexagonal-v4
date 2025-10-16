package api

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

func (r *Router) testRouter() {
	r.HandleFunc("increment", r.testRouterIncrement)
}

var testCounter uint64

func (r *Router) testRouterIncrement(w http.ResponseWriter, rec *http.Request) {
	if rec.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	v := atomic.AddUint64(&testCounter, 1)
	resp := map[string]uint64{"value": v}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
