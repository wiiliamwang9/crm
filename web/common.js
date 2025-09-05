
// 全局变量
let currentSearchTerm = '';
let isSearching = false;
let searchCache = new Map();
let searchHistory = [];
let tagDimensions = [];
let selectedTags = []; // 存储选中的标签

// 环境配置
let API_CONFIG = {};

// DOM元素
let searchInput;
let searchResults;
let historyContainer;
let tagDimensionsContainer;
let cancelBtn;
let clearHistoryBtn;
let trashBtn;

// 初始化DOM元素
function initializeDOMElements() {
    // 兼容不同页面的搜索输入框
    searchInput = document.getElementById('searchInput') || document.querySelector('.search-input');
    searchResults = document.querySelector('.search-results');
    historyContainer = document.querySelector('.history-tags');
    tagDimensionsContainer = document.getElementById('tags-dimensions-container');
    // 兼容不同页面的取消按钮
    cancelBtn = document.getElementById('cancelBtn') || document.querySelector('.cancel-btn');
    clearHistoryBtn = document.querySelector('.clear-history');
    trashBtn = document.querySelector('.trash-btn');
}

// 获取API基础URL
function getApiBaseUrl() {
    if (window.GlobalApiConfig) {
        return window.GlobalApiConfig.BASE_URL;
    }
    // 回退到服务器地址
    return 'https://static.lamdar.cn:9501/crm';
}

// 初始化事件监听器
function initializeEventListeners() {
    if (searchInput) {

        // 添加input事件监听器（带防抖）
        const debouncedSearch = debounce(handleSearch, 300);
        searchInput.addEventListener('input', (e) => {
            debouncedSearch();
        });

        // 添加keypress事件监听器
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                handleSearch();
            }
        });

    } else {
        console.error('searchInput元素不存在，无法绑定事件监听器');
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

    if (trashBtn) {
        trashBtn.addEventListener('click', () => {
            clearSearchHistory();
            alert('搜索历史已清空');
        });
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

        // 如果没有搜索历史，添加一些示例数据用于演示
        if (searchHistory.length === 0) {
            searchHistory = [];
            saveSearchHistory();
        }
        
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

    if (searchHistory.length === 0) {
        // 如果没有搜索历史，显示提示信息
        const emptyTip = document.createElement('span');
        emptyTip.className = 'tag empty-tag';
        emptyTip.textContent = '暂无搜索历史';
        emptyTip.style.backgroundColor = '#f5f5f5';
        emptyTip.style.color = '#999';
        emptyTip.style.cursor = 'default';
        historyContainer.appendChild(emptyTip);
        return;
    }
    
    searchHistory.forEach(term => {
        const tag = document.createElement('span');
        tag.className = 'tag';
        tag.textContent = term;
        tag.addEventListener('click', () => {
            if (searchInput) {
                searchInput.value = term;
                handleSearch();
                // 添加到搜索历史
                addToSearchHistory(term);
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
    if (!searchInput) {
        console.error('searchInput不存在');
        return;
    }

    const term = searchInput.value.trim();

    if (term === currentSearchTerm) {
        return;
    }

    currentSearchTerm = term;

    if (term === '') {
        clearSearchResults();
        return;
    }

    if (term.length < 1) {
        return;
    }

    // 添加到搜索历史
    addToSearchHistory(term);
    
    searchCustomers(term);
}

// 根据标签搜索客户
function searchCustomersByTags(tagNames) {
    if (isSearching) {
        return;
    }

    const tagsParam = tagNames.join(',');
    const cacheKey = `search_tags_${tagsParam}`;
    const cachedResult = getCachedResult(cacheKey);

    if (cachedResult) {
        displaySearchResults(cachedResult, `标签: ${tagNames.join(', ')}`);
        return;
    }

    isSearching = true;

    const apiUrl = `${getApiBaseUrl()}/api/v1/customers/search?tags=${encodeURIComponent(tagsParam)}`;
    console.log('🔍 标签搜索API URL:', apiUrl);

    fetch(apiUrl)
        .then(response => {
            console.log('📡 标签搜索响应状态:', response.status);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('📡 标签搜索响应数据:', data);
            if (data.data) {
                setCachedResult(cacheKey, data.data);
                displaySearchResults(data.data, `标签: ${tagNames.join(', ')}`);
            } else {
                throw new Error(data.message || '搜索失败');
            }
        })
        .catch(error => {
            console.error('❌ 标签搜索失败:', error);
            displaySearchError('搜索服务暂时不可用');
        })
        .finally(() => {
            isSearching = false;
        });
}

// 搜索客户
function searchCustomers(term) {

    if (isSearching) {
        return;
    }

    const cacheKey = `search_${term}`;
    const cachedResult = getCachedResult(cacheKey);

    if (cachedResult) {
        displaySearchResults(cachedResult, term);
        return;
    }


    isSearching = true;

    const apiUrl = `${getApiBaseUrl()}/api/v1/customers/search?keyword=${encodeURIComponent(term)}`;

    const requestFn = () => {
        return fetchWithTimeout(apiUrl, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
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
            if (data.data && Array.isArray(data.data)) {
                setCachedResult(cacheKey, data.data);
                displaySearchResults(data.data, term);
                addToSearchHistory(term);
            } else {
                throw new Error(data.message || '搜索失败');
            }
        })
        .catch(error => {
            console.error('搜索客户失败，尝试使用测试数据:', error);
            // 如果API请求失败，尝试使用本地测试数据
            loadTestCustomers(term, cacheKey);
        })
        .finally(() => {
            isSearching = false;
        });
}

// 加载测试客户数据
function loadTestCustomers(term, cacheKey) {
    fetch('./test-customers.json')
        .then(response => {
            if (!response.ok) {
                throw new Error('无法加载测试数据');
            }
            return response.json();
        })
        .then(data => {
            if (data.code === 200 && data.data) {
                // 根据搜索词过滤客户数据
                const filteredCustomers = data.data.filter(customer => {
                    const contactName = (customer.ContactName || customer.contactName || '').toLowerCase();
                    const name = (customer.name || '').toLowerCase();
                    const searchTerm = term.toLowerCase();

                    return contactName.includes(searchTerm) || name.includes(searchTerm);
                });

                setCachedResult(cacheKey, filteredCustomers);
                displaySearchResults(filteredCustomers, term);
                addToSearchHistory(term);
            } else {
                throw new Error('测试数据格式错误');
            }
        })
        .catch(error => {
            console.error('加载测试数据失败:', error);
            displaySearchError('搜索服务暂时不可用');
        });
}

// 显示搜索结果
function displaySearchResults(customers, searchTerm) {
    // 查找search.html中的客户显示区域
    const contentContainer = document.querySelector('.content-container');
    const customerSection = contentContainer ? contentContainer.querySelector('.search-content-section') : null;

    if (!customerSection) {
        console.error('未找到客户显示区域');
        return;
    }

    // 显示内容容器
    if (contentContainer) {
        contentContainer.classList.add('show');
    }

    const resultsHeader = customerSection.querySelector('h3');
    if (resultsHeader) {
        resultsHeader.textContent = `客户 (${customers.length})`;
    }

    // 清空现有的客户项和"未找到相关客户"提示，但保留h3标题
    const existingItems = customerSection.querySelectorAll('.search-content-item, .no-results');
    existingItems.forEach(item => item.remove());

    if (customers.length === 0) {
        const noResultsDiv = document.createElement('div');
        noResultsDiv.className = 'no-results';
        noResultsDiv.textContent = '未找到相关客户';
        noResultsDiv.style.padding = '20px';
        noResultsDiv.style.textAlign = 'center';
        noResultsDiv.style.color = '#999';
        customerSection.appendChild(noResultsDiv);
    } else {
        customers.forEach(customer => {
            const resultItem = createSearchContentItem(customer);
            customerSection.appendChild(resultItem);
        });
    }
}

// 创建搜索内容项（适配search.html结构）
function createSearchContentItem(customer) {
    const item = document.createElement('div');
    item.className = 'search-content-item';

    // 按照用户要求的字段映射：
    // - 主标题显示ContactName字段
    // - 副标题显示name字段，且保留文字前面的"|"的分隔符号
    // - 标签显示tags字段
    const contactName = customer.contact_name || customer.ContactName || customer.contactName || '未知客户';
    const name = customer.name || '';
    const tags = customer.tags || [];

    // 创建标签HTML
    let tagsHtml = '';
    if (tags && tags.length > 0) {
        // 定义可用的标签样式
        const tagClasses = ['search-content-tag-blue', 'search-content-tag-green', 'search-content-tag-orange'];

        tagsHtml = tags.map(tag => {
            // 随机选择一个样式
            const randomIndex = Math.floor(Math.random() * tagClasses.length);
            const tagClass = tagClasses[randomIndex];

            if (typeof tag === 'string') {
                return `<span class="${tagClass}">${escapeHtml(tag)}</span>`;
            } else if (tag.name) {
                return `<span class="${tagClass}">${escapeHtml(tag.name)}</span>`;
            }
            return '';
        }).join('');
    }

    item.innerHTML = `
        <img src="./image/avatar.png" alt="头像">
        <div class="search-content-info">
            <div class="search-content-name">${escapeHtml(contactName)}</div>
            <div class="search-content-organization"> | </div>
            <div class="search-content-contact">${escapeHtml(name.replace(/\n/g, ' '))}</div>
            <div class="search-content-tags">
                ${tagsHtml}
            </div>
        </div>
    `;

    // 添加点击事件
    item.addEventListener('click', () => {
        if (customer.id) {
            window.location.href = `customer.html?id=${customer.id}`;
        }
    });

    return item;
}

// 创建客户结果项（保留原有函数以兼容其他页面）
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
                <span class="customer-name">${escapeHtml(customer.contact_name || customer.alias || customer.name || '未知客户')}</span>
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
    const name = customer.contact_name || customer.alias || customer.name || '?';
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

    // 隐藏内容容器
    const contentContainer = document.querySelector('.content-container');
    if (contentContainer) {
        contentContainer.classList.remove('show');

        // 清除"未找到相关客户"提示
        const customerSection = contentContainer.querySelector('.search-content-section');
        if (customerSection) {
            const noResultsDiv = customerSection.querySelector('.no-results');
            if (noResultsDiv) {
                noResultsDiv.remove();
            }
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
    console.log('🏷️ 开始加载标签维度数据...');
    const loadingElement = document.getElementById('tags-loading');
    if (loadingElement) {
        loadingElement.style.display = 'block';
        console.log('✅ 显示加载提示');
    }

    // 同时加载标签维度和标签数据
    console.log('📡 开始并行请求标签维度和标签数据');
    Promise.all([
        loadTagDimensionsData(),
        loadTagsData()
    ])
        .then(([dimensions, tags]) => {
            console.log('✅ 标签数据加载成功:', {dimensions: dimensions.length, tags: tags.length});

            // 将标签数据按维度分组
            const tagsByDimension = {};
            tags.forEach(tag => {
                if (!tagsByDimension[tag.dimension_id]) {
                    tagsByDimension[tag.dimension_id] = [];
                }
                tagsByDimension[tag.dimension_id].push(tag);
            });
            console.log('🔗 标签按维度分组完成:', tagsByDimension);

            // 为每个维度添加对应的标签
            tagDimensions = dimensions.map(dimension => ({
                ...dimension,
                tags: tagsByDimension[dimension.id] || []
            }));
            console.log('📋 最终标签维度数据:', tagDimensions);

            renderTagDimensions();
        })
        .catch(error => {
            console.error('❌ 加载标签数据失败:', error);
            displayTagDimensionsError(error.message);
        })
        .finally(() => {
            if (loadingElement) {
                loadingElement.style.display = 'none';
                console.log('✅ 隐藏加载提示');
            }
        });
}

// 加载标签维度数据
function loadTagDimensionsData() {
    const url = `${getApiBaseUrl()}/api/v1/tag-dimensions`;
    console.log('📡 请求标签维度数据:', url);
    
    const requestFn = () => {
        return fetchWithTimeout(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        }, 10000);
    };

    return retryRequest(requestFn, 3, 2000)
        .then(response => {
            console.log('📡 标签维度响应状态:', response.status);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('📡 标签维度响应数据:', data);
            if (data.data) {
                return data.data;
            } else {
                throw new Error(data.message || '加载标签维度失败');
            }
        });
}

// 加载标签数据
function loadTagsData() {
    const url = `${getApiBaseUrl()}/api/v1/tags`;
    console.log('📡 请求标签数据:', url);

    const requestFn = () => {
        return fetchWithTimeout(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        }, 10000);
    };

    return retryRequest(requestFn, 3, 2000)
        .then(response => {
            console.log('📡 标签响应状态:', response.status);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('📡 标签响应数据:', data);
            if (data.data) {
                return data.data;
            } else {
                throw new Error(data.message || '加载标签失败');
            }
        });
}

// 渲染标签维度
function renderTagDimensions() {
    if (!tagDimensionsContainer) {
        console.error('标签维度容器未找到');
        return;
    }

    tagDimensionsContainer.innerHTML = '';

    if (!tagDimensions || tagDimensions.length === 0) {
        tagDimensionsContainer.innerHTML = '<div class="no-tags">暂无标签维度</div>';
        return;
    }

    tagDimensions.forEach(dimension => {
        const dimensionElement = createTagDimensionElement(dimension);
        tagDimensionsContainer.appendChild(dimensionElement);
    });
}

// 创建标签维度元素
function createTagDimensionElement(dimension) {
    const dimensionDiv = document.createElement('div');
    dimensionDiv.className = 'tag-group';

    const titleSpan = document.createElement('span');
    titleSpan.className = 'group-title';
    titleSpan.textContent = dimension.name + '：';
    dimensionDiv.appendChild(titleSpan);

    const tagsDiv = document.createElement('div');
    tagsDiv.className = 'tag-items';

    if (dimension.tags && dimension.tags.length > 0) {
        dimension.tags.forEach(tag => {
            const tagElement = createTagElement(tag);
            tagsDiv.appendChild(tagElement);
        });
    } else {
        tagsDiv.innerHTML = '<span class="no-tags">暂无标签</span>';
    }

    dimensionDiv.appendChild(tagsDiv);
    return dimensionDiv;
}

// 创建标签元素
function createTagElement(tag) {
    const tagSpan = document.createElement('span');
    tagSpan.className = 'tag-item';
    tagSpan.textContent = tag.name;
    tagSpan.dataset.tagName = tag.name; // 存储标签名称用于查找
    tagSpan.dataset.tagId = tag.id; // 存储标签ID

    // 如果标签有特定颜色类型，添加对应的CSS类
    if (tag.color_type) {
        tagSpan.classList.add(tag.color_type);
    }

    // 检查是否已选中
    if (selectedTags.some(selectedTag => selectedTag.name === tag.name)) {
        tagSpan.classList.add('selected');
    }

    tagSpan.addEventListener('click', () => {
        toggleTagSelection(tag, tagSpan);
    });

    return tagSpan;
}

// 切换标签选中状态
function toggleTagSelection(tag, tagElement) {
    const isSelected = selectedTags.some(selectedTag => selectedTag.name === tag.name);

    if (isSelected) {
        // 取消选中
        selectedTags = selectedTags.filter(selectedTag => selectedTag.name !== tag.name);
        tagElement.classList.remove('selected');
        console.log('🏷️ 取消选中标签:', tag.name);
    } else {
        // 选中标签
        selectedTags.push({
            id: tag.id,
            name: tag.name,
            dimension_id: tag.dimension_id
        });
        tagElement.classList.add('selected');
        console.log('🏷️ 选中标签:', tag.name);
    }

    console.log('🏷️ 当前选中的标签:', selectedTags);

    // 触发搜索
    searchBySelectedTags();
}

// 根据选中的标签搜索
function searchBySelectedTags() {
    if (selectedTags.length === 0) {
        // 如果没有选中标签，清空搜索
        if (searchInput) {
            searchInput.value = '';
        }
        clearSearchResults();
        return;
    }

    // 构建标签名称列表用于搜索
    const tagNames = selectedTags.map(tag => tag.name);
    console.log('🔍 使用标签搜索:', tagNames);

    // 调用搜索函数
    searchCustomersByTags(tagNames);
}

// 根据标签搜索（保留原函数用于兼容）
function searchByTag(tagName) {
    if (searchInput) {
        searchInput.value = tagName;
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

        navButtons.forEach((button, index) => {
            const buttonText = button.querySelector('.nav-label-active, .nav-label')?.textContent || '';

            // Skip the button that already has onclick in HTML
            if (buttonText === '客户') {
                return;
            }

            button.style.cursor = 'pointer';
            button.addEventListener('click', (e) => {
                e.preventDefault();
                this.handleNavigation(index, buttonText);
            });
        });
    }

    handleNavigation(index, buttonText) {

        switch (buttonText) {
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
    } else if (typeof window.baseUrl !== 'undefined') {
        API_CONFIG = {
            baseUrl: window.baseUrl
        };
    }

    // 初始化DOM元素
    initializeDOMElements();

    // 初始化事件监听器
    initializeEventListeners();

    // 加载搜索历史
    loadSearchHistory();

    // 加载标签维度（如果容器存在）
    console.log('🔍 检查标签维度容器:', tagDimensionsContainer);
    if (tagDimensionsContainer) {
        console.log('✅ 标签维度容器存在，开始加载标签数据');
        loadTagDimensions();
    } else {
        console.log('❌ 标签维度容器不存在，跳过标签加载');
    }

    // 初始化导航处理
    new NavigationHandler();


    // 如果有搜索输入框，设置焦点
    if (searchInput) {
        searchInput.focus();
    } else {
        console.error('未找到搜索输入框元素');
    }

}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializePage);
} else {
    initializePage();
}


// 搜索历史功能实现完成
// 功能包括：
// 1. 自动保存搜索记录到localStorage
// 2. 页面加载时从localStorage加载历史记录
// 3. 点击历史记录标签可触发搜索
// 4. 限制显示最多10条历史记录
// 5. 垃圾桶按钮可清空所有历史记录