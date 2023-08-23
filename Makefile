
BuildTime=$(shell date "+%Y-%m-%d %H:%M:%S")
GitCommitID=$(shell git rev-parse HEAD)
GitTag=$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
RtPkg=git.qutoutiao.net/gopher/qms/internal/pkg/runtime

TargetName := mybase
OutputDir := outputs

all: default

dev: run

#设置代理
proxy:
	go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,https://mirrors.aliyun.com/goproxy/,direct   #[翻墙需要]设置代理

run: proxy
	go run main.go

clean:
	go clean -v ./...
	go mod tidy -v
	rm -rf ${OutputDir}

proto:
	@protoc -I app/api/proto app/api/proto/*.proto  --qms_out=plugins=grpc,paths=source_relative:app/api/proto/.

#默认的build命令
default: proxy
	# 准备好编译目录结构
	rm -rf $(OutputDir)
	mkdir -p $(OutputDir)
	cp -rf ./conf $(OutputDir)

	# 记录编译时间
	echo `date '+%Y-%m-%d %H:%M:%S'` > $(OutputDir)/build.log

	go build -o ${OutputDir}/${TargetName} -tags jsoniter -a \
	    -ldflags '-X "${RtPkg}.BuildTime=${BuildTime}" -X "${RtPkg}.GitCommitID=${GitCommitID}" -X "${RtPkg}.GitTag=${GitTag}"' \
	    main.go


local: proxy
	# 准备好编译目录结构
	rm -rf $(OutputDir)
	mkdir -p $(OutputDir)
	cp -rf ./conf $(OutputDir)

	# 记录编译时间
	echo `date '+%Y-%m-%d %H:%M:%S'` > $(OutputDir)/build.log

	go build -o ${OutputDir}/${TargetName} main.go