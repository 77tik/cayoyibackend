import { getStructuralResult } from '@/api/request';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { core, ModelContainer } from 'ys-dte';
import { getButtonsByCard } from '../../config';
import { ResourceManager } from '../../ResourceManager';
import { SimulationResult } from '../../types';
import { ModelLoader } from '../ModelLoader';
import VirtualProbe from '../VirtualProbe';
import { MultipleModelWrapper } from './styles';

// import MultiModelSection from '../MultiModelSection';

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
  const [loading1, setLoading1] = useState<boolean>(true);
  const [loading2, setLoading2] = useState<boolean>(true);

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
          setLoading1(true);
          setLoading2(true);

          // 获取工况一的结果
          const result1 = (await getStructuralResult(
            queryParams.condition1,
          )) as {
            data: { data: SimulationResult };
          };

          // 检查是否已被取消
          if (isCancelled) return;
          setSimulationResult1(result1.data.data);

          // 获取工况二的结果
          const result2 = (await getStructuralResult(
            queryParams.condition2,
          )) as {
            data: { data: SimulationResult };
          };

          // 检查是否已被取消
          if (isCancelled) return;
          setSimulationResult2(result2.data.data);

          setActiveCard('位移场');
          setActiveButton('');
        } catch (error) {
          console.error('获取仿真结果失败:', error);
        } finally {
          if (!isCancelled) {
            setLoading1(false);
            setLoading2(false);
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
    camera.position.set(0, 0, 15);

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
  }, [simulationResult1, activeCard, queryParams]);

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
    camera.position.set(0, 0, 15);

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
  }, [simulationResult2, activeCard, queryParams]);

  // 当模型配置变化时加载新模型 - 工况一
  const loadModel1 = useCallback(async () => {
    if (
      !container1Ref.current ||
      !simulationResult1 ||
      !resourceManager1Ref.current
    ) {
      return;
    }
    setLoading1(true);
    await modelLoader1Ref.current.loadModel(
      simulationResult1,
      activeCard,
      activeButton,
    );
    setLoading1(false);
  }, [activeCard, activeButton, simulationResult1]);

  // 当模型配置变化时加载新模型 - 工况二
  const loadModel2 = useCallback(async () => {
    if (
      !container2Ref.current ||
      !simulationResult2 ||
      !resourceManager2Ref.current
    ) {
      return;
    }
    setLoading2(true);
    await modelLoader2Ref.current.loadModel(
      simulationResult2,
      activeCard,
      activeButton,
    );
    setLoading2(false);
  }, [activeCard, activeButton, simulationResult2]);

  useEffect(() => {
    if (container1Ref.current && simulationResult1) {
      loadModel1();
    }
  }, [activeCard, activeButton, simulationResult1, loadModel1]);

  useEffect(() => {
    if (container2Ref.current && simulationResult2) {
      loadModel2();
    }
  }, [activeCard, activeButton, simulationResult2, loadModel2]);

  // 复位相机到初始状态 - 工况一
  const resetCamera1 = useCallback(() => {
    if (container1Ref.current && initialCameraState1Ref.current) {
      const { camera, controls } = container1Ref.current;
      const initialState = initialCameraState1Ref.current;

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

  // 复位相机到初始状态 - 工况二
  const resetCamera2 = useCallback(() => {
    if (container2Ref.current && initialCameraState2Ref.current) {
      const { camera, controls } = container2Ref.current;
      const initialState = initialCameraState2Ref.current;

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
  const handleButtonClick = (buttonKey: string) => {
    if (buttonKey === '复位') {
      resetCamera1();
      resetCamera2();
    } else {
      setActiveButton(activeButton === buttonKey ? '' : buttonKey);
    }
  };

  const buttons = getButtonsByCard();

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

      <div className="model-container">
        {/* 工况一模型 */}
        <div className="model-item">
          <div className="model-item-title">工况一</div>
          <div className="model-item-content">
            <div id={domId1} className="model-container-inner">
              {loading1 && (
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
            <div id={domId2} className="model-container-inner">
              {loading2 && (
                <div className="loading-overlay">
                  <div className="spinner"></div>
                  <div>模型加载中...</div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="info-cards">
        {/* 工况一表单 */}
        <div className="info-card-container">
          {/* 位移场-工况一 */}
          <div
            className={`info-card ${activeCard === '位移场' ? 'active' : ''}`}
            onClick={() => setActiveCard('位移场')}
          >
            <div className="card-title">
              <span>位移场-工况一</span>
              <img src="/images/icons/full.png" alt="位移场" />
            </div>
            <div className="card-content">
              <div className="info-item">
                <span className="label">最大位移</span>
                <span className="value">
                  {simulationResult1?.deplace?.max
                    ? `${simulationResult1.deplace.max.toFixed(4)} mm`
                    : '--'}
                </span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">
                  {simulationResult1?.deplace?.max_displacement_location ||
                    '--'}
                </span>
              </div>
            </div>
          </div>
          {/* 位移场-工况二 */}
          <div
            className={`info-card ${activeCard === '位移场' ? 'active' : ''}`}
            onClick={() => setActiveCard('位移场')}
          >
            <div className="card-title">
              <span>位移场-工况二</span>
              <img src="/images/icons/full.png" alt="位移场" />
            </div>
            <div className="card-content">
              <div className="info-item">
                <span className="label">最大位移</span>
                <span className="value">
                  {simulationResult2?.deplace?.max
                    ? `${simulationResult2.deplace.max.toFixed(4)} mm`
                    : '--'}
                </span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">
                  {simulationResult2?.deplace?.max_displacement_location ||
                    '--'}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* 工况二表单 */}
        <div className="info-card-container">
          {/* 应力场-工况一 */}
          <div
            className={`info-card ${activeCard === '应力场' ? 'active' : ''}`}
            onClick={() => setActiveCard('应力场')}
          >
            <div className="card-title">
              <span>应力场-工况一</span>
              <img src="/images/icons/full.png" alt="应力场" />
            </div>
            <div className="card-content">
              <div className="info-item">
                <span className="label">最大应力</span>
                <span className="value">
                  {simulationResult1?.contrainte?.max
                    ? `${simulationResult1.contrainte.max.toFixed(2)} MPa`
                    : '--'}
                </span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">
                  {simulationResult1?.contrainte?.max_stress_location || '--'}
                </span>
              </div>
            </div>
          </div>
          {/* 应力场-工况二 */}
          <div
            className={`info-card ${activeCard === '应力场' ? 'active' : ''}`}
            onClick={() => setActiveCard('应力场')}
          >
            <div className="card-title">
              <span>应力场-工况二</span>
              <img src="/images/icons/full.png" alt="应力场" />
            </div>
            <div className="card-content">
              <div className="info-item">
                <span className="label">最大应力</span>
                <span className="value">
                  {simulationResult2?.contrainte?.max
                    ? `${simulationResult2.contrainte.max.toFixed(2)} MPa`
                    : '--'}
                </span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">
                  {simulationResult2?.contrainte?.max_stress_location || '--'}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <VirtualProbe
        container={container1Ref.current}
        resourceManager={resourceManager1Ref.current}
        activeCard={activeCard}
        activeButton={activeButton}
      />
      <VirtualProbe
        container={container2Ref.current}
        resourceManager={resourceManager2Ref.current}
        activeCard={activeCard}
        activeButton={activeButton}
      />
    </MultipleModelWrapper>
  );
};

export default MultiModelInfo;
