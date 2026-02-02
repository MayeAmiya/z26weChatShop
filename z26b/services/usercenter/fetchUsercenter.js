import { config } from '../../config/index';
import { getUserInfo, testLogin } from '../auth/auth';

/** 获取个人中心信息 */
function mockFetchUserCenter() {
  const { delay } = require('../_utils/delay');
  const { genUsercenter } = require('../../model/usercenter');
  return delay(200).then(() => genUsercenter());
}

/** 获取个人中心信息 */
export function fetchUserCenter() {
  if (config.useMock) {
    return mockFetchUserCenter();
  }
  return (async () => {
    // 确保已有用户登录（测试账号自动登录）
    let user = getUserInfo();
    if (!user) {
      await testLogin();
      user = getUserInfo();
    }

    const userInfo = {
      avatarUrl: user?.avatar || user?.Avatar || '',
      nickName: user?.nickName || user?.NickName || '微信用户',
      phoneNumber: user?.phoneNumber || '',
    };

    const countsData = [
      {
        type: 'address',
        num: user?.addressCount ?? 0,
      },
    ];

    const customerServiceInfo = {
      servicePhone: '400-000-0000',
      serviceTimeDuration: '09:00-18:00',
    };

    return { userInfo, countsData, customerServiceInfo };
  })();
}

