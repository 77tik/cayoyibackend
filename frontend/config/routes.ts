export default [
  {
    path: '/',
    redirect: '/fluid',
  },
  {
    name: '登录',
    path: '/login',
    component: './Login',
    layout: false,
    hideChildrenInMenu: true,
    hideInMenu: true,
  },
  {
    path: '/fluid',
    component: '@/pages/Fluid',
    name: '流体仿真模拟',
    layout: false,
  },
  {
    path: '/structure',
    component: '@/pages/Structure',
    name: '结构仿真模拟',
    layout: false,
  },
  {
    path: '/safety',
    component: '@/pages/Safety',
    name: '安全水头区域',
    layout: false,
  },
  {
    path: '/monitor',
    component: '@/pages/Monitor',
    name: '应变数据展示',
    layout: false,
  },
];
