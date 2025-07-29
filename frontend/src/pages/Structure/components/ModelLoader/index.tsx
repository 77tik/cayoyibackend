// 所有模型加载相关逻辑
import alovaInstance from '@/api/alovaInstance';
import { core, JsonLoader, ShaderMesh } from 'ys-dte';
import { ResourceManager } from '../../ResourceManager';
import { SimulationResult } from '../../types';

export class ModelLoader {
  private loader: JsonLoader;
  private resourceManager: ResourceManager | null = null;

  constructor() {
    this.loader = new JsonLoader();
  }

  setResourceManager(resourceManager: ResourceManager) {
    this.resourceManager = resourceManager;
  }

  // 获取材质函数 - 根据模型类型返回适当的材质
  private getMaterial(modelType: string) {
    // 几何模型使用半透明蓝色材质
    if (modelType === 'geometry') {
      return new core.MeshBasicMaterial({
        color: 0x000000,
        side: core.DoubleSide,
        wireframe: false,
        transparent: true,
        opacity: 0.1,
      });
    }
    // 默认材质
    return new core.MeshBasicMaterial({ color: 0xffffff });
  }

  // 根据 activeCard 和 activeButton 计算 modelUrl
  private getModelUrl(
    simulationResult: SimulationResult,
    activeCard: string,
    activeButton: string,
  ) {
    // 根据不同卡片和按钮返回对应的 URL
    switch (activeCard) {
      case '位移场':
        if (activeButton === '几何' || activeButton === '') {
          return simulationResult.mesh_json;
        } else if (activeButton === '剖面') {
          return simulationResult.deplace?.deplace;
        }
        break;
      case '应力场':
        if (activeButton === '几何' || activeButton === '') {
          return simulationResult.mesh_json;
        } else if (activeButton === '剖面') {
          return simulationResult.contrainte?.contrainte;
        }
        break;
    }
    return '';
  }

  // 加载单个模型
  private async loadSingleModel(
    modelUrl: string,
    activeCard: string,
    activeButton: string,
  ) {
    if (!modelUrl || !this.resourceManager) return;

    const fullUrl = `${(alovaInstance as any).options.baseURL}${modelUrl}`;

    try {
      const meshData = await new Promise<any>((resolve, reject) => {
        this.loader.loadDataMesh(fullUrl, resolve, reject);
      });

      meshData.rotateX(-Math.PI / 2);

      let mesh: any;

      // 判断是否为几何模式（只有点击几何按钮才显示纯几何模型）
      if (activeButton === '几何') {
        const attrName = '_YAMI';
        mesh = new ShaderMesh(meshData, attrName, {
          materialType: 'MeshLambertMaterial' as any,
          center: false,
          colorList: [0x888888],
        });
      } else {
        // 仿真云图模式：使用着色器材质显示数据属性
        let attrName = '';
        switch (activeCard) {
          case '位移场':
            attrName = 'resu____DEPL';
            break;
          case '应力场':
            attrName = 'resu____SIEQ_NOEU';
            break;
          default:
            attrName = 'resu____DEPL';
        }
        const attrInfo = meshData.userData?.[attrName];
        const geometryAttr = meshData.attributes?.[attrName];
        if (attrInfo && geometryAttr) {
          const minValue = attrInfo.MinValue ?? attrInfo.minValue;
          const maxValue = attrInfo.MaxValue ?? attrInfo.maxValue;
          mesh = new ShaderMesh(meshData, attrName, {
            minValue,
            maxValue,
          });
        } else {
          mesh = new core.Mesh(meshData, this.getMaterial('geometry'));
        }
      }

      // 设置模型的位置、旋转和缩放
      if (
        (activeCard === '位移场' || activeCard === '应力场') &&
        activeButton === '剖面'
      ) {
        mesh.position.set(0, 0, 0);
        mesh.rotation.set(Math.PI / 2, 0, 0);
      } else {
        mesh.position.set(0, 0, 0);
      }
      mesh.scale.set(1, 1, 1);
      this.resourceManager.addResource(mesh);
    } catch (error) {
      console.error('模型加载失败:', error);
    }
  }

  // 主加载函数
  async loadModel(
    simulationResult: SimulationResult,
    activeCard: string,
    activeButton: string,
  ) {
    console.log('🔧 ModelLoader.loadModel 开始执行', {
      simulationResult,
      activeCard,
      activeButton,
    });

    if (!this.resourceManager) {
      console.log('❌ resourceManager 为空，退出');
      return;
    }

    // 清理所有现有资源
    this.resourceManager.clearAll();

    // 结构分析只有普通的单模型加载模式
    const modelUrl = this.getModelUrl(
      simulationResult,
      activeCard,
      activeButton,
    );

    console.log('🔗 计算出的模型URL:', modelUrl);

    if (modelUrl) {
      await this.loadSingleModel(modelUrl, activeCard, activeButton);
    } else {
      console.log('❌ 模型URL为空，无法加载模型');
    }
  }
}
