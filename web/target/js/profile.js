// 资料页面管理器
class ProfileManager {
  constructor() {
    this.currentCustomerId = null;
    this.customerData = null;
    this.isInitialized = false;
    
    this.init();
  }
  
  init() {
    // 等待DOM加载完成
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        this.loadCustomerData();
      });
    } else {
      this.loadCustomerData();
    }
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
        this.customerData = result.data;
        
        console.log('Customer profile data loaded:', this.customerData);
        
        // 如果资料tab已经激活，立即显示数据
        this.displayProfileData();
      } else {
        console.error('Failed to load customer profile data:', response.statusText);
      }
    } catch (error) {
      console.error('Error loading customer profile data:', error);
    }
  }
  
  // 显示资料数据
  displayProfileData() {
    if (!this.customerData) {
      return;
    }
    
    // 更新别名字段
    this.updateContactName();
    
    // 更新电话字段
    this.updatePhones();
    
    // 更新微信字段
    this.updateWechats();
    
    // 更新头像/照片
    this.updatePhotos();
  }
  
  // 更新别名显示
  updateContactName() {
    const contactNameElement = document.getElementById('customerContactName');
    if (contactNameElement) {
      contactNameElement.textContent = this.customerData.contact_name || '-';
    }
  }
  
  // 更新电话显示
  updatePhones() {
    const phoneElement = document.getElementById('customerPhones');
    if (phoneElement) {
      if (this.customerData.phones && this.customerData.phones.length > 0) {
        // 使用逗号和空格分隔，让文本能在合适位置换行
        phoneElement.textContent = this.customerData.phones.join(', ');
      } else {
        phoneElement.textContent = '-';
      }
    }
  }
  
  // 更新微信显示
  updateWechats() {
    const wechatElement = document.getElementById('customerWechats');
    if (wechatElement) {
      if (this.customerData.wechats && this.customerData.wechats.length > 0) {
        // 使用逗号和空格分隔，让文本能在合适位置换行
        wechatElement.textContent = this.customerData.wechats.join(', ');
      } else {
        wechatElement.textContent = '-';
      }
    }
  }
  
  // 更新照片显示
  updatePhotos() {
    const imageElements = document.querySelectorAll('#profile-tab .frame-18 img');
    
    if (this.customerData.photos && this.customerData.photos.length > 0) {
      // 更新已有的照片元素
      imageElements.forEach((img, index) => {
        if (index < this.customerData.photos.length) {
          img.src = this.customerData.photos[index];
          img.style.display = 'block';
        } else {
          img.style.display = 'none';
        }
      });
    } else if (this.customerData.avatar) {
      // 如果没有photos但有avatar，显示avatar
      if (imageElements.length > 0) {
        imageElements[0].src = this.customerData.avatar;
        imageElements[0].style.display = 'block';
        // 隐藏其他图片
        for (let i = 1; i < imageElements.length; i++) {
          imageElements[i].style.display = 'none';
        }
      }
    } else {
      // 没有照片时使用默认头像
      imageElements.forEach(img => {
        img.src = 'https://c.animaapp.com/mepb4kxjUWP5uf/img/---8.png';
      });
    }
  }
  
  // 当切换到资料tab时调用
  onProfileTabActivated() {
    // 如果还没有加载数据，则加载
    if (!this.customerData) {
      this.loadCustomerData();
    } else {
      // 如果已经有数据，直接显示
      this.displayProfileData();
    }
  }
}

// 创建全局实例
window.profileManager = new ProfileManager();

// 为其他脚本提供接口
window.onProfileTabActivated = function() {
  if (window.profileManager) {
    window.profileManager.onProfileTabActivated();
  }
};