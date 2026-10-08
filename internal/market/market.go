// Package market 定义行情数据接入规范：极简 Bar 时间序列 + Provider 窄接口。
//
// 指标引擎与展示层只依赖 Provider 接口，不直接触碰具体数据源，
// go2tdx-lite 是它的一个实现。更换数据源只需再实现 Provider。
package market

import (
	"context"
	"time"
)

// Bar 一根 K 线。
type Bar struct {
	Time                 time.Time
	Open, High, Low, Close float64
	Volume               float64
}

// Period 序列周期。
type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
)

// Series 一段 K 线时间序列。
type Series struct {
	Code   string
	Period Period
	Bars   []Bar
}

// Query 数据查询。
type Query struct {
	Code   string
	Period Period
	Count  int // 取最近 N 根
}

// Provider 数据接入窄接口。
type Provider interface {
	Bars(ctx context.Context, q Query) (Series, error)
}
