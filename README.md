# util

`util` 是 infrago 通用工具包。

## 安装

```bash
go get github.com/infrago/util@latest
```

## 公开 API（摘自源码）

- `func (p hashringNodes) Len() int           { return len(p) }`
- `func (p hashringNodes) Less(i, j int) bool { return p[i].spotValue < p[j].spotValue }`
- `func (p hashringNodes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }`
- `func (p hashringNodes) Sort()              { sort.Sort(p) }`
- `func NewHashRing(weights map[string]int, spotsArgs ...int) *HashRing`
- `func (h *HashRing) Append(nodeKey string, weight int)`
- `func (h *HashRing) Remove(nodeKey string)`
- `func (h *HashRing) Update(nodeKey string, weight int)`
- `func (h *HashRing) Locate(s string) string`
- `func NewWRR(weights map[string]int) *WRR`
- `func (w *WRR) Next() string`
