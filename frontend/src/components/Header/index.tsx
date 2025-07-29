import React from 'react';
import { history } from 'umi';
import { HeaderWrapper } from './style';

interface HeaderProps {
  selectedTitle: string;
  onUnitChange: (value: string) => void;
}

const Header: React.FC<HeaderProps> = ({ selectedTitle, onUnitChange }) => {
  const handleNavigation = (path: string, unit: string) => {
    onUnitChange(unit);
    history.push(path);
  };

  return (
    <HeaderWrapper>
      <nav className="navbar">
        <div className="pageTitle">
          <span>水轮机流体仿真智能同数系统</span>
        </div>
        <ul>
          <li
            onClick={() => handleNavigation('/fluid', 'fluidSimulation')}
            className={selectedTitle === 'fluidSimulation' ? 'selected' : ''}
          >
            <a href="#">流体仿真模拟</a>
          </li>
          <li
            onClick={() =>
              handleNavigation('/structure', 'structureSimulation')
            }
            className={
              selectedTitle === 'structureSimulation' ? 'selected' : ''
            }
          >
            <a href="#">结构仿真模拟</a>
          </li>
          <li
            onClick={() => handleNavigation('/safety', 'safetyWaterHead')}
            className={selectedTitle === 'safetyWaterHead' ? 'selected' : ''}
          >
            <a href="#">安全水头区域</a>
          </li>
          <li
            onClick={() => handleNavigation('/monitor', 'monitorDataDisplay')}
            className={selectedTitle === 'monitorDataDisplay' ? 'selected' : ''}
          >
            <a href="#">应变数据展示</a>
          </li>
        </ul>
        {/* <button className="headerBtn" type="button" /> */}
      </nav>
    </HeaderWrapper>
  );
};

export default Header;
