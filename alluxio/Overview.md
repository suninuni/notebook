# Overview

https://docs.alluxio.io/os/user/stable/en/Overview.html

## Introduction

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

## Architecture

![Architecture](https://docs.alluxio.io/os/user/stable/img/architecture-overview-simple-docs.png)

一个 leading master，备用 masters，一个 job aster，备用 job masters，wokers，job workers。master 和 worker 构成了 Alluxio servers。Client 通过 Spark，MapReduce 的作业，Alluxio Cli 或者 Alluxio FUSE 层与 Alluxio Servers 进行交互。Allxuio Job Masters 和 Job Wokers 可以认为是独立的一套 Job Service，是一个轻量级的任务计划框架，负责为 Job Workers 分配不同的操作，主要是 Alluxio 与 UFS 之间的数据交换。可以不放置在同个 Pod 中，建议放在一起以为 RPC 和数据传输提供低延迟。架构概述

### Masters

![Alluxio Masters](https://docs.alluxio.io/os/user/stable/img/architecture-master-docs.png)

- Leading Master：管理系统的全局元数据，包括文件系统元数据（e.g. 文件系统 inode 树），块元数据（e.g. 块位置）和 worker 容量的元数据（可用空间和已用空间）。只会在自己的存储中查询元数据。將所有文件系统事务的记录存储到分布式持久存储中，以恢复 master 的状态，这些被称为 journal（类似 oplog）。
- Standby Masters：读取 Leading Master 记录的 journal，使自己的主机状态副本保持最新。创建日志检查点，以便快速恢复。但不会处理任何请求，只在 Leading Master 故障后参与选举。
- Secondary Masters (for UFS journal)：当使用 UFS 日志运行单个不具有 HA 模式的 Alluxio Master 时，可以在主服务器上启动一个 Secondary Master，只会创建日志检查点以实现快速恢复，永远无法升级为 Leading Master
- Job Masters：

### Wokers

### Client
