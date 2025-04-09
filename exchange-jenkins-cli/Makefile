PKG := "exchange-jenkins-cli"

install: ## 安装依赖
	@go mod tidy

build: ## 构建二进制程序
	@go build -o ${PKG}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'