// 所有模型加载相关逻辑
import alovaInstance from '@/api/alovaInstance';
import { core, JsonLoader, ShaderMesh } from 'ys-dte';
import { ResourceManager } from '../../ResourceManager';
import { ProfileType, SimulationResult } from '../../types';

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
    // 几何模型使用半透明材质
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
    profileType: ProfileType,
  ) {
    // 根据不同卡片和按钮返回对应的 URL
    switch (activeCard) {
      case '速度场':
        if (activeButton === '几何') {
          return simulationResult.data?.mesh_json;
        } else if (activeButton === '剖面') {
          return profileType === 'h'
            ? simulationResult.data?.velocity?.h
            : simulationResult.data?.velocity?.v;
        } else if (activeButton === '流线') {
          return ''; // 流线模式特殊处理，返回空
        } else if (activeButton === '') {
          // 默认状态显示速度场仿真云图（使用mesh_json）
          return simulationResult.data?.mesh_json;
        }
        break;
      case '压力场':
        if (activeButton === '几何') {
          return simulationResult.data?.mesh_json;
        } else if (activeButton === '剖面') {
          return profileType === 'h'
            ? simulationResult.data?.pressure?.h
            : simulationResult.data?.pressure?.v;
        } else if (activeButton === '') {
          // 默认状态显示压力场仿真云图（使用mesh_json）
          return simulationResult.data?.mesh_json;
        }
        break;
      case '空化分布':
        if (activeButton === '几何') {
          //   return simulationResult.data?.vof?.mesh_json;
          // } else if (activeButton === '剖面') {
          return ''; // 空化剖面模式特殊处理，返回空
        } else if (activeButton === '') {
          // 默认状态显示空化分布仿真云图
          return ''; // 空化默认也是特殊处理
        }
        break;
      case '涡带分布':
        if (activeButton === '几何') {
          //   return simulationResult.data?.vof?.mesh_json;
          // } else if (activeButton === '剖面') {
          return ''; // 涡带剖面模式特殊处理，返回空
        } else if (activeButton === '') {
          // 默认状态显示涡带分布仿真云图
          return ''; // 涡带默认也是特殊处理
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
    profileType: ProfileType,
  ) {
    if (!modelUrl || !this.resourceManager) return;

    const fullUrl = `${(alovaInstance as any).options.baseURL}${modelUrl}`;

    try {
      const meshData = await new Promise<any>((resolve, reject) => {
        this.loader.loadDataMesh(fullUrl, resolve, reject);
      });

      console.log('meshData', meshData);

      meshData.rotateX(-Math.PI / 2);

      let mesh: any;

      // 根据卡片类型确定要使用的属性
      let attrName = '';
      switch (activeCard) {
        case '速度场':
          attrName = 'velocity';
          break;
        case '压力场':
          attrName = 'total_pressure';
          break;
        case '空化分布':
          attrName = 'phase_1_vof';
          break;
        case '涡带分布':
          attrName = 'total_pressure';
          break;
        default:
          attrName = 'velocity';
      }

      const attrInfo = meshData.userData?.[attrName];
      const geometryAttr = meshData.attributes?.[attrName];

      // 判断是否为几何模式（只有点击几何按钮才显示纯几何模型）
      if (activeButton === '几何') {
        const attrName = '_YAMI';
        mesh = new ShaderMesh(meshData, attrName, {
          materialType: 'MeshLambertMaterial' as any,
          colorList: [0x888888],
        });
      } else if (attrInfo && geometryAttr) {
        // 仿真云图模式：使用着色器材质显示数据属性
        const minValue = attrInfo.MinValue ?? attrInfo.minValue;
        const maxValue = attrInfo.MaxValue ?? attrInfo.maxValue;

        mesh = new ShaderMesh(meshData, attrName, {
          minValue,
          maxValue,
          materialType: 'MeshBasicMaterial' as any,
        });
      } else {
        console.warn(`⚠️ 属性 ${attrName} 不存在于模型数据中，使用基础材质`);
        // 如果没有对应属性，则使用半透明材质作为fallback
        mesh = new core.Mesh(meshData, this.getMaterial('geometry'));
      }

      // 设置模型的位置、旋转和缩放
      mesh.position.set(0, 0, 0);

      // 为压力场和速度场的剖面模式设置不同的旋转角度
      if (activeCard === '速度场' || activeCard === '压力场') {
        if (activeButton === '剖面') {
          if (profileType === 'h') {
            // 速度场横剖面
            mesh.rotation.set(Math.PI / 2, Math.PI / 2, 0);
          } else {
            // 速度场纵剖面
            mesh.rotation.set(0, 0, 0);
          }
        } else {
          // 其他模式使用默认旋转
          mesh.rotation.set(0, Math.PI / 2, 0);
        }
      } else {
        // 非剖面模式使用默认旋转
        mesh.rotation.set(0, Math.PI / 2, 0);
      }

      mesh.scale.set(1, 1, 1);
      this.resourceManager.addResource(mesh);
    } catch (error) {
      console.error('模型加载失败:', error);
    }
  }

  // 加载速度场流线模式
  private async loadStreamlineMode(simulationResult: SimulationResult) {
    if (!this.resourceManager) return;

    // 1. 加载基础几何模型（半透明）
    const geometryUrl = simulationResult.data?.mesh_json;
    if (geometryUrl) {
      const fullGeometryUrl = `${
        (alovaInstance as any).options.baseURL
      }${geometryUrl}`;
      try {
        const geometryMeshData = await new Promise<any>((resolve, reject) => {
          this.loader.loadDataMesh(fullGeometryUrl, resolve, reject);
        });
        geometryMeshData.rotateX(-Math.PI / 2);

        const geometryMesh = new core.Mesh(
          geometryMeshData,
          this.getMaterial('geometry'),
        );
        geometryMesh.position.set(-9, 6, 0);
        geometryMesh.rotation.set(0, Math.PI / 2, 0);
        geometryMesh.scale.set(1, 1, 1);
        this.resourceManager.addResource(geometryMesh);
      } catch (error) {
        console.error('几何模型加载失败', error);
      }
    }

    // 2. 加载流线模型（使用自身数据）
    const streamlineUrl = simulationResult.data?.velocity?.stream_line;
    if (streamlineUrl) {
      const streamFullUrl = `${
        (alovaInstance as any).options.baseURL
      }${streamlineUrl}`;
      try {
        const streamMeshData = await new Promise<any>((resolve, reject) => {
          this.loader.loadDataMesh(streamFullUrl, resolve, reject);
        });
        streamMeshData.rotateX(-Math.PI / 2);

        // ————————————————————
        const attrName = 'velocity';
        const attrInfo = streamMeshData.userData?.[attrName];
        const geometryAttr = streamMeshData.attributes?.[attrName];

        let streamlineMesh: any;

        if (attrInfo && geometryAttr) {
          /***************** 关键修改开始 *****************/
          // 使用模型自带的颜色属性（如果存在）
          const hasColorAttribute = !!streamMeshData.attributes.color;

          // 创建线材质
          const lineMaterial = new core.LineBasicMaterial({
            vertexColors: hasColorAttribute, // 如果模型有颜色属性则启用
            // linewidth: 1,
          });

          // 创建线模型（保持模型原有颜色）
          streamlineMesh = new core.LineSegments(streamMeshData, lineMaterial);
          /***************** 关键修改结束 *****************/
        } else {
          /***************** 关键修改开始 *****************/
          // 没有速度属性时使用红色线模型
          streamlineMesh = new core.LineSegments(
            streamMeshData,
            new core.LineBasicMaterial({ color: 0xff0000 }),
          );
          /***************** 关键修改结束 *****************/
        }

        // ————————————————————
        streamlineMesh.position.set(-9, 6, 0);
        streamlineMesh.rotation.set(0, Math.PI / 2, 0);
        streamlineMesh.scale.set(1, 1, 1);
        this.resourceManager.addResource(streamlineMesh);
      } catch (error) {
        console.error('流线模型加载失败', error);
      }
    }
  }

  // 加载空化分布模式
  private async loadCavitationProfileMode(
    simulationResult: SimulationResult,
    onlyGeometry: boolean = false,
  ) {
    if (!this.resourceManager) return;

    // 1. 加载基础几何模型（半透明）
    const geometryUrl = simulationResult.data?.vof?.mesh_json;
    if (geometryUrl) {
      const fullGeometryUrl = `${
        (alovaInstance as any).options.baseURL
      }${geometryUrl}`;
      try {
        const geometryMeshData = await new Promise<any>((resolve, reject) => {
          this.loader.loadDataMesh(fullGeometryUrl, resolve, reject);
        });
        geometryMeshData.rotateX(-Math.PI / 2);

        // 计算几何模型的边界框和中心点
        geometryMeshData.computeBoundingBox();

        let geometryMesh: any;
        if (onlyGeometry) {
          // 几何模型材质为MeshLambertMaterial
          const attrName = '_YAMI';
          geometryMesh = new ShaderMesh(geometryMeshData, attrName, {
            materialType: 'MeshLambertMaterial' as any,
            center: false,
            colorList: [0x888888],
          });
        } else {
          geometryMesh = new core.Mesh(
            geometryMeshData,
            this.getMaterial('geometry'),
          );
        }

        geometryMesh.position.set(0, 4.5, 0);
        geometryMesh.scale.set(5, 5, 5);
        this.resourceManager.addResource(geometryMesh);

        if (onlyGeometry) return;

        // 2. 加载空化分布数据模型
        const vofUrl = simulationResult.data?.vof?.vof;
        if (vofUrl) {
          const vofFullUrl = `${
            (alovaInstance as any).options.baseURL
          }${vofUrl}`;
          const vofMeshData = await new Promise<any>((resolve, reject) => {
            this.loader.loadDataMesh(vofFullUrl, resolve, reject);
          });
          vofMeshData.rotateX(-Math.PI / 2);

          const attrName = 'phase_1_vof';
          const attrInfo = vofMeshData.userData?.[attrName];
          const geometryAttr = vofMeshData.attributes?.[attrName];

          let vofMesh: any;

          if (attrInfo && geometryAttr) {
            const minValue = attrInfo.MinValue ?? attrInfo.minValue;
            const maxValue = attrInfo.MaxValue ?? attrInfo.maxValue;

            vofMesh = new ShaderMesh(vofMeshData, attrName, {
              minValue,
              maxValue,
              materialType: 'MeshBasicMaterial' as any,
              center: false,
            });
          } else {
            console.warn('空化模型没有phase_1_vof属性，使用基础材质');
            vofMesh = new core.Mesh(
              vofMeshData,
              new core.MeshBasicMaterial({
                color: 0xff0000,
                side: core.DoubleSide,
              }),
            );
          }

          vofMesh.position.y -= -4.3;
          vofMesh.scale.set(5, 5, 5);
          this.resourceManager.addResource(vofMesh);
        }
      } catch (error) {
        console.error('空化模型加载失败:', error);
      }
    }
  }

  // 加载涡带分布模式
  private async loadVortexProfileMode(
    simulationResult: SimulationResult,
    onlyGeometry: boolean = false,
  ) {
    if (!this.resourceManager) return;

    // 1. 加载基础几何模型（半透明）
    const geometryUrl = simulationResult.data?.mesh_json;
    if (geometryUrl) {
      const fullGeometryUrl = `${
        (alovaInstance as any).options.baseURL
      }${geometryUrl}`;
      try {
        const geometryMeshData = await new Promise<any>((resolve, reject) => {
          this.loader.loadDataMesh(fullGeometryUrl, resolve, reject);
        });
        geometryMeshData.rotateX(-Math.PI / 2);

        let geometryMesh: any;
        if (onlyGeometry) {
          // 几何模型材质为MeshLambertMaterial
          const attrName = '_YAMI';
          geometryMesh = new ShaderMesh(geometryMeshData, attrName, {
            materialType: 'MeshLambertMaterial' as any,
            center: false,
            colorList: [0x888888],
          });
        } else {
          geometryMesh = new core.Mesh(
            geometryMeshData,
            this.getMaterial('geometry'),
          );
        }

        geometryMesh.rotation.set(0, Math.PI / 2, 0);

        geometryMesh.position.set(-9, 6, 0);

        this.resourceManager.addResource(geometryMesh);
      } catch (error) {
        console.error('涡带几何模型加载失败', error);
      }
    }

    if (onlyGeometry) return;

    // 2. 加载涡带分布数据模型
    const vortexUrl = simulationResult.data?.vortex?.vortex;
    if (vortexUrl) {
      const vortexFullUrl = `${
        (alovaInstance as any).options.baseURL
      }${vortexUrl}`;
      try {
        const vortexMeshData = await new Promise<any>((resolve, reject) => {
          this.loader.loadDataMesh(vortexFullUrl, resolve, reject);
        });
        vortexMeshData.rotateX(-Math.PI / 2);

        const attrName = 'total_pressure';
        const attrInfo = vortexMeshData.userData?.[attrName];
        const geometryAttr = vortexMeshData.attributes?.[attrName];

        let vortexMesh: any;

        if (attrInfo && geometryAttr) {
          const minValue = attrInfo.MinValue ?? attrInfo.minValue;
          const maxValue = attrInfo.MaxValue ?? attrInfo.maxValue;

          vortexMesh = new ShaderMesh(vortexMeshData, attrName, {
            minValue,
            maxValue,
            materialType: 'MeshBasicMaterial' as any,
            center: false,
          });
        } else {
          console.warn('涡带模型没有total_pressure属性，使用基础材质');
          vortexMesh = new core.Mesh(
            vortexMeshData,
            new core.MeshBasicMaterial({
              color: 0xff0000,
              side: core.DoubleSide,
            }),
          );
        }

        vortexMesh.rotation.set(0, Math.PI / 2, 0);

        vortexMesh.position.set(-9, 6, 0);

        this.resourceManager.addResource(vortexMesh);
      } catch (error) {
        console.error('涡带分布模型加载失败', error);
      }
    }
  }

  // 主加载函数
  async loadModel(
    simulationResult: SimulationResult,
    activeCard: string,
    activeButton: string,
    profileType: ProfileType,
  ) {
    if (!this.resourceManager) return;

    // 清理所有现有资源
    this.resourceManager.clearAll();

    // 速度场流线模式需要特殊处理 - 加载两个模型
    if (activeCard === '速度场' && activeButton === '流线') {
      await this.loadStreamlineMode(simulationResult);
    }
    // 空化分布剖面模式需要特殊处理 - 加载两个模型
    else if (activeCard === '空化分布') {
      const onlyGeometry = activeButton === '几何';
      await this.loadCavitationProfileMode(simulationResult, onlyGeometry);
    }
    // 涡带分布剖面模式需要特殊处理 - 加载两个模型
    else if (activeCard === '涡带分布') {
      const onlyGeometry = activeButton === '几何';
      await this.loadVortexProfileMode(simulationResult, onlyGeometry);
    }
    // 其他模式加载单个模型
    else {
      const modelUrl = this.getModelUrl(
        simulationResult,
        activeCard,
        activeButton,
        profileType,
      );
      if (modelUrl) {
        await this.loadSingleModel(
          modelUrl,
          activeCard,
          activeButton,
          profileType,
        );
      }
    }
  }
}
