import { Select as AntSelect } from 'antd';
import styled from 'styled-components';
// 多工况参数面板组件样式
export const Wrapper = styled.div`
  .multiConditionPanel {
    height: 730px;
    position: relative; /* 为绝对定位的子元素提供定位上下文 */
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
    width: 225px;
    height: 1px;
    background-color: #d9d9d9; // 自定义线条颜色
  }

  /* 有效水头下的按钮 */
  .multipleGrid {
    display: flex;
    flex-wrap: wrap;
    /* gap: 10px; // 可调整按钮之间的间距 */
    gap: 16px 10px;
  }

  .multipleItem {
    width: 184px;
    height: 32px;
    /* margin: 8px 0; */
    border-radius: 16px;
    /* background: rgba(167, 179, 205, 0.2); */
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    border: none;
    background: url('/images/椭圆形9蒙版.png') no-repeat center/cover;
    opacity: 1;
    span {
      font-family: JiangChengXieHei-700W;
      font-size: 16px;
      color: #ffffff;
      letter-spacing: 0;
      line-height: 20px;
      text-align: center;
      /* font-weight: 700; */
      text-shadow: 0 4px 4px rgba(2, 166, 200, 0.5);
    }

    box-shadow: inset 0 -1px 3px 0 #0fcfe5,
      inset 0 2px 1px 0 rgba(255, 255, 255, 0.5);

    /* &.multipleSelect {
      background: rgba(167, 179, 205, 0.2);
      box-shadow: none;
      border: 1px solid #d0d6e3;
    } */
  }

  /* 查询按钮样式 */
  .multiQueryButton {
    transition: all 0.2s;
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
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    bottom: 40px;
    left: 50%;
    transform: translateX(-50%);

    &:hover {
      transform: translateX(-50%) translateY(-1px);
      box-shadow: 0 4px 8px rgba(82, 142, 255, 0.3);
      opacity: 0.9;
    }

    &:disabled {
      background: #ccc;
      cursor: not-allowed;
      margin-bottom: 0;
      box-shadow: none;
      opacity: 1;
    }
  }
`;

export const Select = styled(AntSelect)``;
