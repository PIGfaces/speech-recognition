.PHONY: all build clean run check cover lint docker help

SUMMARY_SPEECH_TOOL=summarySpeech
COMMON_SPEECH_TOOL=commonSpeech
RECORD_SPEECH_TOOL=recordSpeech

all: check build

build: build_summary

clean:
	rm -f ${SUMMARY_SPEECH_TOOL} ${COMMON_SPEECH_TOOL} ${RECORD_SPEECH_TOOL}
	@go clean

build_summary: 
	@go build -o "${SUMMARY_SPEECH_TOOL}" cmd/speech/main.go

build_common: 
	@go build -o "${SUMMARY_SPEECH_TOOL}" cmd/speech/main.go

build_record: 
	@go build -o "${SUMMARY_SPEECH_TOOL}" cmd/speech/main.go

check:
	@go fmt ./
	@go vet ./

lint:
	golangci-lint run --enable-all

help:
	@echo "make: 格式化代码 并编译生成二进制文件"
	@echo "make build 编译生成二进制文件, 编译目录为 cmd/speech"
	@echo "make build_common 仅编译生成同步语音工具二进制文件, 编译目录为 cmd/commonspeech"
	@echo "make build_record 仅编译生成异步点播语音工具二进制文件, 编译目录为 cmd/recordspeech"
	@echo "make check 格式化 go 代码"
	@echo "make lint 执行代码检查"