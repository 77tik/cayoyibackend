import { getFluidResult } from '@/api/request';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { core, ModelContainer } from 'ys-dte';
import { getButtonsByCard, getCardConfig } from '../../config';
import { ResourceManager } from '../../ResourceManager';
import { ProfileType, SimulationResult } from '../../types';
import { ModelLoader } from '../ModelLoader';
import VirtualProbe from '../VirtualProbe';
import { MultipleModelWrapper } from './styles';

// 多工况的信息显示组件
interface MultiModelInfoProps {
  activeCard: string;
  setActiveCard: (card: string) => void;
  queryParams: {
    condition1: { effective_head: number; active_power: number };
    condition2: { effective_head: number; active_power: number };
  } | null;
}

const MultiModelInfo: React.FC<MultiModelInfoProps> = ({
  activeCard,
  setActiveCard,
  queryParams,
}) => {
  const [activeButton, setActiveButton] = useState<string>('');
  const [simulationResult1, setSimulationResult1] =
    useState<SimulationResult | null>(null);
  const [simulationResult2, setSimulationResult2] =
    useState<SimulationResult | null>(null);
  const [profileType, setProfileType] = useState<ProfileType>('h');
  const [loading, setLoading] = useState<boolean>(true);
  const [showProfilePanel, setShowProfilePanel] = useState(false);

  // 使用Ref跟踪资源状态 - 工况一
  const container1Ref = useRef<any>(null);
  const resourceManager1Ref = useRef<ResourceManager | null>(null);
  const modelLoader1Ref = useRef<ModelLoader>(new ModelLoader());
  const initialCameraState1Ref = useRef<any>(null);

  // 使用Ref跟踪资源状态 - 工况二
  const container2Ref = useRef<any>(null);
  const resourceManager2Ref = useRef<ResourceManager | null>(null);
  const modelLoader2Ref = useRef<ModelLoader>(new ModelLoader());
  const initialCameraState2Ref = useRef<any>(null);

  // 动态生成唯一domId，确保每次切换都重新挂载
  const domId1 = `model-container-1-${activeCard}-${
    queryParams?.condition1?.effective_head ?? ''
  }-${queryParams?.condition1?.active_power ?? ''}`;
  const domId2 = `model-container-2-${activeCard}-${
    queryParams?.condition2?.effective_head ?? ''
  }-${queryParams?.condition2?.active_power ?? ''}`;

  // 当查询参数变化时获取结果
  useEffect(() => {
    let isCancelled = false;

    const fetchResults = async () => {
      if (queryParams) {
        try {
          setLoading(true);

          // 获取工况一的结果
          const result1 = (await getFluidResult(queryParams.condition1)) as {
            data: SimulationResult;
          };

          // 检查是否已被取消
          if (isCancelled) return;
          setSimulationResult1(result1.data);

          // 获取工况二的结果
          const result2 = (await getFluidResult(queryParams.condition2)) as {
            data: SimulationResult;
          };

          // 检查是否已被取消
          if (isCancelled) return;
          setSimulationResult2(result2.data);

          setActiveCard('速度场');
          setActiveButton('');
        } catch (error) {
          console.error('获取仿真结果失败:', error);
        } finally {
          if (!isCancelled) {
            setLoading(false);
          }
        }
      }
    };

    fetchResults();

    // 清理函数，用于取消正在进行的请求
    return () => {
      isCancelled = true;
    };
  }, [queryParams]);

  // 初始化工况一三维场景
  useEffect(() => {
    if (!simulationResult1) return;

    // 如果容器已存在，先清理资源
    if (container1Ref.current) {
      try {
        resourceManager1Ref.current?.dispose();
        container1Ref.current.dispose();
      } catch (error) {
        console.warn('清理工况一容器时出错:', error);
      }
    }

    const container = new ModelContainer({
      domId: domId1,
      debug: false,
      antialias: true,
    });

    const { scene, camera, renderer, controls } = container;

    // 初始化灯光
    const light1 = new core.AmbientLight(0xffffff, 0.6);
    scene.add(light1);

    const light2 = new core.DirectionalLight(0xffffff, 0.8);
    light2.position.set(-100, -100, 100);
    scene.add(light2);

    const light3 = new core.DirectionalLight(0xffffff, 0.4);
    light3.position.set(100, 100, 100);
    scene.add(light3);

    // 设置相机位置
    camera.position.set(0, 0, 100);

    // 设置控制器限制
    controls.maxPolarAngle = Math.PI;
    controls.minPolarAngle = 0;
    controls.enablePan = false;

    // 设置渲染器
    renderer.autoClear = true;
    renderer.setClearAlpha(0);

    controls.update();

    container1Ref.current = container;
    resourceManager1Ref.current = new ResourceManager(container);
    modelLoader1Ref.current.setResourceManager(resourceManager1Ref.current);

    // 保存初始相机状态
    initialCameraState1Ref.current = {
      position: camera.position.clone(),
      quaternion: camera.quaternion.clone(),
      target: controls.target.clone(),
      ...('fov' in camera && { fov: camera.fov }),
    };

    return () => {
      // 清理资源
      try {
        if (container1Ref.current) {
          // 先手动 detach domElement，防止 ys-dte dispose 时 removeChild 报错
          const domElement = container1Ref.current.renderer?.domElement;
          if (domElement && domElement.parentNode) {
            try {
              domElement.parentNode.removeChild(domElement);
            } catch (e) {
              // 已经被移除则忽略
            }
          }
        }
        resourceManager1Ref.current?.dispose();
        if (container1Ref.current) {
          container1Ref.current.dispose();
        }
      } catch (error) {
        console.warn('清理工况一资源时出错:', error);
      }
      container1Ref.current = null;
      resourceManager1Ref.current = null;
    };
  }, [simulationResult1, domId1]);

  // 初始化工况二三维场景
  useEffect(() => {
    if (!simulationResult2) return;

    // 如果容器已存在，先清理资源
    if (container2Ref.current) {
      try {
        resourceManager2Ref.current?.dispose();
        container2Ref.current.dispose();
      } catch (error) {
        console.warn('清理工况二容器时出错:', error);
      }
    }

    const container = new ModelContainer({
      domId: domId2,
      debug: false,
      antialias: true,
    });

    const { scene, camera, renderer, controls } = container;

    // 初始化灯光
    const light1 = new core.AmbientLight(0xffffff, 0.6);
    scene.add(light1);

    const light2 = new core.DirectionalLight(0xffffff, 0.8);
    light2.position.set(-100, -100, 100);
    scene.add(light2);

    const light3 = new core.DirectionalLight(0xffffff, 0.4);
    light3.position.set(100, 100, 100);
    scene.add(light3);

    // 设置相机位置
    camera.position.set(0, 0, 100);

    // 设置控制器限制
    controls.maxPolarAngle = Math.PI;
    controls.minPolarAngle = 0;
    controls.enablePan = false;

    // 设置渲染器
    renderer.autoClear = true;
    renderer.setClearAlpha(0);

    controls.update();

    container2Ref.current = container;
    resourceManager2Ref.current = new ResourceManager(container);
    modelLoader2Ref.current.setResourceManager(resourceManager2Ref.current);

    // 保存初始相机状态
    initialCameraState2Ref.current = {
      position: camera.position.clone(),
      quaternion: camera.quaternion.clone(),
      target: controls.target.clone(),
      ...('fov' in camera && { fov: camera.fov }),
    };

    return () => {
      // 清理资源
      try {
        if (container2Ref.current) {
          // 先手动 detach domElement，防止 ys-dte dispose 时 removeChild 报错
          const domElement = container2Ref.current.renderer?.domElement;
          if (domElement && domElement.parentNode) {
            try {
              domElement.parentNode.removeChild(domElement);
            } catch (e) {
              // 已经被移除则忽略
            }
          }
        }
        resourceManager2Ref.current?.dispose();
        if (container2Ref.current) {
          container2Ref.current.dispose();
        }
      } catch (error) {
        console.warn('清理工况二资源时出错:', error);
      }
      container2Ref.current = null;
      resourceManager2Ref.current = null;
    };
  }, [simulationResult2, domId2]);

  // 当模型配置变化时加载新模型
  const loadModels = useCallback(async () => {
    if (!simulationResult1 || !simulationResult2) return;

    setLoading(true);

    // 加载工况一模型
    if (container1Ref.current && resourceManager1Ref.current) {
      await modelLoader1Ref.current.loadModel(
        simulationResult1,
        activeCard,
        activeButton,
        profileType,
      );
    }

    // 加载工况二模型
    if (container2Ref.current && resourceManager2Ref.current) {
      await modelLoader2Ref.current.loadModel(
        simulationResult2,
        activeCard,
        activeButton,
        profileType,
      );
    }

    setLoading(false);
  }, [
    activeCard,
    activeButton,
    profileType,
    simulationResult1,
    simulationResult2,
  ]);

  useEffect(() => {
    if (
      container1Ref.current &&
      container2Ref.current &&
      simulationResult1 &&
      simulationResult2
    ) {
      loadModels();
    }
  }, [
    activeCard,
    activeButton,
    profileType,
    simulationResult1,
    simulationResult2,
    loadModels,
  ]);

  // 复位相机到初始状态
  const resetCamera = useCallback(() => {
    // 复位工况一相机
    if (container1Ref.current && initialCameraState1Ref.current) {
      const { camera, controls } = container1Ref.current;
      const initialState = initialCameraState1Ref.current;

      camera.position.copy(initialState.position);
      camera.quaternion.copy(initialState.quaternion);

      if ('fov' in camera && 'fov' in initialState) {
        camera.fov = initialState.fov;
        camera.updateProjectionMatrix();
      }

      controls.target.copy(initialState.target);
      controls.update();
    }

    // 复位工况二相机
    if (container2Ref.current && initialCameraState2Ref.current) {
      const { camera, controls } = container2Ref.current;
      const initialState = initialCameraState2Ref.current;

      camera.position.copy(initialState.position);
      camera.quaternion.copy(initialState.quaternion);

      if ('fov' in camera && 'fov' in initialState) {
        camera.fov = initialState.fov;
        camera.updateProjectionMatrix();
      }

      controls.target.copy(initialState.target);
      controls.update();
    }
  }, []);

  // 处理按钮点击
  const handleButtonClick = (buttonKey: string) => {
    if (buttonKey === '复位') {
      resetCamera();
    } else {
      setActiveButton(activeButton === buttonKey ? '' : buttonKey);
      if (buttonKey === '剖面' && activeButton !== '剖面') {
        setShowProfilePanel(true);
      } else if (buttonKey === '剖面' && activeButton === '剖面') {
        setShowProfilePanel(false);
      }
    }
  };

  const cardConfig1 = getCardConfig(simulationResult1);
  const cardConfig2 = getCardConfig(simulationResult2);
  const buttons = getButtonsByCard(activeCard);

  return (
    <MultipleModelWrapper>
      {/* 右侧功能按钮 */}
      <div className="rightBtn">
        {buttons.map((button) => (
          <button
            key={button.key}
            type="button"
            className={`function-btn ${
              activeButton === button.key ? 'active' : ''
            }`}
            onClick={() => handleButtonClick(button.key)}
          >
            <img
              src={
                activeButton === button.key ? button.activeIcon : button.icon
              }
              alt={button.text}
            />
            {button.text}
          </button>
        ))}
      </div>

      {/* 剖面按钮 */}
      {activeButton === '剖面' &&
        (activeCard === '速度场' || activeCard === '压力场') &&
        showProfilePanel && (
          <div className="profile-panel">
            <div className="profile-panel-header">
              <div className="profile-tabs">
                <div
                  className={`profile-tab ${
                    profileType === 'h' ? 'active' : ''
                  }`}
                  onClick={() => setProfileType('h')}
                >
                  横剖面
                </div>
                <div
                  className={`profile-tab ${
                    profileType === 'v' ? 'active' : ''
                  }`}
                  onClick={() => setProfileType('v')}
                >
                  纵剖面
                </div>
              </div>
              <div
                className="profile-close"
                onClick={() => setShowProfilePanel(false)}
                title="关闭"
              >
                ×
              </div>
            </div>
            <div className="profile-panel-content">
              {profileType === 'h' ? (
                <img
                  src="/images/icons/cross_section.png"
                  alt="横剖面"
                  className="profile-panel-content-h"
                />
              ) : (
                <img
                  src="/images/icons/longisection.png"
                  alt="纵剖面"
                  className="profile-panel-content-v"
                />
              )}
            </div>
          </div>
        )}

      <div className="model-container">
        {/* 工况一模型 */}
        <div className="model-item">
          <div className="model-item-title">工况一</div>
          <div className="model-item-content">
            <div id={domId1} key={domId1} className="model-container-inner">
              {loading && (
                <div className="loading-overlay">
                  <div className="spinner"></div>
                  <div>模型加载中...</div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* 工况二模型 */}
        <div className="model-item">
          <div className="model-item-title">工况二</div>
          <div className="model-item-content">
            <div id={domId2} key={domId2} className="model-container-inner">
              {loading && (
                <div className="loading-overlay">
                  <div className="spinner"></div>
                  <div>模型加载中...</div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* 工况一信息卡片 */}
      <div className="info-cards">
        {Object.entries(cardConfig1).map(([card, items]) => (
          <div
            key={`condition1-${card}`}
            className={`info-card ${activeCard === card ? 'active' : ''}`}
            onClick={() => setActiveCard(card)}
          >
            <div className="card-title">
              <span>{card}-工况一</span>
              <img src="/images/icons/full.png" alt={card} />
            </div>
            <div className="card-content">
              {items.map((item, index) => (
                <div key={index} className="info-item">
                  <span className="label">{item.label}</span>
                  <span className="value">{item.value}</span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* 工况二信息卡片 */}
      <div className="info-cards">
        {Object.entries(cardConfig2).map(([card, items]) => (
          <div
            key={`condition2-${card}`}
            className={`info-card ${activeCard === card ? 'active' : ''}`}
            onClick={() => setActiveCard(card)}
          >
            <div className="card-title">
              <span>{card}-工况二</span>
              <img src="/images/icons/full.png" alt={card} />
            </div>
            <div className="card-content">
              {items.map((item, index) => (
                <div key={index} className="info-item">
                  <span className="label">{item.label}</span>
                  <span className="value">{item.value}</span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* 工况一探针，仅在容器和资源管理器都存在时渲染 */}
      {container1Ref.current && resourceManager1Ref.current && (
        <VirtualProbe
          container={container1Ref.current}
          resourceManager={resourceManager1Ref.current}
          activeCard={activeCard}
          activeButton={activeButton}
        />
      )}
      {/* 工况二探针，仅在容器和资源管理器都存在时渲染 */}
      {container2Ref.current && resourceManager2Ref.current && (
        <VirtualProbe
          container={container2Ref.current}
          resourceManager={resourceManager2Ref.current}
          activeCard={activeCard}
          activeButton={activeButton}
        />
      )}
    </MultipleModelWrapper>
  );
};

export default MultiModelInfo;
