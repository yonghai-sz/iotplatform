
`goctl api -o platform.api`

`goctl api go -api api-dsl/main.api -dir .`  
`goctl api go -api internal.api -dir .`

`goctl api swagger --api api/platform.api --dir api --filename platform`


`goctl api validate -api <file>.api`

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

Create a user management API with CRUD operations

Add rate limiting and circuit breaker to my API

Create a user management API with authentication

Create a user management REST API with CRUD operations
```

# List API routes
`goctl api doc -dir .`  
