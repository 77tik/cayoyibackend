import styled from 'styled-components';

export const ProbePopupWrapper = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  padding: 8px 12px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  border-radius: 4px;
  font-size: 14px;
  pointer-events: none; /* 避免遮挡鼠标事件 */
  transform: translate(-420%, -330%); /* 调整位置，使其在点的上方 */
  white-space: nowrap;
  z-index: 100;
  visibility: hidden; /* 默认隐藏 */
  opacity: 0.6;

  &.visible {
    visibility: visible;
  }
`;
