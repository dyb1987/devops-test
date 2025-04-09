
# 编译构建
```shell
# 安装依赖
make install
# 编译
make build
```

# 使用方式

## 帮助信息
```shell
# 帮助信息
./exchange-jenkins-cli --help

```

## 配置文件
在项目根目录，增加config.toml，用于连接jenkins
```toml
[jenkins]
url = "http://jenkins.domain.com"
username = "admin"
password = "123456789"

```

## 示例
在项目根目录，创建 template/ 目录，然后指定输入源和模板文件

### 创建job
```shell
# 创建后端job
./exchange-jenkins-cli create --input apps.json --template apps.xml
# 创建前端web-ui job
./exchange-jenkins-cli create --input web.json --template web-ui.xml
# 创建前端h5-ui job
./exchange-jenkins-cli create --input web.json --template h5-ui.xml
# 创建前端excloud-ui job
./exchange-jenkins-cli create --input web.json --template excloud-ui.xml
# 创建前端manager-ui job
./exchange-jenkins-cli create --input web.json --template manager-ui.xml
# 创建前端partner-ui job
./exchange-jenkins-cli create --input web.json --template partner-ui.xml
# 创建前端wallet-ui job
./exchange-jenkins-cli create --input web.json --template wallet-ui.xml
```

### 运行项目
根据 apps.json 定义的分支名称、启动顺序、job 名称等，批量构建 job
```shell
# 运行项目，按照顺序批量自动化构建
./exchagne-jenkins-cli run --input apps.json
```


### 修改job
根据输入源 apps.json 和模板文件 config.xml 更新 job
```shell
# 更新多个job
./exchange-jenkins-cli update -n exchange-test-account -n exchagne-test-coin --input apps.json --template apps.xml
# 更新所有job
./exchange-jenkins-cli update --input apps.json --template apps.xml

```

### 删除job
根据视图批量删除 job
```shell
# 删除多个job
./exchange-jenkins-cli delete -n exchange-test-account -n exchagne-test-coin
# 批量删除job
./exchange-jenkins-cli delete --view exchange-test
```