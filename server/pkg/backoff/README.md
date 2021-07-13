# backoff

提供简单的随机递增backoff算法，用于多次重试等场合

## 用法

```go
func retry() {
    min := 50 * time.Millisecond
    max := 5 * time.Second
    backoff := New(
        Min(min),
        Max(max),
    )
    for count := 0; count < 10; count++ {
        d := backoff.Duration()
        t.Logf("count: %d, duration: %v\n", count, d)
        time.Sleep(d)
    }
}
```