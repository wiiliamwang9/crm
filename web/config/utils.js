// 通用工具函数

// 获取API基础URL
function getApiBaseUrl() {
  if (window.GlobalApiConfig) {
    return window.GlobalApiConfig.BASE_URL;
  }
  // 回退到默认值
  return 'http://localhost:8081';
}

// 存储用户列表的全局变量
window.usersList = [];

// 加载用户列表
async function loadUsersList() {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/users/active`);
    if (response.ok) {
      const result = await response.json();
      window.usersList = result.data || [];
      console.log('用户列表加载成功:', window.usersList.length, '个用户');
      
      // 更新所有执行人选择器
      updateExecutorSelectors();
    } else {
      console.error('加载用户列表失败:', response.status);
      // 使用默认用户列表
      window.usersList = [
        { id: 1, name: '王欢欢' },
        { id: 2, name: '张三' },
        { id: 3, name: '李四' },
        { id: 4, name: '王五' }
      ];
      updateExecutorSelectors();
    }
  } catch (error) {
    console.error('Error loading users list:', error);
    // 使用默认用户列表
    window.usersList = [
      { id: 1, name: '王欢欢' },
      { id: 2, name: '张三' },
      { id: 3, name: '李四' },
      { id: 4, name: '王五' }
    ];
    updateExecutorSelectors();
  }
}

// 更新所有执行人选择器
function updateExecutorSelectors() {
  const selectors = [
    'executorPerson',
    'editExecutorPerson'
  ];
  
  selectors.forEach(selectorId => {
    const selector = document.getElementById(selectorId);
    if (selector && window.usersList.length > 0) {
      // 清空现有选项
      selector.innerHTML = '';
      
      // 添加用户选项
      window.usersList.forEach(user => {
        const option = document.createElement('option');
        option.value = user.id;
        option.textContent = user.name;
        selector.appendChild(option);
      });
      
      console.log(`执行人选择器 ${selectorId} 更新完成:`, window.usersList.length, '个选项');
    }
  });
}

// 导出函数到全局作用域
window.getApiBaseUrl = getApiBaseUrl;
window.loadUsersList = loadUsersList;
window.updateExecutorSelectors = updateExecutorSelectors;