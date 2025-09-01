// Homepage filter functionality
class FilterManager {
  constructor() {
    this.init();
  }

  init() {
    this.setupFilterSections();
    this.setupEventListeners();
  }

  setupFilterSections() {
    // Get all filter sections
    this.filterSections = document.querySelectorAll('.view-13');
    
    // Initialize each filter section
    this.filterSections.forEach((section, index) => {
      this.initializeFilterSection(section, index);
    });
  }

  initializeFilterSection(section, sectionIndex) {
    const filterOptions = section.querySelectorAll('.filter-option');
    
    // Add data attributes for identification
    section.setAttribute('data-section', sectionIndex);
    
    // Add click handlers to all filter options
    filterOptions.forEach((option, optionIndex) => {
      option.addEventListener('click', (e) => {
        this.switchFilter(section, option, e);
      });
    });
  }

  switchFilter(section, clickedOption, event) {
    event.preventDefault();
    
    // Get the filter value
    const filterValue = clickedOption.getAttribute('data-filter');
    
    // Don't switch if clicking the same option
    if (clickedOption.classList.contains('active')) return;
    
    // Add transition class
    section.classList.add('switching');
    
    // Remove active class from all options in this section
    const allOptions = section.querySelectorAll('.filter-option');
    allOptions.forEach(option => {
      option.classList.remove('active');
      // Reset text color for inactive options
      const textWrapper = option.querySelector('.text-wrapper-11, .text-wrapper-12');
      if (textWrapper) {
        textWrapper.className = 'text-wrapper-12';
      }
    });
    
    // Add active class to clicked option
    clickedOption.classList.add('active');
    
    // Update text wrapper class for active option
    const activeTextWrapper = clickedOption.querySelector('.text-wrapper-11, .text-wrapper-12');
    if (activeTextWrapper) {
      activeTextWrapper.className = 'text-wrapper-11';
    }
    
    // Remove transition class after animation
    setTimeout(() => {
      section.classList.remove('switching');
      
      // Trigger content update
      this.updateContent(section, filterValue);
    }, 150);
  }

  updateContent(section, selectedFilter) {
    // Get the section index to determine which content to update
    const sectionIndex = parseInt(section.getAttribute('data-section'));
    
    // Find the corresponding content area
    const contentArea = section.parentElement.querySelector('.view-14');
    
    if (contentArea) {
      // Add loading state
      contentArea.classList.add('loading');
      
      // Simulate content loading
      setTimeout(() => {
        // Remove loading state
        contentArea.classList.remove('loading');
        
        // You can add specific content updates here based on the filter
        console.log(`Updated content for section ${sectionIndex} with filter: ${selectedFilter}`);
      }, 300);
    }
  }
}

// Search functionality
class SearchManager {
  constructor() {
    this.init();
  }

  init() {
    this.searchInput = document.querySelector('.view-8');
    this.searchText = document.querySelector('.text-wrapper-7');
    
    if (this.searchInput) {
      this.setupSearch();
    }
  }

  setupSearch() {
    this.searchInput.style.cursor = 'pointer';
    
    this.searchInput.addEventListener('click', () => {
      this.activateSearch();
    });
  }

  activateSearch() {
    // Create actual input element
    const input = document.createElement('input');
    input.type = 'text';
    input.placeholder = '搜索客户...';
    input.className = 'search-input';
    
    // Style the input
    Object.assign(input.style, {
      border: 'none',
      outline: 'none',
      background: 'transparent',
      color: 'var(--lightlinecolor-border-4)',
      fontSize: '14px',
      fontFamily: '"Inter", Helvetica',
      width: '100%',
      padding: '0'
    });
    
    // Replace placeholder text with input
    this.searchText.style.display = 'none';
    this.searchInput.appendChild(input);
    
    // Focus the input
    input.focus();
    
    // Handle blur event
    input.addEventListener('blur', () => {
      if (!input.value.trim()) {
        this.deactivateSearch(input);
      }
    });
    
    // Handle enter key
    input.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') {
        this.performSearch(input.value);
      }
    });
  }

  deactivateSearch(input) {
    input.remove();
    this.searchText.style.display = 'block';
  }

  performSearch(query) {
    console.log('Searching for:', query);
    window.location.href = '../query/query_target.html';
  }

  activateSearch() {
    window.location.href = '../query/query_target.html';
  }
}

// Action buttons functionality
class ActionManager {
  constructor() {
    this.init();
  }

  init() {
    this.setupActionButtons();
  }

  setupActionButtons() {
    // WeChat buttons
    const wechatButtons = document.querySelectorAll('.frame-6');
    wechatButtons.forEach(button => {
      button.style.cursor = 'pointer';
      button.addEventListener('click', (e) => {
        e.preventDefault();
        this.handleWechatAction(button);
      });
    });

    // More action buttons
    const moreButtons = document.querySelectorAll('.frame-7');
    moreButtons.forEach(button => {
      button.style.cursor = 'pointer';
      button.addEventListener('click', (e) => {
        e.preventDefault();
        this.handleMoreActions(button);
      });
    });

    // View all buttons
    const viewAllButtons = document.querySelectorAll('.view-12');
    viewAllButtons.forEach(button => {
      button.style.cursor = 'pointer';
      button.addEventListener('click', (e) => {
        e.preventDefault();
        this.handleViewAll(button);
      });
    });
  }

  handleWechatAction(button) {
    // Add click animation
    button.style.transform = 'scale(0.95)';
    setTimeout(() => {
      button.style.transform = 'scale(1)';
    }, 150);
    
    console.log('WeChat action triggered');
    // Implement WeChat functionality here
  }

  handleMoreActions(button) {
    // Add click animation
    button.style.transform = 'scale(0.95)';
    setTimeout(() => {
      button.style.transform = 'scale(1)';
    }, 150);
    
    console.log('More actions triggered');
    // Implement more actions menu here
  }

  handleViewAll(button) {
    console.log('View all triggered');
    // Implement navigation to full list
  }
}

// User data management
class UserManager {
  constructor() {
    this.currentUser = null;
    this.init();
  }

  async init() {
    await this.loadUserData();
  }

  async loadUserData() {
    try {
      const response = await fetch(window.GlobalApiConfig.getUrl('/api/v1/users/3'));
      if (response.ok) {
        const result = await response.json();
        this.currentUser = result.data;
        this.updateUserDisplay();
      } else {
        console.error('Failed to load user data:', response.status);
      }
    } catch (error) {
      console.error('Error loading user data:', error);
    }
  }

  updateUserDisplay() {
    const userNameElement = document.querySelector('.text-wrapper-3');
    if (userNameElement && this.currentUser) {
      userNameElement.textContent = this.currentUser.name;
    }
  }
}

// Today's follow-up data manager
class TodayFollowUpManager {
  constructor() {
    this.currentFilter = '全部';
    this.todayCustomers = [];
    this.init();
  }

  init() {
    this.setupStatusFilters();
    this.setupViewAllButton();
    this.loadTodayFollowUpData();
  }

  setupViewAllButton() {
    const viewAllButton = document.querySelector('.view-12');
    if (viewAllButton) {
      viewAllButton.style.cursor = 'pointer';
      viewAllButton.addEventListener('click', () => {
        this.showAllCustomers();
      });
    }
  }

  async showAllCustomers() {
    try {
      // 获取所有今日待办数据
      const todoResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/todos?date_type=today&page=1&page_size=1000'), {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      if (todoResponse.ok) {
        const todoResult = await todoResponse.json();
        const todayTodos = todoResult.data || [];
        
        // 获取所有客户详细信息
        await this.loadAllCustomersWithTodos(todayTodos);
        this.filterAndDisplayCustomers();
      }
    } catch (error) {
      console.error('Error loading all customers:', error);
    }
  }

  async loadAllCustomersWithTodos(todos) {
    const customerIds = [...new Set(todos.map(todo => todo.customer_id))];
    const customerPromises = customerIds.map(async (customerId) => {
      try {
        const response = await fetch(window.GlobalApiConfig.getUrl(`/api/v1/customers/${customerId}`));
        if (response.ok) {
          const result = await response.json();
          const customer = result.data;
          
          // 为客户添加今日待办事项
          customer.todayTodos = todos.filter(todo => todo.customer_id === customerId);
          return customer;
        }
      } catch (error) {
        console.error(`Error loading customer ${customerId}:`, error);
      }
      return null;
    });

    const customers = await Promise.all(customerPromises);
    this.todayCustomers = customers.filter(customer => customer !== null);
  }

  setupStatusFilters() {
    const statusSection = document.querySelector('.view-13');
    if (!statusSection) return;

    const statusButtons = statusSection.querySelectorAll('.text-wrapper-11, .text-wrapper-12');
    statusButtons.forEach(button => {
      button.style.cursor = 'pointer';
      button.addEventListener('click', (e) => {
        this.handleStatusFilter(e.target.textContent);
      });
    });
  }

  async loadTodayFollowUpData() {
    try {
      // 获取今日待办数据
      const todoResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/todos?date_type=today&page=1&page_size=100'), {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      if (todoResponse.ok) {
        const todoResult = await todoResponse.json();
        const todayTodos = todoResult.data || [];
        
        // 获取客户详细信息
        await this.loadCustomersWithTodos(todayTodos);
        this.filterAndDisplayCustomers();
      }
    } catch (error) {
      console.error('Error loading today follow-up data:', error);
    }
  }

  async loadCustomersWithTodos(todos) {
    const customerIds = [...new Set(todos.map(todo => todo.customer_id))];
    const customerPromises = customerIds.slice(0, 3).map(async (customerId) => {
      try {
        const response = await fetch(window.GlobalApiConfig.getUrl(`/api/v1/customers/${customerId}`));
        if (response.ok) {
          const result = await response.json();
          const customer = result.data;
          
          // 为客户添加今日待办事项
          customer.todayTodos = todos.filter(todo => todo.customer_id === customerId);
          return customer;
        }
      } catch (error) {
        console.error(`Error loading customer ${customerId}:`, error);
      }
      return null;
    });

    const customers = await Promise.all(customerPromises);
    this.todayCustomers = customers.filter(customer => customer !== null);
  }

  handleStatusFilter(status) {
    this.currentFilter = status;
    this.updateFilterButtons(status);
    this.filterAndDisplayCustomers();
  }

  updateFilterButtons(activeStatus) {
    const statusSection = document.querySelector('.view-13');
    if (!statusSection) return;

    const allButtons = statusSection.querySelectorAll('.text-wrapper-11, .text-wrapper-12');
    allButtons.forEach(button => {
      // 清除所有内联样式，让CSS类来控制样式
      button.style.backgroundColor = '';
      button.style.color = '';
      
      // 使用CSS类而不是内联样式
      const parentOption = button.closest('.filter-option');
      if (parentOption) {
        if (button.textContent === activeStatus) {
          parentOption.classList.add('active');
          button.className = 'text-wrapper-11';
        } else {
          parentOption.classList.remove('active');
          button.className = 'text-wrapper-12';
        }
      }
    });
  }

  filterAndDisplayCustomers() {
    let filteredCustomers = this.todayCustomers;

    if (this.currentFilter !== '全部') {
      filteredCustomers = this.todayCustomers.filter(customer => {
        return this.getCustomerTodoStatus(customer) === this.currentFilter;
      });
    }

    this.displayCustomers(filteredCustomers);
  }

  getCustomerTodoStatus(customer) {
    if (!customer.todayTodos || customer.todayTodos.length === 0) {
      return '待办';
    }

    const latestTodo = customer.todayTodos.sort((a, b) => new Date(b.planned_time) - new Date(a.planned_time))[0];
    
    // 根据待办内容和状态推断客户状态
    const todoContent = (latestTodo.title || latestTodo.content || '').toLowerCase();
    
    if (latestTodo.status === 'completed') {
      if (todoContent.includes('发样') || todoContent.includes('样品')) {
        return '已发样';
      } else if (todoContent.includes('发货') || todoContent.includes('物流')) {
        return '已发货';
      }
    }
    
    if (todoContent.includes('定期') || todoContent.includes('回访')) {
      return '定期';
    } else if (todoContent.includes('发样') || todoContent.includes('样品')) {
      return '已发样';
    } else if (todoContent.includes('发货') || todoContent.includes('物流')) {
      return '已发货';
    }
    
    return '待办';
  }

  displayCustomers(customers) {
    const customerList = document.getElementById('today-customers-list');
    if (!customerList) return;

    customerList.innerHTML = '';

    // 更新计数器
    const countElement = document.getElementById('today-count');
    if (countElement) {
      countElement.textContent = customers.length;
    }

    if (customers.length === 0) {
      customerList.innerHTML = '<div style="text-align: center; padding: 20px; color: #999;">暂无今日待跟进客户</div>';
      return;
    }

    // 只显示前3个客户
    customers.slice(0, 3).forEach(customer => {
      const customerItem = this.createCustomerItem(customer);
      customerList.appendChild(customerItem);
    });
  }

  createCustomerItem(customer) {
    const div = document.createElement('div');
    div.className = 'view-15';
    div.style.cursor = 'pointer';
    
    const customerCategory = customer.category || '茶叶店';
    const customerTags = this.formatCustomerTags(customer.tags);
    const todayTodoContent = this.getTodayTodoContent(customer);
    const avatarHtml = this.createCustomerAvatarDisplay(customer);

    div.innerHTML = `
      <div class="view-16">
        ${avatarHtml}
        <div class="view-17">
          <div class="view-18">
            <div class="text-wrapper-13">${customer.name || '未知客户'}</div>
            <div class="text-wrapper-14">|</div>
            <div class="text-wrapper-15">${customerCategory}</div>
            <img class="img-3" src="https://c.animaapp.com/mewiuc6y0Nbcsg/img/clock.svg" />
          </div>
          <div class="view-19">
            ${customerTags}
          </div>
          <div class="text-wrapper-19">${todayTodoContent}</div>
        </div>
      </div>
      <div class="view-20">
        <div class="frame-6" onclick="this.handleWechatAction(${customer.id})">
          <img class="img-3" src="../statics/wechat.svg" />
          <div class="text-wrapper-6">发微信</div>
        </div>
        <img class="frame-7" src="../statics/black_1_more-one.svg" onclick="this.handleMoreActions(${customer.id})" />
      </div>
    `;

    // 添加点击整个卡片跳转到客户详情
    div.addEventListener('click', (e) => {
      if (!e.target.closest('.view-20')) {
        window.location.href = `../target/target.html?customer_id=${customer.id}`;
      }
    });

    return div;
  }

  createCustomerAvatarDisplay(customer) {
    if (customer.name) {
      const firstChar = customer.name.charAt(0);
      return `<div class="image-3" style="width: 48px; height: 48px; border-radius: 50%; background: #1890ff; color: white; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: bold;">${firstChar}</div>`;
    } else {
      return `<div class="image-3" style="width: 48px; height: 48px; border-radius: 50%; background: #ccc; color: white; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: bold;">?</div>`;
    }
  }

  formatCustomerTags(tags) {
    const tagElements = [];
    
    // 处理数组形式的标签
    if (tags && Array.isArray(tags) && tags.length > 0) {
      tags.slice(0, 2).forEach(tag => {
        tagElements.push(`<div class="frame-3"><div class="text-wrapper-16">${tag}</div></div>`);
      });
    }
    
    // 如果没有标签或标签不足，添加默认标签
    if (tagElements.length === 0) {
      tagElements.push('<div class="frame-4"><div class="text-wrapper-17">普通客户</div></div>');
    }
    
    // 限制最多显示3个标签
    return tagElements.slice(0, 3).join('');
  }

  getTodayTodoContent(customer) {
    if (!customer.todayTodos || customer.todayTodos.length === 0) {
      return '暂无今日待办';
    }
    
    const latestTodo = customer.todayTodos.sort((a, b) => new Date(b.planned_time) - new Date(a.planned_time))[0];
    const content = latestTodo.title || latestTodo.content || '暂无待办内容';
    
    // 限制显示长度
    return content.length > 20 ? content.substring(0, 20) + '...' : content;
  }

  getStatusClass(status) {
    const statusClasses = {
      '待办': 'status-pending',
      '定期': 'status-scheduled',
      '已发样': 'status-sample',
      '已发货': 'status-shipped'
    };
    return statusClasses[status] || '';
  }

  handleStatusFilter(filterText) {
    this.currentFilter = filterText;
    
    // 更新筛选按钮状态
    const statusSection = document.querySelector('.view-13');
    if (statusSection) {
      const filterOptions = statusSection.querySelectorAll('.filter-option');
      filterOptions.forEach(option => {
        option.classList.remove('active');
        if (option.textContent.trim() === filterText) {
          option.classList.add('active');
        }
      });
    }
    
    this.filterAndDisplayCustomers();
  }
}

// 近期待跟进管理器
class UpcomingFollowUpManager {
  constructor() {
    this.currentFilter = '全部';
    this.upcomingCustomers = [];
    this.init();
  }

  init() {
    this.setupStatusFilters();
    this.setupViewAllButton();
    this.loadUpcomingFollowUpData();
  }

  setupViewAllButton() {
    // 查找近期待跟进的查看全部按钮（第二个.view-12）
    const viewAllButtons = document.querySelectorAll('.view-12');
    const upcomingViewAllButton = viewAllButtons[1]; // 第二个是近期待跟进的
    if (upcomingViewAllButton) {
      upcomingViewAllButton.style.cursor = 'pointer';
      upcomingViewAllButton.addEventListener('click', () => {
        this.showAllCustomers();
      });
    }
  }

  setupStatusFilters() {
    // 查找近期待跟进的筛选区域（第二个.view-13）
    const statusSections = document.querySelectorAll('.view-13');
    const upcomingStatusSection = statusSections[1]; // 第二个是近期待跟进的
    if (!upcomingStatusSection) return;

    const statusButtons = upcomingStatusSection.querySelectorAll('.text-wrapper-11, .text-wrapper-12');
    statusButtons.forEach(button => {
      button.style.cursor = 'pointer';
      button.addEventListener('click', (e) => {
        this.handleStatusFilter(e.target.textContent);
      });
    });
  }

  async loadUpcomingFollowUpData() {
    try {
      const filter = this.currentFilter;

      if (filter === '全部') {
        // 加载近期（未来7天）的待办数据
        const tomorrow = new Date();
        tomorrow.setDate(tomorrow.getDate() + 1);
        const nextWeek = new Date();
        nextWeek.setDate(nextWeek.getDate() + 7);
        
        const todoResponse = await fetch(window.GlobalApiConfig.getUrl(`/api/v1/todos?status=pending&start_date=${tomorrow.toISOString()}&end_date=${nextWeek.toISOString()}&page=1&page_size=100`), {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (todoResponse.ok) {
          const todoResult = await todoResponse.json();
          const upcomingTodos = todoResult.data || [];
          await this.loadCustomersWithTodos(upcomingTodos);
        }
      } else if (filter === '半年未下单') {
        // 加载半年未下单客户
        const customerResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/customers/special?type=no_order_half_year&page=1&page_size=100'), {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (customerResponse.ok) {
          const customerResult = await customerResponse.json();
          this.upcomingCustomers = (customerResult.data.customers || []).map(customer => {
            customer.upcomingTodos = [];
            customer.noTodoMessage = '未添加待办！';
            return customer;
          });
        }
      } else if (filter === '一直未下单') {
        // 加载一直未下单客户
        const customerResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/customers/special?type=never_ordered&page=1&page_size=100'), {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (customerResponse.ok) {
          const customerResult = await customerResponse.json();
          this.upcomingCustomers = (customerResult.data.customers || []).map(customer => {
            customer.upcomingTodos = [];
            customer.noTodoMessage = '未添加待办！';
            return customer;
          });
        }
      }

      this.displayCustomers(this.upcomingCustomers);
    } catch (error) {
      console.error('Error loading upcoming follow-up data:', error);
      // 显示错误信息
      const upcomingCustomerList = document.getElementById('upcoming-customers-list');
      if (upcomingCustomerList) {
        upcomingCustomerList.innerHTML = '<div style="text-align: center; padding: 20px; color: #f53f3f;">加载数据失败，请刷新重试</div>';
      }
    }
  }

  async loadCustomersWithTodos(todos) {
    const customerIds = [...new Set(todos.map(todo => todo.customer_id))];
    const customerPromises = customerIds.slice(0, 3).map(async (customerId) => {
      try {
        const response = await fetch(window.GlobalApiConfig.getUrl(`/api/v1/customers/${customerId}`));
        if (response.ok) {
          const result = await response.json();
          const customer = result.data;
          
          // 为客户添加近期待办事项
          customer.upcomingTodos = todos.filter(todo => todo.customer_id === customerId);
          return customer;
        }
      } catch (error) {
        console.error(`Error loading customer ${customerId}:`, error);
      }
      return null;
    });

    const customers = await Promise.all(customerPromises);
    this.upcomingCustomers = customers.filter(customer => customer !== null);
  }

  async showAllCustomers() {
    try {
      const filter = this.currentFilter;
      
      if (filter === '全部') {
        // 获取所有待办数据
        const todoResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/todos?date_type=all&page=1&page_size=1000'), {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (todoResponse.ok) {
          const todoResult = await todoResponse.json();
          const upcomingTodos = todoResult.data || [];
          await this.loadAllCustomersWithTodos(upcomingTodos);
        }
      } else if (filter === '半年未下单') {
        const customerResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/customers/special?type=no_order_half_year&page=1&page_size=1000'), {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (customerResponse.ok) {
          const customerResult = await customerResponse.json();
          this.upcomingCustomers = (customerResult.data.customers || []).map(customer => {
            customer.upcomingTodos = [];
            customer.noTodoMessage = '未添加待办！';
            return customer;
          });
        }
      } else if (filter === '一直未下单') {
        const customerResponse = await fetch(window.GlobalApiConfig.getUrl('/api/v1/customers/special?type=never_ordered&page=1&page_size=1000'), {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });

        if (customerResponse.ok) {
          const customerResult = await customerResponse.json();
          this.upcomingCustomers = (customerResult.data.customers || []).map(customer => {
            customer.upcomingTodos = [];
            customer.noTodoMessage = '未添加待办！';
            return customer;
          });
        }
      }
      
      this.displayCustomers(this.upcomingCustomers);
    } catch (error) {
      console.error('Error loading all upcoming customers:', error);
    }
  }

  async loadAllCustomersWithTodos(todos) {
    const customerIds = [...new Set(todos.map(todo => todo.customer_id))];
    const customerPromises = customerIds.map(async (customerId) => {
      try {
        const response = await fetch(window.GlobalApiConfig.getUrl(`/api/v1/customers/${customerId}`));
        if (response.ok) {
          const result = await response.json();
          const customer = result.data;
          
          // 为客户添加近期待办事项
          customer.upcomingTodos = todos.filter(todo => todo.customer_id === customerId);
          return customer;
        }
      } catch (error) {
        console.error(`Error loading customer ${customerId}:`, error);
      }
      return null;
    });

    const customers = await Promise.all(customerPromises);
    this.upcomingCustomers = customers.filter(customer => customer !== null);
  }

  handleStatusFilter(status) {
    this.currentFilter = status;
    this.updateFilterButtons(status);
    this.loadUpcomingFollowUpData();
  }

  updateFilterButtons(activeStatus) {
    const statusSections = document.querySelectorAll('.view-13');
    const upcomingStatusSection = statusSections[1]; // 第二个是近期待跟进的
    if (!upcomingStatusSection) return;

    const allButtons = upcomingStatusSection.querySelectorAll('.text-wrapper-11, .text-wrapper-12');
    allButtons.forEach(button => {
      // 清除所有内联样式，让CSS类来控制样式
      button.style.backgroundColor = '';
      button.style.color = '';
      
      // 使用CSS类而不是内联样式
      const parentOption = button.closest('.filter-option');
      if (parentOption) {
        if (button.textContent === activeStatus) {
          parentOption.classList.add('active');
          button.className = 'text-wrapper-11';
        } else {
          parentOption.classList.remove('active');
          button.className = 'text-wrapper-12';
        }
      }
    });
  }

  displayCustomers(customers) {
    const upcomingCustomerList = document.getElementById('upcoming-customers-list');
    if (!upcomingCustomerList) return;

    upcomingCustomerList.innerHTML = '';

    // 更新计数器
    const countElement = document.getElementById('upcoming-count');
    if (countElement) {
      countElement.textContent = customers.length;
    }

    if (customers.length === 0) {
      upcomingCustomerList.innerHTML = '<div style="text-align: center; padding: 20px; color: #999;">暂无近期待跟进客户</div>';
      return;
    }

    // 只显示前3个客户
    customers.slice(0, 3).forEach(customer => {
      const customerItem = this.createCustomerItem(customer);
      upcomingCustomerList.appendChild(customerItem);
    });
  }

  createCustomerItem(customer) {
    const div = document.createElement('div');
    div.className = 'view-15';
    
    const customerCategory = customer.category || '茶叶店';
    const customerTags = this.formatCustomerTags(customer.tags);
    const upcomingTodoContent = this.getUpcomingTodoContent(customer);
    const avatarHtml = this.createCustomerAvatarDisplay(customer);

    div.innerHTML = `
      <div class="view-16">
        ${avatarHtml}
        <div class="view-17">
          <div class="view-18">
            <div class="text-wrapper-13">${customer.name || '未知客户'}</div>
            <div class="text-wrapper-14">|</div>
            <div class="text-wrapper-15">${customerCategory}</div>
            <img class="img-3" src="https://c.animaapp.com/mewiuc6y0Nbcsg/img/clock.svg" />
          </div>
          <div class="view-19">
            ${customerTags}
          </div>
          <div class="text-wrapper-19">${upcomingTodoContent}</div>
        </div>
      </div>
      <div class="view-20">
        <div class="frame-6">
          <img class="img-3" src="../statics/wechat.svg" />
          <div class="text-wrapper-6">发微信</div>
        </div>
        <img class="frame-7" src="../statics/black_1_more-one.svg" />
      </div>
    `;

    return div;
  }

  createCustomerAvatarDisplay(customer) {
    if (customer.name) {
      const firstChar = customer.name.charAt(0);
      return `<div class="image-3" style="width: 48px; height: 48px; border-radius: 50%; background: #1890ff; color: white; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: bold;">${firstChar}</div>`;
    } else {
      return `<div class="image-3" style="width: 48px; height: 48px; border-radius: 50%; background: #ccc; color: white; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: bold;">?</div>`;
    }
  }

  formatCustomerTags(tags) {
    const tagElements = [];
    
    // 处理数组形式的标签
    if (tags && Array.isArray(tags) && tags.length > 0) {
      tags.slice(0, 2).forEach(tag => {
        tagElements.push(`<div class="frame-3"><div class="text-wrapper-16">${tag}</div></div>`);
      });
    }
    
    // 如果没有标签或标签不足，添加默认标签
    if (tagElements.length === 0) {
      tagElements.push('<div class="frame-4"><div class="text-wrapper-17">普通客户</div></div>');
    }
    
    // 限制最多显示3个标签
    return tagElements.slice(0, 3).join('');
  }

  getUpcomingTodoContent(customer) {
    // 如果是特殊客户（半年未下单、一直未下单）且没有待办
    if (customer.noTodoMessage) {
      return customer.noTodoMessage;
    }
    
    if (!customer.upcomingTodos || customer.upcomingTodos.length === 0) {
      return '暂无近期待办';
    }
    
    const latestTodo = customer.upcomingTodos.sort((a, b) => new Date(a.planned_time) - new Date(b.planned_time))[0];
    return latestTodo.title || latestTodo.content || '暂无待办内容';
  }
}

// Navigation manager for bottom tabs
class NavigationManager {
  constructor() {
    this.init();
  }

  init() {
    this.setupNavigationHandlers();
  }

  setupNavigationHandlers() {
    const navButtons = document.querySelectorAll('.frame-12');
    console.log('Found navigation buttons:', navButtons.length);
    
    navButtons.forEach((button, index) => {
      const buttonText = button.querySelector('.text-wrapper-21, .text-wrapper-22')?.textContent || '';
      console.log(`Setting up button ${index}: ${buttonText}`);
      
      // Skip the button that already has onclick in HTML
      if (buttonText === '客户') {
        console.log('Skipping 客户 button - handled by HTML onclick');
        return;
      }
      
      button.style.cursor = 'pointer';
      button.addEventListener('click', (e) => {
        e.preventDefault();
        console.log(`Navigation button clicked: ${index} (${buttonText})`);
        this.handleNavigation(index, buttonText);
      });
    });
  }

  handleNavigation(index, buttonText) {
    console.log(`Handling navigation: ${index} - ${buttonText}`);
    
    switch(index) {
      case 0: // 首页
        console.log('首页 clicked');
        break;
      case 1: // 客户
        console.log('客户 clicked - navigating to customer.html');
        window.location.href = '../customer/customer.html';
        break;
      case 2: // 业绩
        console.log('业绩 clicked');
        break;
      case 3: // 设置
        console.log('设置 clicked');
        break;
      default:
        console.log('Unknown navigation button:', index, buttonText);
    }
  }
}

// Initialize all functionality when DOM is loaded
document.addEventListener('DOMContentLoaded', async () => {
  console.log('Homepage DOM loaded, initializing managers...');
  
  // Wait for global config to load
  if (window.GlobalApiConfig) {
    console.log('Global API config found, initializing managers...');
    new UserManager();
    new TodayFollowUpManager();
    new UpcomingFollowUpManager();
  } else {
    console.warn('Global API config not found');
  }
  
  new FilterManager();
  new SearchManager();
  new ActionManager();
  new NavigationManager();
  
  console.log('All managers initialized');
});
