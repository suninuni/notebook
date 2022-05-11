# Cli

## alluxio fs

- alluxio fs leader 查询 Leading Master

## alluxio copyDir conf

每一个Alluxio server进程都会去读取本地安装目录下的`conf/alluxio-site.properties`以及`conf/alluxio-env.sh`文件, 所以对某一项设置做修改后, 需要在所有安装 master 和worker的节点上都更新这一文件. 一个简单的方法是在 master 节点上做过修改以后, 使用我们提供的"bin/alluxio copyDir"这个脚本, 自动的根据conf/workers的内容,  把maste r节点上的配置文件 rsync 到所有的 worker 节点上.

Alluxio 属性由以下列表最早的配置决定：

- JVM 系统参数（-Dproperty=key alluxio fs 中也可使用）
- 环境变量（${ALLUXIO_HOME}/conf/alluxio-env.sh）
- 参数配置文件（ 当Alluxio集群启动时, 每一个Alluxio服务端进程（包括master和worke） 在目录 `${HOME}/.alluxio/`, `/etc/alluxio/` and `${ALLUXIO_HOME}/conf` 下顺序读取 alluxio-site.properties , 当 alluxio-site.properties 文件被找到，将跳过剩余路径的查找.）
- 集群默认值

alluxio getConf --source --master

从 v1.8 开始，每个 Alluxio 客户端都可以使用从 master 节点获取的集群范围的配置值初始化其配置。具体来说，当不同的客户端应用程序（如 Alluxio Shel l命令、Spark 作业或MapReduce 作业）连接到一个 Alluxio 服务时，它们将使用 master 节点提供的默认值初始化自己的 Alluxio 配置属性，这些默认值是基于 master 节点的 `${ALLUXIO_HOME}/conf/alluxio-site.properties` 属性文件。因此，集群管理员可以放置客户端设置（例如，`alluxio.user.*`）或网络传输设置(如 `alluxio.security.authentication.type`)在 master 节点的 `${ALLUXIO_HOME}/conf/alluxio-site.properties`。 它将被分布并成为集群范围内的默认值，用于新的 Aluxio 客户端。

Spark用户可以通过对 Spark executor 的 `spark.executor.extraJavaOptions` 和 Spark drivers 的 `spark.driver.extraJavaOptions` 添加 "-Dproperty=value" 向Spark job 传递 JVM 环境参数。
