# crispy-garbanzo

## 1、开发环境

+ Go 版本

  - Go 1.22.2

+ Web框架

  - Gin

+ ORM库

  - gorm.io/gorm + MySQL

+ 缓存库

  - github.com/go-redis/redis/v8

## 2、基础环境安装与运行

```
cd crispy-garbanzo
go get
go run main.go
```

## 3、配置文件

+ 系统配置

  - config/config.yaml

+ vscode启动配置

```
{
    "name": "crispy-garbanzo",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceRoot}/main.go",
    "env": {},
    "args": []
}
```

## 4 生成swaggo接口文档
- [gin-swagger](http://github.com/swaggo/gin-swagger)
- [swaggo-files](http://github.com/swaggo/files)
- [swaggo-swag](http://github.com/swaggo/swag)

```
http://localhost:8000/swagger/index.html
```