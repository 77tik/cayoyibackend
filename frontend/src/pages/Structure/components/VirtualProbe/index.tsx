import { useEffect, useRef, useState } from 'react';
import { core, ModelContainer } from 'ys-dte';
import { ResourceManager } from '../../ResourceManager';
import { ProbePopupWrapper } from './styles';

interface VirtualProbeProps {
  container: ModelContainer | null;
  resourceManager: ResourceManager | null;
  activeCard: string; // 当前激活的卡片，如 '速度场'
  activeButton: string; // 当前激活的按钮，如 '几何'、'剖面'、'流线'
}

interface ProbeInfo {
  visible: boolean;
  x: number;
  y: number;
  text: string;
}

function VirtualProbe({
  container,
  resourceManager,
  activeCard,
  activeButton,
}: VirtualProbeProps): JSX.Element {
  const [probeInfo, setProbeInfo] = useState<ProbeInfo>({
    visible: false,
    x: 0,
    y: 0,
    text: '',
  });
  const raycaster = useRef(new core.Raycaster());
  const mouse = useRef(new core.Vector2());
  const popupRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!container || !container.renderer || !container.renderer.domElement)
      return;

    const domElement = container.renderer.domElement;

    const handleMouseClick = (event: MouseEvent) => {
      if (!resourceManager || !container || !container.camera) {
        setProbeInfo({ ...probeInfo, visible: false });
        return;
      }

      const rect = domElement.getBoundingClientRect();
      mouse.current.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
      mouse.current.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;

      raycaster.current.setFromCamera(mouse.current, container.camera);

      const meshes = resourceManager.getAllResources();
      const intersects = raycaster.current.intersectObjects(meshes, false);

      if (intersects.length > 0) {
        const intersection = intersects[0];
        const point = intersection.point;
        const face = intersection.face;
        const object = intersection.object as core.Mesh;

        if (!face || !object.geometry) {
          setProbeInfo({ ...probeInfo, visible: false });
          return;
        }

        const geometry = object.geometry as core.BufferGeometry;

        // Find the closest vertex to the intersection point
        const vertices = [face.a, face.b, face.c];
        let closestVertexIndex = -1;
        let minDistance = Infinity;

        const positionAttribute = geometry.getAttribute(
          'position',
        ) as core.BufferAttribute;

        for (const vertexIndex of vertices) {
          const vertexPosition = new core.Vector3().fromBufferAttribute(
            positionAttribute,
            vertexIndex,
          );
          const distance = point.distanceTo(vertexPosition);
          if (distance < minDistance) {
            minDistance = distance;
            closestVertexIndex = vertexIndex;
          }
        }

        if (closestVertexIndex === -1) {
          setProbeInfo({ ...probeInfo, visible: false });
          return;
        }

        let probeText = '';
        let attrName = '';
        let unit = '';

        switch (activeCard) {
          case '速度场':
            attrName = 'velocity';
            unit = 'm/s';
            break;
          case '压力场':
            attrName = 'total_pressure';
            unit = 'Pa';
            break;
          case '空化分布':
            attrName = 'phase_1_vof';
            unit = '';
            break;
          case '涡带分布':
            attrName = 'total_pressure'; // Or other relevant attribute
            // attrName = 'raw_q_criterion'; // Or other relevant attribute
            unit = 'Pa';
            break;
          case '位移场':
            attrName = 'resu____DEPL';
            unit = 'mm';
            break;
          case '应力场':
            attrName = 'resu____SIEQ_NOEU';
            unit = 'MPa';
            break;
          default:
            setProbeInfo({ ...probeInfo, visible: false });
            return;
        }

        const attribute = geometry.getAttribute(
          attrName,
        ) as core.BufferAttribute;

        if (attribute) {
          const value = new core.Vector3().fromBufferAttribute(
            attribute,
            closestVertexIndex,
          );
          const displayValue =
            attribute.itemSize === 1 ? value.x : value.length();
          let formattedValue = displayValue.toFixed(2);
          if (activeCard === '位移场') {
            formattedValue = displayValue.toFixed(6);
          } else if (activeCard === '应力场') {
            formattedValue = (displayValue / 1000000).toFixed(2);
          }
          probeText = `${activeCard
            .replace('分布', '')
            .replace('场', '')}: ${formattedValue} ${unit}`;
        } else {
          setProbeInfo({ ...probeInfo, visible: false });
          return;
        }

        // Convert 3D point to 2D screen coordinates
        const screenPoint = point.clone().project(container.camera);
        const screenX = ((screenPoint.x + 1) / 2) * rect.width + rect.left;
        const screenY = ((-screenPoint.y + 1) / 2) * rect.height + rect.top;

        setProbeInfo({
          visible: true,
          x: screenX,
          y: screenY,
          text: probeText,
        });
      } else {
        setProbeInfo({ ...probeInfo, visible: false });
      }
    };

    domElement.addEventListener('click', handleMouseClick);

    return () => {
      domElement.removeEventListener('click', handleMouseClick);
    };
  }, [container, resourceManager, activeCard, activeButton]);

  useEffect(() => {
    setProbeInfo((info) => ({ ...info, visible: false }));
  }, [container, resourceManager, activeCard, activeButton]);

  return (
    // @ts-ignore
    <ProbePopupWrapper
      ref={popupRef}
      className={probeInfo.visible ? 'visible' : ''}
      style={{
        left: `${probeInfo.x}px`,
        top: `${probeInfo.y}px`,
      }}
    >
      {probeInfo.text}
    </ProbePopupWrapper>
  );
}

export default VirtualProbe;
