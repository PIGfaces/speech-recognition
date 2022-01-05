# 调用语音API工具

## 基本说明

### 目录说明

```
├── cmd                 // 编译程序的入口
│   ├── commonspeech    // 仅含有同步语音
│   ├── recordspeech    // 仅含有异步点播语音
│   └── speech          // 含有同步语音和异步点播语音
├── go.mod
├── go.sum
├── Makefile            // 构建规则
├── pkg                 // 此项目中的公共业务代码
│   ├── handler         // 任务处理入口
│   ├── reader          // 读取任务
│   ├── recognition     // 调用识别 SDK
│   ├── schedule        // 并发管理、读写管理
│   └── writer          // 保存结果
├── README.md
└── rsa_private_key.pem // 测试账号的私钥
```

### 编译   

- 仅编译同步语音: `make build`   
- 仅编译异步点播语音: `make build_record`   
- 仅编译同步语音: `make build_common`   

### 其他   

`make help` 查看帮助   

## 工具使用   

`./speechRecord -h` 查看使用帮助

**Example**:    
```shell
$ ./speechRecord -PrivKeyPath=./rsa_private_key.pem -SecretIds=61cafe78a997180094b52b52,61cc35427d169e001d42479e,61cc35a602399d001da97b7c -TaskFilePath=./TTlastweek.txt
```

### 参数说明   

- `-PrivKeyPath` : 私钥路径   
- `-SecretIds` : 本次任务需要执行的 `secretId` 接口   
- `-TaskFilePath`: 本次任务需要审核的语音 `url` , 一行一条 `url`, 按行区分    
- `-concurrentNum`: 本次任务每秒执行的请求数量