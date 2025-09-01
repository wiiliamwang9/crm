// ==================== 全局API配置 ====================
// 统一管理前端的API配置，避免不同页面配置不一致的问题

// 环境配置定义
const ENVIRONMENTS = {
  // 本地开发环境配置
  LOCAL: {
    API_BASE_URL: 'http://127.0.0.1:8081',
    ENV_NAME: 'LOCAL',
    TIMEOUT: 30000 // 30秒超时
  },
  
  // 服务器环境配置
  SERVER: {
    API_BASE_URL: 'https://static.lamdar.cn:9501/crm',
    ENV_NAME: 'SERVER',
    TIMEOUT: 30000 // 30秒超时
  }
};

// 当前使用的环境配置
// 切换环境：修改这里的配置来切换不同环境
// const CURRENT_ENV = ENVIRONMENTS.LOCAL;        // 使用本地环境
const CURRENT_ENV = ENVIRONMENTS.SERVER;    // 使用服务器环境

// 导出的全局配置对象
window.GlobalApiConfig = {
  // 基础配置
  BASE_URL: CURRENT_ENV.API_BASE_URL,
  ENV_NAME: CURRENT_ENV.ENV_NAME,
  TIMEOUT: CURRENT_ENV.TIMEOUT,
  
  // API端点定义
  ENDPOINTS: {
    // 客户相关API
    CUSTOMERS: {
      LIST: '/api/v1/customers',
      DETAIL: '/api/v1/customers/{id}',
      SEARCH: '/api/v1/customers/search'
    },
    
    // 待办相关API
    TODOS: {
      LIST: '/api/v1/todos',
      DETAIL: '/api/v1/todos/{id}',
      CREATE: '/api/v1/todos',
      UPDATE: '/api/v1/todos/{id}',
      DELETE: '/api/v1/todos/{id}',
      COMPLETE: '/api/v1/todos/{id}/complete',
      CANCEL: '/api/v1/todos/{id}/cancel'
    },
    
    // 用户相关API
    USERS: {
      LIST: '/api/v1/users',
      DETAIL: '/api/v1/users/{id}',
      ACTIVE: '/api/v1/users/active'
    },
    
    // 提醒相关API
    REMINDERS: {
      LIST: '/api/v1/reminders',
      DETAIL: '/api/v1/reminders/{id}',
      CREATE: '/api/v1/reminders',
      UPDATE: '/api/v1/reminders/{id}',
      DELETE: '/api/v1/reminders/{id}'
    }
  },
  
  // 获取完整的API URL
  getUrl: function(endpoint, params = {}) {
    let url = this.BASE_URL + endpoint;
    
    // 替换路径参数，如 {id} -> 实际ID值
    for (const [key, value] of Object.entries(params)) {
      url = url.replace(`{${key}}`, value);
    }
    
    return url;
  },
  
  // 获取基础URL
  getBaseUrl: function() {
    return this.BASE_URL;
  },
  
  // 获取环境信息
  getEnvInfo: function() {
    return {
      name: this.ENV_NAME,
      baseUrl: this.BASE_URL,
      timeout: this.TIMEOUT
    };
  },
  
  // 测试API连接
  testConnection: async function() {
    try {
      const response = await fetch(this.BASE_URL + '/api/v1/users/active', {
        method: 'OPTIONS',
        timeout: this.TIMEOUT
      });
      return {
        success: response.ok,
        status: response.status,
        statusText: response.statusText
      };
    } catch (error) {
      return {
        success: false,
        error: error.message
      };
    }
  }
};

// 为了向后兼容，保持原有的API_CONFIG对象
window.API_CONFIG = {
  BASE_URL: window.GlobalApiConfig.BASE_URL,
  ENV_NAME: window.GlobalApiConfig.ENV_NAME,
  TIMEOUT: window.GlobalApiConfig.TIMEOUT
};

// 为了向后兼容，保持原有的getApiBaseUrl函数
window.getApiBaseUrl = function() {
  return window.GlobalApiConfig.BASE_URL;
};

// 控制台输出配置信息
console.log('=== 全局API配置加载完成 ===');
console.log('当前环境:', window.GlobalApiConfig.ENV_NAME);
console.log('API基础URL:', window.GlobalApiConfig.BASE_URL);
console.log('超时时间:', window.GlobalApiConfig.TIMEOUT + 'ms');
console.log('==============================');