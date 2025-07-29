import Header from '@/components/Header';
import React, { useState } from 'react';
import { Wrapper } from './styles';

// 安全水头区域
const SafetyPage: React.FC = () => {
  const [selectedTitle, setSelectedTitle] = useState('safetyWaterHead');

  const handleDownload = () => {
    const link = document.createElement('a');
    link.href = '/pdf/test.pdf';
    link.download = '研究报告.pdf';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <Wrapper>
      <div className="container">
        <Header selectedTitle={selectedTitle} onUnitChange={setSelectedTitle} />

        <div className="content">
          <div className="body">
            {/* 机组运转特性曲线 */}
            <div className="left">
              <div className="safeTitle"> 机组运转特性曲线</div>
              <div className="safeText">
                水头（m）峰值效率=96.22%，P=227.26MW
              </div>
              <div className="leftChart">图表1</div>
              <div className="chartResult">
                {/* 左边 */}
                <div className="chartResultLeftItem">
                  <div className="chartResultLeftItemIcon"></div>
                  <div className="chartResultLeftItemText">
                    综合分析泸定电站水轮机空化性能与转轮结构应力的变化趋势，确定
                    减小水轮机内部流场空化效应和结构应力集中现象的低有效水头安全
                    区间为：60.2~66.2m。
                  </div>
                </div>
                {/* 右边 */}
                <div className="chartResultRightItem">
                  <div className="chartResultRightItemContent">60.2~66.2m</div>
                  <div className="chartResultRightItemTitle">
                    低有效水头安全区间
                  </div>
                </div>
              </div>
            </div>
            {/* 研究报告 */}
            <div className="right">
              <div className="safeTitle">研究报告</div>
              <div className="downloadBtn">
                {/* 下载 */}
                <span className="downloadText" onClick={handleDownload}>
                  下载
                </span>
              </div>
              <div className="report">
                <div className="reportContent">
                  {/* 泸定电站水轮机空化性能与转轮结构应力分析报 */}
                  <iframe src="/pdf/test.pdf#toolbar=0&view=FitH" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Wrapper>
  );
};

export default SafetyPage;
