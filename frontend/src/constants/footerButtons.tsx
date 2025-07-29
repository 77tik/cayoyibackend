import { Button } from 'antd';


export const footerButtons = [
  { text: '强化功效模拟', key: 'strength', to: '/simulation' },
  { text: '综合功效表现', key: 'performance', to: '/performance' },
  { text: '安全水头反馈', key: 'safety', to: '/safety' },
  {
    key: 'custom',
    render: () => (
      <Button type="dashed" danger onClick={() => alert('自定义按钮')}>
        自定义按钮
      </Button>
    ),
  },
];
