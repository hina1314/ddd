# Runbook: StudyAPIDown

## 告警含义

Prometheus 连续 15 秒无法抓取：

http://127.0.0.1:3000/metrics

这不一定代表程序崩溃，也可能是端口、网络或 metrics 接口故障。

## 1\. 记录现场

记录：

* 告警开始时间
* instance
* environment
* Prometheus Targets 中的 Last Scrape Error
* 服务最后几行日志

不要一看到告警就立刻重启，否则可能破坏用于定位的现场信息。

## 2\. 检查服务入口

执行：

curl.exe --noproxy "*" -i --max-time 3 http://127.0.0.1:3000/livez
curl.exe --noproxy "*" -i --max-time 3 http://127.0.0.1:3000/readyz
curl.exe --noproxy "\*" -i --max-time 3 http://127.0.0.1:3000/metrics

判断：

* 三个接口都连接失败：检查进程和端口。
* livez=200、readyz=503：服务存活，但依赖或停机状态异常。
* livez=200、metrics失败：指标路由或中间件异常。
* 三个接口都正常、up仍为0：检查Prometheus配置和网络。

## 3\. 检查3000端口

执行：

Get-NetTCPConnection -LocalPort 3000 -State Listen -ErrorAction SilentlyContinue

如果有结果，查看占用进程：

$connection = Get-NetTCPConnection -LocalPort 3000 -State Listen |
Select-Object -First 1

Get-Process -Id $connection.OwningProcess

如果没有结果，说明没有程序监听3000端口。

## 4\. 检查启动失败原因

在项目目录执行：

go run .

常见错误：

* invalid configuration：配置错误。
* ping database：数据库不可用。
* address already in use：端口被其他进程占用。
* panic：程序缺陷。
* dependency initialization failed：依赖初始化失败。

先根据错误恢复依赖或配置，不要无条件重复重启。

## 5\. 恢复验证

恢复后必须全部满足：

* /livez 返回200
* /readyz 返回200
* /metrics 返回200
* up{job="study-api"} 返回1
* StudyAPIDown发送resolved
* 核心业务接口通过冒烟测试

仅仅看到进程存在，不代表服务已经恢复。

## 6\. 事后记录

记录：

* 发生时间
* 发现时间
* 恢复时间
* 用户影响
* 根因
* 临时止损措施
* 永久修复措施

