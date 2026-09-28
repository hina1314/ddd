# kit 依赖、发布和业务项目升级

## 仓库边界

| 仓库 | Go 模块 | 职责 |
| --- | --- | --- |
| 业务模板 | `github.com/hina1314/ddd` | 用户功能、SQL、仓储、事务、Token、配置、路由、Wire、部署 |
| [kit](https://github.com/hina1314/kit) | `github.com/hina1314/kit` | 通用错误、响应、日志、国际化、CORS、连接池、指标、健康检查、停机、密码工具 |

kit 已迁到独立仓库，模板不再保存框架源码，也不包含框架发布工作流。业务和 kit 分别测试、提交、发布。公共包不导入模板的业务类型；业务扩展留在用户错误码、校验、响应状态映射、翻译、配置和依赖装配中。

## 当前固定版本

模板直接从远端下载已发布版本，不需要同级 kit 目录：

```go
require github.com/hina1314/kit v0.1.0
```

kit 已发布 `v0.1.0`，模板固定该版本，不使用 `latest` 作为构建依赖。Docker 和 CI 都通过 `go.mod`、`go.sum` 下载相同版本。

## 新项目初始化

不指定 kit 版本时保留模板当前固定版本：

```powershell
./scripts/init-project.ps1 -Module example.com/team/newapp -AppName newapp
```

```bash
./scripts/init-project.sh example.com/team/newapp newapp
```

可以显式选择已发布的 `v0.1.0`：

```powershell
./scripts/init-project.ps1 -Module example.com/team/newapp -AppName newapp -KitVersion v0.1.0
```

```bash
./scripts/init-project.sh example.com/team/newapp newapp v0.1.0
```

脚本接受正式版本、预发布版本和伪版本；PowerShell 保留 `-FrameworkVersion` 参数别名。初始化会重命名业务模块、移除 kit 的本地替换、整理依赖、生成 Wire 并测试，不修改 kit 导入路径。

## kit 发布与项目升级

在 kit 仓库中提交代码并合入 main，等待独立 CI 的格式、依赖、vet、竞态测试与构建检查通过。后续发布使用新的仓库根模块标签，例如：

```bash
git tag v0.1.1
git push origin v0.1.1
```

以上后续发布命令是示例，未在本次修复中执行。kit 的模块在仓库根目录，标签直接使用 `v0.1.1` 这样的版本号，不添加子目录前缀；已发布的 `v0.1.0` 不应重复创建或修改。

业务项目升级到实际已发布的版本：

```bash
go get github.com/hina1314/kit@v0.1.0
go mod tidy
go run github.com/google/wire/cmd/wire ./internal/di
go test ./...
go build ./...
```

发布标签前也可以显式指定已推送的提交，Go 会记录其伪版本。审查模块文件变更与迁移说明，再部署。兼容修复通常只需要更新依赖；配置、装配 API、SQL、数据库迁移和部署文件的变化仍需按迁移说明处理。

升级后先运行 `go mod tidy`，再把 `go.mod`、`go.sum` 一起提交。CI 会重新整理并检查这两个文件没有差异；旧伪版本校验记录未清理就提交，会在此步骤失败。该检查应保留。

## 本地共同开发

本机 kit 路径是 `F:\project\go\kit`。在业务模板根目录可使用本地 Go workspace：

```powershell
go work init . ../kit
```

workspace 文件已加入模板忽略列表，不应提交到业务仓库。它只用于联调，会覆盖 kit 的远端版本；交付前可关闭 workspace 验证固定依赖：

```powershell
$env:GOWORK = 'off'
go test ./...
go build ./...
Remove-Item Env:GOWORK
```

不要把指向个人磁盘路径的 `replace` 提交到模板，否则 Docker、CI 和其他开发者无法独立构建。

## 原模板迁移与验证

迁移已有业务项目时保留业务逻辑、迁移和仓储，接入 kit 包并把定制转成扩展配置：

- 错误、响应、国际化和日志：替换为 kit 包；业务错误码、状态映射、手机号校验和项目翻译留在业务侧。
- 连接池：`database.Open` 提供池初始化，项目注册 SQL 驱动、构造自己的 Store，并管理资源关闭。
- 指标：每个应用创建 `metrics.New`，注册连接池指标；中间件顺序为 request ID → Metrics → Logger → recover → 其他中间件与路由。
- 健康与停机：`health.Probe` 接收项目依赖的检查函数，`lifecycle.Run` 协调摘流和排空；监控响应取消，入口在 HTTP 关闭流程结束后关闭数据库。
- CORS 与密码：使用 kit 的 `middleware.Cors`、`password`；业务密码规则留在项目中。

模板已完成这些接入。验证登录、注册、鉴权的错误码、状态、中文提示及单条日志关联，再验证健康接口、指标和停机。kit 自身测试在独立仓库执行；模板真实 PostgreSQL 测试依赖 `TEST_DATABASE_URL`，未配置时跳过。
