package jqs

import (
	"testing"
	"time"

	"github.com/chainball/mqx-lite/internal/market"
)

func TestIndicators(t *testing.T) {
	var bars []market.Bar
	for i := 0; i < 30; i++ {
		v := float64(10 + i)
		bars = append(bars, market.Bar{Time: time.Now(), Open: v, High: v + 1, Low: v - 1, Close: v, Volume: 100})
	}
	s := market.Series{Code: "000001", Bars: bars}
	e := Engine{}
	_, err := e.Run("return { MA5: MA(close, 5), EMA: EMA(close, 10), BOLL: BOLL(close, 20, 2), ATR: ATR(high, low, close, 14), RSI: RSI(close, 14), MACD: MACD(close) };", s)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
}

func TestBadScript(t *testing.T) {
	e := Engine{}
	s := market.Series{Code: "000001", Bars: []market.Bar{{Time: time.Now(), Close: 1}}}
	if _, err := e.Run("return 42;", s); err == nil {
		t.Fatal("expected error for non-object return")
	}
}
