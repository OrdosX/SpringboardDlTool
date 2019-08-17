# 跳板下载工具

平时下载偶尔会遇到一种国外下载很快、然后国内下载巨慢（只有几kb/s）的资源。这个项目提供了一种方案：借助架设在国外的VPS作为跳板，先将文件下到它上面，再取回本地。虽然是土办法，效率依然比直接本地下载高了不少。

## 服务端的部署

在服务器上安装golang环境，然后输入以下命令：

```
go get github.com/OrdosX/SpringboardDlTool/server
```

之后进入编译完的程序存放的位置（通常是 `~/go/bin` ），输入以下命令：

```
wget https://raw.githubusercontent.com/OrdosX/SpringboardDlTool/master/config.json
```

然后打开，填上相应参数：

```
auth: 密码，客户端与服务端保持一致
host: 服务端此参数无用
port: 监听的端口
```

完成后，记得配置防火墙开放所设置的端口，然后输入 `./server` 打开服务器。要后台运行，可以使用服务或者screen。

## 客户端的使用

安装流程及指令与服务端相同，只是将 `server` 换成 `client` ,并且config.json中 `host` 字段填服务器地址。

下载时，输入以下指令：

```
./client <URL> <要保存为的文件名>

例子：
./client https://example.com/example.zip ex.zip
```

完成后，在与程序同目录下的dl文件夹中可以找到下载的文件。