// Package web 提供 HTTP 服务：K 线查询 + 指标脚本执行。
package web

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/chainball/mqx-lite/internal/jqs"
	"github.com/chainball/mqx-lite/internal/market"
)

//go:embed index.html
var indexHTML []byte

// Server 持有数据源与指标引擎。
type Server struct {
	provider market.Provider
	engine   jqs.Engine
}

// New 构造 Server。
func New(p market.Provider) *Server {
	return &Server{provider: p, engine: jqs.Engine{}}
}

// Handler 返回路由。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
	})
	mux.HandleFunc("/api/bars", s.handleBars)
	mux.HandleFunc("/api/indicators", s.handleIndicators)
	return mux
}

// handleBars GET /api/bars?code=000001&count=100 → ECharts candlestick 数据。
func (s *Server) handleBars(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		code = "000001"
	}
	series, err := s.provider.Bars(r.Context(), market.Query{Code: code, Period: market.PeriodDay, Count: 100})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	times := make([]string, len(series.Bars))
	values := make([][4]float64, len(series.Bars))
	for i, b := range series.Bars {
		times[i] = b.Time.Format("2006-01-02")
		// ECharts candlestick 顺序: [open, close, low, high]
		values[i] = [4]float64{b.Open, b.Close, b.Low, b.High}
	}
	writeJSON(w, map[string]any{"code": series.Code, "times": times, "bars": values})
}

// handleIndicators POST /api/indicators {code, script} → 指标结果。
func (s *Server) handleIndicators(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code   string `json:"code"`
		Script string `json:"script"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Code == "" {
		req.Code = "000001"
	}
	series, err := s.provider.Bars(r.Context(), market.Query{Code: req.Code, Period: market.PeriodDay, Count: 200})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	result, err := s.engine.Run(req.Script, series)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}
