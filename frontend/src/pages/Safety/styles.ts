export const Wrapper = styled.div`
  background-image: url('/images/backgrounds/main-bg.png');
  background-size: cover;
  background-position: center;
  background-position: center top;
  margin: 0px 13px 0px 11px;
  width: 1896px;

  min-height: calc(100vh - 20px);

  position: relative;

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
    margin-top: 0px;
    margin-left: 8px;
    padding: 0 0 26px 21px;
  }

  /* 页面 */
  .content {
    display: flex;
    flex: 1;
  }

  .body {
    height: 926px;
    /* height: 921px; */
    margin-top: 42px;
    position: relative;
    display: flex;
    flex: 1;
    gap: 10px;
  }
  .left,
  .right {
    width: 910px;
    /* width: 905px; */
    display: flex;
    flex-direction: column;
    padding: 20px 16px;
    background: rgba(255, 255, 255, 0.8);
    border: 1px solid rgba(210, 216, 229, 1);
    box-shadow: 0px 4px 6px 0px rgba(82, 142, 255, 0.2);
    border-radius: 6px;

    .safeTitle {
      width: 878px;
      height: 36px;
      line-height: 36px;
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

  /* 左边 */
  /* 文字部分 */
  .safeText {
    margin-top: 20px;
    font-family: PingFangSC-Medium;
    font-size: 16px;
    color: #2c2d31;
    letter-spacing: 0.8px;
    height: 20px;
    line-height: 20px;
    font-weight: 500;
  }
  /* 图表部分 */
  .leftChart {
    margin-top: 12px;
    width: 872px;
    height: 642px;
    background: rgba(82, 142, 255, 0.1);
    /* background: rgba(67, 76, 63, 1); */
    /* background: rgba(0, 115, 100, 1); */
    border: 2px solid rgba(82, 142, 255, 1);
  }
  /* 图表结果 */
  .chartResult {
    margin-top: 20px;
    width: 872px;
    height: 136px;
    padding: 20px;
    background: rgba(82, 142, 255, 0.1);
    border: 2px solid rgba(82, 142, 255, 1);
    display: flex;
    flex-direction: row;
    gap: 12px;

    .chartResultLeftItem {
      width: 618px;
      height: 96px;
      padding: 15px 16px;
      background: rgba(255, 255, 255, 0.2);
      border: 1px solid rgba(255, 255, 255, 1);
      gap: 10px;
      display: flex;
      flex-direction: row;

      .chartResultLeftItemIcon {
        width: 65px;
        height: 66px;
      }
      .chartResultLeftItemText {
        width: 510px;
        font-family: PingFangSC-Medium;
        font-size: 16px;
        color: #6c6f77;
        letter-spacing: 1px;
        line-height: 16px;
        font-weight: 500;
      }
    }
    .chartResultRightItem {
      width: 202px;
      height: 96px;
      background: rgba(255, 255, 255, 0.2);
      border: 1px solid rgba(255, 255, 255, 1);
      gap: 12px;
      padding: 22px 0;
      text-align: center;
      display: flex;
      flex-direction: column;
      justify-content: center;

      .chartResultRightItemContent {
        font-family: JiangChengXieHei-700W;
        font-size: 24px;
        line-height: 24px;
        color: #528eff;
        letter-spacing: 1.2px;
        font-weight: 700;
      }

      .chartResultRightItemTitle {
        font-family: PingFangSC-Medium;
        font-size: 16px;
        line-height: 16px;
        color: #2c2d31;
        letter-spacing: 0;
        font-weight: 500;
      }
    }
  }

  /* 右边 */
  /* 下载按钮 */
  .downloadBtn {
    margin-top: 20px;
    font-family: PingFangSC-Regular;
    font-size: 16px;
    color: #528eff;
    height: 20px;
    line-height: 20px;
    font-weight: 400;
    display: flex;
    justify-content: flex-end;

    .downloadText {
      cursor: pointer;
    }
  }
  /* 报告文件 */
  .report {
    margin-top: 12px;
    width: 876px;
    height: 796px;
    padding: 12px;
    background: rgba(82, 142, 255, 0.1);
    border: 2px solid rgba(82, 142, 255, 1);
    overflow: hidden;

    .reportContent {
      width: 100%;
      height: 100%;
      /* width: 848px; */
      /* height: 768px; */
      /* margin: 12px; */
      background: rgb(255, 255, 255);

      padding: 0;
      margin: 0;

      iframe {
        width: 100%;
        height: 100%;
        border: none;
        clip-path: inset(16px 41px 20px 25px);
      }
    }
  }
`;
