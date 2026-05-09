

When working on your **temporary** branch, you can:  
`git push --force-with-lease`  





`ssh -i /Users/chenyonghai/Desktop/SingaporeTest.pem root@47.84.204.211`    
`uname -m`    



Copies the contents of the id_ed25519.pub file to your clipboard     
`pbcopy < ~/.ssh/id_ed25519.pub`    




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
