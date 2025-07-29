import styled from 'styled-components';

export const LeftWrapper = styled.div`
  /* .left {
    margin-top: 24px;
    width: 420px;
  } */

  .leftPanel {
    width: 420px;
    margin-top: 14px;
    padding: 19px 16px 0 16px;
    color: #fff;
    background: rgba(255, 255, 255, 0.8);
    border: 1px solid rgba(210, 216, 229, 1);
    box-shadow: 0px 4px 6px 0px rgba(82, 142, 255, 0.2);
    border-radius: 6px;

    .panelHeader {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .panelHeaderTitle {
        width: 388px;
        height: 36px;
        background-image: linear-gradient(
          90deg,
          rgba(105, 167, 255, 0.4) 1%,
          rgba(255, 255, 255, 0.6) 100%
        );

        font-family: JiangChengXieHei-700W;
        font-size: 20px;
        color: #2c2d31;
        letter-spacing: 1px;
        font-weight: 700;
      }
    }
  }

  .queryButton {
    width: 140px;
    height: 30px;
    background: #528eff;
    border-radius: 4px;
    border: none;
    padding: 8px 16px;
    font-family: 'PingFangSC-Medium', sans-serif;
    font-size: 16px;
    color: #ffffff;
    letter-spacing: 1px;
    text-align: center;
    font-weight: 500;
    cursor: pointer;
    margin: 0 auto;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 8px rgba(82, 142, 255, 0.3);
      opacity: 0.9;
    }
  }

  /* 选择框 */
  .select-wrapper {
    width: 420px;
    height: 42px;
    margin-top: 24px;
    border-radius: 6px;
    background: linear-gradient(to right, #90a9de, #c3d4f1);
    padding: 2px; /* 渐变边框厚度 */
    display: inline-block;
  }
  .custom-select {
    width: 100%;
    height: 100%;
    border: none;
    border-radius: 6px;
    padding: 0 12px;
    background: rgba(255, 255, 255, 0.8);
    box-shadow: 0px 4px 6px 0px rgba(82, 142, 255, 0.2);
    outline: none;
    appearance: none;

    /* 自定义下拉箭头 */
    background-image: url("data:image/svg+xml,%3Csvg width='10' height='6' viewBox='0 0 10 6' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M0 0l5 6 5-6z' fill='%2390A9DE'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 12px center;
    background-size: 10px 6px;

    /* 文字样式 */
    font-family: PingFangSC-Semibold;
    font-size: 20px;
    color: #2c2d31;
    letter-spacing: 0;
    line-height: 20px;
    font-weight: 600;
  }

  /* 按钮 */
  .conditionSwitcher {
    display: flex;
    border: none;
    margin-bottom: 8px;

    .ant-radio-button-wrapper {
      flex: 1;
      height: 36px;
      line-height: 36px;
      text-align: center;
      font-size: 16px;
      letter-spacing: 0.8px;
      font-family: JiangChengXieHei-500W;
      font-weight: 500;
      color: #6c6f77;
      border: none;
      border-bottom: 2px solid #d9d9d9;
      background: transparent;
      box-shadow: none;

      &::before {
        display: none; // 💥 取消伪元素分隔线
      }

      /* // 🚫 取消按钮之间的边框线
      &:not(:first-child) {
        border-left: none !important;
        margin-left: 0 !important;
      } */

      &.ant-radio-button-wrapper-checked {
        font-family: JiangChengXieHei-700W;
        font-weight: 700;
        color: #2c2d31;
        border-bottom: 2px solid #528eff;
        background: transparent;
      }
    }
  }
`;
