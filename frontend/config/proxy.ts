export default {
  // 开发环境
  dev: {
    '/api/': {
      target: 'http://10.0.4.66:8849', // 本地开发时代理请求过去
      changeOrigin: true,
      pathRewrite: { '^/api': ' ' }, // 保持 /api 前缀
      // 添加一个配置来控制是否使用 mock
      // useMock: true, // 改为 true 时使用Mock接口
      useMock: false, // 改为 false 时使用真实接口
    },
  },
  // 测试环境
  test: {
    '/api/': {
      target: 'http://test-server.com',
      changeOrigin: true,
      pathRewrite: { '^/api': '/api' },
      useMock: false,
    },
  },
  // 生产环境
  prod: {
    '/api/': {
      target: 'http://prod-server.com',
      changeOrigin: true,
      pathRewrite: { '^/api': '/api' },
      useMock: false,
    },
  },
};
