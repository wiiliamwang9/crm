// 工具函数和辅助功能

// 获取API基础URL
function getApiBaseUrl() {
  return window.GlobalApiConfig ? window.GlobalApiConfig.BASE_URL : 'http://localhost:8081';
}

// 获取当前客户ID
function getCurrentCustomerId() {
  const urlParams = new URLSearchParams(window.location.search);
  const customerId = urlParams.get('customer_id');
  
  if (customerId) {
    return parseInt(customerId);
  }
  
  console.warn('未找到customer_id参数，使用默认客户ID: 92');
  return 92;
}

// 等待全局配置加载完成的辅助函数
function waitForGlobalConfig() {
  return new Promise((resolve) => {
    if (window.GlobalApiConfig) {
      resolve();
    } else {
      setTimeout(() => waitForGlobalConfig().then(resolve), 100);
    }
  });
}

// 存储用户列表的全局变量
window.usersList = [];

// 获取执行人名称
function getExecutorName(executorId) {
  if (!executorId) return '未知';
  
  const user = window.usersList.find(u => u.id === parseInt(executorId));
  if (user) {
    return user.name;
  }
  
  const executors = {
    1: '小朋友',
    2: '张三',
    3: '李四',
    4: '王五'
  };
  return executors[executorId] || '未知';
}

// 获取提醒方式名称
function getReminderMethodName(reminderType) {
  const methods = {
    'sms': '短信',
    'wechat': '企业微信',
    'email': '邮件'
  };
  return methods[reminderType] || '未设置';
}

// 格式化待办时间显示
function formatTodoTime(timeStr) {
  if (!timeStr) return '-';
  
  try {
    const date = new Date(timeStr);
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const todoDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    
    const diffDays = Math.floor((todoDate - today) / (1000 * 60 * 60 * 24));
    const timeFormat = date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
    
    if (diffDays === 0) {
      return `今天 ${timeFormat}`;
    } else if (diffDays === 1) {
      return `明天 ${timeFormat}`;
    } else if (diffDays === -1) {
      return `昨天 ${timeFormat}`;
    } else if (diffDays > 1 && diffDays <= 7) {
      return `${diffDays}天后 ${timeFormat}`;
    } else if (diffDays < -1 && diffDays >= -7) {
      return `${Math.abs(diffDays)}天前 ${timeFormat}`;
    } else {
      return date.toLocaleDateString('zh-CN') + ' ' + timeFormat;
    }
  } catch (error) {
    console.error('时间格式化错误:', error);
    return timeStr;
  }
}

// 返回上一页
function goBack() {
  if (document.referrer && document.referrer !== window.location.href) {
    window.history.back();
  } else {
    window.location.href = '../query/query_target.html';
  }
}

// 加载用户列表
async function loadUsersList() {
  try {
    // 先尝试active接口，如果失败则使用通用users接口
    console.log('Loading users from active endpoint');
    let response = await fetch(`${getApiBaseUrl()}/api/v1/users/active`);
    let result = null;
    let users = [];
    
    if (response.ok) {
      result = await response.json();
      users = result.data || [];
      console.log('Active users API response:', result);
    }
    
    // 如果active接口返回空数据，尝试通用users接口
    if (!users || users.length === 0) {
      console.log('Active users empty, trying general users endpoint');
      response = await fetch(`${getApiBaseUrl()}/api/v1/users`);
      if (response.ok) {
        result = await response.json();
        users = result.data || [];
        console.log('General users API response:', result);
        
        // 过滤出在职用户
        users = users.filter(user => user.status === '在职' || user.status === 'active');
      }
    }
    
    console.log('Final users loaded:', users);
    
    // 如果仍然没有数据，使用实际的用户数据作为后备
    if (!users || users.length === 0) {
      console.warn('API returned no users, using fallback users');
      users = [
        {id: 1, name: '王欢欢'},
        {id: 2, name: '冯敏勇'}
      ];
    }
    
    window.usersList = users;
    updateExecutorSelectors(users);
    
  } catch (error) {
    console.error('Error loading users:', error);
    // 使用实际的用户数据作为后备
    const fallbackUsers = [
      {id: 1, name: '王欢欢'},
      {id: 2, name: '冯敏勇'}
    ];
    window.usersList = fallbackUsers;
    updateExecutorSelectors(fallbackUsers);
  }
}

// 存储可搜索下拉框实例
window.searchableSelects = {};

// 获取执行人选择框的值
function getExecutorPersonValue() {
  const hiddenInput = document.getElementById('executorPerson');
  return hiddenInput ? hiddenInput.value : null;
}

// 获取编辑执行人选择框的值
function getEditExecutorPersonValue() {
  const hiddenInput = document.getElementById('editExecutorPerson');
  return hiddenInput ? hiddenInput.value : null;
}

// 设置执行人选择框的值
function setExecutorPersonValue(value, text) {
  if (window.searchableSelects.executorPerson) {
    window.searchableSelects.executorPerson.setValue(value, text);
  }
}

// 设置编辑执行人选择框的值
function setEditExecutorPersonValue(value, text) {
  if (window.searchableSelects.editExecutorPerson) {
    window.searchableSelects.editExecutorPerson.setValue(value, text);
  }
}

// 更新执行人选择框
function updateExecutorSelectors(users) {
  console.log('Updating executor selectors with users:', users);
  
  // 确保至少有一些用户数据
  if (!users || users.length === 0) {
    console.warn('No users provided to updateExecutorSelectors, using fallback');
    users = [
      {id: 1, name: '王欢欢'},
      {id: 2, name: '冯敏勇'}
    ];
  }
  
  // 初始化可搜索下拉框
  const selectorConfigs = [
    {
      containerId: 'executorPersonContainer',
      key: 'executorPerson'
    },
    {
      containerId: 'editExecutorPersonContainer', 
      key: 'editExecutorPerson'
    }
  ];
  
  selectorConfigs.forEach(config => {
    const container = document.getElementById(config.containerId);
    
    if (container) {
      // 如果已存在实例，更新数据
      if (window.searchableSelects[config.key]) {
        window.searchableSelects[config.key].setData(users);
        console.log(`Updated existing ${config.key} with ${users.length} users`);
      } else {
        // 创建新实例
        window.searchableSelects[config.key] = new SearchableSelect(config.containerId, {
          onChange: (value, text) => {
            console.log(`${config.key} changed:`, {value, text});
          }
        });
        window.searchableSelects[config.key].setData(users);
        console.log(`Created new ${config.key} with ${users.length} users`);
      }
    } else {
      console.warn(`Container ${config.containerId} not found in DOM`);
      // 尝试延迟初始化
      setTimeout(() => {
        const delayedContainer = document.getElementById(config.containerId);
        if (delayedContainer) {
          console.log(`Found delayed container ${config.containerId}, initializing now`);
          updateExecutorSelectors(users);
        }
      }, 1000);
    }
  });
}

// 加载客户信息
async function loadCustomerInfo() {
  try {
    const customerId = getCurrentCustomerId();
    const apiUrl = `${getApiBaseUrl()}/api/v1/customers/${customerId}`;
    console.log('Loading customer info from:', apiUrl);
    console.log('Customer ID:', customerId);
    
    const response = await fetch(apiUrl);
    
    if (response.ok) {
      const result = await response.json();
      const customer = result.data || result;
      console.log('Customer API response:', result);
      console.log('Customer data to display:', customer);
      
      updateCustomerDisplay(customer);
    } else {
      console.error('Failed to load customer info:', response.status, response.statusText);
      const errorText = await response.text();
      console.error('Error response:', errorText);
    }
  } catch (error) {
    console.error('Error loading customer info:', error);
  }
}

// 更新页面上的客户信息显示
function updateCustomerDisplay(customer) {
  console.log('updateCustomerDisplay called with:', customer);
  
  // 更新客户姓名
  const customerNameElements = document.querySelectorAll('.text-wrapper-2');
  console.log('Found customer name elements:', customerNameElements.length);
  customerNameElements.forEach((element, index) => {
    console.log(`Updating name element ${index}:`, element, 'with name:', customer.name);
    if (customer.name) {
      element.textContent = customer.name;
    }
  });
  
  // 更新客户标签
  const tagsContainer = document.querySelector('.view-7');
  if (tagsContainer && customer.tags) {
    tagsContainer.innerHTML = '';
    
    let tags = [];
    if (Array.isArray(customer.tags)) {
      tags = customer.tags;
    } else if (typeof customer.tags === 'string' && customer.tags.trim()) {
      tags = customer.tags.split(',').map(tag => tag.trim()).filter(tag => tag);
    }
    
    tags.forEach((tag, index) => {
      const tagElement = document.createElement('div');
      if (index === 0) {
        tagElement.className = 'div-wrapper';
      } else if (index === 1) {
        tagElement.className = 'frame-2';
      } else if (index === 2) {
        tagElement.className = 'frame-3';
      } else {
        tagElement.className = 'customer-tag';
      }
      
      if (index < 3) {
        tagElement.innerHTML = `<div class="text-wrapper-${index + 3}">${tag}</div>`;
      } else {
        tagElement.innerHTML = `<div class="tag-text">${tag}</div>`;
      }
      
      tagsContainer.appendChild(tagElement);
    });
    
    console.log('Customer tags updated:', tags);
  }
  
  // 更新销售字段
  const sellerElement = document.querySelector('.text-wrapper-7');
  console.log('Seller element found:', !!sellerElement);
  console.log('Customer saller_name:', customer.saller_name);
  if (sellerElement && customer.saller_name) {
    console.log('Updating seller element with:', customer.saller_name);
    sellerElement.textContent = customer.saller_name;
  } else {
    console.log('Not updating seller - element:', !!sellerElement, 'saller_name:', customer.saller_name);
  }
  
  // 更新客户头像
  const avatarElement = document.querySelector('.image');
  if (avatarElement && customer.name) {
    if (customer.avatar && customer.avatar.trim()) {
      avatarElement.src = customer.avatar;
      avatarElement.style.display = 'block';
      avatarElement.classList.remove('text-avatar');
    } else {
      const firstChar = customer.name.charAt(0);
      avatarElement.style.display = 'none';
      
      let textAvatar = document.querySelector('.text-avatar');
      if (!textAvatar) {
        textAvatar = document.createElement('div');
        textAvatar.className = 'text-avatar';
        avatarElement.parentNode.insertBefore(textAvatar, avatarElement);
      }
      textAvatar.textContent = firstChar;
    }
  }
  
  // 更新组织信息
  const organizationElement = document.querySelector('.text-wrapper-6');
  const organizationFrame = document.querySelector('.frame-5');
  
  if (organizationElement && organizationFrame) {
    if (customer.organization && customer.organization.trim()) {
      organizationElement.textContent = `组织：${customer.organization}`;
      organizationFrame.style.display = 'block';
    } else {
      organizationFrame.style.display = 'none';
    }
  }
  
  // 更新标题中的客户名
  const titleElements = document.querySelectorAll('.modal-title');
  titleElements.forEach(element => {
    const titleText = element.textContent;
    if (titleText.includes('"') && customer.name) {
      element.textContent = titleText.replace(/"[^"]*"/, `"${customer.name}"`);
    }
  });
  
  console.log('Customer display updated:', customer.name);
}

// 测试API端点
async function testApiEndpoints() {
  const config = window.GlobalApiConfig;
  const testUrls = [
    config.getUrl(config.ENDPOINTS.USERS.ACTIVE),
    config.getUrl(config.ENDPOINTS.CUSTOMERS.DETAIL, { id: 108 }),
    config.getUrl(config.ENDPOINTS.TODOS.LIST) + '?customer_id=108'
  ];
  
  console.log('=== API端点测试 ===');
  for (const url of testUrls) {
    try {
      const response = await fetch(url, { method: 'OPTIONS' });
      console.log(`OPTIONS ${url}:`, response.status, response.statusText);
    } catch (error) {
      console.log(`OPTIONS ${url} 失败:`, error.message);
    }
  }
  console.log('==================');
}