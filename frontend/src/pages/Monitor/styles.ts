export const Wrapper = styled.div`
  background-image: url('/images/backgrounds/main-bg.png');
  background-size: cover;
  background-position: center;
  background-position: center top;
  /* margin: 0px 13px 20px 11px; */
  width: 1896px;

  margin: 0px 13px 0px 11px;

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

  /* 四个机组 */
  .monitor-content {
    /* height: 676px; */
    margin-top: 62px;
    /* margin-bottom: 30px; */
    display: flex;
    flex-wrap: wrap;
    gap: 0 20px;
    position: relative; // 添加相对定位
  }
  .section-item {
    background-color: #fff;
    width: 440px;
    height: 466px;
    margin-bottom: 0;
    /* height: 100%; */
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    border-radius: 6px;

    .section-item-title {
      height: 32px;
      line-height: 32px;
      margin-left: 16px;
      font-family: JiangChengXieHei-700W;
      font-size: 16px;
      color: #2c2d31;
      letter-spacing: 0.8px;
      font-weight: 700;

      display: flex;
      align-items: center;
      justify-content: space-between;
      padding-right: 19px;
    }
    .section-item-content {
      height: 100%;
      background-color: rgba(218, 231, 255);
      margin: 0 1px 1px 1px;
      border-radius: 0 0 6px 6px;
      position: relative;

      .cover-strain,
      .volute-manhole-door-strain {
        height: 46px;
        opacity: 0.8;
        background: #465373;
        border: 1px solid rgba(144, 169, 222, 1);
        border-radius: 4px;
        font-family: PingFangSC-Medium;
        font-size: 20px;
        color: #ffffff;
        line-height: 46px;
        font-weight: 500;
        text-align: center;
        position: absolute;
      }
      .cover-strain {
        width: 212px;
        top: 24px;
        left: 20px;
      }
      .volute-manhole-door-strain {
        width: 262px;
        right: 20px;
        bottom: 24px;
      }

      /* 模型图 */
      .strain-model {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
      }
    }
  }

  .monitor-title {
    display: flex;
    flex-wrap: wrap;
    gap: 0 20px;
  }
  /* 日期选择框 */
  .date-select {
    height: 42px;
    width: 518px;
    display: flex;
    align-items: center;

    font-family: PingFangSC-Semibold;
    font-size: 16px;
    color: #2c2d31;
    letter-spacing: 0;
    line-height: 42px;
    font-weight: 600;
    background: rgba(255, 255, 255, 0.8);
    box-shadow: 0px 4px 6px 0px rgba(82, 142, 255, 0.2);
    border-radius: 6px;

    margin-top: 22px;
    background: linear-gradient(to right, #90a9de, #c3d4f1);
    padding: 1px;

    .date-select-data {
      width: 100%;
      height: 100%;
      border-radius: 6px;
      background: rgba(255, 255, 255, 0.8);
      padding-left: 24px;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }

  .date-select .ant-picker {
    border: none !important;
    background-color: rgba(255, 255, 255, 0.1) !important;
  }
  .date-select .ant-picker-focused {
    border: none !important;
    box-shadow: none !important;
    outline: none !important;
  }

  .date-select .ant-picker-input input {
    text-align: center;
  }

  /* 数据切换 */
  .data-switch {
    height: 42px;
    width: 372px;
    display: flex;
    align-items: center;

    font-family: PingFangSC-Semibold;
    font-size: 16px;
    color: #2c2d31;
    letter-spacing: 0;
    line-height: 42px;
    font-weight: 600;
    background: rgba(255, 255, 255, 0.8);
    box-shadow: 0px 4px 6px 0px rgba(82, 142, 255, 0.2);
    border-radius: 6px;

    margin-top: 22px;
    background: linear-gradient(to right, #90a9de, #c3d4f1);
    padding: 1px;

    .data-switch-content {
      width: 100%;
      height: 100%;
      border-radius: 6px;
      background: rgba(255, 255, 255, 0.8);
      padding: 0 32px;
      display: flex;

      .data-switch-options {
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 6px;

        img {
          height: 16px;
          width: 16px;
          vertical-align: middle;
        }

        span {
          vertical-align: middle;
          line-height: 1;
          transform: translateY(-1px);
        }
      }
    }
  }

  /* 下方区域 */
  .monitor-body {
    /* height: 926px; */
    height: 360px;
    margin-top: 16px;
    position: relative;
    display: flex;
    flex: 1;
    gap: 20px;

    /* 监测列表 */
    .monitor-list,
    .monitor-charts {
      height: 360px;
      width: 910px;
      /* width: 905px; */
      display: flex;
      flex-direction: column;
      padding: 20px 16px;
      background: rgba(255, 255, 255, 0.8);
      border: 1px solid rgba(210, 216, 229, 1);
      box-shadow: 0px 4px 6px 0px rgba(82, 142, 255, 0.2);
      border-radius: 6px;

      .monitor-list-title {
        width: 878px;
        height: 36px;
        line-height: 36px;
        margin-bottom: 12px;
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
        border-radius: 4px;
      }
    }
  }
`;
