// Tab页面切换功能
class TabManager {
  constructor() {
    this.init();
  }

  init() {
    this.switchTab('follow-up'); // 默认显示跟进记录
  }

  switchTab(tabName) {
    console.log('Switching to tab:', tabName);
    
    // Hide all tab contents
    const tabContents = document.querySelectorAll('.tab-content');
    tabContents.forEach(content => {
      content.classList.remove('active');
    });

    // Show selected tab content
    const selectedTab = document.getElementById(tabName + '-tab');
    if (selectedTab) {
      selectedTab.classList.add('active');
      console.log('Tab activated:', tabName);
      
      // 当切换到标签页时，加载标签数据
      if (tabName === 'tags') {
        if (typeof window.loadTagDimensions === 'function') {
          try { window.loadTagDimensions(); } catch (e) { console.error('Error calling loadTagDimensions:', e); }
        } else {
          console.warn('loadTagDimensions is not defined');
        }
      }

      // 当切换到跟进记录页时，加载跟进记录
      if (tabName === 'follow-up' && typeof window.loadFollowUpRecords === 'function') {
        setTimeout(() => { try { window.loadFollowUpRecords(); } catch (e) { console.error('Error calling loadFollowUpRecords:', e); } }, 100);
      }

      // 当切换到资料页时，触发资料页激活钩子
      if (tabName === 'profile' && typeof window.onProfileTabActivated === 'function') {
        setTimeout(() => { try { window.onProfileTabActivated(); } catch (e) { console.error('Error calling onProfileTabActivated:', e); } }, 100);
      }

      // 当切换到备注页时，触发备注页激活钩子
      if (tabName === 'notes' && typeof window.onNotesTabActivated === 'function') {
        setTimeout(() => { try { window.onNotesTabActivated(); } catch (e) { console.error('Error calling onNotesTabActivated:', e); } }, 100);
      }
    } else {
      console.log('Tab not found:', tabName + '-tab');
    }

    // Update tab styles - reset all tabs first
    const allTabs = document.querySelectorAll('.components-tab, .components-tab-2');
    allTabs.forEach(tab => {
      tab.className = 'components-tab';
      const hover = tab.querySelector('.hover, .hover-2');
      if (hover) {
        hover.className = 'hover-2';
      }
      const text = tab.querySelector('.text, .text-2');
      if (text) {
        text.className = 'text';
      }
    });

    // Set active tab style
    const activeTabElement = document.querySelector(`[onclick="switchTab('${tabName}')"]`);
    if (activeTabElement) {
      activeTabElement.className = 'components-tab-2';
      const text = activeTabElement.querySelector('.text');
      if (text) {
        text.className = 'text-2';
      }
      console.log('Tab style updated for:', tabName);
    }
  }
}

// 跟进记录相关功能
class FollowUpTab {
  constructor() {
    this.setupEventListeners();
  }

  setupEventListeners() {
    // 可以在这里添加跟进记录相关的事件监听器
  }

  // 打开添加记录弹窗
  openRecordModal() {
    const modal = document.getElementById('recordModal');
    if (modal) {
      modal.style.display = 'flex';
      document.body.style.overflow = 'hidden';
      console.log('Record modal opened');
    }
  }

  // 关闭添加记录弹窗
  closeRecordModal() {
    const modal = document.getElementById('recordModal');
    if (modal) {
      modal.style.display = 'none';
      document.body.style.overflow = 'auto';
      this.resetRecordForm();
      console.log('Record modal closed');
    }
  }

  resetRecordForm() {
    document.getElementById('recordType').value = 'call';
    document.getElementById('recordContent').value = '';
    document.getElementById('createTodoSwitch').checked = false;
    document.getElementById('planTime').value = 'today';
    document.getElementById('planTimeRow').style.display = 'none';
  }

  createRecord() {
    const recordData = {
      type: document.getElementById('recordType').value,
      content: document.getElementById('recordContent').value.trim(),
      createTodo: document.getElementById('createTodoSwitch').checked,
      planTime: document.getElementById('createTodoSwitch').checked ? document.getElementById('planTime').value : null,
      createTime: new Date().toISOString()
    };
    
    if (!recordData.content) {
      alert('请输入记录内容');
      return;
    }
    
    console.log('Creating record:', recordData);
    
    alert('跟进记录创建成功！' + (recordData.createTodo ? '\n同时已创建待办事项。' : ''));
    this.closeRecordModal();
  }
}

// 标签相关功能已移动到 tags.js

// 偏好相关功能
class PreferencesTab {
  constructor() {
    this.setupEventListeners();
  }

  setupEventListeners() {
    // 可以在这里添加偏好相关的事件监听器
  }

  openPreferenceModal() {
    const modal = document.getElementById('preferenceModal');
    if (modal) {
      modal.style.display = 'flex';
      document.body.style.overflow = 'hidden';
      console.log('Preference modal opened');
    }
  }

  closePreferenceModal() {
    const modal = document.getElementById('preferenceModal');
    if (modal) {
      modal.style.display = 'none';
      document.body.style.overflow = 'auto';
      
      document.getElementById('preferenceSearch').value = '';
      document.getElementById('newProductInputRow').style.display = 'none';
      document.getElementById('newProductInput').value = '';
      
      console.log('Preference modal closed');
    }
  }
}

// 熟人相关功能
class ContactsTab {
  constructor() {
    this.setupEventListeners();
  }

  setupEventListeners() {
    // 可以在这里添加熟人相关的事件监听器
  }

  openAddContactModal() {
    document.getElementById('addContactModal').style.display = 'flex';
    document.body.style.overflow = 'hidden';
    console.log('Add contact modal opened');
  }

  closeAddContactModal() {
    document.getElementById('addContactModal').style.display = 'none';
    document.body.style.overflow = 'auto';
    
    const notesInput = document.getElementById('contactNotes');
    if (notesInput) {
      notesInput.value = '';
    }
    
    console.log('Add contact modal closed');
  }

  createContact() {
    const notes = document.getElementById('contactNotes').value.trim();
    
    const contactData = {
      name: '徐小二',
      relationship: '好友',
      notes: notes,
      createTime: new Date().toISOString()
    };
    
    console.log('Creating contact:', contactData);
    
    alert('熟人关系创建成功！');
    this.closeAddContactModal();
  }
}

// 全局函数，供HTML调用
window.switchTab = function(tabName) {
  if (window.tabManager) {
    window.tabManager.switchTab(tabName);
  }
};

window.openRecordModal = function() {
  if (window.followUpTab) {
    window.followUpTab.openRecordModal();
  }
};

window.closeRecordModal = function() {
  if (window.followUpTab) {
    window.followUpTab.closeRecordModal();
  }
};


// openTagModal and closeTagModal functions are defined in tags.js

// toggleTag function is now defined in tags.js

window.openPreferenceModal = function() {
  if (window.preferencesTab) {
    window.preferencesTab.openPreferenceModal();
  }
};

window.closePreferenceModal = function() {
  if (window.preferencesTab) {
    window.preferencesTab.closePreferenceModal();
  }
};

window.openAddContactModal = function() {
  if (window.contactsTab) {
    window.contactsTab.openAddContactModal();
  }
};

window.closeAddContactModal = function() {
  if (window.contactsTab) {
    window.contactsTab.closeAddContactModal();
  }
};


// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
  window.tabManager = new TabManager();
  window.followUpTab = new FollowUpTab();
  window.preferencesTab = new PreferencesTab();
  window.contactsTab = new ContactsTab();
});