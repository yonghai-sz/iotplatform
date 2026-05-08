
`goctl api -o shorturl.api`

`goctl api go -api api-dsl/main.api -dir .`

`goctl api swagger --api api/shorturl.api --dir api --filename shorturl`


```
		val := r.Header.Get("User-Agent")
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, "User-Agent", val)
		newReq := r.WithContext(ctx)
		next(w, newReq)
```


要修改的地方：  
* etc/example.yaml
* internal/config/config.go
* internal/svc/servicecontext.go
* internal/logic/*.go






```text
Generate code from the userservice.api file in ./userservice
```

