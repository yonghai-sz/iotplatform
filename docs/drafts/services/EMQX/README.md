
<https://docs.emqx.com/en/emqx/v5.8/>

```sh
docker container exec -it my-single-emqx-container /bin/bash
docker container stop my-single-emqx-container
docker container rm my-single-emqx-container
```




Config precedence order:  
etc/base.hocon < cluster.hocon < emqx.conf < environment variables
```sh
docker container cp -a my-single-emqx-container:/opt/emqx/etc/emqx.conf /Users/chenyonghai/Desktop/
docker container cp -a my-single-emqx-container:/opt/emqx/data/configs/cluster.hocon /Users/chenyonghai/Desktop/
docker container cp -a my-single-emqx-container:/opt/emqx/etc/base.hocon /Users/chenyonghai/Desktop/
```


`docker container cp -a my-single-emqx-container:/opt/emqx/etc/acl.conf /Users/chenyonghai/Desktop/` 
`docker container cp -a my-single-emqx-container:/opt/emqx/data/authz/acl.conf /Users/chenyonghai/Desktop/` 





<https://docs.emqx.com/en/emqx/v5.8/configuration/configuration.html>
