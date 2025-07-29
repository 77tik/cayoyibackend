# 3D模型文件说明

## 目录结构
```
public/static/fluid/
├── case1/
│   ├── mesh.json          # 主要几何网格文件
│   ├── velocity_h.json    # 速度场水平剖面
│   ├── velocity_v.json    # 速度场垂直剖面  
│   ├── pressure_h.json    # 压力场水平剖面
│   ├── pressure_v.json    # 压力场垂直剖面
│   ├── vof.json          # VOF空化分布
│   ├── vortex.json       # 涡带分布
│   └── stream_line.json  # 流线数据
├── case2/
└── ...
```

## 文件说明

### mesh.json
主要的3D几何网格文件，包含：
- 顶点坐标 (position)
- 法向量 (normal) 
- 各种物理量数据 (velocity, total_pressure, phase_1_vof, raw_q_criterion)
- 物理量的最值信息 (userData)

### 文件大小
- 真实的CFD模型文件通常很大（几十MB到几百MB）
- 包含大量顶点和物理量数据
- 建议压缩存储或使用二进制格式

## 数据格式
所有JSON文件都应该遵循Three.js BufferGeometry格式：

```json
{
  "type": "BufferGeometry",
  "uuid": "unique-identifier",
  "metadata": {
    "type": "BufferGeometry", 
    "version": 4.5
  },
  "attributes": {
    "position": {
      "itemSize": 3,
      "type": "Float32Array",
      "array": [x1,y1,z1, x2,y2,z2, ...]
    },
    "物理量名称": {
      "itemSize": 1,
      "type": "Float32Array", 
      "array": [v1, v2, v3, ...]
    }
  },
  "userData": {
    "物理量名称": {
      "MaxValue": 最大值,
      "MinValue": 最小值
    }
  }
}
```

## 注意事项
1. 所有标量物理量的数组长度必须与顶点数量一致
2. MaxValue和MinValue用于颜色映射，必须准确
3. 文件路径要与API返回的路径对应
4. 大文件建议放在CDN或专门的文件服务器上 