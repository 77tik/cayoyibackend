import { getStructuralConditions } from '@/api/request';
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

// 定义工况数据类型
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

interface ConditionPanelProps {
  onConditionSelect: (condition: {
    effective_head: number;
    active_power: number;
  }) => void;
}

const ConditionPanel: React.FC<ConditionPanelProps> = ({
  onConditionSelect,
}) => {
  const [activeIndex, setActiveIndex] = useState<number | null>(null);
  const [activePowerIndex, setActivePowerIndex] = useState<number | null>(null);
  const [conditions, setConditions] = useState<Condition[]>([]);

  // 获取唯一的有效水头值
  const uniqueHeads = React.useMemo(() => {
    return Array.from(new Set(conditions.map((c) => c.effective_head))).sort(
      (a, b) => a - b,
    );
  }, [conditions]);

  // 获取唯一的有功功率值
  const uniquePowers = React.useMemo(() => {
    return Array.from(new Set(conditions.map((c) => c.active_power))).sort(
      (a, b) => a - b,
    );
  }, [conditions]);

  // 获取当前选中的工况数据
  const currentCondition = React.useMemo(() => {
    if (activeIndex === null || activePowerIndex === null) return null;
    const selectedHead = uniqueHeads[activeIndex];
    const selectedPower = uniquePowers[activePowerIndex];
    return (
      conditions.find(
        (c) =>
          c.effective_head === selectedHead && c.active_power === selectedPower,
      ) || null
    );
  }, [activeIndex, activePowerIndex, uniqueHeads, uniquePowers, conditions]);

  // 模拟获取数据
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response =
          (await getStructuralConditions()) as Promise<FluidConditionResponse>;
        const fluData = (await response).data.data.conditions;

        setConditions(fluData);

        if (fluData.length > 0) {
          setActiveIndex(0);
          setActivePowerIndex(0);
          const firstCondition = fluData[0];
          // 只设置状态，不自动触发查询
          onConditionSelect({
            effective_head: firstCondition.effective_head,
            active_power: firstCondition.active_power,
          });
        }
      } catch (error) {
        console.error('获取流体工况数据失败:', error);
      }
    };
    fetchData();
  }, [onConditionSelect]);

  // 修改点击处理函数 - 不再自动触发查询
  const handleConditionSelect = (index: number) => {
    setActiveIndex(index);
    const selectedHead = uniqueHeads[index];
    const selectedPower =
      activePowerIndex !== null ? uniquePowers[activePowerIndex] : null;

    if (selectedPower !== null) {
      const condition = conditions.find(
        (c) =>
          c.effective_head === selectedHead && c.active_power === selectedPower,
      );
      if (condition) {
        // 只更新状态，不触发查询
        onConditionSelect({
          effective_head: condition.effective_head,
          active_power: condition.active_power,
        });
      }
    }
  };

  // 修改有功功率点击处理函数 - 不再自动触发查询
  const handlePowerSelect = (index: number, power: number) => {
    setActivePowerIndex(index);
    if (activeIndex !== null) {
      const selectedHead = uniqueHeads[activeIndex];
      const condition = conditions.find(
        (c) => c.effective_head === selectedHead && c.active_power === power,
      );
      if (condition) {
        // 只更新状态，不触发查询
        onConditionSelect({
          effective_head: condition.effective_head,
          active_power: condition.active_power,
        });
      }
    }
  };

  return (
    <Wrapper>
      <div className="conditionPanel">
        <div className="waterPoints">
          <div className="titleRow">
            <img src="/images/矩形备份19.png" className="titleIcon" />
            <div className="titleText">有效水头(m)</div>
            <div className="titleLine" />
          </div>
          <div className="pointsGrid">
            {uniqueHeads.map((head, i) => (
              <div
                key={head}
                className={`pointItem ${activeIndex === i ? 'active' : ''}`}
                onClick={() => handleConditionSelect(i)}
              >
                <span>{`${head}m`}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="powerSection">
          <div className="titleRow">
            <img src="/images/矩形备份19.png" className="titleIcon" />
            <div className="titleText">有功功率(MW)</div>
            <div className="titleLine" />
          </div>
          <div className="powerGrid">
            {uniquePowers.map((power, index) => (
              <div
                key={power}
                className={`powerItem ${
                  activePowerIndex === index ? 'active' : ''
                }`}
                onClick={() => handlePowerSelect(index, power)}
              >
                <span className="powerItemLeft">有功功率：{power}MW</span>
                <div className="divider" />
                <span className="powerItemRight">负荷-出力比：80%</span>
              </div>
            ))}
          </div>
        </div>

        <div className="envParams">
          <div className="titleRow">
            <img src="/images/矩形备份19.png" className="titleIcon" />
            <div className="titleText">运行环境参数</div>
            <div className="titleLine" />
          </div>
          <div className="paramGrid">
            {currentCondition &&
              [
                {
                  label: '入口流量',
                  value: currentCondition.inlet_flow,
                  unit: 'm³/s',
                  image: '/images/编组31.png',
                },
                {
                  label: '出口压力',
                  value: currentCondition.outlet_pressure,
                  unit: 'MPa',
                  image: '/images/编组42.png',
                },
                {
                  label: '转轮角速度',
                  value: currentCondition.runner_angular_velocity,
                  unit: 'Rev/min',
                  image: '/images/编组40.png',
                },
                {
                  label: '导叶开度',
                  value: currentCondition.guide_vane_opening,
                  unit: '%',
                  image: '/images/编组41.png',
                },
              ].map((param, index) => (
                <div key={index} className="paramItem">
                  <div className="paramImage">
                    <img src={param.image} alt={param.label} />
                  </div>
                  <div className="paramContent">
                    <div className="paramValue">
                      {param.value || 0}
                      <span className="paramUnit">{param.unit}</span>
                    </div>
                    <div className="paramLabel">{param.label}</div>
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
    </Wrapper>
  );
};

export default ConditionPanel;
