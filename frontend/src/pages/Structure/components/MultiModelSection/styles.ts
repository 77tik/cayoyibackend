import styled from 'styled-components';

export const SectionViewWrapper = styled.div`
  .section-container {
    height: 676px;
    margin-top: 82px;
    margin-bottom: 30px;
    display: flex;
    flex-wrap: wrap;
    gap: 12px 16px;
    position: relative; // 添加相对定位
  }

  .section-item {
    background-color: #fff;
    width: 320px;
    height: 160px;
    margin-bottom: 0;
    /* height: 100%; */
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    border-radius: 6px;

    &.hidden {
      visibility: hidden; // 当处于全屏模式时隐藏原始项
    }

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

      img {
        filter: invert(67%) sepia(8%) saturate(240%) hue-rotate(182deg)
          brightness(90%) contrast(86%); // 将图片颜色改为 #A4A7AD
        width: 10px; // 设置图片大小
        height: 10px;
        /* position: absolute; */
        /* right: 16px; */
        /* top: 50%; */
        /* transform: translateY(-50%); */
        cursor: pointer; // 添加指针样式
      }
    }
    .section-item-content {
      height: 100%;
      background-color: gray;
      margin: 0 1px 1px 1px;
      border-radius: 0 0 6px 6px;
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
          border-bottom: 1px solid #d0d6e3;

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

  // 添加全屏覆盖层样式
  .fullscreen-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 1330px;
    height: 676px;
    /* width: 100%; */
    /* height: 100%; */
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 8px;
    z-index: 100;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .fullscreen-content {
    width: 100%;
    height: 100%;
    background-color: #fff;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
  }

  .fullscreen-header {
    height: 40px;
    line-height: 40px;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-family: JiangChengXieHei-700W;
    font-size: 18px;
    font-weight: 700;
    color: #2c2d31;

    img {
      width: 12px;
      height: 13px;
      filter: invert(67%) sepia(8%) saturate(240%) hue-rotate(182deg)
        brightness(90%) contrast(86%);
      cursor: pointer;
    }
  }

  .fullscreen-body {
    flex: 1;
    background-color: gray;
    margin: 0 1px 1px 1px;
    border-radius: 0 0 8px 8px;
  }
`;
