// 卡片和按钮配置、格式化函数
import { ButtonConfig, CardConfig } from './types';

// 格式化数值为两位小数的函数
export const formatValue = (value: any) => {
  if (typeof value === 'number') {
    return value.toFixed(2);
  }
  return value ?? '';
};

// 生成卡片配置
export const getCardConfig = (simulationResult: any): CardConfig => ({
  速度场: [
    {
      label: '最大速度',
      // value: formatValue(simulationResult?.data?.velocity?.max_velocity),
      value: simulationResult?.data?.velocity?.max
        ? `${formatValue(simulationResult?.data?.velocity?.max)} m/s`
        : '',
    },
    {
      label: '蜗壳出口平均速度',
      value: simulationResult?.data?.velocity?.volute_average_velocity
        ? `${formatValue(
            simulationResult?.data?.velocity?.volute_average_velocity,
          )}m/s`
        : ' ',
    },
  ],
  压力场: [
    {
      label: '最大压力',
      value: simulationResult?.data?.pressure?.max
        ? `${formatValue(simulationResult?.data?.pressure?.max)}  MPa`
        : '',
    },
    {
      label: '蜗壳出口压力',
      value: simulationResult?.data?.pressure?.volute_pressure
        ? `${formatValue(
            simulationResult?.data?.pressure?.volute_pressure,
          )} MPa`
        : '',
    },
  ],
  空化分布: [
    {
      label: '转轮空化数',
      value: formatValue(
        simulationResult?.data?.vof?.runner_cavitation_bubble_count,
      ),
    },
    {
      label: '叶片空化面积',
      value: simulationResult?.data?.vof?.blade_cavitation_area
        ? `${formatValue(
            simulationResult?.data?.vof?.blade_cavitation_area,
          )} m²`
        : '',
    },
  ],
  涡带分布: [
    {
      label: '集中部位',
      value: formatValue(
        simulationResult?.data?.vortex?.vortex_concentration_location,
      ),
    },
  ],
});

// 获取按钮配置
export const getButtonsByCard = (activeCard: string): ButtonConfig[] => {
  const baseButtons = [
    {
      key: '复位',
      icon: '/images/icons/reset.png',
      activeIcon: '/images/icons/reset.png',
      text: '复位',
    },
    {
      key: '几何',
      icon: '/images/icons/geometry_a.png',
      activeIcon: '/images/icons/geometry_b.png',
      text: '几何',
    },
    {
      key: '剖面',
      icon: '/images/icons/section_a.png',
      activeIcon: '/images/icons/section_b.png',
      text: '剖面',
    },
  ];

  if (activeCard === '速度场') {
    return [
      ...baseButtons,
      {
        key: '流线',
        icon: '/images/icons/flow_line_a.png',
        activeIcon: '/images/icons/flow_line_b.png',
        text: '流线',
      },
    ];
  } else if (activeCard === '空化分布' || activeCard === '涡带分布') {
    return [
      {
        key: '复位',
        icon: '/images/icons/reset.png',
        activeIcon: '/images/icons/reset.png',
        text: '复位',
      },
      {
        key: '几何',
        icon: '/images/icons/geometry_a.png',
        activeIcon: '/images/icons/geometry_b.png',
        text: '几何',
      },
    ];
  }
  return baseButtons;
};
