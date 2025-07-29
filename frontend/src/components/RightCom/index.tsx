import { useSafeState, useUpdateEffect } from 'ahooks';
import styled from 'styled-components';
/**
 * children legend组件
 * contentProp
 */
interface RightComProps {
  /**
   * 是否在收起时销毁传进去的子组件
   */
  destroyOnClose?: boolean; // 收起的时候是否销毁children元素
  /**
   * 最外层样式修正
   */
  outerStyle?: React.CSSProperties;
  /**
   * 内容部分的组件
   */
  contentProp: {
    /**
     * 内容部分宽度
     */
    width?: number;
    /**
     * 内容部分高度
     */
    height?: number;
    /**
     * 内容组件
     */
    children: JSX.Element | null;
    /**
     * content 样式
     */
    style?: React.CSSProperties;
  };
  /**
   * 左侧legend
   */
  legendProp?: {
    /**
     * legend距离内容的距离
     */
    marginRight: number;
    /**
     * 内容组件
     */
    content: {
      /**
       * jsx
       */
      children: JSX.Element | null;
      /**
       * 其他样式
       */
      style?: React.CSSProperties;
    }[];
  };
  /**
   * 侧边栏展开关闭操作按钮
   */
  hiddenSideBar?: {
    /**
     * 是否展示
     */
    isShow?: boolean;
    /**
     * 是否展示装饰条
     */
    isShowDecoratorBar?: boolean;
    /**
     * 按钮组件
     */
    decoratorChildren?: JSX.Element | null;
    /**
     * 开启按钮图片
     */
    btnOpenImgUrl?: string;
    /**
     * 关闭按钮图片
     */
    btnCloseImgUrl?: string;
  };
  /**
   * 是否隐藏侧边栏
   */
  isHideRight?: boolean;
  /**
   * 侧边栏展开或者显示的回调函数
   */
  handleRightShowChange?: Function;
}
// 左侧样式框架
const RightCom = ({
  outerStyle = {},
  destroyOnClose = false,
  contentProp = {
    width: 400,
    height: 800,
    children: null,
    style: {},
  },
  legendProp = {
    marginRight: 0,
    content: [],
  },
  hiddenSideBar = {
    isShow: false,
    decoratorChildren: null,
    isShowDecoratorBar: true,
  },
  isHideRight,
  handleRightShowChange,
}: RightComProps) => {
  const [isOpen, setIsOpen] = useSafeState(true);
  useUpdateEffect(() => {
    if (isHideRight != undefined) setIsOpen(!isHideRight);
  }, [isHideRight]);

  const imgObj = {
    open: require('./images/open.png'),
    close: require('./images/close.png'),
  };

  return (
    <Wrapper
      style={outerStyle}
      marginRight={legendProp?.marginRight || 0}
      isOpen={isOpen}
      contentWidth={contentProp.width}
      contentHeight={contentProp.height || 0}
    >
      {/* 隐藏按钮 */}
      {hiddenSideBar.isShow && (
        <div
          className={[
            'content-operation transform-x_time',
            isOpen ? 'transform_in' : 'right-transform_out_optBtn',
          ].join(' ')}
        >
          <div
            className="content-operation-btn"
            onClick={() => {
              setIsOpen(!isOpen);
              handleRightShowChange && handleRightShowChange(!isOpen);
            }}
          >
            {!isOpen ? (
              <img src={imgObj.open} alt="开启" />
            ) : (
              <img src={imgObj.close} alt="关闭" />
            )}
          </div>
          {hiddenSideBar.isShowDecoratorBar &&
            (hiddenSideBar.decoratorChildren || (
              <div className="content-decorator-line color-side-line" />
            ))}
        </div>
      )}
      <div
        style={contentProp.style}
        className={[
          `content transform-x_time`,
          isOpen ? 'transform_in' : 'right-transform_out',
        ].join(' ')}
      >
        {destroyOnClose ? isOpen && contentProp.children : contentProp.children}
      </div>
      {/* 左侧左下角图例区域 */}
      {legendProp?.content.length > 0 &&
        legendProp.content.map((item) => (
          <div
            style={item?.style}
            className={[
              'content_right_legend transform-x_time',
              isOpen ? '' : 'transform_in',
            ].join(' ')}
          >
            {item.children}
          </div>
        ))}
    </Wrapper>
  );
};

const Wrapper = styled.div<{
  marginRight: number;
  contentWidth: number;
  contentHeight: number;
  isOpen: boolean;
}>`
  position: fixed;
  z-index: 3;
  top: 90rem;
  display: flex;
  color: #fff;
  right: 20rem;
  .content {
    width: ${(props) => (props.isOpen ? `${props.contentWidth}rem` : 0)};
    height: ${(props) => `${props.contentHeight}rem`};
    overflow-x: hidden;
    overflow-y: hidden;
    ::-webkit-scrollbar {
      /*整体样式*/
      width: 5rem;
    }
    ::-webkit-scrollbar-thumb {
      /*滚动条小方块*/
      border-radius: 10rem;
      background: rgba(154, 167, 186, 0.8);
      border-radius: 3rem;
    }
  }
  /* 左侧的legend储物架 */
  .content_right_legend {
    width: fit-content;
    position: absolute;
    right: ${(props) =>
      `${props.isOpen ? props.contentWidth + props.marginRight : 0}rem`};
  }
  .content-operation {
    width: 1;
    height: ${(props) => `${props.contentHeight}rem`};
    display: flex;
    align-items: center;
    .content-decorator-line {
      margin-right: 1rem;
      width: 1rem;
      height: 100%;
    }
    img {
      width: 20rem;
      height: 70rem;
    }
    .content-operation-btn {
      width: 1;
      height: 56rem;
      cursor: pointer;
      z-index: 99;
    }
  }
  /* 左右两个仿抽屉--不删除DOM */
  .transform-x_time {
    transition: all 200ms;
    transition-timing-function: linear;
  }
  /* 内容展开 */
  .transform_in {
    transform: translateX(0);
  }
  /* 内容缩回去的时候动作 */
  .right-transform_out {
    transform: translateX(420rem);
  }
  /* 关闭按钮 缩回去的时候多缩一点点 */
  .right-transform_out_optBtn {
    transform: translateX(20rem);
  }
  // 侧边竖线
  .color-side-line {
    background-image: linear-gradient(
      180deg,
      rgba(188, 224, 255, 0.1) 0%,
      #b9d1ff 51%,
      rgba(188, 224, 255, 0.1) 100%
    );
  }
`;

export { RightCom };
