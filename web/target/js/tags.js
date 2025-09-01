// 标签相关功能

// 全局标签数据
let tagDimensionsData = [];
let selectedTagDimension = null;
let selectedCreateDimensionId = null; // 用于创建新标签时选择的维度
let currentCustomerData = null; // 当前客户数据
let customerSystemTags = []; // 客户的系统标签ID数组

// 清除搜索框
function clearTagSearch() {
  const searchInput = document.getElementById('tag-search-input');
  const clearBtn = document.querySelector('.clear-search');
  if (searchInput) {
    searchInput.value = '';
    searchInput.focus();
  }
  if (clearBtn) {
    clearBtn.style.display = 'none';
  }
}

// 添加新标签
async function addNewTag() {
  const searchInput = document.getElementById('tag-search-input');
  const tagName = searchInput ? searchInput.value.trim() : '';
  
  if (!tagName) {
    alert('请输入标签名称');
    return;
  }
  
  try {
    // 首先检查标签是否已存在
    const existingTag = await searchTagByName(tagName);
    
    if (existingTag) {
      // 标签存在，直接添加到客户
      await addTagToCustomer(existingTag.id);
      alert(`标签"${tagName}"已添加到客户`);
      // 清空搜索框
      searchInput.value = '';
      onTagSearchInput();
      // 刷新标签显示
      await loadTagModal();
    } else {
      // 标签不存在，显示维度选择弹窗
      showDimensionSelectionModal(tagName);
    }
  } catch (error) {
    console.error('添加标签时出错:', error);
    alert('添加标签失败，请稍后重试');
  }
}

// 标签搜索输入处理
function onTagSearchInput() {
  const searchInput = document.getElementById('tag-search-input');
  const clearBtn = document.querySelector('.clear-search');
  const addTagBtn = document.querySelector('.add-tag-btn');
  const createHint = document.querySelector('.create-tag-hint');
  
  if (searchInput && clearBtn) {
    const hasValue = searchInput.value.trim().length > 0;
    clearBtn.style.display = hasValue ? 'block' : 'none';
    
    // 更新添加标签按钮和创建提示的文本
    if (addTagBtn && createHint) {
      const searchValue = searchInput.value.trim();
      if (searchValue) {
        addTagBtn.innerHTML = `
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M12 5v14M5 12h14" stroke="#165DFF" stroke-width="2" stroke-linecap="round"/>
          </svg>
          ${searchValue}
        `;
        createHint.innerHTML = `
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M20.59 13.41l-7.17 7.17a2 2 0 01-2.83 0L2 12V2h10l8.59 8.59a2 2 0 010 2.82z" stroke="#165DFF" stroke-width="2"/>
            <path d="M7 7h.01" stroke="#165DFF" stroke-width="2"/>
          </svg>
          创建新标签"${searchValue}"
        `;
      } else {
        addTagBtn.innerHTML = `
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M12 5v14M5 12h14" stroke="#165DFF" stroke-width="2" stroke-linecap="round"/>
          </svg>
          添加标签
        `;
        createHint.innerHTML = `
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M20.59 13.41l-7.17 7.17a2 2 0 01-2.83 0L2 12V2h10l8.59 8.59a2 2 0 010 2.82z" stroke="#165DFF" stroke-width="2"/>
            <path d="M7 7h.01" stroke="#165DFF" stroke-width="2"/>
          </svg>
          创建新标签
        `;
      }
    }
  }
}

// 简化的标签维度数据加载
async function loadTagDimensions() {
  try {
    showTagsLoading(true);
    console.log('=== 开始加载标签维度数据 ===');
    
    // 等待全局配置加载完成
    if (typeof waitForGlobalConfig === 'function') {
      await waitForGlobalConfig();
      console.log('全局配置已加载');
    }
    
    const baseUrl = getApiBaseUrl();
    console.log('API基础URL:', baseUrl);
    
    const apiUrl = `${baseUrl}/api/v1/tag-dimensions`;
    console.log('完整API URL:', apiUrl);
    
    const response = await fetch(apiUrl);
    console.log('API响应状态:', response.status, response.statusText);
    
    if (response.ok) {
      const result = await response.json();
      console.log('=== API原始响应数据 ===');
      console.log(JSON.stringify(result, null, 2));
      
      tagDimensionsData = result.data || [];
      console.log('=== 解析后的标签维度数据 ===');
      console.log('数据数量:', tagDimensionsData.length);
      console.log('详细数据:', JSON.stringify(tagDimensionsData, null, 2));
      
      // 直接渲染标签，不需要复杂的客户数据处理
      renderTagDimensions(tagDimensionsData);
      showTagsError(''); // 清除错误信息
      console.log('=== 标签数据加载完成 ===');
    } else {
      console.error('加载标签维度失败:', response.status, response.statusText);
      showTagsError(`加载失败 (${response.status})`);
    }
  } catch (error) {
    console.error('=== 加载标签维度时出错 ===');
    console.error('错误详情:', error);
    console.error('错误堆栈:', error.stack);
    showTagsError('网络错误，请检查连接');
  } finally {
    showTagsLoading(false);
  }
}

// 简化版本：不需要加载客户数据
// 原来的 loadCurrentCustomerData 函数已被移除，因为我们只展示标签，不处理客户选择状态

// 显示/隐藏标签加载状态
function showTagsLoading(show) {
  const loadingElement = document.getElementById('tags-loading');
  if (loadingElement) {
    loadingElement.style.display = show ? 'block' : 'none';
  }
}

// 显示标签错误信息
function showTagsError(message) {
  const container = document.getElementById('tags-dimensions-container');
  if (container) {
    container.innerHTML = `<div class="error-message">${message}</div>`;
  }
}

// 简化的标签维度渲染
function renderTagDimensions(dimensions) {
  const container = document.getElementById('tags-dimensions-container');
  if (!container) {
    console.warn('标签容器未找到');
    return;
  }
  
  console.log('=== renderTagDimensions 被调用 ===');
  console.log('开始渲染标签维度:', dimensions);
  console.log('容器元素:', container);
  console.log('容器当前innerHTML:', container.innerHTML);
  console.log('容器样式display:', window.getComputedStyle(container).display);
  console.log('容器样式visibility:', window.getComputedStyle(container).visibility);
  
  if (!dimensions || dimensions.length === 0) {
    container.innerHTML = '<div class="no-data">暂无标签数据</div>';
    console.log('没有标签数据，显示暂无数据提示');
    return;
  }
  
  container.innerHTML = '';
  console.log('清空容器，开始渲染', dimensions.length, '个维度');
  
  dimensions.forEach((dimension, index) => {
    console.log(`渲染第${index + 1}个维度:`, dimension.name, '标签数量:', dimension.tags?.length || 0);
    
    const dimensionElement = document.createElement('div');
    dimensionElement.className = 'tag-category';
    
    // 维度标题
    const labelElement = document.createElement('div');
    labelElement.className = 'tag-label';
    labelElement.textContent = `${dimension.name}：`;
    
    // 标签列表
    const listElement = document.createElement('div');
    listElement.className = 'tag-list';
    
    // 渲染该维度下的所有标签
    if (dimension.tags && dimension.tags.length > 0) {
      dimension.tags.forEach((tag, tagIndex) => {
        console.log(`  渲染标签${tagIndex + 1}:`, tag.name);
        const tagElement = document.createElement('span');
        tagElement.className = 'tag';
        tagElement.textContent = tag.name;
        
        // 如果有自定义颜色，使用内联样式
        if (tag.color && tag.color !== '#e8f3ff') {
          tagElement.style.backgroundColor = tag.color;
        }
        
        if (tag.description) {
          tagElement.title = tag.description;
        }
        
        listElement.appendChild(tagElement);
      });
    } else {
      const emptyElement = document.createElement('span');
      emptyElement.className = 'empty-tags';
      emptyElement.textContent = '该维度暂无标签';
      emptyElement.style.color = '#999';
      emptyElement.style.fontStyle = 'italic';
      listElement.appendChild(emptyElement);
    }
    
    dimensionElement.appendChild(labelElement);
    dimensionElement.appendChild(listElement);
    container.appendChild(dimensionElement);
    
    console.log(`第${index + 1}个维度渲染完成`);
  });
  
  console.log('所有标签渲染完成，容器内容:', container.innerHTML.substring(0, 200) + '...');
  console.log('渲染后容器样式display:', window.getComputedStyle(container).display);
  console.log('渲染后容器offsetHeight:', container.offsetHeight);
  console.log('渲染后容器offsetWidth:', container.offsetWidth);
  console.log('渲染后容器子元素数量:', container.children.length);
}

// 简化版本：不需要标签切换功能
// 原来的 toggleCustomerTag 函数已被移除，因为我们只展示标签，不处理交互

// 打开标签编辑弹窗
// 打开标签弹窗
function openTagModal() {
  selectedTagDimension = null;
  
  // 显示弹窗
  const overlay = document.getElementById('tag-modal-overlay');
  if (overlay) {
    overlay.style.display = 'flex';
    // 初始化弹窗内容
    loadTagModal();
    // 添加点击外部区域关闭弹窗的事件监听器
    setTimeout(() => {
      overlay.addEventListener('click', handleOverlayClick);
    }, 100);
  } else {
    console.error('Tag modal overlay not found!');
  }
}


// 加载标签弹窗
async function loadTagModal() {
  // 获取现有的弹窗元素
  const overlay = document.getElementById('tag-modal-overlay');
  if (!overlay) {
    console.error('Tag modal overlay not found!');
    return;
  }
  
  // 首先加载全部标签展示，这会更新tagDimensionsData
  await loadModalAllTags();
  
  // 加载搜索相关组件
  loadSearchDimensionOptions();
}

// 加载维度选项
// 已移除不需要的维度选择和标签管理相关函数

// 已移除新增维度和标签表单相关函数

// 已移除标签创建和维度选择相关函数
// 已移除维度选择按钮样式更新和隐藏表单相关函数

// 已移除创建维度和标签相关函数

// 已移除更新和删除标签相关函数

// 处理点击弹窗外部区域
function handleOverlayClick(event) {
  // 只有点击的是overlay本身（不是弹窗内容）时才关闭弹窗
  if (event.target.id === 'tag-modal-overlay') {
    closeTagModal();
  }
}

// 关闭标签弹窗
function closeTagModal() {
  const overlay = document.getElementById('tag-modal-overlay');
  if (overlay) {
    overlay.style.display = 'none';
    // 移除事件监听器
    overlay.removeEventListener('click', handleOverlayClick);
    
    // 刷新标签tab页面
    if (typeof loadTagDimensions === 'function') {
      loadTagDimensions();
    }
  }
}


// 加载弹窗中的全部标签
async function loadModalAllTags() {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/tag-dimensions`);
    if (response.ok) {
      const result = await response.json();
      const dimensions = result.data || [];
      // 更新全局变量，确保搜索功能可以使用最新数据
      tagDimensionsData = dimensions;
      renderModalAllTags(dimensions);
    } else {
      console.error('加载标签维度失败:', response.status);
      const container = document.getElementById('modal-all-tags-container');
      if (container) {
        container.innerHTML = '<div class="error-message">加载标签失败</div>';
      }
    }
  } catch (error) {
    console.error('加载标签维度时出错:', error);
    const container = document.getElementById('modal-all-tags-container');
    if (container) {
      container.innerHTML = '<div class="error-message">网络错误，请检查连接</div>';
    }
  }
}

// 渲染弹窗中的全部标签
function renderModalAllTags(dimensions) {
  const container = document.getElementById('modal-all-tags-container');
  if (!container) return;
  
  container.innerHTML = '';
  
  dimensions.forEach(dimension => {
    const categoryDiv = document.createElement('div');
    categoryDiv.className = 'tag-category';
    
    const labelSpan = document.createElement('span');
    labelSpan.className = 'category-label';
    labelSpan.textContent = dimension.name + '：';
    
    const listDiv = document.createElement('div');
    listDiv.className = 'tag-list';
    
    if (dimension.tags && dimension.tags.length > 0) {
      dimension.tags.forEach(tag => {
        const tagSpan = document.createElement('span');
        tagSpan.className = 'tag-item';
        tagSpan.textContent = tag.name;
        tagSpan.style.color = tag.color || '#333'; // 使用标签的颜色
        listDiv.appendChild(tagSpan);
      });
    } else {
      const emptySpan = document.createElement('span');
      emptySpan.className = 'empty-tags';
      emptySpan.textContent = '该维度暂无标签';
      emptySpan.style.color = '#999';
      emptySpan.style.fontStyle = 'italic';
      listDiv.appendChild(emptySpan);
    }
    
    categoryDiv.appendChild(labelSpan);
    categoryDiv.appendChild(listDiv);
    container.appendChild(categoryDiv);
  });
}

// 加载搜索区域的维度选项
function loadSearchDimensionOptions() {
  const selector = document.getElementById('search-dimension-selector');
  if (!selector) return;
  
  selector.innerHTML = '<option value="">请选择维度</option>';
  
  tagDimensionsData.forEach(dimension => {
    const option = document.createElement('option');
    option.value = dimension.id;
    option.textContent = dimension.name;
    selector.appendChild(option);
  });
}

// 标签搜索功能
function onTagSearch() {
  const searchInput = document.getElementById('tag-search-input');
  const dimensionSelector = document.getElementById('search-dimension-selector');
  const resultsContainer = document.getElementById('search-results');
  
  if (!searchInput || !dimensionSelector || !resultsContainer) return;
  
  const searchTerm = searchInput.value.trim();
  const selectedDimensionId = dimensionSelector.value;
  
  if (!searchTerm) {
    resultsContainer.style.display = 'none';
    return;
  }
  
  if (!selectedDimensionId) {
    resultsContainer.innerHTML = '<div class="search-message">请先选择维度</div>';
    resultsContainer.style.display = 'block';
    return;
  }
  
  // 在选定维度中搜索标签
  const selectedDimension = tagDimensionsData.find(d => d.id == selectedDimensionId);
  if (!selectedDimension || !selectedDimension.tags) {
    showCreateTagOption(searchTerm, selectedDimensionId);
    return;
  }
  
  const matchingTags = selectedDimension.tags.filter(tag => 
    tag.name.toLowerCase().includes(searchTerm.toLowerCase())
  );
  
  if (matchingTags.length > 0) {
    displaySearchResults(matchingTags);
  } else {
    showCreateTagOption(searchTerm, selectedDimensionId);
  }
}

// 显示搜索结果
function displaySearchResults(tags) {
  const resultsContainer = document.getElementById('search-results');
  if (!resultsContainer) return;
  
  resultsContainer.innerHTML = '';
  
  tags.forEach(tag => {
    const tagDiv = document.createElement('div');
    tagDiv.className = 'search-result-item';
    tagDiv.innerHTML = `
      <span class="tag-name">${tag.name}</span>
      <span class="tag-desc">${tag.description || ''}</span>
    `;
    resultsContainer.appendChild(tagDiv);
  });
  
  resultsContainer.style.display = 'block';
}

// 显示创建标签选项
function showCreateTagOption(tagName, dimensionId) {
  const resultsContainer = document.getElementById('search-results');
  if (!resultsContainer) return;
  
  const dimension = tagDimensionsData.find(d => d.id == dimensionId);
  const dimensionName = dimension ? dimension.name : '未知维度';
  
  const messageDiv = document.createElement('div');
  messageDiv.className = 'search-message';
  messageDiv.textContent = `未找到标签 "${tagName}"`;
  
  const createDiv = document.createElement('div');
  createDiv.className = 'create-tag-option';
  createDiv.innerHTML = `
    <i class="icon-plus"></i>
    在 "${dimensionName}" 维度中创建标签 "${tagName}"
  `;
  createDiv.onclick = () => createTagFromSearch(tagName, dimensionId);
  
  resultsContainer.innerHTML = '';
  resultsContainer.appendChild(messageDiv);
  resultsContainer.appendChild(createDiv);
  resultsContainer.style.display = 'block';
}

// 从搜索创建标签
async function createTagFromSearch(tagName, dimensionId) {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/tags`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: tagName,
        dimension_id: parseInt(dimensionId),
        color: '#2196F3',
        description: `通过搜索创建的标签`,
        sort_order: 0
      })
    });
    
    if (response.ok) {
      alert(`标签 "${tagName}" 创建成功！`);
      // 重新加载标签数据
      await loadTagDimensions();
      // 重新加载弹窗中的全部标签
      loadModalAllTags();
      // 清空搜索框
      const searchInput = document.getElementById('tag-search-input');
      const resultsContainer = document.getElementById('search-results');
      if (searchInput) searchInput.value = '';
      if (resultsContainer) resultsContainer.style.display = 'none';
    } else {
      const error = await response.text();
      alert(`创建标签失败: ${error}`);
    }
  } catch (error) {
    console.error('Error creating tag:', error);
    alert('创建标签时出错');
  }
}

// 搜索标签（按名称）
async function searchTagByName(tagName) {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/tags?name=${encodeURIComponent(tagName)}`);
    if (response.ok) {
      const result = await response.json();
      const tags = result.data?.list || [];
      // 返回完全匹配的标签
      return tags.find(tag => tag.name === tagName) || null;
    }
  } catch (error) {
    console.error('搜索标签失败:', error);
  }
  return null;
}

// 添加标签到客户
async function addTagToCustomer(tagId) {
  const customerId = getCurrentCustomerId();
  try {
    // 先获取当前客户的系统标签
    const customerResponse = await fetch(`${getApiBaseUrl()}/api/v1/customers/${customerId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (!customerResponse.ok) {
      throw new Error('获取客户信息失败');
    }
    
    const customerData = await customerResponse.json();
    const currentSystemTags = customerData.data.system_tags || [];
    
    // 添加新标签ID到现有标签列表中（如果不存在的话）
    const updatedSystemTags = [...currentSystemTags];
    if (!updatedSystemTags.includes(tagId)) {
      updatedSystemTags.push(tagId);
    }
    
    const response = await fetch(`${getApiBaseUrl()}/api/v1/customers/${customerId}/system-tags`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        system_tags: updatedSystemTags
      })
    });
    
    if (!response.ok) {
      throw new Error('添加标签到客户失败');
    }
  } catch (error) {
    console.error('添加标签到客户失败:', error);
    throw error;
  }
}

// 显示维度选择弹窗
async function showDimensionSelectionModal(tagName) {
  try {
    // 获取所有维度
    const response = await fetch(`${getApiBaseUrl()}/api/v1/tag-dimensions`);
    if (!response.ok) {
      throw new Error('获取维度列表失败');
    }
    
    const result = await response.json();
    const dimensions = result.data || [];
    
    // 创建维度选择弹窗
    createDimensionSelectionModal(dimensions, tagName);
  } catch (error) {
    console.error('获取维度列表失败:', error);
    alert('获取维度列表失败，请稍后重试');
  }
}

// 创建维度选择弹窗HTML
function createDimensionSelectionModal(dimensions, tagName) {
  // 移除已存在的弹窗
  const existingModal = document.getElementById('dimension-selection-modal');
  if (existingModal) {
    existingModal.remove();
  }
  
  const modalHtml = `
    <div id="dimension-selection-modal" class="modal-overlay" style="display: flex;">
      <div class="modal-container">
        <div class="modal-header">
          <h3 class="modal-title">选择标签维度</h3>
          <button class="modal-close" onclick="closeDimensionSelectionModal()">×</button>
        </div>
        <div class="modal-content">
          <p>标签"${tagName}"不存在，请选择要创建的维度：</p>
          <div class="dimension-list">
            ${dimensions.map(dimension => `
              <label class="dimension-item">
                <input type="radio" name="dimension" value="${dimension.id}">
                <span class="dimension-name">${dimension.name}</span>
                <span class="dimension-desc">${dimension.description || ''}</span>
              </label>
            `).join('')}
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" onclick="closeDimensionSelectionModal()">取消</button>
          <button class="btn-primary" onclick="createTagWithDimension('${tagName}')">创建标签</button>
        </div>
      </div>
    </div>
  `;
  
  document.body.insertAdjacentHTML('beforeend', modalHtml);
  
  // 添加点击背景关闭事件
  document.getElementById('dimension-selection-modal').addEventListener('click', function(e) {
    if (e.target === this) {
      closeDimensionSelectionModal();
    }
  });
}

// 关闭维度选择弹窗
function closeDimensionSelectionModal() {
  const modal = document.getElementById('dimension-selection-modal');
  if (modal) {
    modal.remove();
  }
}

// 使用选定的维度创建标签
async function createTagWithDimension(tagName) {
  const selectedDimension = document.querySelector('input[name="dimension"]:checked');
  
  if (!selectedDimension) {
    alert('请选择一个维度');
    return;
  }
  
  const dimensionId = parseInt(selectedDimension.value);
  
  try {
    // 创建新标签
    const response = await fetch(`${getApiBaseUrl()}/api/v1/tags`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        dimension_id: dimensionId,
        name: tagName,
        color: getRandomTagColor(),
        description: `用户创建的标签: ${tagName}`
      })
    });
    
    if (!response.ok) {
      throw new Error('创建标签失败');
    }
    
    const result = await response.json();
    const newTag = result.data;
    
    // 关闭维度选择弹窗
    closeDimensionSelectionModal();
    
    // 添加新标签到客户
    await addTagToCustomer(newTag.id);
    
    alert(`标签"${tagName}"创建成功并已添加到客户`);
    
    // 清空搜索框并刷新
    const searchInput = document.getElementById('tag-search-input');
    if (searchInput) {
      searchInput.value = '';
      onTagSearchInput();
    }
    
    // 刷新标签显示
    await loadTagModal();
    
  } catch (error) {
    console.error('创建标签失败:', error);
    alert('创建标签失败，请稍后重试');
  }
}

// 获取随机标签颜色
function getRandomTagColor() {
  const colors = [
    '#2196F3', '#4CAF50', '#FF9800', '#9C27B0', 
    '#F44336', '#009688', '#795548', '#607D8B',
    '#3F51B5', '#8BC34A', '#FF5722', '#E91E63'
  ];
  return colors[Math.floor(Math.random() * colors.length)];
}

// 切换标签选中状态
function toggleTag(tagElement) {
  tagElement.classList.toggle('selected');
  console.log('Tag toggled:', tagElement.textContent, tagElement.classList.contains('selected'));
}

// 将函数暴露到全局作用域
window.loadTagDimensions = loadTagDimensions;
window.renderTagDimensions = renderTagDimensions;
window.onTagSearchInput = onTagSearchInput;
window.onTagSearch = onTagSearch;
window.createTagFromSearch = createTagFromSearch;
window.openTagModal = openTagModal;
window.closeTagModal = closeTagModal;
window.loadTagModal = loadTagModal;
window.handleOverlayClick = handleOverlayClick;
window.toggleTag = toggleTag;
window.clearTagSearch = clearTagSearch;
window.addNewTag = addNewTag;
window.closeDimensionSelectionModal = closeDimensionSelectionModal;
window.createTagWithDimension = createTagWithDimension;