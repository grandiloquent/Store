# Store

## 引用库

* https://github.com/gorilla
* https://github.com/jackc
* https://github.com/russross/blackfriday
* https://github.com/css/csso
* https://github.com/thebird/Swipe
* https://github.com/tweenjs/tween.js
* https://github.com/golang

## 部署

`settings/settings.json`：设置服务器相关参数

* [打包](./scripts/publish.go)
* [SSH](./scripts/ssh/server.go)
* [数据库](./SQL.md)

## Troubleshooting

### 无法安装 `golang.org/x`的包
 
```

go get github.com/jackc/pgx

mkdir -p $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/xerrors.git

```

### 无法侦听 `5050` 端口 
```shell script
# 打印占用5050端口的程序
$ netstat -ltnp | grep -w ':5050' 

```