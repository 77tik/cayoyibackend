## Description
基于[UMI4](https://umijs.org/?ref=nav.poetries.top)搭建的前端脚手架，封装了一些通用的组件、方法，应用于数字孪生+格物CAE。

## 环境准备
node >= 14.0

## 开始

Install dependencies,

```bash
$ yarn
```

Start the dev server,

```bash
$ yarn start
```

Build,

```bash
$ yarn build
```

## 组件

+ Axios
/src/utils/http.ts
脚手架没有采用UMI自带的request，封装了Axios，并提供了接口验证，接口错误统一处理的功能。

+ Rem
/src/utils/flexible.ts
脚手架封装了rem，实现了适配不同屏幕分辨率的功能。

+ Mobx
脚手架封装了Mobx作为状态管理

+ styled-components
脚手架使用styled-components作为css in js，主要因为styled-components支持动态样式。



