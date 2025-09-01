function openTodoModal() {
  const modal = document.getElementById('todoModal');
  if (modal) {
    // 设置标题为默认的"新增待办"
    const title = modal.querySelector('.modal-title');
    if (title) {
      title.textContent = '新增待办';
    }
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';
    
    // 初始化时间选择器
    initDatetimePicker('todoDatetimePicker');
    
    console.log('Todo modal opened');
  }
}

// 打开SOS求助弹窗（使用待办弹窗但修改标题）
function openSOSModal() {
  const modal = document.getElementById('todoModal');
  if (modal) {
    // 设置标题为"SOS 求助"
    const title = modal.querySelector('.modal-title');
    if (title) {
      title.textContent = 'SOS 求助';
    }
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';
    
    // 初始化时间选择器
    initDatetimePicker('todoDatetimePicker');
    
    console.log('SOS modal opened');
  }
}

// 关闭待办弹窗
function closeTodoModal() {
  const modal = document.getElementById('todoModal');
  if (modal) {
    modal.style.display = 'none';
    document.body.style.overflow = 'auto';
    
    // 清除编辑模式状态
    modal.removeAttribute('data-todo-id');
    
    // 清空表单
    resetTodoForm();
    console.log('Todo modal closed');
  }
}

// 重置表单
function resetTodoForm() {
  // 重置时间选择器为当前时间
  initDatetimePicker('todoDatetimePicker');
  const todoTypeEl = document.getElementById('todoType');
  if (todoTypeEl) todoTypeEl.value = 'call';
  const todoContentEl = document.getElementById('todoContent');
  if (todoContentEl) todoContentEl.value = '';
  const reminderSwitchEl = document.getElementById('reminderSwitch');
  if (reminderSwitchEl) reminderSwitchEl.checked = false;

  // 设置默认执行人（兼容可搜索下拉框）
  const users = Array.isArray(window.usersList) ? window.usersList : [];
  if (users.length > 0) {
    if (typeof setExecutorPersonValue === 'function') {
      setExecutorPersonValue(users[0].id, users[0].name);
    } else {
      const hiddenInput = document.getElementById('executorPerson');
      if (hiddenInput) hiddenInput.value = users[0].id;
    }
  } else {
    if (typeof setExecutorPersonValue === 'function') {
      setExecutorPersonValue('', '');
    } else {
      const hiddenInput = document.getElementById('executorPerson');
      if (hiddenInput) hiddenInput.value = '';
    }
  }

  const reminderMethodEl = document.getElementById('reminderMethod');
  if (reminderMethodEl) reminderMethodEl.value = 'wechat';
  // 隐藏提醒方式选择
  const reminderMethodRow = document.getElementById('reminderMethodRow');
  if (reminderMethodRow) reminderMethodRow.style.display = 'none';
}

// 创建或更新待办
async function createTodo() {
  // 从时间选择器获取选择的时间
  const plannedTime = getPickerDateTime('todoDatetimePicker');
  if (!plannedTime) {
    alert('请选择计划时间');
    return;
  }

  const isReminder = document.getElementById('reminderSwitch').checked;
  const reminderMethodElement = document.getElementById('reminderMethod');
  
  const todoData = {
    type: document.getElementById('todoType').value,
    content: document.getElementById('todoContent').value.trim(),
    reminder: isReminder,
    executorPerson: document.getElementById('executorPerson').value,
    reminderMethod: isReminder && reminderMethodElement ? reminderMethodElement.value : 'wechat',
    plannedTime: plannedTime
  };
  
  // 验证必填字段
  if (!todoData.content) {
    alert('请输入待办内容');
    return;
  }
  
  // 检查是否为编辑模式
  const modal = document.getElementById('todoModal');
  const todoId = modal.getAttribute('data-todo-id');
  const isEdit = !!todoId;
  
  console.log(isEdit ? 'Updating todo:' : 'Creating todo:', todoData);
  console.log('Reminder checked:', isReminder);
  if (isReminder) {
    console.log('Reminder method:', todoData.reminderMethod);
  }
  
  try {
    let todoPayload;
    let url;
    let method;
    
    if (isEdit) {
      // 编辑模式 - 只发送需要更新的字段
      todoPayload = {
        title: todoData.type === 'call' ? '电话回访' : todoData.type === 'visit' ? '上门拜访' : '其他待办',
        content: todoData.content,
        planned_time: plannedTime.toISOString(),
        is_reminder: todoData.reminder,
        reminder_type: todoData.reminder ? todoData.reminderMethod : null,
        reminder_user_id: todoData.reminder ? parseInt(todoData.executorPerson) : null,
        reminder_time: todoData.reminder ? plannedTime.toISOString() : null
      };
      url = `${getApiBaseUrl()}/api/v1/todos/${todoId}`;
      method = 'PUT';
    } else {
      // 创建模式 - 发送完整数据
      todoPayload = {
        customer_id: getCurrentCustomerId(), // 从URL参数获取客户ID
        executor_id: parseInt(todoData.executorPerson),
        title: todoData.type === 'call' ? '电话回访' : todoData.type === 'visit' ? '上门拜访' : '其他待办',
        content: todoData.content,
        planned_time: plannedTime.toISOString(),
        is_reminder: todoData.reminder,
        reminder_type: todoData.reminder ? todoData.reminderMethod : null,
        reminder_user_id: todoData.reminder ? parseInt(todoData.executorPerson) : null,
        reminder_time: todoData.reminder ? plannedTime.toISOString() : null,
        priority: 'medium'
      };
      url = `${getApiBaseUrl()}/api/v1/todos`;
      method = 'POST';
    }
    
    // 发送到后端API
    console.log('Sending request to:', url);
    console.log('Request method:', method);
    console.log('Request payload:', JSON.stringify(todoPayload, null, 2));
    
    const response = await fetch(url, {
      method: method,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(todoPayload)
    });
    
    if (response.ok) {
      let result = null;
      try {
        const contentType = response.headers.get('Content-Type') || '';
        if (response.status !== 204 && contentType.includes('application/json')) {
          result = await response.json();
        } else {
          const text = await response.text().catch(() => '');
          if (text && text.trim()) {
            try { result = JSON.parse(text); } catch (e) { /* non-JSON body */ }
          }
        }
      } catch (e) {
        console.warn('Successful response without JSON body for create/update todo.');
      }
      alert(isEdit ? '待办更新成功！' : '待办创建成功！');
      closeTodoModal();
      // 刷新待办列表
      loadTodoList();
    } else {
      const errorText = await response.text();
      console.error('API Error Response:', errorText);
      console.error('Request Payload:', JSON.stringify(todoPayload, null, 2));
      
      let errorMessage = errorText;
      try {
        const errorJson = JSON.parse(errorText);
        errorMessage = errorJson.message || errorJson.error || errorText;
      } catch (e) {
        // 如果不是JSON格式，使用原始文本
      }
      
      alert((isEdit ? '更新失败：' : '创建失败：') + errorMessage);
    }
  } catch (error) {
    console.error('Error saving todo:', error);
    alert(isEdit ? '更新失败，请稍后重试' : '创建失败，请稍后重试');
  }
}

// 加载待办列表
async function loadTodoList() {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/todos?customer_id=${getCurrentCustomerId()}`);
    if (response.ok) {
      const result = await response.json();
      const todos = result.data || [];
      console.log('Todos loaded:', todos);
      
      // 存储所有待办数据到全局变量
      window.allTodos = todos;
      
      // 默认显示今日待办
      filterTodosByStatus('today');
    }
  } catch (error) {
    console.error('Error loading todos:', error);
  }
}

// 更新待办列表显示
function updateTodoListDisplay(todos) {
  const todoContainer = document.getElementById('todo-list-container');
  if (!todoContainer) return;
  
  // 存储所有待办数据到全局变量
  if (!window.allTodos && todos) {
    window.allTodos = todos;
  }
  
  // 清空现有内容
  todoContainer.innerHTML = '';
  
  // 如果没有待办事项，显示提示信息
  if (!todos || todos.length === 0) {
    todoContainer.innerHTML = '<div class="no-todos">暂无待办事项</div>';
    return;
  }
  
  // 添加新的待办元素
  todos.forEach(todo => {
    const todoElement = createTodoElement(todo);
    todoContainer.appendChild(todoElement);
  });
}

// 存储用户列表的全局变量
window.usersList = [];

// 获取Executor名称
function getExecutorName(executorId) {
  if (!executorId) return '未知';
  
  // 从用户列表中查找
  const user = window.usersList.find(u => u.id === parseInt(executorId));
  if (user) {
    return user.name;
  }
  
  // 如果用户列表还没有加载，返回默认值
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

// 创建待办元素
function createTodoElement(todo) {
  const todoDiv = document.createElement('div');
  todoDiv.className = todo.status === 'completed' ? 'frame-17' : 'frame-16';
  todoDiv.setAttribute('data-status', todo.status || 'pending');
  
  // 格式化时间显示
  const timeDisplay = formatTodoTime(todo.planned_time);
  
  // 获取执行人名称
  const executorName = getExecutorName(todo.executor_id);
  
  // 获取创建人名称
  const creatorName = getExecutorName(todo.creator_id);
  
  // 格式化提醒方式显示
  const reminderDisplay = todo.is_reminder ? 
    getReminderMethodName(todo.reminder_type) : '-';
  
  todoDiv.innerHTML = `
    <div style="width: 100%;">
      <p class="p"><span class="span">待办时间：</span> <span class="text-wrapper-18">${timeDisplay}</span></p>
      <p class="div-2">
        <span class="span">待办内容：</span> <span class="text-wrapper-18">${todo.content}</span>
      </p>
      <div class="frame-6">
        <div class="text-wrapper-8">创建待办：</div>
        <img class="image-2" src="https://c.animaapp.com/mepb4kxjUWP5uf/img/---8.png" />
        <div class="text-wrapper-7">${creatorName}</div>
      </div>
      <div class="frame-6">
        <div class="text-wrapper-8">&nbsp;&nbsp; Executor：</div>
        <img class="image-2" src="https://c.animaapp.com/mepb4kxjUWP5uf/img/---8.png" />
        <div class="text-wrapper-7">${executorName}</div>
      </div>
      <div class="text-wrapper-19">提醒方式：${reminderDisplay}</div>
      <div style="width: 100%; display: flex; justify-content: flex-end; gap: 10px; margin-top: 10px;">
        ${todo.status !== 'completed' ? `<div class="view-12" onclick="completeTodo(${todo.id})"><div class="text-wrapper-10">完成</div></div>` : ''}
        <div class="view-12" onclick="openChangeTimeModal(${todo.id})"><div class="text-wrapper-10">改时间</div></div>
        <div class="view-12" onclick="openEditTodoModal(${todo.id})"><div class="text-wrapper-10">编辑</div></div>
      </div>
    </div>
  `;
  return todoDiv;
}

// 编辑待办
async function editTodo(todoId) {
  try {
    // 获取待办详情
    const response = await fetch(`http://localhost:8081/api/v1/todos/${todoId}`);
    if (response.ok) {
      const result = await response.json();
      const todo = result.data;
      
      // 填充表单数据
      document.getElementById('todoContent').value = todo.content;
      const plannedTime = new Date(todo.planned_time);
      setPickerDateTime('todoDatetimePicker', plannedTime);
      document.getElementById('reminderSwitch').checked = todo.is_reminder;
      
      // 设置编辑模式
      const modal = document.getElementById('todoModal');
      const title = modal.querySelector('.modal-title');
      if (title) {
        title.textContent = '编辑待办';
      }
      
      // 保存待办ID用于更新
      modal.setAttribute('data-todo-id', todoId);
      
      // 显示弹窗
      modal.style.display = 'flex';
      document.body.style.overflow = 'hidden';
      
      console.log('Editing todo:', todo);
    }
  } catch (error) {
    console.error('Error loading todo for edit:', error);
    alert('加载待办信息失败');
  }
}

// 删除待办
async function deleteTodo(todoId) {
  if (!confirm('确定要删除这个待办吗？')) {
    return;
  }
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/todos/${todoId}`, {
      method: 'DELETE'
    });
    
    if (response.ok) {
      alert('待办删除成功！');
      // 刷新待办列表
      loadTodoList();
    } else {
      const error = await response.text();
      alert('删除失败：' + error);
    }
  } catch (error) {
    console.error('Error deleting todo:', error);
    alert('删除失败，请稍后重试');
  }
}

// 完成待办
async function completeTodo(todoId) {
  if (!confirm('确定要完成这个待办吗？')) {
    return;
  }
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/todos/${todoId}/complete`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (response.ok) {
      alert('待办已完成！');
      // 刷新待办列表
      loadTodoList();
    } else {
      const error = await response.text();
      alert('完成失败：' + error);
    }
  } catch (error) {
    console.error('Error completing todo:', error);
    alert('完成失败，请稍后重试');
  }
}

// 删除当前编辑的待办
async function deleteCurrentTodo() {
  const modal = document.getElementById('editTodoModal');
  const todoId = modal ? modal.getAttribute('data-todo-id') : null;
  
  if (!todoId) {
    alert('无法获取待办ID');
    return;
  }
  
  if (!confirm('确定要删除这个待办吗？')) {
    return;
  }
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/todos/${todoId}`, {
      method: 'DELETE'
    });
    
    if (response.ok) {
      alert('待办删除成功！');
      closeEditTodoModal();
      // 刷新待办列表
      loadTodoList();
    } else {
      const error = await response.text();
      alert('删除失败：' + error);
    }
  } catch (error) {
    console.error('Error deleting todo:', error);
    alert('删除失败，请稍后重试');
  }
}

// switchTodoTab 统一由 js/todo-area.js 提供全局实现，这里删除重复定义

// 根据状态筛选待办
function filterTodosByStatus(status) {
  // 根据状态筛选待办数据并重新渲染
  if (!window.allTodos) {
    loadTodoList();
    return;
  }
  
  let filteredTodos = [];
  const now = new Date();
  
  switch(status) {
    case 'today':
      // 先筛选今日待办（未完成的）
      const todayTodos = window.allTodos.filter(todo => {
        const plannedDate = new Date(todo.planned_time);
        return plannedDate.toDateString() === now.toDateString() && todo.status !== 'completed';
      });
      
      // 如果今日有待办，则显示今日待办
      if (todayTodos.length > 0) {
        filteredTodos = todayTodos;
      } else {
        // 如果今日无待办，则显示未来的待办（所有未完成且时间在今天之后的）
        filteredTodos = window.allTodos.filter(todo => {
          const plannedDate = new Date(todo.planned_time);
          return todo.status !== 'completed' && plannedDate > now;
        }).sort((a, b) => new Date(a.planned_time) - new Date(b.planned_time)); // 按时间排序
      }
      break;
    case 'completed':
      filteredTodos = window.allTodos.filter(todo => todo.status === 'completed');
      break;
    case 'pending':
      filteredTodos = window.allTodos.filter(todo => todo.status === 'pending');
      break;
    case 'overdue':
      filteredTodos = window.allTodos.filter(todo => {
        const plannedDate = new Date(todo.planned_time);
        return todo.status === 'pending' && plannedDate < now;
      });
      break;
    case 'all':
      filteredTodos = window.allTodos;
      break;
    default:
      filteredTodos = window.allTodos;
  }
  
  updateTodoListDisplay(filteredTodos);
}

// 点击弹窗背景关闭
document.addEventListener('click', function(event) {
  const modal = document.getElementById('todoModal');
  if (event.target === modal) {
    closeTodoModal();
  }
});

// ESC键关闭弹窗
document.addEventListener('keydown', function(event) {
  if (event.key === 'Escape') {
    closeTodoModal();
    closeRecordModal();
    closePreferenceModal();
    closeTagModal();
    closeAddContactModal();
    closeEditTodoModal();
    closeChangeTimeModal();
  }
});

// 添加记录弹窗相关函数由 tabs.js 中的 FollowUpTab 负责实现，这里仅保留创建逻辑
function resetRecordForm() {
  document.getElementById('recordType').value = 'call';
  document.getElementById('recordContent').value = '';
  document.getElementById('createTodoSwitch').checked = false;
  document.getElementById('planTime').value = 'today';
  document.getElementById('planTimeRow').style.display = 'none';
}

async function createRecord() {
  const recordData = {
    customer_id: getCurrentCustomerId(),
    kind: document.getElementById('recordType').value,
    title: getRecordTitle(document.getElementById('recordType').value),
    content: document.getElementById('recordContent').value.trim(),
    result: '',
    amount: 0,
    cost: 0,
    feedback: '',
    satisfaction: 0,
    remark: '',
    create_todo: document.getElementById('createTodoSwitch').checked,
    todo_planned_time: null,
    todo_content: ''
  };
  
  // 验证必填字段
  if (!recordData.content) {
    alert('请输入记录内容');
    return;
  }
  
  // 处理待办事项相关
  if (recordData.create_todo) {
    const planTime = document.getElementById('planTime').value;
    recordData.todo_planned_time = calculatePlannedTime(planTime);
    recordData.todo_content = `跟进：${recordData.content}`;
  }
  
  console.log('Creating record:', recordData);
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/activities`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json; charset=utf-8'
      },
      body: JSON.stringify(recordData)
    });
    
    if (response.ok) {
      const result = await response.json();
      alert('跟进记录创建成功！' + (recordData.create_todo ? '\n同时已创建待办事项。' : ''));
      closeRecordModal();
      
      // 刷新跟进记录列表
      if (window.refreshFollowUpRecords) {
        refreshFollowUpRecords();
      }
    } else {
      const errorText = await response.text();
      console.error('创建跟进记录失败:', errorText);
      alert('创建失败，请稍后重试');
    }
  } catch (error) {
    console.error('Error creating record:', error);
    alert('网络错误，请稍后重试');
  }
}

// 获取记录标题
function getRecordTitle(kind) {
  const titleMap = {
    'call': '电话沟通',
    'visit': '实地拜访',
    'meeting': '会议洽谈',
    'order': '下单记录',
    'sample': '发样记录',
    'feedback': '客户反馈',
    'complaint': '客户投诉',
    'payment': '付款记录',
    'other': '其他跟进'
  };
  
  return titleMap[kind] || '跟进记录';
}

// 计算计划时间
function calculatePlannedTime(planTimeType) {
  const now = new Date();
  let plannedTime = new Date();
  
  switch (planTimeType) {
    case 'today':
      plannedTime.setHours(now.getHours() + 2); // 2小时后
      break;
    case 'tomorrow':
      plannedTime.setDate(now.getDate() + 1);
      plannedTime.setHours(9, 0, 0, 0); // 明天上午9点
      break;
    case 'this-week':
      plannedTime.setDate(now.getDate() + 3);
      plannedTime.setHours(9, 0, 0, 0); // 3天后上午9点
      break;
    case 'next-week':
      plannedTime.setDate(now.getDate() + 7);
      plannedTime.setHours(9, 0, 0, 0); // 一周后上午9点
      break;
    case 'this-month':
      plannedTime.setDate(now.getDate() + 14);
      plannedTime.setHours(9, 0, 0, 0); // 两周后上午9点
      break;
    default:
      plannedTime.setHours(now.getHours() + 2); // 默认2小时后
  }
  
  return plannedTime.toISOString();
}

// 监听创建待办开关的变化
document.addEventListener('DOMContentLoaded', function() {
  const createTodoSwitch = document.getElementById('createTodoSwitch');
  const planTimeRow = document.getElementById('planTimeRow');
  
  if (createTodoSwitch && planTimeRow) {
    createTodoSwitch.addEventListener('change', function() {
      if (this.checked) {
        planTimeRow.style.display = 'block';
      } else {
        planTimeRow.style.display = 'none';
      }
    });
  }
});

// 点击弹窗背景关闭记录弹窗
document.addEventListener('click', function(event) {
  const recordModal = document.getElementById('recordModal');
  if (event.target === recordModal) {
    closeRecordModal();
  }
  
  const preferenceModal = document.getElementById('preferenceModal');
  if (event.target === preferenceModal) {
    closePreferenceModal();
  }

  const tagModal = document.getElementById('tagModal');
  if (event.target === tagModal) {
    closeTagModal();
  }
  
  const addContactModal = document.getElementById('addContactModal');
  if (event.target === addContactModal) {
    closeAddContactModal();
  }
});

// 偏好弹窗相关函数由 tabs.js 中的 PreferencesTab 负责实现
// 添加偏好标签
function addPreferenceTag(text) {
  if (!text || text.trim() === '') return;
  
  const container = document.getElementById('preferenceTagsContainer');
  const colors = ['blue', 'green', 'yellow'];
  const randomColor = colors[Math.floor(Math.random() * colors.length)];
  
  const tag = document.createElement('span');
  tag.className = `preference-tag preference-tag-${randomColor}`;
  tag.textContent = text.trim();
  tag.onclick = function() {
    removePreferenceTag(this);
  };
  
  container.appendChild(tag);
  
  // 更新创建新产品按钮的文本
  updateCreateProductButton(text.trim());
}

// 删除偏好标签
function removePreferenceTag(tagElement) {
  tagElement.remove();
}

// 更新创建新产品按钮的文本
function updateCreateProductButton(productName) {
  const btn = document.getElementById('createProductBtn');
  btn.innerHTML = `创建新产品 "${productName}"`;
  btn.setAttribute('data-product-name', productName);
}

// 显示新产品输入框
function showNewProductInput() {
  const inputRow = document.getElementById('newProductInputRow');
  const input = document.getElementById('newProductInput');
  
  inputRow.style.display = 'block';
  input.focus();
}

// 创建新产品
function createNewProduct() {
  const productName = document.getElementById('newProductInput').value.trim();
  if (!productName) {
    alert('请输入产品名称');
    return;
  }
  
  // 添加到产品列表
  const productsList = document.querySelector('.products-list');
  const productItem = document.createElement('div');
  productItem.className = 'product-item';
  productItem.innerHTML = `
    <div class="product-name">${productName}:</div>
    <div class="product-preferences">
      <span class="product-preference">新产品</span>
    </div>
  `;
  
  productsList.appendChild(productItem);
  
  // 隐藏输入框并清空
  document.getElementById('newProductInputRow').style.display = 'none';
  document.getElementById('newProductInput').value = '';
  
  alert(`产品 "${productName}" 创建成功！`);
}

// 监听偏好搜索框回车事件
document.addEventListener('DOMContentLoaded', function() {
  const preferenceSearch = document.getElementById('preferenceSearch');
  if (preferenceSearch) {
    preferenceSearch.addEventListener('keypress', function(event) {
      if (event.key === 'Enter') {
        event.preventDefault();
        const text = this.value.trim();
        if (text) {
          addPreferenceTag(text);
          this.value = '';
        }
      }
    });
  }
  
  // 监听新产品输入框回车事件
  const newProductInput = document.getElementById('newProductInput');
  if (newProductInput) {
    newProductInput.addEventListener('keypress', function(event) {
      if (event.key === 'Enter') {
        event.preventDefault();
        createNewProduct();
      }
    });
  }
  
  // 监听创建新产品按钮点击
  const createProductBtn = document.getElementById('createProductBtn');
  if (createProductBtn) {
    createProductBtn.addEventListener('click', function() {
      const productName = this.getAttribute('data-product-name');
      if (productName) {
        // 如果有产品名称，直接创建
        document.getElementById('newProductInput').value = productName;
        createNewProduct();
      } else {
        // 否则显示输入框
        showNewProductInput();
      }
    });
  }
});

// 标签弹窗相关函数已移动到 tags.js

// 标签相关函数已移动到 tags.js，避免重复定义
/*
function clearTagSearch() {
  const searchInput = document.getElementById('tagSearchInput');
  if (searchInput) {
    searchInput.value = '';
    updateTagSections('');
    searchInput.focus();
  }
}

// 切换标签选中状态
function toggleTag(tagElement) {
  tagElement.classList.toggle('selected');
  
  // 这里可以添加保存选中状态的逻辑
  console.log('Tag toggled:', tagElement.textContent, tagElement.classList.contains('selected'));
}

// 添加新标签（选中现有标签）
function addNewTag() {
  const searchInput = document.getElementById('tagSearchInput');
  const tagText = searchInput.value.trim();
  
  if (!tagText) return;
  
  // 查找是否已存在该标签
  const allTags = document.querySelectorAll('#tagModal .tag');
  let existingTag = null;
  
  for (let tag of allTags) {
    if (tag.textContent.toLowerCase() === tagText.toLowerCase()) {
      existingTag = tag;
      break;
    }
  }
  
  if (existingTag) {
    // 如果标签已存在，则选中它
    existingTag.classList.add('selected');
    existingTag.scrollIntoView({ behavior: 'smooth', block: 'center' });
    
    // 清空搜索框
    searchInput.value = '';
    updateTagSections('');
    
    console.log('Existing tag selected:', tagText);
  } else {
    // 如果不存在，提示创建新标签
    alert('该标签不存在，请点击"创建新标签"来创建');
  }
}

// 创建新标签
function createNewTag() {
  const searchInput = document.getElementById('tagSearchInput');
  const tagText = searchInput.value.trim();
  
  if (!tagText) return;
  
  // 创建新标签元素
  const newTag = document.createElement('span');
  newTag.className = 'tag selected';
  newTag.textContent = tagText;
  newTag.onclick = function() { toggleTag(this); };
  
  // 添加到"其他"分类（如果没有则创建）
  let otherCategory = document.querySelector('.tag-category:last-child .tag-list');
  if (!otherCategory) {
    // 创建"其他"分类
    const allTagsSection = document.querySelector('.all-tags-section');
    const newCategory = document.createElement('div');
    newCategory.className = 'tag-category';
    newCategory.innerHTML = `
      <div class="category-title">其他：</div>
      <div class="tag-list"></div>
    `;
    allTagsSection.appendChild(newCategory);
    otherCategory = newCategory.querySelector('.tag-list');
  }
  
  otherCategory.appendChild(newTag);
  
  // 清空搜索框
  searchInput.value = '';
  updateTagSections('');
  
  // 滚动到新创建的标签
  newTag.scrollIntoView({ behavior: 'smooth', block: 'center' });
  
  console.log('New tag created:', tagText);
  alert(`标签"${tagText}"创建成功！`);
}

// 更新标签相关显示
function updateTagSections(searchText) {
  const addTagSection = document.getElementById('addTagSection');
  const createTagSection = document.getElementById('createTagSection');
  const addTagText = document.getElementById('addTagText');
  const createTagText = document.getElementById('createTagText');
  
  if (searchText.trim()) {
    addTagText.textContent = searchText;
    createTagText.textContent = `创建新标签"${searchText}"`;
    addTagSection.style.display = 'flex';
    createTagSection.style.display = 'flex';
  } else {
    addTagSection.style.display = 'none';
    createTagSection.style.display = 'none';
  }
}
*/

// 监听搜索框输入（已移动到 tags.js）
/*
document.addEventListener('DOMContentLoaded', function() {
  const tagSearchInput = document.getElementById('tagSearchInput');
  if (tagSearchInput) {
    tagSearchInput.addEventListener('input', function() {
      updateTagSections(this.value);
    });
    
    tagSearchInput.addEventListener('keypress', function(event) {
      if (event.key === 'Enter') {
        event.preventDefault();
        addNewTag();
      }
    });
  }
});
*/

// 添加熟人弹窗相关函数由 tabs.js 中的 PreferencesTab/ContactsTab 提供，这里仅保留业务逻辑
function createContact() {
  const notes = document.getElementById('contactNotes').value.trim();
  
  const contactData = {
    name: '徐小二',
    relationship: '好友',
    notes: notes,
    createTime: new Date().toISOString()
  };
  
  console.log('Creating contact:', contactData);
  
  // 这里可以添加保存逻辑，比如发送到后端API
  // 模拟保存成功
  alert('熟人关系创建成功！');
  closeAddContactModal();
  
  // 可以在这里更新页面上的熟人列表
  // addContactToList(contactData);
}

// 添加熟人弹窗背景点击关闭事件
document.addEventListener('DOMContentLoaded', function() {
  const addContactModal = document.getElementById('addContactModal');
  if (addContactModal) {
    addContactModal.addEventListener('click', function(e) {
      if (e.target === this) {
        closeAddContactModal();
      }
    });
  }
});

// 编辑待办弹窗相关函数
function openEditTodoModal(todoId) {
  const modal = document.getElementById('editTodoModal');
  if (modal) {
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';
    
    // 设置待办ID用于删除功能
    modal.setAttribute('data-todo-id', todoId);
    
    // 初始化时间选择器
    initDatetimePicker('editTodoDatetimePicker');
    
    // 根据todoId加载现有数据
    loadTodoData(todoId);
    
    console.log('Edit todo modal opened for:', todoId);
  }
}

function closeEditTodoModal() {
  const modal = document.getElementById('editTodoModal');
  if (modal) {
    modal.style.display = 'none';
    document.body.style.overflow = 'auto';
    
    // 清空表单
    resetEditTodoForm();
    console.log('Edit todo modal closed');
  }
}

// 加载待办数据到编辑表单
async function loadTodoData(todoId) {
  try {
    const response = await fetch(`http://localhost:8081/api/v1/todos/${todoId}`);
    if (response.ok) {
      const result = await response.json();
      const todo = result.data;
      
      // 填充表单数据
      if (todo.planned_time) {
        const plannedDateTime = new Date(todo.planned_time);
        setPickerDateTime('editTodoDatetimePicker', plannedDateTime);
      }
      document.getElementById('editTodoType').value = todo.type || 'call';
      document.getElementById('editTodoContent').value = todo.content || '';
      document.getElementById('editReminderSwitch').checked = todo.is_reminder || false;
      document.getElementById('editExecutorPerson').value = todo.executor_id || '1';
      document.getElementById('editReminderMethod').value = todo.reminder_type || 'wechat';
      
      // 根据提醒状态显示/隐藏提醒方式
      const editReminderMethodRow = document.getElementById('editReminderMethodRow');
      if (todo.is_reminder) {
        editReminderMethodRow.style.display = 'block';
      } else {
        editReminderMethodRow.style.display = 'none';
      }
      
      // 存储todoId用于保存时使用
      document.getElementById('editTodoModal').setAttribute('data-todo-id', todoId);
      
      console.log('Todo data loaded:', todo);
    } else {
      console.error('Failed to load todo data');
      alert('加载待办数据失败');
      closeEditTodoModal();
    }
  } catch (error) {
    console.error('Error loading todo data:', error);
    alert('加载待办数据失败');
    closeEditTodoModal();
  }
}

// 格式化时间用于编辑表单
function formatTimeForEdit(plannedTime) {
  if (!plannedTime) return 'today';
  
  const planned = new Date(plannedTime);
  const today = new Date();
  const tomorrow = new Date(today);
  tomorrow.setDate(today.getDate() + 1);
  
  // 比较日期（忽略时间）
  const plannedDate = planned.toDateString();
  const todayDate = today.toDateString();
  const tomorrowDate = tomorrow.toDateString();
  
  if (plannedDate === todayDate) {
    return 'today';
  } else if (plannedDate === tomorrowDate) {
    return 'tomorrow';
  } else {
    return 'next-week'; // 默认为下周
  }
}

// 重置编辑表单
function resetEditTodoForm() {
  // 重置时间选择器为当前时间
  initDatetimePicker('editTodoDatetimePicker');
  const typeEl = document.getElementById('editTodoType');
  if (typeEl) typeEl.value = 'call';
  const contentEl = document.getElementById('editTodoContent');
  if (contentEl) contentEl.value = '';
  const switchEl = document.getElementById('editReminderSwitch');
  if (switchEl) switchEl.checked = false;

  // 设置默认执行人（兼容可搜索下拉框）
  const users = Array.isArray(window.usersList) ? window.usersList : [];
  if (users.length > 0) {
    if (typeof setEditExecutorPersonValue === 'function') {
      setEditExecutorPersonValue(users[0].id, users[0].name);
    } else {
      const hidden = document.getElementById('editExecutorPerson');
      if (hidden) hidden.value = users[0].id;
    }
  } else {
    if (typeof setEditExecutorPersonValue === 'function') {
      setEditExecutorPersonValue('', '');
    } else {
      const hidden = document.getElementById('editExecutorPerson');
      if (hidden) hidden.value = '';
    }
  }
  
  const methodEl = document.getElementById('editReminderMethod');
  if (methodEl) methodEl.value = 'wechat';
  const rowEl = document.getElementById('editReminderMethodRow');
  if (rowEl) rowEl.style.display = 'none';
  const modalEl = document.getElementById('editTodoModal');
  if (modalEl) modalEl.removeAttribute('data-todo-id');
}

// 保存编辑的待办
async function saveTodo() {
  const modal = document.getElementById('editTodoModal');
  const todoId = modal.getAttribute('data-todo-id');
  
  if (!todoId) {
    alert('无法获取待办ID');
    return;
  }
  
  // 从时间选择器获取选择的时间
  const plannedTime = getPickerDateTime('editTodoDatetimePicker');
  if (!plannedTime) {
    alert('请选择计划时间');
    return;
  }
  
  const content = document.getElementById('editTodoContent').value.trim();
  
  // 验证必填字段
  if (!content) {
    alert('请输入待办内容');
    return;
  }
  
  const isReminder = document.getElementById('editReminderSwitch').checked;
  const executorId = parseInt(document.getElementById('editExecutorPerson').value);
  const reminderMethod = document.getElementById('editReminderMethod').value;
  
  const todoData = {
    type: document.getElementById('editTodoType').value,
    content: content,
    planned_time: plannedTime.toISOString(),
    is_reminder: isReminder,
    reminder_type: isReminder ? reminderMethod : null,
    reminder_user_id: isReminder ? executorId : null,
    reminder_time: isReminder ? plannedTime.toISOString() : null,
    executor_id: executorId
  };
  
  console.log('Saving todo:', todoData);
  
  try {
    // 发送到后端API
    const response = await fetch(`${getApiBaseUrl()}/api/v1/todos/${todoId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(todoData)
    });
    
    if (response.ok) {
      let result = null;
      try {
        const contentType = response.headers.get('Content-Type') || '';
        if (response.status !== 204 && contentType.includes('application/json')) {
          result = await response.json();
        } else {
          const text = await response.text().catch(() => '');
          if (text && text.trim()) {
            try { result = JSON.parse(text); } catch (e) { /* non-JSON body */ }
          }
        }
      } catch (e) {
        console.warn('Successful response without JSON body for save todo.');
      }
      alert('待办保存成功！');
      closeEditTodoModal();
      // 刷新待办列表
      loadTodoList();
    } else {
      const error = await response.text();
      alert('保存失败：' + error);
    }
  } catch (error) {
    console.error('Error saving todo:', error);
    alert('保存失败，请稍后重试');
  }
}

// 编辑待办弹窗背景点击关闭事件
document.addEventListener('DOMContentLoaded', function() {
  const editTodoModal = document.getElementById('editTodoModal');
  if (editTodoModal) {
    editTodoModal.addEventListener('click', function(e) {
      if (e.target === this) {
        closeEditTodoModal();
      }
    });
  }
  
  // 改时间弹窗背景点击关闭事件
  const changeTimeModal = document.getElementById('changeTimeModal');
  if (changeTimeModal) {
    changeTimeModal.addEventListener('click', function(e) {
      if (e.target === this) {
        closeChangeTimeModal();
      }
    });
  }
});

// 更多操作下拉菜单相关函数
function toggleMoreActions() {
  const dropdown = document.getElementById('moreActionsDropdown');
  if (dropdown) {
    const isVisible = dropdown.style.display === 'block';
    if (isVisible) {
      dropdown.style.display = 'none';
    } else {
      // 先隐藏其他可能打开的下拉菜单
      hideAllDropdowns();
      
      // 获取三个点按钮的位置
      const moreButton = document.querySelector('.more-actions-container .frame');
      if (moreButton) {
        const rect = moreButton.getBoundingClientRect();
        dropdown.style.left = (rect.right - 120) + 'px'; // 120是下拉菜单的min-width
        dropdown.style.top = (rect.bottom + 5) + 'px'; // 5px间距
      }
      
      dropdown.style.display = 'block';
    }
    console.log('More actions dropdown toggled:', !isVisible);
  }
}

// 隐藏所有下拉菜单
function hideAllDropdowns() {
  const dropdown = document.getElementById('moreActionsDropdown');
  if (dropdown) {
    dropdown.style.display = 'none';
  }
}

// 打电话功能
function makePhoneCall() {
  // 隐藏下拉菜单
  hideAllDropdowns();
  
  // 这里可以添加实际的打电话逻辑
  // 比如使用 tel: 协议调用系统拨号器
  // window.location.href = 'tel:' + phoneNumber;
  
  console.log('Making phone call...');
  alert('正在拨打电话...');
}

// 点击页面其他地方关闭下拉菜单
document.addEventListener('click', function(event) {
  const moreActionsContainer = document.querySelector('.more-actions-container');
  if (moreActionsContainer && !moreActionsContainer.contains(event.target)) {
    hideAllDropdowns();
  }
});

// 滚动时间选择器相关函数

// 创建通用的时间选择器HTML结构
function createDatetimePickerHTML(containerId) {
  // 根据容器ID生成对应的子选择器ID前缀
  let idPrefix;
  if (containerId === 'todoDatetimePicker') {
    idPrefix = '';  // 新增待办使用原来的ID（没有前缀）
  } else if (containerId === 'editTodoDatetimePicker') {
    idPrefix = 'edit';  // 编辑待办使用edit前缀
  } else if (containerId === 'changeTimeDatetimePicker') {
    idPrefix = 'changeTime';  // 改时间使用changeTime前缀
  } else {
    idPrefix = containerId.replace('DatetimePicker', '').replace('Picker', '');
  }
  
  const yearId = idPrefix ? `${idPrefix}YearPicker` : 'yearPicker';
  const monthId = idPrefix ? `${idPrefix}MonthPicker` : 'monthPicker';  
  const dayId = idPrefix ? `${idPrefix}DayPicker` : 'dayPicker';
  const hourId = idPrefix ? `${idPrefix}HourPicker` : 'hourPicker';
  
  return `
    <div class="picker-container">
      <div class="picker-column">
        <div class="picker-header">年</div>
        <div class="picker-scroll" id="${yearId}">
          <!-- 年份选项将通过JavaScript生成 -->
        </div>
      </div>
      <div class="picker-column">
        <div class="picker-header">月</div>
        <div class="picker-scroll" id="${monthId}">
          <!-- 月份选项将通过JavaScript生成 -->
        </div>
      </div>
      <div class="picker-column">
        <div class="picker-header">日</div>
        <div class="picker-scroll" id="${dayId}">
          <!-- 日期选项将通过JavaScript生成 -->
        </div>
      </div>
      <div class="picker-column">
        <div class="picker-header">时</div>
        <div class="picker-scroll" id="${hourId}">
          <!-- 小时选项将通过JavaScript生成 -->
        </div>
      </div>
    </div>
  `;
}

// 初始化时间选择器（确保HTML结构正确）
function ensureDatetimePickerStructure(containerId) {
  const container = document.getElementById(containerId);
  if (!container) return false;
  
  // 检查是否已有正确的结构
  const pickerContainer = container.querySelector('.picker-container');
  if (!pickerContainer || !pickerContainer.querySelector('[id*="YearPicker"], [id*="yearPicker"]')) {
    // 重新创建结构
    container.innerHTML = createDatetimePickerHTML(containerId);
  }
  return true;
}

function initDatetimePicker(containerId) {
  // 首先确保HTML结构正确
  if (!ensureDatetimePickerStructure(containerId)) return;
  
  const container = document.getElementById(containerId);
  
  const yearPicker = container.querySelector('[id*="YearPicker"], [id*="yearPicker"]');
  const monthPicker = container.querySelector('[id*="MonthPicker"], [id*="monthPicker"]');
  const dayPicker = container.querySelector('[id*="DayPicker"], [id*="dayPicker"]');
  const hourPicker = container.querySelector('[id*="HourPicker"], [id*="hourPicker"]');

  // 获取当前时间
  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth() + 1;
  const currentDay = now.getDate();
  const currentHour = now.getHours();

  // 初始化年份选择器（当前年份前后各5年）
  if (yearPicker) {
    yearPicker.innerHTML = '';
    for (let year = currentYear - 5; year <= currentYear + 5; year++) {
      const option = createPickerOption(year, year === currentYear);
      yearPicker.appendChild(option);
    }
    scrollToSelected(yearPicker);
  }

  // 初始化月份选择器
  if (monthPicker) {
    monthPicker.innerHTML = '';
    for (let month = 1; month <= 12; month++) {
      const option = createPickerOption(month + '月', month === currentMonth);
      option.dataset.value = month;
      monthPicker.appendChild(option);
    }
    scrollToSelected(monthPicker);
  }

  // 初始化日期选择器
  if (dayPicker) {
    updateDayPicker(dayPicker, currentYear, currentMonth, currentDay);
  }

  // 初始化小时选择器
  if (hourPicker) {
    hourPicker.innerHTML = '';
    for (let hour = 0; hour <= 23; hour++) {
      const option = createPickerOption(hour.toString().padStart(2, '0') + '时', hour === currentHour);
      option.dataset.value = hour;
      hourPicker.appendChild(option);
    }
    scrollToSelected(hourPicker);
  }

  // 添加事件监听器
  if (yearPicker && monthPicker && dayPicker) {
    yearPicker.addEventListener('click', function(e) {
      if (e.target.classList.contains('picker-option')) {
        selectPickerOption(yearPicker, e.target);
        const selectedYear = parseInt(e.target.textContent);
        const selectedMonth = getSelectedValue(monthPicker);
        const selectedDay = getSelectedValue(dayPicker);
        updateDayPicker(dayPicker, selectedYear, selectedMonth, selectedDay);
      }
    });

    monthPicker.addEventListener('click', function(e) {
      if (e.target.classList.contains('picker-option')) {
        selectPickerOption(monthPicker, e.target);
        const selectedYear = getSelectedValue(yearPicker);
        const selectedMonth = parseInt(e.target.dataset.value);
        const selectedDay = getSelectedValue(dayPicker);
        updateDayPicker(dayPicker, selectedYear, selectedMonth, selectedDay);
      }
    });
  }

  // 为日期和小时添加点击事件
  if (dayPicker) {
    dayPicker.addEventListener('click', function(e) {
      if (e.target.classList.contains('picker-option')) {
        selectPickerOption(dayPicker, e.target);
      }
    });
  }

  if (hourPicker) {
    hourPicker.addEventListener('click', function(e) {
      if (e.target.classList.contains('picker-option')) {
        selectPickerOption(hourPicker, e.target);
      }
    });
  }
}

// 创建选择器选项
function createPickerOption(text, selected = false) {
  const option = document.createElement('div');
  option.className = 'picker-option' + (selected ? ' selected' : '');
  option.textContent = text;
  if (typeof text === 'number') {
    option.dataset.value = text;
  }
  return option;
}

// 选择选项
function selectPickerOption(picker, option) {
  // 检查picker是否存在
  if (!picker || !option) return;
  
  // 移除其他选中状态
  picker.querySelectorAll('.picker-option').forEach(opt => {
    opt.classList.remove('selected');
  });
  // 选中当前选项
  option.classList.add('selected');
  // 滚动到选中项
  scrollToSelected(picker);
}

// 滚动到选中项
function scrollToSelected(picker) {
  // 检查picker是否存在
  if (!picker) return;
  
  const selected = picker.querySelector('.picker-option.selected');
  if (selected) {
    const pickerHeight = picker.clientHeight;
    const optionHeight = selected.clientHeight;
    const optionTop = selected.offsetTop;
    const scrollTop = optionTop - (pickerHeight / 2) + (optionHeight / 2);
    picker.scrollTop = scrollTop;
  }
}

// 获取选中的值
function getSelectedValue(picker) {
  // 检查picker是否存在
  if (!picker) return null;
  
  const selected = picker.querySelector('.picker-option.selected');
  if (selected) {
    return selected.dataset.value ? parseInt(selected.dataset.value) : parseInt(selected.textContent);
  }
  return null;
}

// 更新日期选择器
function updateDayPicker(dayPicker, year, month, selectedDay) {
  const daysInMonth = new Date(year, month, 0).getDate();
  const validSelectedDay = selectedDay && selectedDay <= daysInMonth ? selectedDay : 1;
  
  dayPicker.innerHTML = '';
  for (let day = 1; day <= daysInMonth; day++) {
    const option = createPickerOption(day + '日', day === validSelectedDay);
    option.dataset.value = day;
    dayPicker.appendChild(option);
  }
  scrollToSelected(dayPicker);
}

// 获取选择器的完整日期时间
function getPickerDateTime(containerId) {
  const container = document.getElementById(containerId);
  if (!container) return null;

  const yearPicker = container.querySelector('[id*="YearPicker"], [id*="yearPicker"]');
  const monthPicker = container.querySelector('[id*="MonthPicker"], [id*="monthPicker"]');
  const dayPicker = container.querySelector('[id*="DayPicker"], [id*="dayPicker"]');
  const hourPicker = container.querySelector('[id*="HourPicker"], [id*="hourPicker"]');

  const year = getSelectedValue(yearPicker);
  const month = getSelectedValue(monthPicker);
  const day = getSelectedValue(dayPicker);
  const hour = getSelectedValue(hourPicker);

  if (year && month && day && hour !== null) {
    return new Date(year, month - 1, day, hour, 0, 0);
  }
  return null;
}

// 打开改时间弹窗
async function openChangeTimeModal(todoId) {
  const modal = document.getElementById('changeTimeModal');
  if (!modal) return;
  
  try {
    // 获取待办详情
    const response = await fetch(`http://localhost:8081/api/v1/todos/${todoId}`);
    if (response.ok) {
      const result = await response.json();
      const todo = result.data;
      
      // 显示当前时间
      const currentTimeDisplay = document.getElementById('currentTimeDisplay');
      if (currentTimeDisplay) {
        currentTimeDisplay.textContent = formatTodoTime(todo.planned_time);
      }
      
      // 初始化时间选择器
      initDatetimePicker('changeTimeDatetimePicker');
      
      // 设置当前时间为初始值
      if (todo.planned_time) {
        const currentTime = new Date(todo.planned_time);
        setPickerDateTime('changeTimeDatetimePicker', currentTime);
      }
      
      // 保存todoId用于后续保存
      modal.setAttribute('data-todo-id', todoId);
      
      // 显示弹窗
      modal.style.display = 'flex';
      document.body.style.overflow = 'hidden';
      
      console.log('Change time modal opened for todo:', todoId);
    } else {
      alert('获取待办信息失败');
    }
  } catch (error) {
    console.error('Error loading todo for time change:', error);
    alert('获取待办信息失败');
  }
}

// 关闭改时间弹窗
function closeChangeTimeModal() {
  const modal = document.getElementById('changeTimeModal');
  if (modal) {
    modal.style.display = 'none';
    document.body.style.overflow = 'auto';
    modal.removeAttribute('data-todo-id');
    console.log('Change time modal closed');
  }
}

// 保存新时间
async function saveNewTime() {
  const modal = document.getElementById('changeTimeModal');
  const todoId = modal.getAttribute('data-todo-id');
  
  if (!todoId) {
    alert('无法获取待办ID');
    return;
  }
  
  // 从时间选择器获取新时间
  const newTime = getPickerDateTime('changeTimeDatetimePicker');
  if (!newTime) {
    alert('请选择新的时间');
    return;
  }
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/todos/${todoId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        planned_time: newTime.toISOString()
      })
    });
    
    if (response.ok) {
      alert('时间修改成功！');
      closeChangeTimeModal();
      // 刷新待办列表
      loadTodoList();
    } else {
      const error = await response.text();
      alert('时间修改失败：' + error);
    }
  } catch (error) {
    console.error('Error changing todo time:', error);
    alert('时间修改失败，请稍后重试');
  }
}

// 设置选择器的日期时间
function setPickerDateTime(containerId, datetime) {
  const container = document.getElementById(containerId);
  if (!container || !datetime) return;

  const year = datetime.getFullYear();
  const month = datetime.getMonth() + 1;
  const day = datetime.getDate();
  const hour = datetime.getHours();

  const yearPicker = container.querySelector('[id*="YearPicker"], [id*="yearPicker"]');
  const monthPicker = container.querySelector('[id*="MonthPicker"], [id*="monthPicker"]');
  const dayPicker = container.querySelector('[id*="DayPicker"], [id*="dayPicker"]');
  const hourPicker = container.querySelector('[id*="HourPicker"], [id*="hourPicker"]');

  // 选中对应的年份
  if (yearPicker) {
    const yearOption = Array.from(yearPicker.querySelectorAll('.picker-option')).find(opt => 
      parseInt(opt.textContent) === year
    );
    if (yearOption) {
      selectPickerOption(yearPicker, yearOption);
    }
  }

  // 选中对应的月份
  if (monthPicker) {
    const monthOption = Array.from(monthPicker.querySelectorAll('.picker-option')).find(opt => 
      parseInt(opt.dataset.value) === month
    );
    if (monthOption) {
      selectPickerOption(monthPicker, monthOption);
    }
  }

  // 更新并选中对应的日期
  if (dayPicker) {
    updateDayPicker(dayPicker, year, month, day);
  }

  // 选中对应的小时
  if (hourPicker) {
    const hourOption = Array.from(hourPicker.querySelectorAll('.picker-option')).find(opt => 
      parseInt(opt.dataset.value) === hour
    );
    if (hourOption) {
      selectPickerOption(hourPicker, hourOption);
    }
  }
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
      
      // 更新页面上的客户信息
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
  
  // 更新客户标签 - 第二行
  const tagsContainer = document.querySelector('.view-7');
  if (tagsContainer && customer.tags) {
    // 清空现有标签
    tagsContainer.innerHTML = '';
    
    // 处理tags字段 - 支持数组和逗号分隔的字符串
    let tags = [];
    if (Array.isArray(customer.tags)) {
      tags = customer.tags;
    } else if (typeof customer.tags === 'string' && customer.tags.trim()) {
      tags = customer.tags.split(',').map(tag => tag.trim()).filter(tag => tag);
    }
    
    // 创建标签元素，使用与原有标签相同的结构
    tags.forEach((tag, index) => {
      const tagElement = document.createElement('div');
      // 使用与原有标签相同的类名结构
      if (index === 0) {
        tagElement.className = 'div-wrapper';
      } else if (index === 1) {
        tagElement.className = 'frame-2';
      } else if (index === 2) {
        tagElement.className = 'frame-3';
      } else {
        tagElement.className = 'customer-tag';
      }
      
      // 使用与原有标签相同的内部结构
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
      // 有头像图片，显示图片
      avatarElement.src = customer.avatar;
      avatarElement.style.display = 'block';
      avatarElement.classList.remove('text-avatar');
    } else {
      // 没有头像图片，显示文字头像
      const firstChar = customer.name.charAt(0);
      avatarElement.style.display = 'none';
      
      // 创建或更新文字头像
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
  const organizationFrame = document.querySelector('.frame-5'); // 包含组织信息的父容器
  
  if (organizationElement && organizationFrame) {
    if (customer.organization && customer.organization.trim()) {
      // 有组织数据，显示组织信息
      organizationElement.textContent = `组织：${customer.organization}`;
      organizationFrame.style.display = 'block';
    } else {
      // 没有组织数据，隐藏整个组织字段
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

// 注意：loadUsersList 和 updateExecutorSelectors 函数已移动到 utils.js 中

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

// 返回上一页
function goBack() {
  // 如果有历史记录，返回上一页
  if (document.referrer && document.referrer !== window.location.href) {
    window.history.back();
  } else {
    // 如果没有历史记录，跳转到搜索页面
    window.location.href = '../query/query_target.html';
  }
}