// ProtectedComponent.js
import React, { useState, useEffect } from 'react';
import Modal from './Modal';
import { isClientAuthenticated } from '../utils/auth';  // 导入你的认证检查函数

const ProtectedComponent = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);  // 初始状态为 false
    const [loading, setLoading] = useState(true);  // 初始状态为 true，表示正在加载

    const [isModalOpen, setIsModalOpen] = useState(false);

    // 检查用户的认证状态
    useEffect(() => {
        const checkAuthentication = async () => {
            try {
                const authStatus = await isClientAuthenticated();  // 获取认证状态
                setIsAuthenticated(authStatus);  // 更新认证状态
            } catch (error) {
                console.error("认证检查失败:", error);
            } finally {
                setLoading(false);  // 完成加载，设置 loading 为 false
            }
        };

        checkAuthentication();  // 调用认证检查函数
    }, []);

    useEffect(() => {
        if (!loading && !isAuthenticated) {
            setIsModalOpen(true);  // 如果未认证且加载完成，则打开 Modal
        }
    }, [loading, isAuthenticated]);  // 依赖 loading 和 isAuthenticated

    if (loading) {
        return <div>加载中...</div>;  // 加载中显示 loading 信息
    }

    return (
        <div>
            {isModalOpen && (
                <Modal
                    isOpen={isModalOpen}
                    onClose={() => setIsModalOpen(false)}  // 关闭 Modal 的方法
                />
            )}
            {!isModalOpen && isAuthenticated && children}  {/* 如果认证并且 Modal 已关闭，展示子组件 */}
        </div>
    );
};

export default ProtectedComponent;
