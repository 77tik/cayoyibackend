import type { EChartsOption } from 'echarts';
import * as echarts from 'echarts';
import React, { useCallback, useEffect, useMemo, useRef } from 'react';
import styled from 'styled-components';
import { StrainPoint } from '../../index';

const ChartContainer = styled.div`
  height: 300px;
  .chart-container {
    width: 100%;
    height: 100%;
  }
`;

interface MonitoringTrendChartProps {
  monitorData?: StrainPoint[];
}

const MonitoringTrendChart: React.FC<MonitoringTrendChartProps> = ({
  monitorData = [],
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstanceRef = useRef<echarts.ECharts | null>(null);

  // 数据抽样函数 - 当数据点过多时进行抽样
  const sampleData = useCallback(
    (data: StrainPoint[], maxPoints: number = 200) => {
      if (data.length <= maxPoints) return data;

      const step = Math.ceil(data.length / maxPoints);
      return data.filter((_, index) => index % step === 0);
    },
    [],
  );

  // 使用 useMemo 缓存处理后的数据
  const sampledData = useMemo(() => {
    return sampleData(monitorData);
  }, [monitorData, sampleData]);

  // 格式化日期
  const formatDate = useCallback((timestamp: number) => {
    const date = new Date(timestamp);
    // return `${date.getFullYear()}\n${String(date.getMonth() + 1).padStart(
    //   2,
    //   '0',
    // )}/${String(date.getDate()).padStart(2, '0')}`;
    return `${String(date.getMonth() + 1).padStart(2, '0')}/${String(
      date.getDate(),
    ).padStart(2, '0')}\n${String(date.getHours()).padStart(2, '0')}:${String(
      date.getMinutes(),
    ).padStart(2, '0')}`;
  }, []);

  // 使用 useMemo 缓存图表配置
  const chartOption = useMemo((): EChartsOption => {
    if (!sampledData?.length) return {};

    // 准备图表数据
    const xAxisData = sampledData.map(
      (point) => formatDate(point.timestamp * 1000),
      // new Date(point.timestamp * 1000).toLocaleString(),
    );
    const seriesData = [
      {
        name: '1#机组顶盖',
        data: sampledData.map((point) => point.one_upper),
      },
      {
        name: '2#机组顶盖',
        data: sampledData.map((point) => point.two_cover),
      },
      {
        name: '3#机组顶盖',
        data: sampledData.map((point) => point.three_cover),
      },
      {
        name: '4#机组顶盖',
        data: sampledData.map((point) => point.four_cover),
      },
    ];

    // 计算所有数据的最大值和最小值
    const allValues = seriesData.flatMap((series) => series.data);
    const maxValue = Math.max(...allValues);
    const minValue = Math.min(...allValues);

    // 设置y轴的范围，留出一些边距
    const yAxisMin = Math.floor(minValue * 0.9);
    const yAxisMax = Math.ceil(maxValue * 1.1);

    return {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          label: {
            backgroundColor: '#6a7985',
          },
        },
        formatter: function (params: any) {
          let result = `时间: ${params[0].axisValue}<br/>`;
          params.forEach((item: any) => {
            const value =
              item.value !== null ? Number(item.value).toFixed(2) : '--';
            result += `${item.seriesName}: ${value}<br/>`;
          });
          return result;
        },
      },
      legend: {
        data: ['1#机组顶盖', '2#机组顶盖', '3#机组顶盖', '4#机组顶盖'],
        itemWidth: 16,
        itemHeight: 0,
        top: 12,
        textStyle: {
          fontSize: 14,
        },
      },
      grid: {
        left: '1%',
        right: '1%',
        bottom: '1%',
        top: '20%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        data: xAxisData,
        boundaryGap: true,
        axisLabel: {
          fontSize: 12,
          color: '#666',
          interval: (index: number) => {
            const dataLength = xAxisData.length;
            if (dataLength > 30) {
              return index % Math.floor(dataLength / 10) === 0;
            } else if (dataLength > 15) {
              return index % Math.floor(dataLength / 8) === 0;
            } else if (dataLength > 8) {
              return index % Math.floor(dataLength / 6) === 0;
            }
            return true;
          },
          margin: 15,
          formatter: (value: string) => {
            return value.length > 12 ? value.slice(0, 10) + '...' : value;
          },
        },
        axisTick: {
          alignWithLabel: true,
        },
        axisLine: {
          lineStyle: {
            color: '#999',
          },
        },
      },
      yAxis: {
        type: 'value',
        min: yAxisMin,
        max: yAxisMax,
        splitNumber: 5,
      },
      series: seriesData.map((item) => ({
        name: item.name,
        type: 'line',
        smooth: true,
        data: item.data,
        symbol: 'none',
        sampling: 'lttb',
        emphasis: {
          focus: 'series',
        },
      })),
      animation: false, // 关闭动画以提升性能
    };
  }, [sampledData, formatDate]);

  useEffect(() => {
    if (!chartRef.current) return;

    // 如果图表实例不存在，创建新实例
    if (!chartInstanceRef.current) {
      chartInstanceRef.current = echarts.init(chartRef.current);
    }

    const chart = chartInstanceRef.current;

    // 设置图表配置
    if (Object.keys(chartOption).length > 0) {
      chart.setOption(chartOption, true); // 使用 notMerge: true 提升性能
    }

    // 监听窗口大小变化，调整图表大小
    const handleResize = () => {
      chart.resize();
    };

    window.addEventListener('resize', handleResize);

    // 清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [chartOption]);

  // 组件卸载时清理图表实例
  useEffect(() => {
    return () => {
      if (chartInstanceRef.current) {
        chartInstanceRef.current.dispose();
        chartInstanceRef.current = null;
      }
    };
  }, []);

  return (
    <ChartContainer>
      <div className="chart-container" ref={chartRef} />
      {sampledData.length !== monitorData.length && (
        <div
          style={{
            fontSize: 12,
            color: '#666',
            textAlign: 'center',
            marginTop: 8,
          }}
        >
          {/* 显示 {sampledData.length} / {monitorData.length} 个数据点 */}
        </div>
      )}
    </ChartContainer>
  );
};

export default React.memo(MonitoringTrendChart);
