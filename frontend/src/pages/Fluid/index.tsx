import Header from '@/components/Header';
import React, { useState } from 'react';
import Left from './components/Left';
import MultiModelInfo from './components/MultiModelInfo';
import SingleModelInfo from './components/SingleModelInfo';
import { Wrapper } from './styles';

const FluidPage: React.FC = () => {
  const [selectedUnit, setSelectedUnit] = useState('1号机组');
  const [selectedTitle, setSelectedTitle] = useState('fluidSimulation');
  const [conditionType, setConditionType] = useState<'single' | 'multiple'>(
    'single',
  );
  const [activeCard, setActiveCard] = useState('速度场');

  // 单工况查询参数状态
  const [queryParams, setQueryParams] = useState<{
    effective_head: number;
    active_power: number;
  } | null>(null);

  // 多工况查询参数状态
  const [multiQueryParams, setMultiQueryParams] = useState<{
    condition1: { effective_head: number; active_power: number };
    condition2: { effective_head: number; active_power: number };
  } | null>(null);

  // 多工况表单状态 - 提升到父组件管理
  const [multiConditionState, setMultiConditionState] = useState<{
    selectedHead: number | null;
    selectedPower: number | null;
    selectedHead2: number | null;
    selectedPower2: number | null;
  }>({
    selectedHead: null,
    selectedPower: null,
    selectedHead2: null,
    selectedPower2: null,
  });

  // 处理单工况查询结果
  const handleQueryResult = (params: {
    effective_head: number;
    active_power: number;
  }) => {
    setQueryParams(params);
  };

  // 处理多工况查询结果
  const handleMultiQueryResult = (params: {
    condition1: { effective_head: number; active_power: number };
    condition2: { effective_head: number; active_power: number };
  }) => {
    setMultiQueryParams(params);
  };

  // 处理多工况状态更新
  const handleMultiConditionStateChange = (newState: {
    selectedHead?: number | null;
    selectedPower?: number | null;
    selectedHead2?: number | null;
    selectedPower2?: number | null;
  }) => {
    setMultiConditionState((prev) => ({ ...prev, ...newState }));
  };

  return (
    <Wrapper>
      <div className="container">
        <Header selectedTitle={selectedTitle} onUnitChange={setSelectedTitle} />

        <div className="content">
          <Left
            selectedUnit={selectedUnit}
            onUnitChange={setSelectedUnit}
            conditionType={conditionType}
            setConditionType={setConditionType}
            onQueryResult={handleQueryResult}
            onMultiQueryResult={handleMultiQueryResult}
            multiConditionState={multiConditionState}
            onMultiConditionStateChange={handleMultiConditionStateChange}
          />

          <div className="body">
            {/* 中间3D模型展示区 */}
            <div className="modelView">
              {conditionType === 'single' ? (
                <SingleModelInfo
                  activeCard={activeCard}
                  setActiveCard={setActiveCard}
                  queryParams={queryParams}
                />
              ) : (
                <MultiModelInfo
                  key={JSON.stringify(multiQueryParams)}
                  activeCard={activeCard}
                  setActiveCard={setActiveCard}
                  queryParams={multiQueryParams}
                />
              )}
            </div>
          </div>
        </div>
      </div>
    </Wrapper>
  );
};

export default FluidPage;
