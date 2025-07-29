import { getStructuralResult } from '@/api/request';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { core, ModelContainer } from 'ys-dte';
import { getButtonsByCard, getCardConfig } from '../../config';
import { ResourceManager } from '../../ResourceManager';
import { SimulationResult, SingleModelInfoProps } from '../../types';
import { ModelLoader } from '../ModelLoader';
import VirtualProbe from '../VirtualProbe';
import { SingleModelWrapper } from './styles';

const SingleModelInfo: React.FC<SingleModelInfoProps> = ({
  activeCard,
  setActiveCard,
  queryParams,
}) => {
  const [activeButton, setActiveButton] = useState<string>('');
  const [simulationResult, setSimulationResult] =
    useState<SimulationResult | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  // 使用Ref跟踪资源状态
  const containerRef = useRef<any>(null);
  const resourceManagerRef = useRef<ResourceManager | null>(null);
  const modelLoaderRef = useRef<ModelLoader>(new ModelLoader());
  const initialCameraStateRef = useRef<any>(null);

  // 当查询参数变化时获取结果
  useEffect(() => {
    let isCancelled = false;

    const fetchResult = async () => {
      if (queryParams) {
        try {
          setLoading(true);
          const result = (await getStructuralResult(queryParams)) as {
            code: number;
            msg: string;
            data: { data: SimulationResult };
          };

          // 检查是否已被取消
          if (isCancelled) return;

          // console.log('result', result.data.data);

          setSimulationResult(result.data.data);
          setActiveCard('位移场');
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

    fetchResult();

    // 清理函数，用于取消正在进行的请求
    return () => {
      isCancelled = true;
    };
  }, [queryParams]);

  // 初始化三维场景
  useEffect(() => {
    if (!simulationResult) return;

    // 如果容器已存在，先清理资源
    if (containerRef.current) {
      try {
        resourceManagerRef.current?.dispose();
        containerRef.current.dispose();
      } catch (error) {
        console.warn('清理现有容器时出错:', error);
      }
    }

    const container = new ModelContainer({
      domId: 'model-container',
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
    camera.position.set(0, 0, 15);

    // 设置控制器限制
    controls.maxPolarAngle = Math.PI;
    controls.minPolarAngle = 0;
    controls.enablePan = false;

    // 设置渲染器
    renderer.autoClear = true;
    renderer.setClearAlpha(0);

    controls.update();

    containerRef.current = container;
    resourceManagerRef.current = new ResourceManager(container);
    modelLoaderRef.current.setResourceManager(resourceManagerRef.current);

    // 保存初始相机状态
    initialCameraStateRef.current = {
      position: camera.position.clone(),
      quaternion: camera.quaternion.clone(),
      target: controls.target.clone(),
      ...('fov' in camera && { fov: camera.fov }),
    };

    return () => {
      // 清理资源
      try {
        resourceManagerRef.current?.dispose();
        if (containerRef.current) {
          containerRef.current.dispose();
        }
      } catch (error) {
        console.warn('清理资源时出错:', error);
      }
      containerRef.current = null;
      resourceManagerRef.current = null;
    };
  }, [simulationResult]);

  // 当模型配置变化时加载新模型
  const loadModel = useCallback(async () => {
    if (
      !containerRef.current ||
      !simulationResult ||
      !resourceManagerRef.current
    ) {
      return;
    }
    setLoading(true);
    await modelLoaderRef.current.loadModel(
      simulationResult,
      activeCard,
      activeButton,
    );
    setLoading(false);
  }, [activeCard, activeButton, simulationResult]);

  useEffect(() => {
    if (containerRef.current && simulationResult) {
      loadModel();
    }
  }, [activeCard, activeButton, simulationResult, loadModel]);

  // 复位相机到初始状态
  const resetCamera = useCallback(() => {
    if (containerRef.current && initialCameraStateRef.current) {
      const { camera, controls } = containerRef.current;
      const initialState = initialCameraStateRef.current;

      // 重置相机位置和朝向
      camera.position.copy(initialState.position);
      camera.quaternion.copy(initialState.quaternion);

      // 只有透视相机才有fov属性
      if ('fov' in camera && 'fov' in initialState) {
        camera.fov = initialState.fov;
        camera.updateProjectionMatrix();
      }

      // 重置控制器
      controls.target.copy(initialState.target);
      controls.update();
    }
  }, []);

  // 处理按钮点击
  const handleButtonClick = useCallback(
    (buttonKey: string) => {
      if (buttonKey === '复位') {
        resetCamera();
      } else {
        setActiveButton(activeButton === buttonKey ? '' : buttonKey);
      }
    },
    [resetCamera, activeButton],
  );

  const cardConfig = getCardConfig(simulationResult);
  const buttons = getButtonsByCard();

  return (
    <SingleModelWrapper>
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

      <div id="model-container" className="model-container">
        {loading && (
          <div className="loading-overlay">
            <div className="spinner"></div>
            <div>模型加载中...</div>
          </div>
        )}
      </div>

      {containerRef.current && resourceManagerRef.current && (
        <VirtualProbe
          container={containerRef.current}
          resourceManager={resourceManagerRef.current}
          activeCard={activeCard}
          activeButton={activeButton}
        />
      )}

      <div className="info-cards">
        {Object.entries(cardConfig).map(([card, items]) => (
          <div
            key={card}
            className={`info-card ${activeCard === card ? 'active' : ''}`}
            onClick={() => setActiveCard(card)}
          >
            <div className="card-title">
              <span>{card}</span>
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
    </SingleModelWrapper>
  );
};

export default SingleModelInfo;
