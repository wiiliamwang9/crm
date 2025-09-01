// 待办区域相关功能
class TodoArea {
  constructor() {
    this.init();
  }

  init() {
    this.loadTodoList();
    this.setupEventListeners();
  }

  setupEventListeners() {
    // 初始化待办tab为今日
    this.switchTodoTab('today');
  }

  // 切换待办状态tab
  switchTodoTab(status) {
    console.log('Switching todo tab to:', status);
    
    // 移除所有tab的active状态
    const todoTabs = document.querySelectorAll('.todo-tab');
    todoTabs.forEach(tab => {
      tab.classList.remove('active');
    });
    
    // 给当前选中的tab添加active状态
    const activeTab = document.querySelector(`[data-status="${status}"]`);
    if (activeTab) {
      activeTab.classList.add('active');
    }
    
    // 根据状态筛选和显示待办列表
    this.filterTodosByStatus(status);
  }

  // 根据状态筛选待办
  filterTodosByStatus(status) {
    if (!window.allTodos) {
      this.loadTodoList();
      return;
    }
    
    let filteredTodos = [];
    const now = new Date();
    
    switch(status) {
      case 'today':
        const todayTodos = window.allTodos.filter(todo => {
          const plannedDate = new Date(todo.planned_time);
          return plannedDate.toDateString() === now.toDateString() && todo.status !== 'completed';
        });
        
        if (todayTodos.length > 0) {
          filteredTodos = todayTodos;
        } else {
          filteredTodos = window.allTodos.filter(todo => {
            const plannedDate = new Date(todo.planned_time);
            return todo.status !== 'completed' && plannedDate > now;
          }).sort((a, b) => new Date(a.planned_time) - new Date(b.planned_time));
        }
        break;
      case 'all':
        filteredTodos = window.allTodos;
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
      default:
        filteredTodos = window.allTodos;
    }
    
    this.updateTodoListDisplay(filteredTodos);
  }

  // 加载待办列表
  async loadTodoList() {
    try {
      const response = await fetch(`${getApiBaseUrl()}/api/v1/todos?customer_id=${getCurrentCustomerId()}`);
      if (response.ok) {
        const result = await response.json();
        const todos = result.data || [];
        console.log('Todos loaded:', todos);
        
        window.allTodos = todos;
        // 默认显示今日待办
        this.filterTodosByStatus('today');
      }
    } catch (error) {
      console.error('Error loading todos:', error);
    }
  }

  // 更新待办列表显示
  updateTodoListDisplay(todos) {
    const todoContainer = document.getElementById('todo-list-container');
    if (!todoContainer) return;
    
    if (!window.allTodos && todos) {
      window.allTodos = todos;
    }
    
    todoContainer.innerHTML = '';
    
    if (!todos || todos.length === 0) {
      todoContainer.innerHTML = '<div class="no-todos">暂无待办事项</div>';
      return;
    }
    
    todos.forEach(todo => {
      const todoElement = this.createTodoElement(todo);
      todoContainer.appendChild(todoElement);
    });
  }

  // 创建待办元素
  createTodoElement(todo) {
    const todoDiv = document.createElement('div');
    todoDiv.className = todo.status === 'completed' ? 'frame-17' : 'frame-16';
    todoDiv.setAttribute('data-status', todo.status || 'pending');
    
    const timeDisplay = formatTodoTime(todo.planned_time);
    const executorName = getExecutorName(todo.executor_id);
    const creatorName = getExecutorName(todo.creator_id);
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
          <div class="text-wrapper-8">&nbsp;&nbsp; 执行人：</div>
          <img class="image-2" src="https://c.animaapp.com/mepb4kxjUWP5uf/img/---8.png" />
          <div class="text-wrapper-7">${executorName}</div>
        </div>
        <div class="text-wrapper-19">提醒方式：${reminderDisplay}</div>
        <div style="width: 100%; display: flex; justify-content: flex-end; gap: 10px; margin-top: 10px;">
          ${todo.status !== 'completed' ? `<div class="view-12" onclick="completeTodo(${todo.id})"><div class="text-wrapper-10">完成</div></div>` : ''}
          ${todo.status !== 'completed' ? `<div class="view-12" onclick="openChangeTimeModal(${todo.id})"><div class="text-wrapper-10">改时间</div></div>` : ''}
          <div class="view-12" onclick="openEditTodoModal(${todo.id})"><div class="text-wrapper-10">编辑</div></div>
        </div>
      </div>
    `;
    return todoDiv;
  }
}

// 全局函数，供HTML调用
window.switchTodoTab = function(status) {
  if (window.todoArea) {
    window.todoArea.switchTodoTab(status);
  }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
  window.todoArea = new TodoArea();
});