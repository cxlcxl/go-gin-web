### Go-Gin-web

gin 框架用于 web 快速开发骨架，只接入了 mysql+redis 的简单服务

#### 使用的服务
- jwt
- casbin
- gorm
- redis
- logrus
- validator
- air

#### 搭建
1. 下载后配置 mysql，redis 参数
2. go mod init [app-name]
3. go mod tidy
4. 项目根目录运行
```shell
air
```

运行成功后直接访问参数配置中的端口即可