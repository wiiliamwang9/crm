// 备注功能管理
class NotesManager {
  constructor() {
    this.currentCustomerId = null;
    this.saveTimeout = null;
    this.isInitialized = false;
    
    this.init();
  }
  
  init() {
    // 等待DOM加载完成
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        this.setupEventListeners();
        this.loadCustomerData();
      });
    } else {
      this.setupEventListeners();
      this.loadCustomerData();
    }
  }
  
  setupEventListeners() {
    // 获取textarea元素
    const notesTextarea = document.getElementById('customerNotes');
    const charCountElement = document.getElementById('notesCharCount');
    
    if (!notesTextarea || !charCountElement) {
      // 如果元素不存在，延迟重试
      setTimeout(() => this.setupEventListeners(), 100);
      return;
    }
    
    // 监听内容变化
    notesTextarea.addEventListener('input', (e) => {
      this.updateCharacterCount();
      this.handleContentChange();
    });
    
    // 监听失焦事件
    notesTextarea.addEventListener('blur', () => {
      this.saveNotes();
    });
    
    this.isInitialized = true;
  }
  
  // 获取当前客户ID
  getCurrentCustomerId() {
    if (!this.currentCustomerId) {
      const urlParams = new URLSearchParams(window.location.search);
      this.currentCustomerId = urlParams.get('customer_id');
      
      if (!this.currentCustomerId) {
        console.warn('未找到customer_id参数，使用默认客户ID: 92');
        this.currentCustomerId = '92';
      }
    }
    return this.currentCustomerId;
  }
  
  // 加载客户数据
  async loadCustomerData() {
    try {
      const customerId = this.getCurrentCustomerId();
      const baseUrl = window.GlobalApiConfig ? window.GlobalApiConfig.BASE_URL : 'http://localhost:8081';
      
      const response = await fetch(`${baseUrl}/api/v1/customers/${customerId}`);
      
      if (response.ok) {
        const result = await response.json();
        const customer = result.data;
        
        console.log('Customer data loaded:', customer);
        
        // 显示备注内容
        this.displayCustomerNotes(customer.remark || '');
      } else {
        console.error('Failed to load customer data:', response.statusText);
        this.updateSaveStatus('加载失败', 'error');
      }
    } catch (error) {
      console.error('Error loading customer data:', error);
      this.updateSaveStatus('加载失败', 'error');
    }
  }
  
  // 显示客户备注
  displayCustomerNotes(remark) {
    const notesTextarea = document.getElementById('customerNotes');
    if (notesTextarea) {
      notesTextarea.value = remark;
      this.updateCharacterCount();
      this.updateSaveStatus('已保存', 'saved');
    }
  }
  
  // 更新字符计数
  updateCharacterCount() {
    const notesTextarea = document.getElementById('customerNotes');
    const charCountElement = document.getElementById('notesCharCount');
    
    if (notesTextarea && charCountElement) {
      const currentLength = notesTextarea.value.length;
      charCountElement.textContent = currentLength;
      
      // 如果接近限制，改变颜色
      if (currentLength > 900) {
        charCountElement.style.color = '#f53f3f';
      } else if (currentLength > 800) {
        charCountElement.style.color = '#ff7d00';
      } else {
        charCountElement.style.color = 'var(--color-text-3)';
      }
    }
  }
  
  // 处理内容变化
  handleContentChange() {
    this.updateSaveStatus('未保存', 'saving');
    
    // 清除之前的定时器
    if (this.saveTimeout) {
      clearTimeout(this.saveTimeout);
    }
    
    // 设置新的自动保存定时器（2秒后保存）
    this.saveTimeout = setTimeout(() => {
      this.saveNotes();
    }, 2000);
  }
  
  // 保存备注
  async saveNotes() {
    const notesTextarea = document.getElementById('customerNotes');
    if (!notesTextarea) return;
    
    try {
      this.updateSaveStatus('保存中...', 'saving');
      
      const customerId = this.getCurrentCustomerId();
      const baseUrl = window.GlobalApiConfig ? window.GlobalApiConfig.BASE_URL : 'http://localhost:8081';
      const notes = notesTextarea.value;
      
      // 构造更新数据（只更新remark字段）
      const updateData = {
        remark: notes
      };
      
      const response = await fetch(`${baseUrl}/api/v1/customers/${customerId}/remark`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateData)
      });
      
      if (response.ok) {
        this.updateSaveStatus('已保存', 'saved');
        console.log('Notes saved successfully');
      } else {
        throw new Error('保存失败');
      }
    } catch (error) {
      console.error('Error saving notes:', error);
      this.updateSaveStatus('保存失败', 'error');
    }
  }
  
  // 更新保存状态显示
  updateSaveStatus(text, status) {
    const saveStatusElement = document.getElementById('saveStatus');
    if (saveStatusElement) {
      saveStatusElement.textContent = text;
      saveStatusElement.className = `save-status ${status}`;
    }
  }
  
  // 当切换到备注tab时调用
  onNotesTabActivated() {
    // 确保初始化完成
    if (!this.isInitialized) {
      this.setupEventListeners();
    }
    
    // 如果还没有加载数据，则加载
    const notesTextarea = document.getElementById('customerNotes');
    if (notesTextarea && notesTextarea.value === '') {
      this.loadCustomerData();
    }
  }
}

// 创建全局实例
window.notesManager = new NotesManager();

// 为其他脚本提供接口
window.onNotesTabActivated = function() {
  if (window.notesManager) {
    window.notesManager.onNotesTabActivated();
  }
};