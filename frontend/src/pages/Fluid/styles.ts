export const Wrapper = styled.div`
  background-image: url('/images/backgrounds/main-bg.png');
  background-size: cover;
  background-position: center;
  background-position: center top;
  /* margin: 0px 13px 20px 11px; */
  margin: 0px 13px 0px 11px;
  width: 1896px;
  /* box-sizing: content-box; */
  min-height: calc(100vh - 20px);
  /* background-color: rgba(218, 231, 255); */

  position: relative; // 新增，确保伪元素定位基于Wrapper

  &::before {
    content: '';
    position: fixed; // 用fixed保证全屏
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: -1; // 保证在内容下方
    background-color: rgba(218, 231, 255);
  }

  /* 容器 */
  .container {
    display: flex;
    flex-direction: column;
    /* margin: 6px 8px 8px 8px; */
    margin-top: 0px;
    margin-left: 8px;
    padding: 0 0 26px 21px;
  }

  /* 页面 */
  .content {
    display: flex;
    flex: 1;
    /* overflow: hidden; */
    /* gap: 20px; */
  }

  .left {
    margin-top: 24px;
    width: 420px;
  }
  .body {
    /* width: 100%; */
    /* height: 100%; */
    margin-top: 48px;
    margin-right: 49px;
    position: relative;
  }

  /* 右侧按钮 */
  .rightBtn {
    height: 42px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 10px;
    gap: 10px;
    background: rgba(255, 255, 255, 0.8);
    z-index: 1;
    border-radius: 6px;

    position: absolute; // 使用绝对定位
    right: -1px; // 距离右侧49px

    .function-btn {
      width: 86px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 4px;
      border: none;
      background: transparent;
      border-radius: 4px;
      font-family: PingFangSC-Medium;
      font-size: 14px;
      color: #a4a7ad;
      letter-spacing: 0;
      text-align: center;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.3s;

      img {
        width: 16px;
        height: 16px;
      }

      &:hover,
      &.active {
        background: rgba(82, 142, 255, 0.2);
        color: rgb(137, 143, 155);
      }
    }
  }

  .modelView {
    /* height: 640px; */
    /* width: 1390px; */
    height: 100%;
    width: 100%;
    min-height: 60vh;
  }
  .footer {
    display: flex;
    justify-content: center;
    gap: 20px;
    /* padding: 20px 0; */
    border-top: 1px solid #1890ff;
    margin-top: 20px;
  }
`;
