package jqs

import (
	"math"

	"github.com/dop251/goja"
)

// injectIndicators 把内置算子注入 goja 运行时。
func injectIndicators(vm *goja.Runtime) {
	vm.Set("MA", MA)
	vm.Set("EMA", EMA)
	vm.Set("BOLL", BOLL)
	vm.Set("ATR", ATR)
	vm.Set("RSI", RSI)
	vm.Set("MACD", MACD)
}

// MA 简单移动平均，前 n-1 个位置为 0。
func MA(arr []float64, n int) []float64 {
	out := make([]float64, len(arr))
	var sum float64
	for i, v := range arr {
		sum += v
		if i >= n {
			sum -= arr[i-n]
		}
		if i >= n-1 {
			out[i] = sum / float64(n)
		}
	}
	return out
}

// EMA 指数移动平均。
func EMA(arr []float64, n int) []float64 {
	out := make([]float64, len(arr))
	k := 2.0 / float64(n+1)
	prev := 0.0
	for i, v := range arr {
		if i == 0 {
			prev = v
		} else {
			prev = v*k + prev*(1-k)
		}
		out[i] = prev
	}
	return out
}

// BOLL 布林带，返回 {upper, mid, lower} 三条序列。
func BOLL(cls []float64, n int, k float64) map[string][]float64 {
	mid := MA(cls, n)
	upper := make([]float64, len(cls))
	lower := make([]float64, len(cls))
	for i := n - 1; i < len(cls); i++ {
		var sum float64
		for j := i - n + 1; j <= i; j++ {
			d := cls[j] - mid[i]
			sum += d * d
		}
		std := math.Sqrt(sum / float64(n))
		upper[i] = mid[i] + k*std
		lower[i] = mid[i] - k*std
	}
	return map[string][]float64{"upper": upper, "mid": mid, "lower": lower}
}

// ATR 平均真实波幅（Wilder 平滑）。
func ATR(high, low, cls []float64, n int) []float64 {
	tr := make([]float64, len(cls))
	for i := range cls {
		if i == 0 {
			tr[i] = high[i] - low[i]
			continue
		}
		h := high[i] - low[i]
		hc := math.Abs(high[i] - cls[i-1])
		lc := math.Abs(low[i] - cls[i-1])
		tr[i] = math.Max(h, math.Max(hc, lc))
	}
	out := make([]float64, len(tr))
	prev := 0.0
	for i, v := range tr {
		if i == 0 {
			prev = v
		} else {
			prev = (prev*float64(n-1) + v) / float64(n)
		}
		out[i] = prev
	}
	return out
}

// RSI 相对强弱指标（Wilder 平滑），前 n 个位置为 0。
func RSI(cls []float64, n int) []float64 {
	out := make([]float64, len(cls))
	var avgGain, avgLoss float64
	for i := 1; i < len(cls); i++ {
		ch := cls[i] - cls[i-1]
		var gain, loss float64
		if ch > 0 {
			gain = ch
		} else {
			loss = -ch
		}
		if i <= n {
			avgGain += gain / float64(n)
			avgLoss += loss / float64(n)
		} else {
			avgGain = (avgGain*float64(n-1) + gain) / float64(n)
			avgLoss = (avgLoss*float64(n-1) + loss) / float64(n)
		}
		if i >= n {
			if avgLoss == 0 {
				out[i] = 100
			} else {
				out[i] = 100 - 100/(1+avgGain/avgLoss)
			}
		}
	}
	return out
}

// MACD 异同移动平均线，返回 {dif, dea, macd} 三条序列。
func MACD(cls []float64) map[string][]float64 {
	ema12 := EMA(cls, 12)
	ema26 := EMA(cls, 26)
	dif := make([]float64, len(cls))
	for i := range cls {
		dif[i] = ema12[i] - ema26[i]
	}
	dea := EMA(dif, 9)
	macd := make([]float64, len(cls))
	for i := range cls {
		macd[i] = 2 * (dif[i] - dea[i])
	}
	return map[string][]float64{"dif": dif, "dea": dea, "macd": macd}
}
