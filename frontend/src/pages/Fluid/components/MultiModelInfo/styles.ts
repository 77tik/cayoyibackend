import styled from 'styled-components';
// 多工况的信息显示组件样式
export const MultipleModelWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  margin-left: 60px;
  position: relative; // 添加相对定位，作为浮动元素的参考

  // 添加剖面面板的浮动样式
  .profile-panel {
    position: absolute;
    top: 50px;
    right: 0px;
    z-index: 1000;
    background: #f7faff;
    border-radius: 6px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.08);
    width: 392px;
    height: 240px;
    border: 1px solid #e3e8f0;
    display: flex;
    flex-direction: column;
    padding: 0;
  }

  .profile-panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    height: 36px;
    border-bottom: 1px solid #e3e8f0;
    background: transparent;
  }

  .profile-tabs {
    display: flex;
    align-items: center;
    height: 100%;
  }

  .profile-tab {
    font-size: 16px;
    color: #222;
    padding: 0 16px;
    height: 48px;
    line-height: 48px;
    cursor: pointer;
    position: relative;
    background: none;
    border: none;
    outline: none;
    transition: color 0.2s;
  }
  .profile-tab.active {
    color: #528eff;
  }
  .profile-tab.active::after {
    content: '';
    position: absolute;
    left: 8px;
    right: 8px;
    bottom: 6px;
    height: 3px;
    border-radius: 2px;
    background: #528eff;
  }

  .profile-close {
    font-size: 22px;
    color: #999;
    cursor: pointer;
    margin-left: 12px;
    transition: color 0.2s;
    user-select: none;
  }
  .profile-close:hover {
    color: #528eff;
  }

  .profile-panel-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 120px;

    .profile-panel-content-h {
      max-width: calc(100% - 20px);
      max-height: calc(100% - 40px);
    }
    .profile-panel-content-v {
      max-width: 100%;
      max-height: 100%;
    }
  }

  .model-container {
    margin-top: 82px;
    margin-bottom: 32px;
    display: flex; /* 启用 Flex 布局 */
    gap: 20px;
    height: 536px;
  }

  .model-item {
    background-color: #fff;
    width: 650px;
    height: 100%;
    display: flex; /* 启用 Flex 布局 */
    flex-direction: column;
    border-radius: 6px;

    .model-item-title {
      height: 32px;
      line-height: 32px;
      margin-left: 16px;
      font-family: JiangChengXieHei-700W;
      font-size: 16px;
      color: #2c2d31;
      letter-spacing: 0.8px;
      font-weight: 700;
    }
    .model-item-content {
      height: 100%;
      /* background: rgba(108, 111, 119, 0.4); */
      background: rgba(82, 142, 255, 0.3);

      margin: 0 1px 1px 1px;
      border-radius: 0 0 6px 6px;
      position: relative;

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
        background: rgb(204, 204, 204);

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        z-index: 999;

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
  }

  /* 信息表单 */
  .info-cards {
    display: flex;
    gap: 56px;
    margin-bottom: 20px;
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
  .info-cards:last-child {
    margin-bottom: 0;
  }
`;
