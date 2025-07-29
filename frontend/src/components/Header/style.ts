import styled from 'styled-components';

export const HeaderWrapper = styled.div`
  display: flex;
  height: 36px;
  width: 100%;
  /* align-items: center; */
  /* justify-content: space-between; */
  margin: 14px 0 10px 3px;

  .pageTitle {
    margin-top: 14px;
    font-family: YouSheBiaoTiHei;
    font-size: 36px;
    color: #ffffff;
    letter-spacing: 6px;
    line-height: 36px;
    text-shadow: 2px 2px 2px rgba(10, 48, 162, 0.4);
    font-weight: 400;
  }
  .navbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .navbar ul {
    display: flex;
    list-style: none;
    padding: 0;
    margin: 0;
    /* margin-left: 72px; */
    margin-left: 46px;
  }

  .navbar li {
    height: 41px;
    width: 162px;
    margin-top: 27px;
    /* margin-right: 30px; */
    margin-right: 10px;

    display: flex;
    justify-content: center;
    align-items: center;
  }

  .navbar a {
    font-family: JiangChengXieHei-500W;
    width: 162px;
    height: 41px;
    font-size: 20px;
    line-height: 41px;
    color: #6c6f77;
    letter-spacing: 1px;
    text-align: center;
    font-weight: 500;
    text-decoration: none;
  }

  .navbar li.selected a {
    font-family: JiangChengXieHei-700W;
    font-size: 20px;
    line-height: 41px;
    color: #2c2d31;
    letter-spacing: 1px;
    text-align: center;
    font-weight: 700;
    background-image: url('/images/title_bg.png');
    background-size: cover;
    position: relative; /* 为伪元素提供定位基准 */
  }
  /* 添加自定义下划线 */
  .navbar li.selected a::after {
    content: '';
    display: block;
    width: 152px;
    height: 4px;
    background: #528eff;
    position: absolute;
    left: 2px;
    bottom: -4px;
  }
  /* .navbar li:hover a {
    color: #e6f2ff;
  } */

  /*应用监测按钮 */
  .headerBtn {
    width: 124px;
    height: 46px;
    position: absolute;
    top: 25px;
    right: 46px; /* 距离右侧 31px */
    background-image: url('/images/编组38.png');
    background-size: cover;
    /* background-position: center; */
    background-color: transparent;
    border: none;
    cursor: pointer;
  }
`;
