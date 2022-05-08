# chaoyue-golib

超越专用，golang lib 库

## 一、使用案例
### 1、反向代理
```go
func main() {
	r := gin.New()
	r.Use(core())
	r.GET("/api/v1/info", func(ctx *gin.Context) {
		http.RedirectHandler(ctx, urlAddress, "/api/v1")
	})
	r.Run(":8080")
}
```