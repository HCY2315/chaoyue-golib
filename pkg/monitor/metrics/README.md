# metrics

> 监控指标客户端，引用自[`go-metrics`](https://github.com/armon/go-metrics)，

## 增加的功能有：

* `statsd/statsite`客户端增加`set`类型的指标

削减的功能有：

* `circonus`：由于依赖包编译出错，且不常用

---

## 注意事项

* **本包中**statsd用udp传输metrics，statsite用tcp传输metrics
