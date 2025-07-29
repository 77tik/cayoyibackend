import { Pagination, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import React, { useCallback, useMemo, useState } from 'react';
import styled from 'styled-components';

const FormWrapper = styled.div`
  height: 100%;
  overflow: auto;

  /* ✅ 隐藏滚动条但可滚动 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */

  &::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
  }

  .ant-table {
    background: #f7faff;
    width: 100%;
    table-layout: fixed;
    border-right: 1px solid #e2e7f5;
  }

  /* ✅ 第 1 行表头：sticky 在顶部 */
  .ant-table-thead > tr:nth-child(1) > th {
    position: sticky;
    top: 0;
    z-index: 3;
    background: rgb(235, 239, 249);
  }

  /* ✅ 第 2 行表头：sticky 在第 1 行下方 */
  .ant-table-thead > tr:nth-child(2) > th {
    position: sticky;
    top: 32px;
    z-index: 2;
    background: rgb(235, 239, 249);
  }

  /* 通用样式（可以保留你的原样式） */
  .ant-table-thead > tr > th {
    font-weight: 600;
    text-align: center;
    border: 1px solid #e2e7f5;
    height: 30px;
    line-height: 30px;
    padding: 0 !important;
  }

  .ant-table-tbody > tr > td {
    text-align: center;
    border: 1px solid #e2e7f5;
    height: 32px;
    line-height: 32px;
    padding: 0 !important;
  }

  .ant-table-cell {
    font-size: 14px;
  }
`;

interface StrainPoint {
  timestamp: number;
  one_upper: number;
  one_door: number;
  two_cover: number;
  two_door: number;
  three_cover: number;
  three_door: number;
  four_cover: number;
  four_door: number;
}

interface MonitorFormProps {
  monitorData: StrainPoint[];
}

interface MonitorRow {
  key: number;
  time: string;
  unit1: number[];
  unit2: number[];
  unit3: number[];
  unit4: number[];
}

const MonitorForm: React.FC<MonitorFormProps> = ({ monitorData }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(50); // 每页显示50条数据

  // 使用 useMemo 缓存转换后的表格数据（倒序显示，最新的在前面）
  const tableData: MonitorRow[] = useMemo(() => {
    return monitorData
      .map((point, index) => ({
        key: index, // 添加 key
        time: new Date(point.timestamp * 1000).toLocaleString(),
        unit1: [Math.round(point.one_upper), Math.round(point.one_door)],
        unit2: [Math.round(point.two_cover), Math.round(point.two_door)],
        unit3: [Math.round(point.three_cover), Math.round(point.three_door)],
        unit4: [Math.round(point.four_cover), Math.round(point.four_door)],
      }))
      .reverse(); // 倒序，最新数据在前面
  }, [monitorData]);

  // 分页数据
  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return tableData.slice(startIndex, endIndex);
  }, [tableData, currentPage, pageSize]);

  // 使用 useMemo 缓存列配置
  const columns: ColumnsType<MonitorRow> = useMemo(
    () => [
      {
        title: '监测时间',
        dataIndex: 'time',
        key: 'time',
        // rowSpan: 2,
        width: 160,
        fixed: 'left',
      },
      {
        title: '1#机组',
        children: [
          {
            title: '顶盖',
            dataIndex: ['unit1', 0],
            key: 'unit1_0',
            width: 100,
          },
          {
            title: '涡壳入孔门',
            dataIndex: ['unit1', 1],
            key: 'unit1_1',
            width: 120,
          },
        ],
      },
      {
        title: '2#机组',
        children: [
          {
            title: '顶盖',
            dataIndex: ['unit2', 0],
            key: 'unit2_0',
            width: 100,
          },
          {
            title: '涡壳入孔门',
            dataIndex: ['unit2', 1],
            key: 'unit2_1',
            width: 120,
          },
        ],
      },
      {
        title: '3#机组',
        children: [
          {
            title: '顶盖',
            dataIndex: ['unit3', 0],
            key: 'unit3_0',
            width: 100,
          },
          {
            title: '涡壳入孔门',
            dataIndex: ['unit3', 1],
            key: 'unit3_1',
            width: 120,
          },
        ],
      },
      {
        title: '4#机组',
        children: [
          {
            title: '顶盖',
            dataIndex: ['unit4', 0],
            key: 'unit4_0',
            width: 100,
          },
          {
            title: '涡壳入孔门',
            dataIndex: ['unit4', 1],
            key: 'unit4_1',
            width: 120,
          },
        ],
      },
    ],
    [],
  );

  // 分页改变处理函数
  const handlePageChange = useCallback((page: number, size: number) => {
    setCurrentPage(page);
    setPageSize(size);
  }, []);

  return (
    <FormWrapper>
      <div className="monitor-form">
        <Table
          columns={columns}
          dataSource={paginatedData}
          pagination={false}
          bordered
          size="middle"
          className="custom-table"
          // scroll={{ y: 400, x: 800 }}
          rowKey="key"
        />
        <div style={{ marginTop: 16, textAlign: 'center' }}>
          <Pagination
            current={currentPage}
            pageSize={pageSize}
            total={tableData.length}
            showSizeChanger
            showQuickJumper
            showTotal={(total, range) =>
              `第 ${range[0]}-${range[1]} 条/共 ${total} 条`
            }
            pageSizeOptions={['20', '50', '100', '200']}
            onChange={handlePageChange}
            onShowSizeChange={handlePageChange}
          />
        </div>
      </div>
    </FormWrapper>
  );
};

export default React.memo(MonitorForm);
