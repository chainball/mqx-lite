# mqx-lite

K 线展示 + jqs 指标脚本引擎。Go 后端，数据源走 go2tdx-lite 实时行情（无数据库），指标用 goja 执行 JS 脚本，前端 ECharts 展示。

## 特性

- **通用数据接入规范**：`market.Provider` 窄接口，指标引擎与展示层只依赖它，更换数据源只需再实现一个 Provider
- **jqs 指标脚本引擎**：goja 执行 JS，内置 MA / EMA / BOLL / ATR / RSI / MACD
- **Web 展示**：ECharts K 线图 + JS 脚本编辑器，指标实时叠加
- **无数据库**：数据源直连 go2tdx-lite 行情服务器，零持久化

## 架构

```
浏览器（ECharts K线 + JS 编辑器）
        │  HTTP
   Go 后端（net/http）
   ├── jqs 引擎：goja + 内置算子
   ├── market.Provider 接口（数据接入规范）
   └── tdx 适配器：go2tdx-lite → market.Series
        │  7709 TCP
   通达信行情服务器
```

## 快速开始

```bash
go run ./cmd/server
# 打开 http://localhost:8080，输入代码（如 000001）加载 K 线
```

环境变量：

| 变量 | 默认 | 说明 |
|---|---|---|
| `ADDR` | `:8080` | HTTP 监听地址 |

## jqs 脚本

脚本在浏览器编辑器里编写，包成 IIFE 执行，直接用 `return` 返回指标对象。内置数据数组：`open` `high` `low` `close` `volume`（等长）。

内置算子：

| 算子 | 签名 | 返回 |
|---|---|---|
| `MA` | `MA(arr, n)` | `[]float64` |
| `EMA` | `EMA(arr, n)` | `[]float64` |
| `BOLL` | `BOLL(close, n, k)` | `{upper, mid, lower}` |
| `ATR` | `ATR(high, low, close, n)` | `[]float64` |
| `RSI` | `RSI(close, n)` | `[]float64` |
| `MACD` | `MACD(close)` | `{dif, dea, macd}` |

示例：

```js
return {
  MA5: MA(close, 5),
  MA20: MA(close, 20),
  BOLL: BOLL(close, 20, 2),
  MACD: MACD(close),
};
```

## API

- `GET /api/bars?code=000001` → K 线数据（ECharts candlestick 格式）
- `POST /api/indicators` `{code, script}` → 指标结果

## 数据接入规范

```go
type Provider interface {
	Bars(ctx context.Context, q Query) (Series, error)
}
```

任何数据源实现这个接口即可接入。`internal/tdx` 是基于 go2tdx-lite 的实现。

## 免责声明

本项目的行情数据来自 go2tdx-lite 所接入的通达信 7709 协议服务器，仅供学习、研究用途。使用者需自行确认数据来源的合法授权，见 [go2tdx-lite](https://github.com/chainball/go2tdx-lite) 的免责声明。

## License

[MIT](LICENSE)
