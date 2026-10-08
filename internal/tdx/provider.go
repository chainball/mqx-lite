// Package tdx 用 go2tdx-lite 实现 market.Provider：通达信行情 → market.Series。
package tdx

import (
	"context"

	"github.com/chainball/go2tdx-lite"
	"github.com/chainball/mqx-lite/internal/market"
)

// Provider 用 go2tdx-lite 客户端提供 K 线数据。
type Provider struct {
	client *go2tdx.Client
}

// New 建客户端（DialBest 从内置服务器列表测速选优）。
func New() (*Provider, error) {
	c, err := go2tdx.DialBest(32)
	if err != nil {
		return nil, err
	}
	return &Provider{client: c}, nil
}

// Close 释放连接。
func (p *Provider) Close() error {
	return p.client.Close()
}

// Bars 实现 market.Provider，取最近 count 根 K 线。
func (p *Provider) Bars(_ context.Context, q market.Query) (market.Series, error) {
	count := uint16(q.Count)
	if count == 0 {
		count = 100
	}
	bars, err := p.client.Bars().Get(q.Code, periodOf(q.Period), 0, count)
	if err != nil {
		return market.Series{}, err
	}
	out := market.Series{Code: q.Code, Period: q.Period, Bars: make([]market.Bar, 0, len(bars))}
	for _, b := range bars {
		out.Bars = append(out.Bars, market.Bar{
			Time:   b.Time,
			Open:   b.Open,
			High:   b.High,
			Low:    b.Low,
			Close:  b.Close,
			Volume: b.Volume,
		})
	}
	return out, nil
}

func periodOf(p market.Period) go2tdx.Period {
	switch p {
	case market.PeriodWeek:
		return go2tdx.Week
	case market.PeriodMonth:
		return go2tdx.Month
	default:
		return go2tdx.Day
	}
}
