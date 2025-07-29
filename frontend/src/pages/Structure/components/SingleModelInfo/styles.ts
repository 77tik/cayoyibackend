import styled from 'styled-components';

// 单工况的信息显示组件样式
export const SingleModelWrapper = styled.div`
  display: flex;
  flex-direction: column;
  /* width: 100%; */
  height: 100%;
  margin-left: 60px;
  position: relative; // 添加相对定位，作为浮动元素的参考

  .model-container {
    margin-top: 82px;
    margin-bottom: 66px;
    height: 640px;
    position: relative;
    width: 100%;

    display: flex; /* 启用 Flex 布局 */
    gap: 20px;

    /* .speed-tag {
      position: absolute;
      padding: 4px 8px;
      background: rgba(0, 0, 0, 0.6);
      border-radius: 4px;
      color: #fff;
      font-size: 14px;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
    } */

    .model-container-inner {
      width: 100%;
      height: 100%;
      position: relative;
    }

    .loading-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(255, 255, 255, 0.8);
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      z-index: 1000;

      .spinner {
        width: 40px;
        height: 40px;
        border: 4px solid #f3f3f3;
        border-top: 4px solid #528eff;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 16px;
      }

      @keyframes spin {
        0% {
          transform: rotate(0deg);
        }
        100% {
          transform: rotate(360deg);
        }
      }
    }
  }

  /* 卡片样式 */
  .info-cards {
    gap: 56px;
    display: flex;
    justify-content: center; // 水平居中

    .info-card {
      height: 118px;
      width: 288px;
      background: #fff;
      border-radius: 6px 6px 0px 0px;

      .card-title {
        height: 32px;
        line-height: 32px;
        background: rgba(108, 111, 119, 0.4);
        border-radius: 6px 6px 0px 0px;
        padding: 0 16px;

        font-family: JiangChengXieHei-700W;
        font-size: 16px;
        color: #ffffff;
        letter-spacing: 0.8px;
        font-weight: 700;
        position: relative;
        img {
          position: absolute;
          right: 16px;
          top: 50%;
          transform: translateY(-50%);
        }
      }

      &.active {
        .card-title {
          background-image: linear-gradient(270deg, #528eff, #4b6ce6);
        }

        .value {
          color: #528eff !important;
        }
      }

      .card-content {
        padding: 0 16px;
        font-size: 16px;

        .info-item {
          line-height: 42px;
          display: flex;
          justify-content: space-between;
          border-bottom: 1px solid #d0d6e3; /* 横线样式 */

          &:last-child {
            margin-bottom: 0;
            border-bottom: none;
          }
        }
      }
    }
  }
`;
