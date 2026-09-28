# Go Fiber API Template

可直接复制作为新项目起点的 Go API 模板。核心只保留用户注册、登录、资料维护和通用基础设施，不包含商城业务。

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
internal/app/           应用用例编排
internal/domain/        领域实体和接口
internal/infra/         PostgreSQL 仓储实现
internal/di/            Wire 依赖注入
token/                  PASETO
ops/                    Docker、Prometheus、Nginx 运维文件
```

新增业务模块时，先写迁移和查询，再运行 `make sqlc`；补齐领域、仓储、应用服务和 handler 后运行 `make wire`。不要手改 `db/model` 和 `internal/di/wire_gen.go`。

## 质量检查

```bash
make check
```

CI 会检查格式、`go vet`、sqlc/Wire 生成结果、竞态测试、数据库迁移、构建和真实 PostgreSQL 仓储测试。本机没有 PostgreSQL 时，集成测试会自动跳过。

## 管理后台建议

通用后台接口建议放在同一仓库的独立 `/v1/admin` 路由和应用模块中，共享领域与事务；只有当团队、发布周期或安全边界已经独立时再拆服务。管理前端建议作为单独前端项目，使用成熟组件库生成 CRUD 页面。认证、RBAC、审计日志、菜单权限、批量操作和导出属于框架下一阶段能力，不应和具体商城业务绑定。
