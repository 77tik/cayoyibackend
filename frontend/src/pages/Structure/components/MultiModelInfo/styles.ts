import styled from 'styled-components';
// 多工况的信息显示组件样式
export const MultipleModelWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  margin-left: 60px;

  .model-container {
    margin-top: 82px;
    margin-bottom: 32px;
    display: flex; /* 启用 Flex 布局  */
    gap: 20px;
    height: 674px;
  }

  .model-item {
    background-color: #fff;
    width: 650px;
    /* height: 674px; */
    /* height: 100%; */
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
      /* background-color: gray; */
      /* background: rgb(204, 204, 204); */
      background: rgba(82, 142, 255, 0.3);

      margin: 0 1px 1px 1px;

      border-radius: 0 0 6px 6px;

      .model-container-inner {
        height: 100%;
      }

      .loading-overlay {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        background: rgb(204, 204, 204);

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
    gap: 96px;
    display: flex;
    flex-direction: row;
    justify-content: center;

    .info-card-container {
      display: flex;
      flex-direction: row;
      gap: 16px;
    }

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
