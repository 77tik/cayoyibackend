import { Button, Checkbox, Form, FormProps, Input } from 'antd';
import { history } from '@umijs/max';
import forge from 'node-forge';

const Wrapper = styled.div`
  font-size: 20rem;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
`;

type FieldType = {
  username?: string;
  password?: string;
};
const LoginPage = () => {
  const onFinish: FormProps<FieldType>['onFinish'] = async (values) => {
    console.log('Success:', values);
    if (values.password) {
      // 导入公钥 公钥可以通过API请求从服务器获取
      let publicKeyValue = '';
      if (publicKeyValue) {
        // 获取输入的密码
        const pw = values.password;

        // 将Base64编码的公钥解码为DER格式。
        const publicKeyDer = forge.util.decode64(publicKeyValue);

        // 将DER格式的公钥解析成ASN.1格式。
        const asn1 = forge.asn1.fromDer(publicKeyDer);

        // 从ASN.1格式中提取出公钥对象。
        const publicKey = forge.pki.publicKeyFromAsn1(asn1);

        // 使用RSAES-OAEP作为加密方案，指定SHA-256作为消息摘要算法（md）和MGF1掩码生成函数（mgf1）
        const encrypted = publicKey.encrypt(pw, 'RSAES-OAEP', {
          md: forge.md.sha256.create(),
          mgf1: {
            md: forge.md.sha256.create(),
          },
        });

        //将加密后的数据编码为Base64
        const encryptedBase64 = forge.util.encode64(encrypted);

        //调用登录接口将encryptedBase64作为入参
        // await login({password:encryptedBase64})
      }
    }
    // 跳转首页
    // history.replace('/home');
    history.replace('/fluid');
  };

  return (
    <Wrapper>
      <Form
        name="basic"
        labelCol={{ span: 4 }}
        wrapperCol={{ span: 20 }}
        style={{ width: 400 }}
        initialValues={{ remember: true }}
        onFinish={onFinish}
        autoComplete="off"
      >
        <Form.Item<FieldType>
          label="用户名"
          name="username"
          rules={[{ required: true, message: '请输入用户名!' }]}
        >
          <Input />
        </Form.Item>

        <Form.Item<FieldType>
          label="密码"
          name="password"
          rules={[{ required: true, message: '请输入密码!' }]}
        >
          <Input.Password />
        </Form.Item>

        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            style={{ marginLeft: '200px' }}
          >
            登录
          </Button>
        </Form.Item>
      </Form>
    </Wrapper>
  );
};

export default LoginPage;
