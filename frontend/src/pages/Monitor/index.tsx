import { getStrainMonitoring } from '@/api/request';
import Header from '@/components/Header';
import useDataCache from '@/hooks/useDataCache';
import { DatePicker } from 'antd';
import type { Dayjs } from 'dayjs';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import MonitorForm from './components/MonitorForm';
import MonitoringTrendChart from './components/MonitoringTrendChart';
import UnitModelCard from './components/UnitModelCard';
import { Wrapper } from './styles';

// 应变监测
export interface StrainPoint {
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
interface StrainMonitoringResponse {
  data: {
    data: {
      points: StrainPoint[];
    };
  };
}

// 日期选择器
const { RangePicker } = DatePicker;

const MonitorPage: React.FC = () => {
  const [selectedTitle, setSelectedTitle] = useState('monitorDataDisplay');
  const [selectedDataOption, setSelectedDataOption] = useState('realtime');
  const [monitorData, setMonitorData] = useState<StrainPoint[]>([]);
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs] | null>(null);
  const [loading, setLoading] = useState(false);
  const [lastUpdateTime, setLastUpdateTime] = useState<string>('');

  // 使用数据缓存 Hook
  const { getCachedData } = useDataCache<StrainMonitoringResponse>();

  // 使用 useMemo 缓存单元数据
  const units = useMemo(
    () =>
      Array.from({ length: 4 }, (_, index) => ({
        id: index + 1,
        title: `${index + 1}号机组`,
      })),
    [],
  );

  // 使用 useCallback 优化数据获取函数
  const fetchData = useCallback(async () => {
    if (loading) return; // 防止重复请求

    try {
      setLoading(true);
      const now = Math.floor(Date.now() / 1000);

      let startTime = now;
      let endTime = now;

      // 如果日期选择器有值，优先使用日期选择器的时间范围
      if (dateRange) {
        startTime = Math.floor(dateRange[0].valueOf() / 1000);
        endTime = Math.floor(dateRange[1].valueOf() / 1000);
      } else {
        // 否则根据选择的数据选项设置时间范围
        if (selectedDataOption === 'realtime') {
          startTime = now - 24 * 60 * 60; // 24小时前
        } else if (selectedDataOption === '3days') {
          startTime = now - 3 * 24 * 60 * 60;
        } else if (selectedDataOption === '7days') {
          startTime = now - 7 * 24 * 60 * 60;
        }
      }

      // 使用缓存获取数据
      const response = await getCachedData(
        [startTime, endTime],
        () =>
          getStrainMonitoring(
            startTime,
            endTime,
          ) as Promise<StrainMonitoringResponse>,
      );

      setMonitorData(response.data.data.points);
      setLastUpdateTime(new Date().toLocaleString());
      console.log('监测数据:', response.data.data.points.length, '条记录');
    } catch (error) {
      console.error('获取应变监测数据失败:', error);
      setMonitorData([]);
    } finally {
      setLoading(false);
    }
  }, [selectedDataOption, dateRange, getCachedData]); // ✅ 移除 loading 依赖

  // 防抖处理数据获取
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchData();
    }, 300); // 300ms 防抖

    return () => clearTimeout(timer);
  }, [fetchData]);

  // 实时数据自动更新 - 每分钟更新一次
  useEffect(() => {
    let intervalId: NodeJS.Timeout | null = null;

    if (selectedDataOption === 'realtime' && !dateRange) {
      // 立即执行一次数据获取
      fetchData();

      // 设置每分钟更新一次
      intervalId = setInterval(() => {
        console.log('定时更新实时数据...');
        fetchData();
      }, 60000); // 60秒 = 1分钟
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [selectedDataOption, dateRange, fetchData]);

  // 使用 useMemo 缓存最新的应变数据
  const latestStrainData = useMemo(() => {
    return monitorData?.[monitorData.length - 1];
  }, [monitorData]);

  // 使用 useCallback 优化事件处理函数
  const handleTitleChange = useCallback((title: string) => {
    setSelectedTitle(title);
  }, []);

  const handleDataOptionChange = useCallback((option: string) => {
    setSelectedDataOption(option);
  }, []);

  const handleDateRangeChange = useCallback((dates: any) => {
    setDateRange(dates);
  }, []);

  return (
    <Wrapper>
      <div className="container">
        <Header
          selectedTitle={selectedTitle}
          onUnitChange={handleTitleChange}
        />

        {/* 四个机组展示 */}
        <div className="monitor-content">
          {units.map((unit) => (
            <UnitModelCard
              key={unit.id}
              unit={unit}
              strainData={latestStrainData}
            />
          ))}
        </div>

        {/* 日期选择器与数据切换 */}
        <div className="monitor-title">
          {/* 日期选择 */}
          <div className="date-select">
            <div className="date-select-data">
              日期选择
              <RangePicker
                showTime
                style={{ width: 380, height: 40 }}
                placeholder={['开始时间', '结束时间']}
                allowClear
                onChange={handleDateRangeChange}
              />
            </div>
          </div>

          {/* 数据切换 */}
          <div className="data-switch">
            <div className="data-switch-content">
              <div
                className="data-switch-options"
                onClick={() => handleDataOptionChange('realtime')}
              >
                <img
                  src={
                    selectedDataOption === 'realtime'
                      ? '/images/编组12.png'
                      : '/images/椭圆形.png'
                  }
                  alt=""
                />
                <span>实时数据</span>
              </div>
              <div
                className="data-switch-options"
                onClick={() => handleDataOptionChange('3days')}
              >
                <img
                  src={
                    selectedDataOption === '3days'
                      ? '/images/编组12.png'
                      : '/images/椭圆形.png'
                  }
                  alt=""
                />
                <span>近3日</span>
              </div>
              <div
                className="data-switch-options"
                onClick={() => handleDataOptionChange('7days')}
              >
                <img
                  src={
                    selectedDataOption === '7days'
                      ? '/images/编组12.png'
                      : '/images/椭圆形.png'
                  }
                  alt=""
                />
                <span>近7日</span>
              </div>
            </div>
          </div>
        </div>

        {/* 列表与趋势图 */}
        <div className="monitor-body">
          {/* 监测列表 */}
          <div className="monitor-list">
            {/* <div className="monitor-list-title">
              监测列表 {loading && '(加载中...)'}
              {selectedDataOption === 'realtime' && !dateRange && (
                <span
                  style={{
                    fontSize: '12px',
                    color: '#666',
                    marginLeft: '10px',
                  }}
                >
                  自动更新中 | 最后更新: {lastUpdateTime}
                </span>
              )}
            </div> */}
            <div className="monitor-list-title">监测列表</div>
            <MonitorForm monitorData={monitorData} />
          </div>
          {/* 应急监测趋势图 */}
          <div className="monitor-charts">
            {/* <div className="monitor-list-title">
              应急监测趋势图
               {loading && '(加载中...)'}
              {selectedDataOption === 'realtime' && !dateRange && (
                <span
                  style={{
                    fontSize: '12px',
                    color: '#666',
                    marginLeft: '10px',
                  }}
                >
                  自动更新中 | 最后更新: {lastUpdateTime}
                </span>
              )}
            </div> */}
            <div className="monitor-list-title">应急监测趋势图</div>
            <MonitoringTrendChart monitorData={monitorData} />
          </div>
        </div>
      </div>
    </Wrapper>
  );
};

export default MonitorPage;
