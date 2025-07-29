import alovaInstance from './alovaInstance';

// 1. 查询所有流体仿真模拟工况
export const getFluidCondition = () =>
  alovaInstance.Get('/api/fluid/conditions');

// 2. 查询流体仿真模拟结果
export const getFluidResult = (params: {
  effective_head: number;
  active_power: number;
}) => alovaInstance.Get('/api/fluid/result', { params });

// 3. 查询所有结构仿真模拟工况
export const getStructuralConditions = () =>
  alovaInstance.Get('/api/structural/conditions');

// 4. 查询结构仿真模拟结果
export const getStructuralResult = (params: {
  effective_head: number;
  active_power: number;
}) => alovaInstance.Get('/api/structural/result', { params });

// 5. 应变监测（无参）
// export const getStrainMonitoring = () =>
//   alovaInstance.Get('/strain_monitoring');

// 5. 应变监测（有参）
export const getStrainMonitoring = (start_time: number, stop_time: number) =>
  alovaInstance.Get('/api/strain_monitoring', {
    params: {
      start_time,
      stop_time,
    },
  });
