- 使用的端口
  - http 7001
  - tcp 9001
- control-center-server 订阅了所有服务的 key
- 服务节点信息在 redis 中的的格式
  - key = serviceName
  - value = serviceVersion + IP + port

