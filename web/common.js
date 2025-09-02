
// 全局变量
let currentSearchTerm = '';
let isSearching = false;
let searchCache = new Map();
let searchHistory = [];
let tagDimensions = [];

// 环境配置
let API_CONFIG = {};

// DOM元素
let searchInput;
let searchResults;
let historyContainer;
let tagDimensionsContainer;
let cancelBtn;
let clearHistoryBtn;

// 初始化DOM元素
function initializeDOMElements() {
    searchInput = document.querySelector('.search-input');
    searchResults = document.querySelector('.search-results');
    historyContainer = document.querySelector('.history-tags');
    tagDimensionsContainer = document.getElementById('tags-dimensions-container');
    cancelBtn = document.querySelector('.cancel-btn');
    clearHistoryBtn = document.querySelector('.clear-history');
}

// 获取API基础URL
function getApiBaseUrl() {
    return API_CONFIG.baseUrl || 'http://localhost:8080';
}

// 初始化事件监听器
function initializeEventListeners() {
    if (searchInput) {
        searchInput.addEventListener('input', debounce(handleSearch, 300));
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                handleSearch();
            }
        });
    }

    if (cancelBtn) {
        cancelBtn.addEventListener('click', () => {
            if (searchInput) searchInput.value = '';
            clearSearchResults();
            currentSearchTerm = '';
        });
    }

    if (clearHistoryBtn) {
        clearHistoryBtn.addEventListener('click', clearSearchHistory);
    }
}

// 防抖函数
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 缓存管理
function getCachedResult(key) {
    const cached = searchCache.get(key);
    if (cached && Date.now() - cached.timestamp < 300000) { // 5分钟缓存
        return cached.data;
    }
    return null;
}

function setCachedResult(key, data) {
    searchCache.set(key, {
        data: data,
        timestamp: Date.now()
    });
}

// 带超时的fetch函数
function fetchWithTimeout(url, options = {}, timeout = 10000) {
    return Promise.race([
        fetch(url, options),
        new Promise((_, reject) => 
            setTimeout(() => reject(new Error('请求超时')), timeout)
        )
    ]);
}

// 搜索历史管理
function loadSearchHistory() {
    try {
        const history = localStorage.getItem('searchHistory');
        searchHistory = history ? JSON.parse(history) : [];
        renderSearchHistory();
    } catch (error) {
        console.error('加载搜索历史失败:', error);
        searchHistory = [];
    }
}

function saveSearchHistory() {
    try {
        localStorage.setItem('searchHistory', JSON.stringify(searchHistory));
    } catch (error) {
        console.error('保存搜索历史失败:', error);
    }
}

function addToSearchHistory(term) {
    if (!term || term.trim() === '') return;
    
    const trimmedTerm = term.trim();
    searchHistory = searchHistory.filter(item => item !== trimmedTerm);
    searchHistory.unshift(trimmedTerm);
    searchHistory = searchHistory.slice(0, 10); // 只保留最近10个
    
    saveSearchHistory();
    renderSearchHistory();
}

function renderSearchHistory() {
    if (!historyContainer) return;
    
    historyContainer.innerHTML = '';
    searchHistory.forEach(term => {
        const tag = document.createElement('span');
        tag.className = 'history-tag';
        tag.textContent = term;
        tag.addEventListener('click', () => {
            if (searchInput) {
                searchInput.value = term;
                handleSearch();
            }
        });
        historyContainer.appendChild(tag);
    });
}

function clearSearchHistory() {
    searchHistory = [];
    saveSearchHistory();
    renderSearchHistory();
}

// 重试机制
function retryRequest(requestFn, maxRetries = 3, delay = 1000) {
    return new Promise((resolve, reject) => {
        let retries = 0;
        
        function attempt() {
            requestFn()
                .then(resolve)
                .catch(error => {
                    retries++;
                    if (retries < maxRetries) {
                        console.log(`请求失败，${delay}ms后重试 (${retries}/${maxRetries})`);
                        setTimeout(attempt, delay);
                    } else {
                        reject(error);
                    }
                });
        }
        
        attempt();
    });
}

// 搜索处理函数
function handleSearch() {
    if (!searchInput) return;
    
    const term = searchInput.value.trim();
    if (term === currentSearchTerm) return;
    
    currentSearchTerm = term;
    
    if (term === '') {
        clearSearchResults();
        return;
    }
    
    if (term.length < 2) {
        return;
    }
    
    searchCustomers(term);
}

// 搜索客户
function searchCustomers(term) {
    if (isSearching) return;
    
    const cacheKey = `search_${term}`;
    const cachedResult = getCachedResult(cacheKey);
    
    if (cachedResult) {
        displaySearchResults(cachedResult, term);
        return;
    }
    
    isSearching = true;
    
    const requestFn = () => {
        return fetchWithTimeout(`${getApiBaseUrl()}/api/v1/customers/search`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                keyword: term,
                page: 1,
                pageSize: 20
            })
        }, 8000);
    };
    
    retryRequest(requestFn, 2, 1500)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.code === 200 && data.data) {
                setCachedResult(cacheKey, data.data);
                displaySearchResults(data.data, term);
                addToSearchHistory(term);
            } else {
                throw new Error(data.message || '搜索失败');
            }
        })
        .catch(error => {
            console.error('搜索客户失败:', error);
            displaySearchError(error.message);
        })
        .finally(() => {
            isSearching = false;
        });
}

// 显示搜索结果
function displaySearchResults(customers, searchTerm) {
    if (!searchResults) return;
    
    const resultList = searchResults.querySelector('.result-list');
    const resultsHeader = searchResults.querySelector('.results-header h3');
    
    if (!resultList || !resultsHeader) return;
    
    resultsHeader.textContent = `客户 (${customers.length})`;
    resultList.innerHTML = '';
    
    if (customers.length === 0) {
        resultList.innerHTML = '<div class="no-results">未找到相关客户</div>';
    } else {
        customers.forEach(customer => {
            const resultItem = createCustomerResultItem(customer);
            resultList.appendChild(resultItem);
        });
    }
    
    searchResults.classList.remove('hidden');
}

// 创建客户结果项
function createCustomerResultItem(customer) {
    const item = document.createElement('div');
    item.className = 'result-item';
    
    const avatar = createCustomerAvatarDisplay(customer);
    const contactInfo = getCustomerContactInfo(customer);
    const tags = createCustomerTags(customer);
    
    item.innerHTML = `
        ${avatar}
        <div class="result-content">
            <div class="result-header">
                <span class="customer-name">${escapeHtml(customer.alias || customer.name || '未知客户')}</span>
                <span class="customer-level">${getCustomerLevelText(customer.level)}</span>
            </div>
            <div class="customer-contact">${contactInfo}</div>
            <div class="customer-tags">${tags}</div>
            <div class="customer-status">${getCustomerStatusText(customer.status)}</div>
        </div>
    `;
    
    item.addEventListener('click', () => {
        window.location.href = `target.html?id=${customer.id}`;
    });
    
    return item;
}

// 创建客户头像显示
function createCustomerAvatarDisplay(customer) {
    if (customer.avatar && customer.avatar.trim() !== '') {
        return `<div class="avatar"><img src="${escapeHtml(customer.avatar)}" alt="头像" onerror="this.parentElement.innerHTML='${getCustomerInitials(customer)}'"></div>`;
    } else {
        return `<div class="text-avatar">${getCustomerInitials(customer)}</div>`;
    }
}

// 获取客户姓名首字母
function getCustomerInitials(customer) {
    const name = customer.alias || customer.name || '?';
    if (name.length === 0) return '?';
    if (name.length === 1) return name.toUpperCase();
    return name.substring(0, 2).toUpperCase();
}

// 创建客户标签
function createCustomerTags(customer) {
    if (!customer.tags || customer.tags.length === 0) {
        return '<span class="no-tags">无标签</span>';
    }
    
    return customer.tags.map(tag => 
        `<span class="tag" style="background-color: ${tag.color || '#e0e0e0'}">${escapeHtml(tag.name)}</span>`
    ).join('');
}

// 获取客户联系信息
function getCustomerContactInfo(customer) {
    const contacts = [];
    if (customer.phone) contacts.push(customer.phone);
    if (customer.wechat) contacts.push(`微信: ${customer.wechat}`);
    return contacts.length > 0 ? contacts.join(' | ') : '无联系方式';
}

// 获取客户等级文字
function getCustomerLevelText(level) {
    const levelMap = {
        1: '普通客户',
        2: '重要客户', 
        3: 'VIP客户',
        4: '核心客户'
    };
    return levelMap[level] || '普通客户';
}

// 获取客户状态文字
function getCustomerStatusText(status) {
    const statusMap = {
        1: '正常',
        2: '暂停合作',
        3: '已流失'
    };
    return statusMap[status] || '正常';
}

// HTML转义函数
function escapeHtml(text) {
    if (typeof text !== 'string') return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 清空搜索结果
function clearSearchResults() {
    if (searchResults) {
        searchResults.classList.add('hidden');
        const resultList = searchResults.querySelector('.result-list');
        if (resultList) {
            resultList.innerHTML = '';
        }
    }
}

// 显示搜索错误
function displaySearchError(message) {
    if (!searchResults) return;
    
    const resultList = searchResults.querySelector('.result-list');
    const resultsHeader = searchResults.querySelector('.results-header h3');
    
    if (resultsHeader) {
        resultsHeader.textContent = '搜索出错';
    }
    
    if (resultList) {
        resultList.innerHTML = `<div class="error-message">搜索失败: ${escapeHtml(message)}</div>`;
    }
    
    searchResults.classList.remove('hidden');
}

// 标签维度数据加载
function loadTagDimensions() {
    const loadingElement = document.getElementById('tags-loading');
    if (loadingElement) {
        loadingElement.style.display = 'block';
    }
    
    const requestFn = () => {
        return fetchWithTimeout(`${getApiBaseUrl()}/api/v1/tag-dimensions`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        }, 10000);
    };
    
    retryRequest(requestFn, 3, 2000)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.code === 200 && data.data) {
                tagDimensions = data.data;
                renderTagDimensions();
            } else {
                throw new Error(data.message || '加载标签维度失败');
            }
        })
        .catch(error => {
            console.error('加载标签维度失败:', error);
            displayTagDimensionsError(error.message);
        })
        .finally(() => {
            if (loadingElement) {
                loadingElement.style.display = 'none';
            }
        });
}

// 渲染标签维度
function renderTagDimensions() {
    if (!tagDimensionsContainer) return;
    
    tagDimensionsContainer.innerHTML = '';
    
    tagDimensions.forEach(dimension => {
        const dimensionElement = createTagDimensionElement(dimension);
        tagDimensionsContainer.appendChild(dimensionElement);
    });
}

// 创建标签维度元素
function createTagDimensionElement(dimension) {
    const element = document.createElement('div');
    element.className = 'tag-dimension';
    
    const header = document.createElement('h4');
    header.textContent = dimension.name;
    element.appendChild(header);
    
    const tagsContainer = document.createElement('div');
    tagsContainer.className = 'dimension-tags';
    
    if (dimension.tags && dimension.tags.length > 0) {
        dimension.tags.forEach(tag => {
            const tagElement = createTagElement(tag);
            tagsContainer.appendChild(tagElement);
        });
    } else {
        tagsContainer.innerHTML = '<span class="no-tags">暂无标签</span>';
    }
    
    element.appendChild(tagsContainer);
    return element;
}

// 创建标签元素
function createTagElement(tag) {
    const element = document.createElement('span');
    element.className = 'filter-tag';
    element.textContent = tag.name;
    element.style.backgroundColor = tag.color || '#e0e0e0';
    
    element.addEventListener('click', () => {
        searchByTag(tag);
    });
    
    return element;
}

// 根据标签搜索
function searchByTag(tag) {
    if (searchInput) {
        searchInput.value = tag.name;
        handleSearch();
    }
}

// 显示标签维度错误
function displayTagDimensionsError(message) {
    if (!tagDimensionsContainer) return;
    
    tagDimensionsContainer.innerHTML = `
        <div class="error-message">
            <p>加载标签失败: ${escapeHtml(message)}</p>
            <button onclick="loadTagDimensions()" class="retry-btn">重试</button>
        </div>
    `;
}

// 导航处理类
class NavigationHandler {
    constructor() {
        this.setupNavigationHandlers();
    }

    setupNavigationHandlers() {
        const navButtons = document.querySelectorAll('.nav-item');
        console.log('Found navigation buttons:', navButtons.length);
        
        navButtons.forEach((button, index) => {
            const buttonText = button.querySelector('.nav-label-active, .nav-label')?.textContent || '';
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
        console.log(`Handling navigation to: ${buttonText}`);
        
        switch(buttonText) {
            case '首页':
                window.location.href = 'index.html';
                break;
            case '业绩':
                console.log('业绩页面暂未实现');
                break;
            case '设置':
                console.log('设置页面暂未实现');
                break;
            default:
                console.log(`未知的导航项: ${buttonText}`);
        }
    }
}

// 页面初始化
function initializePage() {
    // 加载API配置
    if (typeof window.API_CONFIG !== 'undefined') {
        API_CONFIG = window.API_CONFIG;
    }
    
    // 初始化DOM元素
    initializeDOMElements();
    
    // 初始化事件监听器
    initializeEventListeners();
    
    // 加载搜索历史
    loadSearchHistory();
    
    // 加载标签维度（如果容器存在）
    if (tagDimensionsContainer) {
        loadTagDimensions();
    }
    
    // 初始化导航处理
    new NavigationHandler();
    
    console.log('页面初始化完成');
}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializePage);
} else {
    initializePage();
}