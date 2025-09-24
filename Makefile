SHELL := /bin/bash

# 自动查找所有 gen.sh
API_SERVICES := $(wildcard api/*/gen.sh)
APP_SERVICES := $(wildcard app/*/gen.sh)

.PHONY: all api rpc docker kube build gen clean help

# 默认同时生成 api + rpc
all: api rpc

# 生成所有 API 服务
api:
	@echo "==> Generating API code for all services…"
	@for f in $(API_SERVICES); do \
		svc_dir=$$(dirname $$f); \
		echo "----> $$svc_dir"; \
		cd $$svc_dir && chmod +x gen.sh && ./gen.sh || exit $$?; \
	done

# 生成所有 RPC 服务
rpc:
	@echo "==> Generating RPC/proto code for all services…"
	@for f in $(APP_SERVICES); do \
		svc_dir=$$(dirname $$f); \
		echo "----> $$svc_dir"; \
		cd $$svc_dir && chmod +x gen.sh && ./gen.sh || exit $$?; \
	done

docker:
	@echo "==> Generating Dockerfile for all services…"
	@# 在 api 下找 main.go
	@for main in $$(find api -maxdepth 2 -type f -name main.go); do \
		svc_dir=$$(dirname $$main); \
		echo "----> $$svc_dir (Dockerfile)"; \
		cd $$svc_dir && goctl docker -go main.go || exit $$?; \
	done
	@# 在 app 下找 main.go
	@for main in $$(find app -maxdepth 2 -type f -name main.go); do \
		svc_dir=$$(dirname $$main); \
		echo "----> $$svc_dir (Dockerfile)"; \
		cd $$svc_dir && goctl docker -go main.go || exit $$?; \
	done

kube:
	@echo "==> Generating Kubernetes YAML for all services…"
	# 遍历 api 下的 main.go
	@for main in $$(find api -maxdepth 2 -type f -name main.go); do \
		svc_dir=$$(dirname $$main); \
		svc_name=$$(basename $$svc_dir); \
		echo "----> Generating kube YAML for $$svc_name"; \
		(cd $$svc_dir && \
			goctl kube deploy \
				-name $$svc_name \
				-namespace kline \
				-image $$svc_name:v1 \
				-o $$svc_name.yaml \
				-port 10001 \
		) || exit $$?; \
	done
	# 遍历 app 下的 main.go
	@for main in $$(find app -maxdepth 2 -type f -name main.go); do \
		svc_dir=$$(dirname $$main); \
		svc_name=$$(basename $$svc_dir); \
		echo "----> Generating kube YAML for $$svc_name"; \
		(cd $$svc_dir && \
			goctl kube deploy \
				-name $$svc_name \
				-namespace kline \
				-image $$svc_name:v1 \
				-o $$svc_name.yaml \
				-port 10001 \
		) || exit $$?; \
	done

# 别名
gen: all

build:
	@echo "==> Building all services (GOOS=darwin GOARCH=amd64 CGO_ENABLED=0)…"

	# 遍历 api 下的 main.go
	@for main in $$(find api -maxdepth 2 -type f -name main.go); do \
		svc_dir=$$(dirname $$main); \
		svc_name=$$(basename $$svc_dir); \
		echo "----> Tidying and building $$svc_name in $$svc_dir/bin"; \
		( \
			cd $$svc_dir && \
			go mod tidy && \
			mkdir -p bin && \
			GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$$svc_name main.go \
		) || exit $$?; \
	done

	# 遍历 app 下的 main.go
	@for main in $$(find app -maxdepth 2 -type f -name main.go); do \
		svc_dir=$$(dirname $$main); \
		svc_name=$$(basename $$svc_dir); \
		echo "----> Tidying and building $$svc_name in $$svc_dir/bin"; \
		( \
			cd $$svc_dir && \
			go mod tidy && \
			mkdir -p bin && \
			GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$$svc_name main.go \
		) || exit $$?; \
	done

# 清理（可按需定制，示例删除所有 pb 目录）
clean:
	@echo "==> Cleaning generated pb directories…"
	@find app -type d -name '*pb' -exec rm -rf {} +
	@echo "Done."

help:
	@echo "Usage:"
	@echo "  make            # run all api and rpc generators"
	@echo "  make api        # run all api/*/gen.sh"
	@echo "  make rpc        # run all app/*/gen.sh"
	@echo "  make clean      # clean pb directories under app/"
