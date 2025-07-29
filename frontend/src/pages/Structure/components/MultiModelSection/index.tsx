import React, { useState } from 'react';
import { SectionViewWrapper } from './styles';
// 剖面视图组件
interface SectionViewProps {
  activeCard: string;
  setActiveCard: (card: string) => void;
}

const SectionView: React.FC<SectionViewProps> = ({
  activeCard,
  setActiveCard,
}) => {
  // 生成12个工况的数组
  const conditions = Array.from({ length: 13 }, (_, index) => ({
    id: index + 1,
    title: `${index + 1}#叶片`,
  }));

  // 添加全屏状态
  const [fullScreenItem, setFullScreenItem] = useState<number | null>(null);

  // 处理全屏切换
  const handleFullScreen = (id: number, e: React.MouseEvent) => {
    e.stopPropagation(); // 阻止事件冒泡
    setFullScreenItem(fullScreenItem === id ? null : id);
  };

  return (
    <SectionViewWrapper>
      <div className="section-container">
        {fullScreenItem !== null && (
          <div className="fullscreen-overlay">
            <div className="fullscreen-content">
              <div className="fullscreen-header">
                {conditions.find((c) => c.id === fullScreenItem)?.title || ''}
                <img
                  src="/images/icons/full.png"
                  alt="退出全屏"
                  onClick={(e) => handleFullScreen(fullScreenItem, e)}
                />
              </div>
              <div className="fullscreen-body">
                {/* 这里将来会放3D模型渲染内容 */}
              </div>
            </div>
          </div>
        )}

        {conditions.map((condition) => (
          <div
            key={condition.id}
            className={`section-item ${
              fullScreenItem === condition.id ? 'hidden' : ''
            }`}
          >
            <div className="section-item-title">
              {condition.title}
              <img
                src="/images/icons/full.png"
                alt="全屏"
                onClick={(e) => handleFullScreen(condition.id, e)}
              />
            </div>
            <div className="section-item-content">
              {/* 这里将来会放3D模型渲染内容 */}
            </div>
          </div>
        ))}
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
                <span className="value">2.5 mm</span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">转轮叶片</span>
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
                <span className="value">2.5 mm</span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">转轮叶片</span>
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
                <span className="value">150 MPa</span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">转轮叶片</span>
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
                <span className="value">150 MPa</span>
              </div>
              <div className="info-item">
                <span className="label">出现部位</span>
                <span className="value">转轮叶片</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </SectionViewWrapper>
  );
};

export default SectionView;
