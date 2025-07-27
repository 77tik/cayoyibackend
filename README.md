# cayoyibackend

+ 我的后端学习之路
+ Weedfilesys 开发进度：
  + ChunkCacheVolume
  + OnDiskCacheLayer
+ 1. 存储层（storage/）
     先实现最基础的数据存储与读写功能，比如单机的卷（volume）文件读写、数据块管理。
     这样可以保证你有了最基本的文件存储能力。
2. 元数据与管理（master/server/ 和 filer/）
   实现Master Server，负责卷的分配、管理和节点注册。
   实现Filer，用于管理文件和目录的元数据（如文件名、目录结构、权限等）。
3. 集群与拓扑（cluster/ 和 topology/）
   加入集群管理，实现节点注册、心跳、卷分布、拓扑结构（如数据中心、机架、节点等）。
   这样可以支持多节点和分布式部署。
4. 副本与一致性（replication/ 和 raft/）
   实现数据副本机制，保证数据可靠性。
   可以引入 Raft 等一致性协议，确保主节点选举和元数据一致性。
5. 对外接口（server/、s3api/、shell/）
   实现HTTP API，让用户可以上传、下载、删除文件。
   可以扩展 S3 兼容接口，方便与云原生生态集成。
   实现命令行工具，方便调试和管理。
6. 挂载与集成（mount/）
   实现 FUSE 挂载，让用户可以像本地文件系统一样访问分布式存储。
7. 安全与监控（security/、admin/、telemetry/）
   增加认证、授权、TLS 加密等安全特性。
   实现监控、管理界面等运维功能。
   建议实现顺序：

storage → master/server → filer → cluster/topology → replication/raft → server/s3api/shell → mount → security/admin