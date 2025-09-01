// 工具函数库

// 格式化日期
function formatDate(dateString, format = 'YYYY-MM-DD') {
    if (!dateString) return '';
    
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return '';
    
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    switch (format) {
        case 'YYYY-MM-DD':
            return `${year}-${month}-${day}`;
        case 'YYYY-MM-DD HH:mm':
            return `${year}-${month}-${day} ${hours}:${minutes}`;
        case 'YYYY-MM-DD HH:mm:ss':
            return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        case 'MM-DD':
            return `${month}-${day}`;
        case 'HH:mm':
            return `${hours}:${minutes}`;
        default:
            return `${year}-${month}-${day}`;
    }
}

// 格式化相对时间
function formatRelativeTime(dateString) {
    if (!dateString) return '';
    
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    
    const minutes = Math.floor(diff / (1000 * 60));
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    
    if (minutes < 1) return '刚刚';
    if (minutes < 60) return `${minutes}分钟前`;
    if (hours < 24) return `${hours}小时前`;
    if (days < 30) return `${days}天前`;
    
    return formatDate(dateString);
}

// 防抖函数
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 节流函数
function throttle(func, limit) {
    let inThrottle;
    return function(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}

// 深拷贝
function deepClone(obj) {
    if (obj === null || typeof obj !== 'object') return obj;
    if (obj instanceof Date) return new Date(obj.getTime());
    if (obj instanceof Array) return obj.map(item => deepClone(item));
    if (typeof obj === 'object') {
        const clonedObj = {};
        for (const key in obj) {
            if (obj.hasOwnProperty(key)) {
                clonedObj[key] = deepClone(obj[key]);
            }
        }
        return clonedObj;
    }
}

// 验证电话号码
function validatePhone(phone) {
    const phoneRegex = /^1[3-9]\d{9}$/;
    return phoneRegex.test(phone);
}

// 验证邮箱
function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// 生成随机ID
function generateId(prefix = 'id') {
    return `${prefix}_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

// 显示消息提示
function showMessage(message, type = 'info', duration = 3000) {
    // 移除已存在的消息
    const existingMessage = document.querySelector('.message-toast');
    if (existingMessage) {
        existingMessage.remove();
    }
    
    // 创建消息元素
    const messageEl = document.createElement('div');
    messageEl.className = `message-toast message-${type}`;
    messageEl.textContent = message;
    
    // 添加样式
    messageEl.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 12px 24px;
        border-radius: 4px;
        color: white;
        font-size: 14px;
        z-index: 10000;
        opacity: 0;
        transform: translateX(100%);
        transition: all 0.3s ease;
    `;
    
    // 根据类型设置颜色
    switch (type) {
        case 'success':
            messageEl.style.backgroundColor = '#52c41a';
            break;
        case 'error':
            messageEl.style.backgroundColor = '#ff4d4f';
            break;
        case 'warning':
            messageEl.style.backgroundColor = '#faad14';
            break;
        default:
            messageEl.style.backgroundColor = '#1890ff';
    }
    
    document.body.appendChild(messageEl);
    
    // 显示动画
    setTimeout(() => {
        messageEl.style.opacity = '1';
        messageEl.style.transform = 'translateX(0)';
    }, 10);
    
    // 自动隐藏
    setTimeout(() => {
        messageEl.style.opacity = '0';
        messageEl.style.transform = 'translateX(100%)';
        setTimeout(() => {
            if (messageEl.parentNode) {
                messageEl.parentNode.removeChild(messageEl);
            }
        }, 300);
    }, duration);
}

// 确认对话框
function showConfirm(message, onConfirm, onCancel) {
    const overlay = document.createElement('div');
    overlay.className = 'confirm-overlay';
    overlay.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 10001;
    `;
    
    const dialog = document.createElement('div');
    dialog.className = 'confirm-dialog';
    dialog.style.cssText = `
        background: white;
        border-radius: 8px;
        padding: 24px;
        min-width: 300px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    `;
    
    dialog.innerHTML = `
        <div class="confirm-message" style="margin-bottom: 20px; font-size: 16px; color: #333;">${message}</div>
        <div class="confirm-actions" style="display: flex; justify-content: flex-end; gap: 12px;">
            <button class="btn-cancel" style="padding: 8px 16px; border: 1px solid #d9d9d9; background: white; color: #333; border-radius: 4px; cursor: pointer;">取消</button>
            <button class="btn-confirm" style="padding: 8px 16px; border: none; background: #1890ff; color: white; border-radius: 4px; cursor: pointer;">确定</button>
        </div>
    `;
    
    overlay.appendChild(dialog);
    document.body.appendChild(overlay);
    
    const btnCancel = dialog.querySelector('.btn-cancel');
    const btnConfirm = dialog.querySelector('.btn-confirm');
    
    const cleanup = () => {
        if (overlay.parentNode) {
            overlay.parentNode.removeChild(overlay);
        }
    };
    
    btnCancel.addEventListener('click', () => {
        cleanup();
        if (onCancel) onCancel();
    });
    
    btnConfirm.addEventListener('click', () => {
        cleanup();
        if (onConfirm) onConfirm();
    });
    
    overlay.addEventListener('click', (e) => {
        if (e.target === overlay) {
            cleanup();
            if (onCancel) onCancel();
        }
    });
}

// 加载中状态
function showLoading(container, message = '加载中...') {
    const loading = document.createElement('div');
    loading.className = 'loading-overlay';
    loading.innerHTML = `
        <div class="loading-spinner"></div>
        <div class="loading-text">${message}</div>
    `;
    loading.style.cssText = `
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(255, 255, 255, 0.8);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    `;
    
    const spinner = loading.querySelector('.loading-spinner');
    spinner.style.cssText = `
        width: 32px;
        height: 32px;
        border: 3px solid #f3f3f3;
        border-top: 3px solid #1890ff;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 12px;
    `;
    
    // 添加动画样式
    if (!document.querySelector('#loading-styles')) {
        const style = document.createElement('style');
        style.id = 'loading-styles';
        style.textContent = `
            @keyframes spin {
                0% { transform: rotate(0deg); }
                100% { transform: rotate(360deg); }
            }
        `;
        document.head.appendChild(style);
    }
    
    if (container) {
        container.style.position = 'relative';
        container.appendChild(loading);
    }
    
    return {
        hide: () => {
            if (loading.parentNode) {
                loading.parentNode.removeChild(loading);
            }
        }
    };
}

// 导出工具函数
export {
    formatDate,
    formatRelativeTime,
    debounce,
    throttle,
    deepClone,
    validatePhone,
    validateEmail,
    generateId,
    showMessage,
    showConfirm,
    showLoading
};