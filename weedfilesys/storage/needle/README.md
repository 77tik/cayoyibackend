| 字段名             | 说明                                     |
| --------------- | -------------------------------------- |
| Cookie          | 4 字节，随机值，用于安全目的                        |
| Id              | 8 字节，Needle 的唯一标识符                     |
| Size            | 4 字节，表示 Needle 数据段的总大小（不含头部、校验、时间戳、填充） |
| DataSize        | 4 字节，Data 字段的长度                        |
| Data            | 可变长度，实际的文件数据                           |
| Flags           | 1 字节，位掩码，指示是否存在压缩、名称、MIME 类型等字段        |
| NameSize/Name   | 可选字段，1 字节长度 + 可变长度名称                   |
| MimeSize/Mime   | 可选字段，1 字节长度 + 可变长度 MIME 类型             |
| LastModified    | 5 字节，文件的最后修改时间                         |
| Ttl             | 2 字节，生命周期（Time-To-Live）                |
| PairsSize/Pairs | 可选的键值对（元数据）                            |
| Checksum        | 4 字节，Data 数据部分的 CRC32 校验值              |
| Timestamp       | 8 字节，仅 v3 存在，用于记录写入时间（纳秒）              |
| Padding         | 0–7 字节，用于补齐 Needle 长度至 8 的整数倍          |

| 位值   | 名称                      | 含义                            |
| ---- | ----------------------- | ----------------------------- |
| 0x01 | FlagIsCompressed        | 数据已压缩（通常用 gzip）               |
| 0x02 | FlagHasName             | 包含 Name 字段                    |
| 0x04 | FlagHasMime             | 包含 Mime 字段                    |
| 0x08 | FlagHasLastModifiedDate | 包含 LastModified 字段            |
| 0x10 | FlagHasTtl              | 包含 Ttl 字段                     |
| 0x20 | FlagHasPairs            | 包含 Pairs 字段                   |
| 0x80 | FlagIsChunkManifest     | 数据为 chunk manifest（用于大文件分块索引） |
