# Go Fiber API Template

可直接复制作为新项目起点的 Go API 模板。核心只保留用户注册、登录、资料维护和通用基础设施，不包含商城业务。

业务模板通过 Go 模块依赖使用独立仓库 [kit](https://github.com/hina1314/kit)。kit 维护错误、响应、日志、国际化、CORS、连接池、指标、健康检查、停机和密码工具；业务项目维护用户规则、SQL、仓储、Token 与装配。本仓库不再包含框架源码，也不使用本地 `replace`。目前固定 kit 已发布版本 `v0.1.0`，后续可以更新模块版本获得框架修复。详见 [框架升级与迁移](docs/framework-upgrades.md)。

## 技术栈

- Go 1.26、Fiber v3
- PostgreSQL、`database/sql`、pgx、sqlc
- Wire 依赖注入、PASETO v2 访问令牌
- Prometheus 指标、结构化日志、健康检查和优雅停机
- Docker Compose、GitHub Actions、版本发布脚本

sqlc 生成类型安全的 SQL 调用；连接池由 pgx 的 `database/sql` 驱动管理。模板已配置连接上限、空闲连接、连接寿命和超时。

## 开始新项目

从模板创建仓库后，运行初始化脚本替换 Go 模块路径和应用名：

```powershell
./scripts/init-project.ps1 -Module github.com/your-name/your-project -AppName your-project
```

Linux/macOS：

```bash
./scripts/init-project.sh github.com/your-name/your-project your-project
```

上述命令保留模板 `go.mod` 中固定的 kit 版本，也可以在初始化时显式选择已发布版本：

```powershell
./scripts/init-project.ps1 -Module github.com/your-name/your-project -AppName your-project -KitVersion v0.1.0
```

```bash
./scripts/init-project.sh github.com/your-name/your-project your-project v0.1.0
```

默认版本为已发布的 `v0.1.0`；脚本也接受伪版本。初始化脚本只改业务模块导入，保留 `github.com/hina1314/kit` 引用；PowerShell 的旧参数 `-FrameworkVersion` 作为 `-KitVersion` 的别名兼容。

复制配置并修改数据库连接和 32 字节令牌密钥：

```bash
cp app.env.example app.env
make generate
DB_SOURCE='postgres://postgres:password@127.0.0.1:5432/app?sslmode=disable' make migrate_up
go run .
```

`app.env` 是可选文件，生产环境可以完全使用环境变量。全部配置项见 `app.env.example`。

## API

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/v1/signup` | 手机号或邮箱注册 |
| POST | `/v1/login` | 登录并签发 PASETO |
| GET | `/v1/user/info` | 获取当前用户 |
| PATCH | `/v1/user/profile` | 部分更新用户资料 |
| GET | `/livez` | 进程存活检查 |
| GET | `/readyz` | 数据库就绪检查 |
| GET | `/metrics` | Prometheus 指标，可用 `METRICS_TOKEN` 保护 |

注册和登录接口有按客户端限流。请求包含 request ID；未处理的 panic 会被恢复，服务端错误会记录具体原因，但响应不会泄露内部错误。

## 目录约定

```text
config/                 配置和翻译
db/migration/           数据库迁移
db/query/               手写 SQL
db/model/               sqlc 生成代码
internal/api/           HTTP handler、中间件、路由
internal/api/presentation/  业务错误状态映射及响应装配
internal/api/validation/    业务校验规则与校验错误转换
internal/app/           应用用例编排
internal/domain/        领域实体和接口
internal/domain/user/usererrors/  用户错误码（保留原 API 值）
internal/infra/         PostgreSQL 仓储实现
internal/di/            Wire 依赖注入
token/                  PASETO
ops/                    Docker、Prometheus、Nginx 运维文件
```

新增业务模块时，先写迁移和查询，再运行 `make sqlc`；补齐领域、仓储、应用服务和 handler 后运行 `make wire`。不要手改 `db/model` 和 `internal/di/wire_gen.go`。

Wire 通过带 `tools` 构建标签的 `tools.go` 保留在模块依赖中，版本由 `go.mod` 固定。使用 `go run github.com/google/wire/cmd/wire ./internal/di` 运行，兼容尚不支持 `go.mod` 中 `tool` 指令的 IDE。其 `x/tools` 版本也在项目中固定；不要给生成命令额外加 `@版本`，以免绕开项目依赖选择并触发旧生成器与当前 Go 工具链的兼容问题。

## 质量检查

```bash
make check
```

CI 会检查格式、`go vet`、sqlc/Wire 生成结果、竞态测试、数据库迁移、构建和真实 PostgreSQL 仓储测试。本机没有 PostgreSQL 时，集成测试会自动跳过。

模板的 `make test`、`make vet` 和 CI 检查业务模块；kit 的自身测试由独立仓库 CI 运行。模板测试会使用固定版本的 kit 验证业务兼容性。也可以手动运行：

```bash
go test ./...
```

## 管理后台建议

通用后台接口建议放在同一仓库的独立 `/v1/admin` 路由和应用模块中，共享领域与事务；只有当团队、发布周期或安全边界已经独立时再拆服务。管理前端建议作为单独前端项目，使用成熟组件库生成 CRUD 页面。认证、RBAC、审计日志、菜单权限、批量操作和导出属于框架下一阶段能力，不应和具体商城业务绑定。
