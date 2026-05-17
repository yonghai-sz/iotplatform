

`git submodule update --remote .cursorrules`  
`git submodule update --remote .github/ai-context`  


Updates:  
`git submodule update --remote --recursive`  




### Endpoints
- **Shorten**:
  - `GET http://localhost:8080/shorten?url=https://example.com`
- **Expand**:
  - `GET http://localhost:8080/expand?shorten=<key>`


test the API Gateway service
`curl -i "http://localhost:8080/shorten?url=http://www.example.cn"`
`curl -i "http://localhost:8080/expand?shorten=fb5cd9"`


### Useful ports
- **API**: `8080`
- **MySQL**: `3307` 
- **Redis**: `6379`
- **Etcd**: `2379`




# 系统架构图或数据流图
# Mermaid 图表




`docker build -f services/internal-api/Dockerfile --target runtime -t <namespace>/<repo>-internal-api:<tag> .`

```text
Why am I getting "http: named cookie not present" error?

Review my handler code for go-zero anti-patterns
```




# Check project file structure
`find . -name "*.api" -o -name "*.proto" -o -name "*.go" | head -50`  
