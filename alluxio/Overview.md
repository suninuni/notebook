# Overview

## [Introduction](https://docs.alluxio.io/os/user/stable/en/Overview.html)

- 首个开源的用于云分析和 AI 的数据编排技术。使存储层中的数据更容易更快速的被驱动层访问。
- 优点
    - 内存速度 I/O，支持分层存储，更低成本的利用内存，SSD 及磁盘。
    - 简化的云对象存储。由于云对象存储以 object 为基础，在对常见的文件系统操作（目录查询，文件重命名等）上会有很大的性能开销，使用 Alluxio 可以缓解这些问题。
    - 简化的数据管理。对多个数据源的单点访问，以及同时支持同一存储系统的不同版本而不需要驱动进行复杂的配置。
    - 易于部署使用。將来自应用程序的数据访问请求转换为基础存储接口，与 Hadoop 兼容，现有的数据分析应用如 Spark、MapReduce 等不需要代码更改就可以在 Alluxio 上运行。
- 特点
    - 全局命名空间。即多个独立存储系统的单点访问。
    - 智能多层缓存。即配置策略可自动优化数据放置，以提高内存和磁盘的性能和可靠性。
    - 服务端 API 转换。即支持行业通用的 API，如 HDFS API，S3 API，FUSE API，REST API。

## [Architecture](https://docs.alluxio.io/os/user/stable/en/overview/Architecture.html)

![Architecture](https://docs.alluxio.io/os/user/stable/img/architecture-overview-simple-docs.png)

- 一个 leading master，备用 masters，一个 job aster，备用 job masters，wokers，job workers。master 和 worker 构成了 Alluxio servers。
- Client 通过 Spark，MapReduce 的作业，Alluxio Cli 或者 Alluxio FUSE 层与 Alluxio Servers 进行交互。
- Allxuio Job Masters 和 Job Wokers 可以认为是独立的一套 Job Service，是一个轻量级的任务计划框架，Allxuio Job Masters 负责为 Job Workers 分配不同的操作，主要是 Alluxio 与 UFS 之间的数据交换。
- Master 和 Job Master 可以不放置在同个 Pod 中，建议放在一起以为 RPC 和数据传输提供低延迟。

### Masters

![Masters](https://docs.alluxio.io/os/user/stable/img/architecture-master-docs.png)

#### Leading Master

- 管理系统的全局元数据，包括文件系统元数据（e.g. 文件系统 inode 树），块元数据（e.g. 块位置）和 worker 容量的元数据（可用空间和已用空间）。
- 只会在自己的存储中查询元数据。將所有文件系统事务的记录存储到分布式持久存储中，以恢复 master 的状态，这些被称为 journal（类似 oplog）。

#### Standby Masters

- 读取 Leading Master 记录的 journal，使自己的主机状态副本保持最新。
- 创建日志检查点，以便快速恢复。
- 不会处理任何请求，只在 Leading Master 故障后参与选举。

#### Secondary Masters (for UFS journal)

当使用 UFS 日志运行单个不具有 HA 模式的 Alluxio Master 时，可以在主服务器上启动一个 Secondary Master，只会创建日志检查点以实现快速恢复，永远无法升级为 Leading Master

#### Job Masters

- 对 Leading Master 的需要执行的操作进行调度并分配到 Alluxio Job Workers 中执行，使 Alluxio Master 可以使用更少的资源并更快的响应。
- 提供了可扩展的框架。

### Wokers

![Workers](https://docs.alluxio.io/os/user/stable/img/architecture-worker-docs.png)

#### Alluxio Workers

- 负责管理分配用户配置的本地存储（memory，SSD，HDD 等）。
- 將数据存储为 blocks，通过读取或创建新的 blocks 来处理 client 的读写请求。
- 只管理 blocks，blocks 中存数据的映射关系只存储在 master 中。
- 可以理解为对本地存储的抽象，使 client 可以更快更方便的读取数据，且不需要处理与底层存储的连接。
- 因为 RAM 空间有限，当空间满时，Workers 通过驱逐策略决定將哪些数据保留在 Alluxio 中。

#### Alluxio Job Workers

- 是 Alluxio 文件系统的客户端，包含 Alluxio Client。
- 接收 Alluxio Job Master 的命令包括运行加载，持久化，复制，移动和在给定文件系统上复制操作。

### Client

- 为用户提供一个与 Alluxio 交互的网关。
- 与 Leading Master 交互以执行元数据操作，与 workers 交互以读取喝写入存储在 Alluxio 中的数据。

### 数据流：读

#### Local Cache Hit

![dataflow-local-cache-hit](https://docs.alluxio.io/os/user/stable/img/dataflow-local-cache-hit.gif)

- Application -请求数据-> Alluxio Client -查询数据位置-> Alluxio Master -> Alluxio Client read file from RAM（short-circuit reads 短路读，绕过了 Alluxio Worker，需要文件系统的权限）。
- 短路不可行的情况下，Alluxio 提供基于 domain socket 的短路读，即通过预创建的 socket 从 Alluxio Worker 传输数据到 client。

#### Remote Cache Hit

![dataflow-remote-cache-hit](https://docs.alluxio.io/os/user/stable/img/dataflow-remote-cache-hit.gif)

- Application -> Alluxio Client -> Alluxio Master -> Alluxio Client -> 其他机器上的 Alluxio Worker 从 RAM 中读取文件 -> Alluxio Client 返回数据给 Application 并通过当前机器上 Alluxio Worker 创建数据副本（block）进行缓存。

#### Cache Miss

![dataflow-cache-miss](https://docs.alluxio.io/os/user/stable/img/dataflow-cache-miss.gif)

- Application -> Alluxio Client -> Alluxio Master -> Alluxio Client -> Alluxio Worker 从 Under Storage（底层存储）读取数据返回并缓存（异步）-> Alluxio Client。

### 数据流：写

#### Write to Alluxio only (MUST_CACHE)

![dataflow-must-cache](https://docs.alluxio.io/os/user/stable/img/dataflow-must-cache.gif)

- 不会持久存储到 UFS 中，所以数据可能会丢失，只适合临时数据。

#### Write through to UFS (CACHE_THROUGH)

![dataflow-cache-through](https://docs.alluxio.io/os/user/stable/img/dataflow-cache-through.gif)

- 同时写入到 UFS 和 RAM 中。
- 如果写 UFS 是强一致的（如 HDFS），那么如果数据不在 UFS 层有任何改动，两边的数据就是一致的。
- 如果写 UFS 的是事件一致的（如 S3，即只有操作成功，文件状态才会改变），那么存在不一致的窗口

#### Write back to UFS (ASYNC_THROUGH)

![dataflow-async-through](https://docs.alluxio.io/os/user/stable/img/dataflow-async-through.gif)

- 先写到 RAM 中，通过后台任务持久化到 UFS 中（可能丢失）。

#### Write to UFS Only (THROUGH)

- 只写到 UFS 中。但 Alluxio 知道文件及其状态。因此数据仍然是一致的。

## [Job Service](https://docs.alluxio.io/os/user/stable/en/overview/JobService.html)
