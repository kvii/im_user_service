# im_user_service

open im 用户服务 demo。演示了已上线的用户系统如何对接 open im server。

## 运行

在工程根目录执行 `go run .` 运行工程。

## 用户

工程在 [db.go](./db.go) 中用 go “数组”模拟了一个用户数据库。在运行工程前，可以在该文件中添加用户。

## 配置文件

若要自定义 open im server 的 ip，请在当前目录创建 `config.json` 文件，内容如下：

```json
{
    "imAddr": "http://{你的 open im server ip}:10002"
}
```