# 概览

分布式kv数据库，基于raft协议实现，针对简单kv进行存储并同步到各个节点。

部署要求`实例`为基数个。


# 使用方式

```shell
go build -o praetor .
```

- 启动node1: ./praetor console --http 127.0.0.1:7001 --raft 127.0.0.1:7000 --myid 1 --cluster 1/127.0.0.1:7000,2/127.0.0.1:8000,3/127.0.0.1:9000
- 启动node2: ./praetor console --http 127.0.0.1:8001 --raft 127.0.0.1:8000 --myid 2 --cluster 1/127.0.0.1:7000,2/127.0.0.1:8000,3/127.0.0.1:9000
- 启动node3: ./praetor console --http 127.0.0.1:9001 --raft 127.0.0.1:9000 --myid 3 --cluster 1/127.0.0.1:7000,2/127.0.0.1:8000,3/127.0.0.1:9000
- 添加：curl http://127.0.0.1:7001/set?key=test_key&value=test_value
- 获取：curl http://127.0.0.1:7001/get?key=test_key
- 删除：curl http://127.0.0.1:7001/get?key=test_key
- 节点情况: curl http://127.0.0.1:7001/state

*备注：添加和删除只能在leader上进行操作*

# 后续规划

- [ ] 日志文件记录
- [ ] 支持配置文件启动
- [ ] 任意节点都可以做添加和删除
- [ ] 存储算法继续优化目前只是map
