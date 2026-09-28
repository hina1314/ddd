# 实验服务器全新部署

在服务器上以 `liao` 登录执行。GitHub 的 `DEPLOY_USER` 也设置为 `liao`。新目录为 `/opt/apps/app`，新辅助脚本为 `/usr/local/sbin/app-deploy`。

以下清理步骤会删除旧服务、旧配置和旧 PostgreSQL 数据卷，适用于已经决定放弃旧数据的全新部署。不要在保留数据的项目中使用。

## 清理旧部署

先保存旧仓库地址，再停止旧 Compose 项目并删除其数据卷：

```bash
set -e
cd /opt/apps
repo_url=$(git -C /opt/apps/ddd remote get-url origin)
cd /opt/apps/ddd/ops/docker
sudo docker compose down --volumes --remove-orphans
cd /opt/apps
sudo rm -rf -- /opt/apps/ddd
sudo rm -f -- /usr/local/sbin/study-deploy /etc/sudoers.d/study-deploy
```

不要执行全局 Docker 清理；上述操作只针对旧 Compose 项目。若旧 sudo 权限写在其他文件中，可通过 `sudo visudo` 移除旧 `study-deploy` 的条目。

## 创建新部署

继续在同一个 SSH 终端执行，使用前面保存的 `repo_url`：

```bash
sudo install -d -o liao -g "$(id -gn liao)" -m 0755 /opt/apps/app
git clone --branch main "$repo_url" /opt/apps/app
cd /opt/apps/app
sha=$(git rev-parse HEAD)
sudo bash ops/docker/prepare.sh
sudo sed -i "s|^APP_IMAGE=.*|APP_IMAGE=ghcr.io/hina1314/ddd:sha-$sha|" ops/docker/.env
sudo install -o root -g root -m 0755 ops/docker/deploy.sh /usr/local/sbin/app-deploy
printf '%s\n' 'liao ALL=(root) NOPASSWD: /usr/local/sbin/app-deploy' | sudo tee /etc/sudoers.d/app-deploy >/dev/null
sudo chmod 0440 /etc/sudoers.d/app-deploy
sudo visudo -cf /etc/sudoers.d/app-deploy
```

`prepare.sh` 自动生成新的数据库密码和令牌密钥，不会打印它们。不要在原配置上重新运行；当前新克隆目录没有旧配置，适合运行一次。

新服务使用 `API_HOST_PORT=3001`，只绑定服务器回环地址。需要其他端口时通过 `sudoedit ops/docker/.env` 修改。Docker Engine、Compose 插件、Git、OpenSSL、curl 和 flock 应已安装。

## 启动与验证

确认该提交的 GitHub Actions 镜像发布步骤已经成功。若 GHCR 镜像私有，先在服务器使用具有镜像读取权限的凭据执行 `sudo docker login ghcr.io`，让部署辅助脚本使用的 root Docker 客户端能够拉取镜像。

```bash
cd /opt/apps/app/ops/docker
sudo docker compose config --quiet
sudo docker compose up -d db
sudo -n /usr/local/sbin/app-deploy "$sha"
sudo docker compose ps
curl --fail --include http://127.0.0.1:3001/readyz
```

辅助脚本会拉取相同提交的镜像、执行迁移、启动 API，并检查 `X-App-Version: sha-<提交号>`。镜像不存在时先修复或完成镜像发布，不要换用不对应当前提交的标签。

确认 GitHub 的 `DEPLOY_HOST` 指向此服务器、`DEPLOY_USER=liao`，已有 SSH 密钥能登录该账号，且该账号能无交互执行 `git fetch origin main`。然后从 main 重跑 Deploy 工作流。以后正常部署无需再次清理数据卷或生成密码。
