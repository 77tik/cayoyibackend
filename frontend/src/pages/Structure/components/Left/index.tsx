import { Radio } from 'antd';
import React, { useEffect } from 'react';
import ConditionPanel from '../ConditionPanel';
import MultiConditionPanel from '../MultiConditionPanel';
import { LeftWrapper } from './styles';

// const { Option } = Select;

interface LeftProps {
  selectedUnit: string;
  onUnitChange: (value: string) => void;
  conditionType: 'single' | 'multiple';
  setConditionType: (value: 'single' | 'multiple') => void;
  //
  onQueryResult: (params: {
    effective_head: number;
    active_power: number;
  }) => void;
  onMultiQueryResult: (params: {
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

const Left: React.FC<LeftProps> = ({
  selectedUnit,
  onUnitChange,
  conditionType,
  setConditionType,
  onQueryResult,
  onMultiQueryResult,
  multiConditionState,
  onMultiConditionStateChange,
}) => {
  const [selectedCondition, setSelectedCondition] = React.useState<{
    effective_head: number;
    active_power: number;
  } | null>(null);
  // 添加状态标记首次查询是否已完成
  const [initialQueryDone, setInitialQueryDone] = React.useState(false);

  // 处理查询按钮点击
  const handleQuery = () => {
    if (selectedCondition) {
      onQueryResult(selectedCondition);
    }
  };

  // 当selectedCondition变化且首次查询未完成时自动触发查询
  useEffect(() => {
    if (selectedCondition && !initialQueryDone) {
      onQueryResult(selectedCondition);
      setInitialQueryDone(true); // 标记首次查询已完成
    }
  }, [selectedCondition, initialQueryDone, onQueryResult]);

  // 当切换机组或工况类型时重置首次查询状态
  useEffect(() => {
    setInitialQueryDone(false);
  }, [selectedUnit, conditionType]);

  return (
    <LeftWrapper>
      <div className="left">
        {/* 机组选择 */}
        <div className="select-wrapper">
          <select
            value={selectedUnit}
            onChange={(e) => onUnitChange(e.target.value)}
            className="custom-select"
          >
            <option value="1号机组">1号机组</option>
            <option value="2号机组">2号机组</option>
            <option value="3号机组">3号机组</option>
            <option value="4号机组">4号机组</option>
          </select>
        </div>
        {/* 左侧参数面板 */}
        <div className="leftPanel">
          <div className="panelHeader">
            <div className="panelHeaderTitle">低有效水头工况</div>
            <Radio.Group
              value={conditionType}
              onChange={(e) => setConditionType(e.target.value)}
              className="conditionSwitcher"
            >
              <Radio.Button value="single">单工况模拟</Radio.Button>
              <Radio.Button value="multiple">多工况模拟</Radio.Button>
            </Radio.Group>
          </div>
          {conditionType === 'single' ? (
            <ConditionPanel onConditionSelect={setSelectedCondition} />
          ) : (
            <MultiConditionPanel
              onQueryResult={onMultiQueryResult}
              multiConditionState={multiConditionState}
              onMultiConditionStateChange={onMultiConditionStateChange}
            />
          )}
          {conditionType === 'single' && (
            <button
              className="queryButton"
              type="button"
              onClick={handleQuery}
              disabled={!selectedCondition}
            >
              结果查询
            </button>
          )}
        </div>
      </div>
    </LeftWrapper>
  );
};

export default Left;
