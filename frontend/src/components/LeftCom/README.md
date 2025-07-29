## 左侧推拉框组件

该组件默认展示为左侧可拉伸框形式，可修改为左侧固定展示栏，应用于更多业务场景

## 使用方式以及注意事项

### 代码示例

```javascript
<LeftCom
  destroyOnClose={false}
  contentProp={{
    children: <Component />
  }}
  isHideLeft={false}
  handleLeftShowChange={(e: boolean) => {
    message.info(e ? '左侧展开了' : '左侧隐藏了')
  }}
/>
```

### 参数说明

| 参数                 | 说明                           | 类型     | 默认值                                                                                                                               | 必填项 |
| -------------------- | ------------------------------ | -------- | ------------------------------------------------------------------------------------------------------------------------------------ | ------ |
| destroyOnClose       | 收起的时候是否销毁children元素 | boolean  | false                                                                                                                                |        |
| contentProps         | 内容部分                       | object   | `{width?: number; height?: number; children: JSX.Element; style?: React.CSSProperties; }`                                            | √      |
| legendProp           | 外部浮动展示框（如：图例等）   | object   | `{ marginLeft: number; content:{children:JSX.Element; style?:React.CSSProperties}[]}`                                                |        |
| hiddenSideBar        | 隐藏侧边栏                     | object   | `{isShow?: boolean; isShowDecoratorBar?: boolean;decoratorChildren?: JSX.Element; btnOpenImgUrl?: string; btnCloseImgUrl?: string }` |        |
| isHideLeft           | 是否隐藏侧边栏                 | boolean  | false                                                                                                                                |        |
| handleLeftShowChange | 侧边栏展开或者显示的回调函数   | Function |                                                                                                                                      |        |
