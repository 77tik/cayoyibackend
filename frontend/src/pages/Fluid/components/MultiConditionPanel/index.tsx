import { getFluidCondition } from '@/api/request';
import { Select } from 'antd';
import React, { useEffect, useState } from 'react';
import { Wrapper } from './styles';

// 定义API响应类型
interface FluidConditionResponse {
  data: {
    data: {
      conditions: Condition[];
    };
  };
}

interface Condition {
  id: number;
  name: string;
  effective_head: number;
  active_power: number;
  inlet_flow: number;
  outlet_pressure: number;
  runner_angular_velocity: number;
  guide_vane_opening: number;
  load_to_output_ratio: number;
}

interface MultiConditionPanelProps {
  onQueryResult: (params: {
    condition1: { effective_head: number; active_power: number };
    condition2: { effective_head: number; active_power: number };
  }) => void;
  // 多工况状态管理
  multiConditionState: {
    selectedHead: number | null;
    selectedPower: number | null;
    selectedHead2: number | null;
    selectedPower2: number | null;
  };
  onMultiConditionStateChange: (newState: {
    selectedHead?: number | null;
    selectedPower?: number | null;
    selectedHead2?: number | null;
    selectedPower2?: number | null;
  }) => void;
}

const MultiConditionPanel: React.FC<MultiConditionPanelProps> = ({
  onQueryResult,
  multiConditionState,
  onMultiConditionStateChange,
}) => {
  const [conditions, setConditions] = useState<Condition[]>([]);
  // 使用外部传入的状态，不再使用内部状态
  const { selectedHead, selectedPower, selectedHead2, selectedPower2 } =
    multiConditionState;

  // 获取唯一的有效水头值
  const uniqueHeads = React.useMemo(() => {
    return Array.from(new Set(conditions.map((c) => c.effective_head))).sort(
      (a, b) => a - b,
    );
  }, [conditions]);

  // 获取唯一的有功功率值
  const uniquePowers = React.useMemo(() => {
    return Array.from(new Set(conditions.map((c) => c.active_power))).sort(
      (a, b) => b - a,
    );
  }, [conditions]);

  // 获取当前选中的工况数据
  const currentCondition = React.useMemo(() => {
    if (selectedHead === null || selectedPower === null) return null;
    return (
      conditions.find(
        (c) =>
          c.effective_head === selectedHead && c.active_power === selectedPower,
      ) || null
    );
  }, [selectedHead, selectedPower, conditions]);

  // 获取工况二当前选中的工况数据
  const currentCondition2 = React.useMemo(() => {
    if (selectedHead2 === null || selectedPower2 === null) return null;
    return (
      conditions.find(
        (c) =>
          c.effective_head === selectedHead2 &&
          c.active_power === selectedPower2,
      ) || null
    );
  }, [selectedHead2, selectedPower2, conditions]);

  // 获取数据
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = (await getFluidCondition()) as FluidConditionResponse;
        const fluData = response.data.data.conditions;
        setConditions(fluData);
        if (fluData.length > 0 && !selectedHead && !selectedPower) {
          // 只在状态为空时设置默认值
          const firstHead = Array.from(
            new Set(fluData.map((c) => c.effective_head)),
          ).sort((a, b) => a - b)[0];
          const firstPower = Array.from(
            new Set(fluData.map((c) => c.active_power)),
          ).sort((a, b) => b - a)[0];
          onMultiConditionStateChange({
            selectedHead: firstHead,
            selectedPower: firstPower,
            selectedHead2: firstHead,
            selectedPower2: firstPower,
          });
        }
      } catch (e) {
        console.error('接口请求失败', e);
      }
    };
    fetchData();
  }, [selectedHead, selectedPower, onMultiConditionStateChange]);

  // 处理有效水头选择
  const handleHeadChange = (value: number) => {
    onMultiConditionStateChange({ selectedHead: value });
  };

  // 处理有功功率选择
  const handlePowerSelect = (power: number) => {
    onMultiConditionStateChange({ selectedPower: power });
  };

  // 处理工况二有效水头选择
  const handleHeadChange2 = (value: number) => {
    onMultiConditionStateChange({ selectedHead2: value });
  };

  // 处理工况二有功功率选择
  const handlePowerSelect2 = (power: number) => {
    onMultiConditionStateChange({ selectedPower2: power });
  };

  // 处理查询按钮点击
  const handleQuery = () => {
    if (
      selectedHead !== null &&
      selectedPower !== null &&
      selectedHead2 !== null &&
      selectedPower2 !== null
    ) {
      onQueryResult({
        condition1: {
          effective_head: selectedHead,
          active_power: selectedPower,
        },
        condition2: {
          effective_head: selectedHead2,
          active_power: selectedPower2,
        },
      });
    }
  };

  return (
    <Wrapper>
      <div className="multiConditionPanel">
        {/* 工况设置-工况一 */}
        <div className="titleRow">
          <img src="/images/矩形备份19.png" className="titleIcon" />
          <div className="titleText">工况设置-工况一</div>
          <div className="titleLine" />
        </div>

        {/* 工况设置-工况一-内容 */}
        <div className="multipleGrid">
          {/* 有效水头选择框 */}
          <div className="multipleItem multipleSelect">
            <span>有效水头：</span>
            <Select
              bordered={false}
              placeholder="请选择"
              value={selectedHead}
              onChange={handleHeadChange}
              style={{
                width: 90,
                color: 'white',
              }}
              dropdownStyle={{
                color: 'white',
              }}
            >
              {uniqueHeads.map((head) => (
                <Select.Option key={head} value={head}>
                  {head}m
                </Select.Option>
              ))}
            </Select>
          </div>

          {/* 有功功率选择 */}
          <div className="multipleItem multipleSelect">
            <span>有功功率：</span>
            <Select
              bordered={false}
              placeholder="请选择"
              value={selectedPower}
              onChange={handlePowerSelect}
              style={{
                width: 90,
                color: 'white',
              }}
              dropdownStyle={{
                color: 'white',
              }}
            >
              {uniquePowers.map((power) => (
                <Select.Option key={power} value={power}>
                  {power}MW
                </Select.Option>
              ))}
            </Select>
          </div>

          {/* 运行环境参数 */}
          {currentCondition &&
            [
              {
                label: '负荷-出力比',
                value: currentCondition.load_to_output_ratio,
                unit: '%',
              },
              {
                label: '入口流量',
                value: currentCondition.inlet_flow,
                unit: 'm³/s',
              },
              {
                label: '导叶开度',
                value: currentCondition.guide_vane_opening,
                unit: '%',
              },
              {
                label: '出口压力',
                value: currentCondition.outlet_pressure,
                unit: 'MPa',
              },
              {
                label: '转轮角速度',
                value: currentCondition.runner_angular_velocity,
                unit: 'Rev/min',
              },
            ].map((param, index) => (
              <div key={index} className="multipleItem">
                <span>{`${param.label}: ${param.value || 0}${
                  param.unit
                }`}</span>
              </div>
            ))}
        </div>

        {/* 工况设置-工况二  */}
        <div className="titleRow">
          <img src="/images/矩形备份19.png" className="titleIcon" />
          <div className="titleText">工况设置-工况二</div>
          <div className="titleLine" />
        </div>

        {/* 工况设置-工况二-内容 */}
        <div className="multipleGrid">
          {/* 有效水头选择框 */}
          <div className="multipleItem multipleSelect">
            <span>有效水头：</span>
            <Select
              bordered={false}
              placeholder="请选择"
              value={selectedHead2}
              onChange={handleHeadChange2}
              style={{
                width: 90,
                color: 'white',
              }}
              dropdownStyle={{
                color: 'white',
              }}
            >
              {uniqueHeads.map((head) => (
                <Select.Option key={head} value={head}>
                  {head}m
                </Select.Option>
              ))}
            </Select>
          </div>

          {/* 有功功率选择 */}
          <div className="multipleItem multipleSelect">
            <span>有功功率：</span>
            <Select
              bordered={false}
              placeholder="请选择"
              value={selectedPower2}
              onChange={handlePowerSelect2}
              style={{
                width: 90,
                color: 'white',
              }}
              dropdownStyle={{
                color: 'white',
              }}
            >
              {uniquePowers.map((power) => (
                <Select.Option key={power} value={power}>
                  {power}MW
                </Select.Option>
              ))}
            </Select>
          </div>

          {/* 运行环境参数 */}
          {currentCondition2 &&
            [
              {
                label: '负荷-出力比',
                value: currentCondition2.load_to_output_ratio,
                unit: '%',
              },
              {
                label: '入口流量',
                value: currentCondition2.inlet_flow,
                unit: 'm³/s',
              },
              {
                label: '导叶开度',
                value: currentCondition2.guide_vane_opening,
                unit: '%',
              },
              {
                label: '出口压力',
                value: currentCondition2.outlet_pressure,
                unit: 'MPa',
              },
              {
                label: '转轮角速度',
                value: currentCondition2.runner_angular_velocity,
                unit: 'Rev/min',
              },
            ].map((param, index) => (
              <div key={index} className="multipleItem">
                <span>{`${param.label}: ${param.value || 0}${
                  param.unit
                }`}</span>
              </div>
            ))}
        </div>

        {/* 查询按钮 */}
        <button
          className="multiQueryButton"
          type="button"
          onClick={handleQuery}
          disabled={
            !selectedHead || !selectedPower || !selectedHead2 || !selectedPower2
          }
        >
          结果查询
        </button>
      </div>
    </Wrapper>
  );
};

export default MultiConditionPanel;
