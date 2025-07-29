/**
 * Author: harvey_li
 * Date: 2022-08-09 09:23:55
 * LastEditors: harvey_li
 * LastEditTime: 2022-09-02 13:17:26
 * FilePath: \frontend\src\pages\Home\index.tsx
 * Description:
 */
import Guide from '@/components/Guide';
import { trim } from '@/utils/format';
import { PageContainer } from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import styles from './index.less';

// 无须引入styled 即可使用
const Wrapper = styled.div`
  font-size: 20rem;
`;

const HomePage: React.FC = () => {
  const { name } = useModel('global');
  return (
    <PageContainer ghost>
      <Wrapper className={styles.container}>
        <Guide name={trim(name)} />
      </Wrapper>
    </PageContainer>
  );
};

export default HomePage;
