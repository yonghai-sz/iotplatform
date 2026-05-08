
<https://docs.emqx.com/zh/neuron/v2.3/south-devices/modbus-tcp/modbus-tcp.html>    


```sh
docker run \
    -d \
    --name my-single-neuron-container \
    --log-opt max-size=100m \
    -p 7000:7000 \
    --privileged=true \
    -v /home/chen/somevolumes/neuron/persistence:/opt/neuron/persistence \
    --restart=always \
    emqx/neuron:2.3.0-alpine
```

