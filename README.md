# DogHandler

监听Prometheus WatchDog告警消息。


## 思路

```
dog handler 流程
http服务
接受post请求：
1.按配置创建timer
2a1.已存在timerID
2a2.重置counter
2a3.重置timer
2b1.不存在ID
2b2.创建timerID
2b3.创建timer
3.定时器结束，counter计数+1
4.counter计数达到配置，进行报警，重置counter

conf:
global:
  interval
  maxcounter
  alert:
    alertinfo
job:
  - id
    describe(option)
    interval(option)
    maxcounter(option)
    alert:
      alertinfo
```