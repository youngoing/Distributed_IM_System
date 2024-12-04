import { useEffect } from 'react';
import { useRouter } from 'next/router';
import Msg from './msg';
import { isClientAuthenticated } from '../utils/auth';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    const checkAuthentication = async () => {
      const isAuthenticated = await isClientAuthenticated();
      console.log('isAuthenticated:', isAuthenticated);
      if (!isAuthenticated) {
        
        // 客户端重定向到登录页面
        router.push('/login');
      }
    };

    checkAuthentication();
  }, [router]);

  return (
    <Msg>
      {/* 这里可以放置其他内容 */}
    </Msg>
  );
}