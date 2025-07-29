// 卡片和按钮配置、格式化函数
import { CardConfig } from './types';

// 格式化数值为两位小数的函数
export const formatValue = (value: any) => {
  if (typeof value === 'number') {
    return value.toFixed(2);
  }
  return value ?? '';
};

// 生成卡片配置
export const getCardConfig = (simulationResult: any): CardConfig => ({
  位移场: [
    {
      label: '最大位移',
      value: simulationResult?.deplace?.max
        ? `${simulationResult.deplace.max.toFixed(4)} mm`
        : '--',
    },
    {
      label: '出现部位',
      value: simulationResult?.deplace?.max_displacement_location || '--',
    },
  ],
  应力场: [
    {
      label: '最大应力',
      value: simulationResult?.contrainte?.max
        ? `${simulationResult.contrainte.max.toFixed(2)} MPa`
        : '--',
    },
    {
      label: '出现部位',
      value: simulationResult?.contrainte?.max_stress_location || '--',
    },
  ],
});

// 获取按钮配置 - 结构分析只有几何、剖面、复位
export const getButtonsByCard = () => {
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

  return baseButtons;
};
