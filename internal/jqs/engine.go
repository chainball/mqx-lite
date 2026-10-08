// Package jqs 实现 java-quant-script 指标脚本引擎（精简版）。
//
// 用户用 JS 写指标脚本，引擎注入 K 线数据（open/high/low/close/volume 数组）
// 与一组内置算子（MA/EMA/BOLL/ATR/RSI/MACD），执行后返回指标序列供展示层叠加到 K 线。
package jqs

import (
	"fmt"

	"github.com/dop251/goja"

	"github.com/chainball/mqx-lite/internal/market"
)

// Result 指标输出：指标名 → 单序列（[]float64）或多序列（map[string][]float64）。
type Result map[string]any

// Engine 指标脚本执行引擎。
type Engine struct{}

// Run 注入 K 线数据与算子，执行脚本，返回指标结果。
//
// 脚本约定：return 一个对象，键为指标名，值为序列，例如
//
//	MA5: MA(close, 5),
//	BOLL: BOLL(close, 20, 2),   // {upper, mid, lower}
func (Engine) Run(script string, s market.Series) (Result, error) {
	vm := goja.New()
	injectBars(vm, s)
	injectIndicators(vm)

	// 脚本包成 IIFE，使顶层 return 合法。
	v, err := vm.RunString("(function(){\n" + script + "\n})()")
	if err != nil {
		return nil, fmt.Errorf("jqs: script failed: %w", err)
	}
	exported, ok := v.Export().(map[string]any)
	if !ok {
		return nil, fmt.Errorf("jqs: script must return an object {name: series}, got %T", v.Export())
	}
	return Result(exported), nil
}

// injectBars 把 K 线各字段注入为等长 JS 数组。
func injectBars(vm *goja.Runtime, s market.Series) {
	n := len(s.Bars)
	open := make([]float64, n)
	high := make([]float64, n)
	low := make([]float64, n)
	cls := make([]float64, n)
	volume := make([]float64, n)
	for i, b := range s.Bars {
		open[i] = b.Open
		high[i] = b.High
		low[i] = b.Low
		cls[i] = b.Close
		volume[i] = b.Volume
	}
	vm.Set("open", open)
	vm.Set("high", high)
	vm.Set("low", low)
	vm.Set("close", cls)
	vm.Set("volume", volume)
}
