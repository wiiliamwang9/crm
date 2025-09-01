// 可搜索下拉框组件

class SearchableSelect {
  constructor(containerId, options = {}) {
    this.containerId = containerId;
    this.container = document.getElementById(containerId);
    this.options = options;
    this.data = [];
    this.selectedValue = null;
    this.selectedText = '';
    this.isOpen = false;
    
    if (!this.container) {
      console.error(`SearchableSelect: Container ${containerId} not found`);
      return;
    }
    
    this.init();
  }
  
  init() {
    this.input = this.container.querySelector('.searchable-select-input');
    this.dropdown = this.container.querySelector('.searchable-select-dropdown');
    this.searchInput = this.container.querySelector('.searchable-select-search-input');
    this.optionsContainer = this.container.querySelector('.searchable-select-options');
    this.hiddenInput = this.container.querySelector('input[type="hidden"]');
    
    this.bindEvents();
  }
  
  bindEvents() {
    // 点击主输入框显示/隐藏下拉列表
    this.input.addEventListener('click', () => {
      this.toggle();
    });
    
    // 搜索输入事件
    this.searchInput.addEventListener('input', (e) => {
      this.filterOptions(e.target.value);
    });
    
    // 点击外部关闭下拉框
    document.addEventListener('click', (e) => {
      if (!this.container.contains(e.target)) {
        this.close();
      }
    });
    
    // 阻止搜索框点击时关闭下拉框
    this.searchInput.addEventListener('click', (e) => {
      e.stopPropagation();
    });
  }
  
  setData(data) {
    this.data = data;
    this.renderOptions();
  }
  
  renderOptions() {
    this.optionsContainer.innerHTML = '';
    
    if (this.data.length === 0) {
      const emptyOption = document.createElement('div');
      emptyOption.className = 'searchable-select-option';
      emptyOption.textContent = '暂无数据';
      emptyOption.style.color = '#86909c';
      this.optionsContainer.appendChild(emptyOption);
      return;
    }
    
    this.data.forEach(item => {
      const option = document.createElement('div');
      option.className = 'searchable-select-option';
      option.textContent = item.name;
      option.dataset.value = item.id;
      option.dataset.name = item.name;
      
      if (this.selectedValue == item.id) {
        option.classList.add('selected');
      }
      
      option.addEventListener('click', () => {
        this.selectOption(item.id, item.name);
      });
      
      this.optionsContainer.appendChild(option);
    });
  }
  
  filterOptions(searchText) {
    const options = this.optionsContainer.querySelectorAll('.searchable-select-option');
    const lowerSearchText = searchText.toLowerCase();
    
    options.forEach(option => {
      const name = option.dataset.name;
      if (name && name.toLowerCase().includes(lowerSearchText)) {
        option.classList.remove('hidden');
      } else {
        option.classList.add('hidden');
      }
    });
  }
  
  selectOption(value, text) {
    this.selectedValue = value;
    this.selectedText = text;
    this.input.value = text;
    this.hiddenInput.value = value;
    
    // 更新选中状态
    const options = this.optionsContainer.querySelectorAll('.searchable-select-option');
    options.forEach(option => {
      if (option.dataset.value == value) {
        option.classList.add('selected');
      } else {
        option.classList.remove('selected');
      }
    });
    
    this.close();
    
    // 触发change事件
    if (this.options.onChange) {
      this.options.onChange(value, text);
    }
  }
  
  setValue(value, text) {
    this.selectedValue = value;
    this.selectedText = text;
    this.input.value = text || '';
    this.hiddenInput.value = value || '';
    
    // 更新UI中的选中状态
    const options = this.optionsContainer.querySelectorAll('.searchable-select-option');
    options.forEach(option => {
      if (option.dataset.value == value) {
        option.classList.add('selected');
      } else {
        option.classList.remove('selected');
      }
    });
  }
  
  getValue() {
    return this.selectedValue;
  }
  
  getText() {
    return this.selectedText;
  }
  
  open() {
    this.isOpen = true;
    this.container.classList.add('open');
    this.searchInput.value = '';
    this.filterOptions('');
    
    // 聚焦搜索框
    setTimeout(() => {
      this.searchInput.focus();
    }, 100);
  }
  
  close() {
    this.isOpen = false;
    this.container.classList.remove('open');
  }
  
  toggle() {
    if (this.isOpen) {
      this.close();
    } else {
      this.open();
    }
  }
  
  clear() {
    this.selectedValue = null;
    this.selectedText = '';
    this.input.value = '';
    this.hiddenInput.value = '';
    
    const options = this.optionsContainer.querySelectorAll('.searchable-select-option');
    options.forEach(option => {
      option.classList.remove('selected');
    });
  }
}