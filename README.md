# gox

gox 是一个 Go 语言工具库，提供了一系列常用的功能模块，涵盖并发、网络、加密、字符串处理、时间处理、文件操作等多个领域。

## 包列表

### 并发与同步

| 包 | 描述 |
|---|---|
| [conc](conc) | 并发工具集，包含原子操作、互斥锁、通道、协程池、障碍器、懒加载等 |
| [barrier](conc/barrier) | 同步障碍器 |
| [chanx](conc/chanx) | 通道操作工具，包含扇入扇出、管道、过滤等功能 |
| [goroutinex](conc/goroutinex) | 协程工具 |
| [mutexx](conc/mutexx) | 互斥锁扩展 |
| [oncex](conc/oncex) | 一次性执行组 |
| [poolx](conc/poolx) | 对象池和缓冲区池 |
| [waiter](conc/waiter) | 等待组扩展 |
| [asyncbatch](conc/asyncbatch) | 异步批量处理 |
| [gofer](conc/gofer) | 协程池抽象接口 |
| [lazyload](conc/lazyload) | 懒加载组 |
| [brave](conc/brave) | 带恢复的协程执行 |

### 网络与地址

| 包 | 描述 |
|---|---|
| [addrx](addrx) | 网络地址工具，IP 获取、端口选择、IP 转换等 |
| [httpx](httpx) | HTTP 客户端/服务端工具 |
| [netx](netx) | 网络工具（预留） |

### 加密与安全

| 包 | 描述 |
|---|---|
| [cryptox](cryptox) | 加密工具集 |
| [aesx](cryptox/aesx) | AES 加密 |
| [hmacx](cryptox/hmacx) | HMAC 签名 |
| [md5x](cryptox/md5x) | MD5 哈希 |
| [rsax](cryptox/rsax) | RSA 加密与签名 |
| [shax](cryptox/shax) | SHA 哈希家族 |
| [tlsx](cryptox/tlsx) | TLS 配置 |
| [x509](cryptox/x509) | X.509 证书工具 |
| [auth](auth) | 认证工具（Basic Auth） |

### 数据库

| 包 | 描述 |
|---|---|
| [databasex](databasex) | 数据库工具集 |
| [pagex](databasex/pagex) | 分页工具 |
| [sqls](databasex/sqls) | SQL 安全检查 |
| [unsafesql](databasex/unsafesql) | SQL 拼接构造器 |

### 字符串与文本

| 包 | 描述 |
|---|---|
| [stringx](stringx) | 字符串工具集 |
| [strconvx](strconvx) | 字符串转换工具 |
| [textx](textx) | 文本编码转换（中文、日文、韩文） |
| [fmtx](fmtx) | 格式化工具 |

### 时间与日期

| 包 | 描述 |
|---|---|
| [timex](timex) | 时间工具集，包含日期计算、时间比较等 |

### 反射与泛型

| 包 | 描述 |
|---|---|
| [reflectx](reflectx) | 反射工具，字段访问、类型操作等 |
| [constraintx](constraintx) | 泛型约束工具 |
| [protox](protox) | Protobuf 工具，消息克隆、切片转换等 |

### 文件与 IO

| 包 | 描述 |
|---|---|
| [filex](filex) | 文件操作工具，复制、解压、遍历等 |
| [iox](iox) | IO 工具，复制、关闭、长度计算等 |

### 图像处理

| 包 | 描述 |
|---|---|
| [imagex](imagex) | 图片处理工具，格式转换、缩放、旋转等 |

### 错误处理

| 包 | 描述 |
|---|---|
| [errorx](errorx) | 错误处理工具，错误链、多错误合并等 |

### 代码生成

| 包 | 描述 |
|---|---|
| [gen](gen) | 代码生成工具 |

### 数据结构

| 包 | 描述 |
|---|---|
| [heapx](heapx) | 堆结构 |
| [listx](listx) | 双向链表、无锁队列 |
| [mapx](mapx) | 映射工具 |
| [ringx](ringx) | 环形缓冲区 |
| [slicex](slicex) | 切片工具 |
| [sortx](sortx) | 排序工具 |
| [distributed](distributed) | 分布式工具，子集选择、序列生成等 |

### 日志

| 包 | 描述 |
|---|---|
| [slogx](slogx) | 日志扩展，结构化日志、上下文等 |

### 重试与退避

| 包 | 描述 |
|---|---|
| [backoff](backoff) | 退避算法，指数退避、斐波那契、线性等 |
| [retry](retry) | 重试策略 |

### 运行时

| 包 | 描述 |
|---|---|
| [runtimex](runtimex) | 运行时工具，堆栈获取等 |
| [osx](osx) | 操作系统工具，环境变量、信号处理等 |

### 运算符

| 包 | 描述 |
|---|---|
| [operator](operator) | 运算符工具，指针、三元运算符等 |

### 随机数

| 包 | 描述 |
|---|---|
| [randx](randx) | 随机数工具，支持多种随机数生成器 |

### 上下文

| 包 | 描述 |
|---|---|
| [contextx](contextx) | 上下文工具 |

### 工具集

| 包 | 描述 |
|---|---|
| [tools](tools) | 通用工具集 |

### 数学

| 包 | 描述 |
|---|---|
| [mathx](mathx) | 数学工具，数值比较、四舍五入等 |

### 二进制

| 包 | 描述 |
|---|---|
| [bytex](bytex) | 二进制数据处理工具，Diff、编辑等 |

## 安装

```bash
go get github.com/soyacen/gox
```

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/soyacen/gox/slicex"
    "github.com/soyacen/gox/strconvx"
)

func main() {
    // 切片工具
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
    max := slicex.Max(nums)
    fmt.Println("Max:", max)

    // 字符串转换
    i, _ := strconvx.ParseInt[int]("42", 10, 0)
    fmt.Println("Parsed:", i)
}
```

## 许可证

MIT License
