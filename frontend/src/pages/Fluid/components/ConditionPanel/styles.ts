import styled from 'styled-components';
// 单工况参数面板组件样式
export const Wrapper = styled.div`
  .conditionPanel {
    height: 680px;
  }
  .pointItem {
    padding: 8px;
    border: 1px solid #1890ff;
    border-radius: 4px;
    text-align: center;
    cursor: pointer;
    color: #2c2d31;
    &:hover {
      background: rgba(24, 144, 255, 0.2);
    }
  }
  /* 小标题样式 */
  .titleRow {
    margin: 16px 0 16px 0;
    display: flex;
    align-items: center;
    height: 16px;
    position: relative; // 为右侧横线定位准备
  }

  .titleIcon {
  }

  .titleText {
    margin: 0px 11px 0 8px;
    font-family: PingFangSC-Medium;
    font-size: 16px;
    color: #2c2d31;
    letter-spacing: 0.89px;
    line-height: 16px;
    font-weight: 500;
  }

  .titleLine {
    position: absolute;
    right: 0;
    bottom: 0;
    width: 262px;
    height: 1px;
    background-color: #d9d9d9; // 自定义线条颜色
  }
  /* 有效水头下的按钮 */
  .pointsGrid {
    display: flex;
    flex-wrap: wrap;
    /* gap: 10px; // 可调整按钮之间的间距 */
    gap: 16px 10px;
  }

  .pointItem {
    width: 120px;
    height: 32px;
    /* margin: 8px 0; */
    border-radius: 16px;
    background: rgba(167, 179, 205, 0.2);
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    border: none;
    span {
      font-family: JiangChengXieHei-400W;
      font-size: 16px;
      color: #6c6f77;
      letter-spacing: 0;
      line-height: 20px;
      font-weight: 400;
      text-align: center;
    }

    &.active {
      background: url('/images/椭圆形9蒙版.png') no-repeat center/cover;
      opacity: 1;

      span {
        font-family: JiangChengXieHei-700W;
        font-size: 16px;
        color: #ffffff;
        letter-spacing: 0;
        line-height: 20px;
        text-align: center;
        font-weight: 700;
        text-shadow: 0 4px 4px rgba(2, 166, 200, 0.5);
      }
    }
  }

  /* 有功功率 */
  .powerItem {
    width: 380px;
    height: 32px;
    margin-top: 16px;
    border-radius: 16px;
    background: rgba(167, 179, 205, 0.2);
    align-items: center;
    border: none;
    display: flex;
    justify-content: space-between;
    /* cursor: pointer; */
    .powerItemLeft,
    .powerItemRight {
      width: 150px;
      height: 20px;
      flex: 1;
      font-family: JiangChengXieHei-400W;
      font-size: 16px;
      color: #6c6f77;
      line-height: 20px;
      text-align: center;
      font-weight: 400;
    }

    .divider {
      width: 1px;
      height: 24px;
      background-color: #d9d9d9;
    }

    &.active {
      background: url('/images/椭圆形9蒙版.png') no-repeat center/cover;
      opacity: 1;
      /* border: none; */
      span {
        font-family: JiangChengXieHei-700W;
        color: #fff;
        text-shadow: 0 4px 4px rgba(2, 166, 200, 0.5);
        font-weight: 700;
      }
      .divider {
        background-color: #fff; /* 添加这一行 */
      }
    }
  }

  /* 运行环境参数 */
  .paramGrid {
    display: grid;
    grid-template-columns: repeat(2, 180px);
    background: #f8fafd;
    gap: 12px 16px;
  }
  .paramItem {
    width: 180px;
    height: 80px;
    padding: 16px 0 16px 20px;
    background: linear-gradient(180deg, #edf2f8 0%, #ffffff 100%);
    box-shadow: 0 2px 4px rgba(22, 90, 152, 0.2);
    border-radius: 6px;
    display: flex;
  }

  /* 图片容器 */
  .paramImage {
    margin-right: 10.5px;
    flex-shrink: 0;
  }
  .paramImage img {
    object-fit: contain;
  }

  .paramContent {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 100%;
    padding: 2px 0;
  }
  .paramValue {
    font-family: JiangChengXieHei-700W;
    font-size: 20px;
    color: #528eff;
    line-height: 1;
    margin-bottom: 10px;
  }
  .paramUnit {
    font-size: 16px;
    margin-left: 6px;
  }
  .paramLabel {
    font-family: PingFangSC-Regular;
    font-size: 16px;
    color: #6c6f77;
    line-height: 1;
  }
`;
