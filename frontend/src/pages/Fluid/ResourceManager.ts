// 资源管理器类
export class ResourceManager {
  private container: any;
  private resources: any[] = [];

  constructor(container: any) {
    this.container = container;
  }

  addResource(resource: any) {
    this.resources.push(resource);
    this.container.scene.add(resource);
  }

  getAllResources() {
    return this.resources;
  }

  clearAll() {
    this.resources.forEach((resource) => {
      this.container.scene.remove(resource);
      if (resource.geometry) resource.geometry.dispose();
      if (resource.material) {
        if (Array.isArray(resource.material)) {
          resource.material.forEach((mat: any) => mat.dispose());
        } else {
          resource.material.dispose();
        }
      }
    });
    this.resources = [];
  }

  dispose() {
    this.clearAll();
    this.container = null;
  }
}
