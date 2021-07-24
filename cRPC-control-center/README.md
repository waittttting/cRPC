- 使用的端口
  - http 7001
  - tcp 9001
- control-center-server 订阅了所有服务的 key
- 服务节点信息在 redis 中的的格式
  - key = serviceName
  - value = serviceVersion + IP + port

- 服务的 IP 等节点信息，存储在 redis 内
- 当订阅的服务有节点变更时，control-center 向各个节点推订阅的节点变化信号时，下推全量的数据