// 跟进记录管理模块

// 全局变量
window.activitiesModule = {
  currentPage: 1,
  pageSize: 5,
  totalRecords: 0,
  isLoading: false,
  allActivities: [],
  currentEditingActivity: null,
  displayedCount: 0  // 当前显示的记录数量
};

// 初始化跟进记录模块
function initActivitiesModule() {
  console.log('初始化跟进记录模块');
}

// 加载跟进记录列表
async function loadFollowUpRecords(reset = true) {
  if (window.activitiesModule.isLoading) return;
  
  window.activitiesModule.isLoading = true;
  const customerId = getCurrentCustomerId();
  
  // 如果是重置加载，重置相关状态
  if (reset) {
    window.activitiesModule.currentPage = 1;
    window.activitiesModule.allActivities = [];
    window.activitiesModule.displayedCount = 0;
  }
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/activities/customer/${customerId}?page=${window.activitiesModule.currentPage}&page_size=${window.activitiesModule.pageSize}`);
    
    if (response.ok) {
      const result = await response.json();
      const data = result.data || {};
      const activities = data.list || [];
      
      console.log('跟进记录加载成功:', activities);
      
      window.activitiesModule.totalRecords = data.total || 0;
      
      if (reset) {
        // 首次加载，替换所有记录
        window.activitiesModule.allActivities = activities;
        window.activitiesModule.displayedCount = activities.length;
        updateFollowUpListDisplay(activities);
      } else {
        // 加载更多，追加记录
        window.activitiesModule.allActivities.push(...activities);
        window.activitiesModule.displayedCount += activities.length;
        appendFollowUpRecords(activities);
      }
      
      updateLoadMoreButton();
    } else {
      console.error('加载跟进记录失败:', response.status, response.statusText);
      showFollowUpError('加载跟进记录失败');
    }
  } catch (error) {
    console.error('Error loading follow-up records:', error);
    showFollowUpError('网络错误，请稍后重试');
  } finally {
    window.activitiesModule.isLoading = false;
  }
}

// 更新跟进记录列表显示
function updateFollowUpListDisplay(activities) {
  const container = document.getElementById('follow-up-list-container');
  if (!container) return;
  
  container.innerHTML = '';
  
  if (!activities || activities.length === 0) {
    container.innerHTML = '<div class="no-activities">暂无跟进记录</div>';
    return;
  }
  
  console.log('Displaying activities:', activities);
  
  activities.forEach(activity => {
    console.log('Creating element for activity:', activity.id, activity.kind);
    const activityElement = createFollowUpElement(activity);
    container.appendChild(activityElement);
  });
}

// 创建跟进记录元素
function createFollowUpElement(activity) {
  const div = document.createElement('div');
  div.className = 'follow-up-item';
  div.setAttribute('data-activity-id', activity.id);
  
  // 确定图标类型和样式
  const iconInfo = getActivityIconInfo(activity.kind);
  
  // 构建详情内容
  let detailsHtml = '';
  if (activity.content) {
    detailsHtml += `<div class="follow-up-detail">内容：${activity.content}</div>`;
  }
  if (activity.amount > 0) {
    detailsHtml += `<div class="follow-up-detail">金额：${activity.amount}</div>`;
  }
  if (activity.cost > 0) {
    detailsHtml += `<div class="follow-up-detail">成本：${activity.cost}</div>`;
  }
  if (activity.feedback) {
    detailsHtml += `<div class="follow-up-detail">反馈：${activity.feedback}</div>`;
  }
  if (activity.satisfaction > 0) {
    detailsHtml += `<div class="follow-up-detail">满意度：${activity.satisfaction}分</div>`;
  }
  
  // 构建操作按钮
  let actionHtml = '';
  if (activity.kind === 'order' || activity.kind === 'sample' || activity.kind === 'feedback' || activity.kind === 'complaint') {
    if (activity.feedback && activity.feedback.trim() !== '') {
      actionHtml = `<button class="action-btn" onclick="openEditFeedbackModal(${activity.id})">编辑反馈</button>`;
    } else {
      actionHtml = `<button class="action-btn" onclick="openAddFeedbackModal(${activity.id})">添加反馈</button>`;
    }
  }
  
  console.log('Activity:', activity.id, 'Kind:', activity.kind, 'Feedback:', activity.feedback, 'ActionHTML:', actionHtml);
  
  div.innerHTML = `
    <div class="follow-up-header">
      <div class="follow-up-icon ${iconInfo.className}">
        <span>${iconInfo.symbol}</span>
      </div>
      <div class="follow-up-title">${activity.title || getActivityKindDisplayName(activity.kind)}</div>
      <div class="follow-up-time">${activity.time_ago || formatTimeAgo(activity.created_at)}</div>
    </div>
    ${detailsHtml ? `<div class="follow-up-details">${detailsHtml}</div>` : ''}
    ${actionHtml ? `<div class="follow-up-action">${actionHtml}</div>` : ''}
  `;
  
  return div;
}

// 获取活动图标信息
function getActivityIconInfo(kind) {
  const iconMap = {
    'call': { className: 'info', symbol: '📞' },
    'visit': { className: 'info', symbol: '🏢' },
    'email': { className: 'info', symbol: '📧' },
    'wechat': { className: 'info', symbol: '💬' },
    'meeting': { className: 'info', symbol: '🤝' },
    'order': { className: 'success', symbol: '✓' },
    'sample': { className: 'warning', symbol: '📦' },
    'feedback': { className: 'info', symbol: 'i' },
    'complaint': { className: 'error', symbol: '⚠' },
    'payment': { className: 'success', symbol: '💰' },
    'other': { className: 'info', symbol: 'i' }
  };
  
  return iconMap[kind] || iconMap['other'];
}

// 获取活动类型显示名称
function getActivityKindDisplayName(kind) {
  const kindNames = {
    'call': '电话沟通',
    'visit': '实地拜访',
    'email': '邮件',
    'wechat': '微信沟通',
    'meeting': '会议洽谈',
    'order': '下单记录',
    'sample': '发样记录',
    'feedback': '客户反馈',
    'complaint': '客户投诉',
    'payment': '付款记录',
    'other': '其他'
  };
  
  return kindNames[kind] || '未知类型';
}

// 格式化时间显示
function formatTimeAgo(dateString) {
  if (!dateString) return '';
  
  try {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffHours = diffMs / (1000 * 60 * 60);
    
    if (diffHours < 1) {
      return '刚刚';
    } else if (diffHours < 24) {
      return '今天';
    } else if (diffHours < 48) {
      return '昨天';
    } else if (diffHours < 72) {
      return '2天前';
    } else {
      const diffDays = Math.floor(diffHours / 24);
      return `${diffDays}天前`;
    }
  } catch (error) {
    console.error('时间格式化错误:', error);
    return dateString;
  }
}

// 追加跟进记录到列表
function appendFollowUpRecords(activities) {
  const container = document.getElementById('follow-up-list-container');
  if (!container) return;
  
  activities.forEach(activity => {
    const activityElement = createFollowUpElement(activity);
    container.appendChild(activityElement);
  });
}

// 更新加载更多按钮状态
function updateLoadMoreButton() {
  const button = document.getElementById('load-more-btn');
  if (!button) return;
  
  const hasMoreRecords = window.activitiesModule.displayedCount < window.activitiesModule.totalRecords;
  
  if (hasMoreRecords) {
    button.style.display = 'block';
    const span = button.querySelector('span');
    if (span) {
      span.textContent = `查看更多 (已显示 ${window.activitiesModule.displayedCount}/${window.activitiesModule.totalRecords} 条)`;
    }
  } else {
    button.style.display = 'none';
  }
}

// 加载更多跟进记录
async function loadMoreActivities() {
  if (window.activitiesModule.isLoading) return;
  
  // 检查是否还有更多记录
  if (window.activitiesModule.displayedCount >= window.activitiesModule.totalRecords) {
    return;
  }
  
  window.activitiesModule.currentPage++;
  
  // 调用loadFollowUpRecords，传入reset=false表示追加加载
  await loadFollowUpRecords(false);
}

// 显示跟进记录错误信息
function showFollowUpError(message) {
  const container = document.getElementById('follow-up-list-container');
  if (container) {
    container.innerHTML = `<div class="error-message">${message}</div>`;
  }
}

// 打开添加反馈弹窗
function openAddFeedbackModal(activityId) {
  console.log('Opening add feedback modal for activity:', activityId);
  const activity = window.activitiesModule.allActivities.find(a => a.id == activityId);
  if (!activity) {
    console.error('Activity not found:', activityId);
    alert('找不到跟进记录，请刷新页面重试');
    return;
  }
  
  window.activitiesModule.currentEditingActivity = activity;
  
  // 创建反馈弹窗HTML（如果不存在）
  if (!document.getElementById('feedbackModal')) {
    createFeedbackModal();
  }
  
  // 重置表单
  document.getElementById('feedbackContent').value = '';
  document.getElementById('feedbackSatisfaction').value = '5';
  
  // 设置标题
  const modal = document.getElementById('feedbackModal');
  const title = modal.querySelector('.modal-title');
  if (title) {
    title.textContent = `为"${activity.title || getActivityKindDisplayName(activity.kind)}"添加反馈`;
  }
  
  // 显示弹窗
  modal.style.display = 'flex';
  document.body.style.overflow = 'hidden';
}

// 打开编辑反馈弹窗
function openEditFeedbackModal(activityId) {
  console.log('Opening edit feedback modal for activity:', activityId);
  const activity = window.activitiesModule.allActivities.find(a => a.id == activityId);
  if (!activity) {
    console.error('Activity not found:', activityId);
    alert('找不到跟进记录，请刷新页面重试');
    return;
  }
  
  window.activitiesModule.currentEditingActivity = activity;
  
  // 创建反馈弹窗HTML（如果不存在）
  if (!document.getElementById('feedbackModal')) {
    createFeedbackModal();
  }
  
  // 填充现有数据
  document.getElementById('feedbackContent').value = activity.feedback || '';
  document.getElementById('feedbackSatisfaction').value = activity.satisfaction || '5';
  
  // 设置标题
  const modal = document.getElementById('feedbackModal');
  const title = modal.querySelector('.modal-title');
  if (title) {
    title.textContent = `编辑"${activity.title || getActivityKindDisplayName(activity.kind)}"的反馈`;
  }
  
  // 显示弹窗
  modal.style.display = 'flex';
  document.body.style.overflow = 'hidden';
}

// 创建反馈弹窗HTML
function createFeedbackModal() {
  const modalHtml = `
    <div id="feedbackModal" class="modal-overlay" style="display: none;">
      <div class="modal-container">
        <div class="modal-header">
          <h3 class="modal-title">添加反馈</h3>
          <button class="modal-close" onclick="closeFeedbackModal()">×</button>
        </div>
        <div class="modal-content">
          <div class="form-row">
            <label class="form-label">反馈内容</label>
            <textarea class="form-textarea" id="feedbackContent" placeholder="请输入客户反馈内容..." rows="4"></textarea>
          </div>
          <div class="form-row">
            <label class="form-label">客户满意度</label>
            <select class="form-select" id="feedbackSatisfaction">
              <option value="1">1分 - 非常不满意</option>
              <option value="2">2分 - 不满意</option>
              <option value="3">3分 - 一般</option>
              <option value="4">4分 - 满意</option>
              <option value="5" selected>5分 - 非常满意</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" onclick="closeFeedbackModal()">取消</button>
          <button class="btn-primary" onclick="saveFeedback()">保存反馈</button>
        </div>
      </div>
    </div>
  `;
  
  document.body.insertAdjacentHTML('beforeend', modalHtml);
  
  // 添加点击背景关闭事件
  document.getElementById('feedbackModal').addEventListener('click', function(e) {
    if (e.target === this) {
      closeFeedbackModal();
    }
  });
}

// 关闭反馈弹窗
function closeFeedbackModal() {
  const modal = document.getElementById('feedbackModal');
  if (modal) {
    modal.style.display = 'none';
    document.body.style.overflow = 'auto';
    window.activitiesModule.currentEditingActivity = null;
  }
}

// 保存反馈
async function saveFeedback() {
  const activity = window.activitiesModule.currentEditingActivity;
  if (!activity) return;
  
  const feedback = document.getElementById('feedbackContent').value.trim();
  const satisfaction = parseInt(document.getElementById('feedbackSatisfaction').value);
  
  if (!feedback) {
    alert('请输入反馈内容');
    return;
  }
  
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/v1/activities/${activity.id}/feedback`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        feedback: feedback,
        satisfaction: satisfaction
      })
    });
    
    if (response.ok) {
      alert('反馈保存成功！');
      closeFeedbackModal();
      
      // 刷新跟进记录列表
      window.activitiesModule.currentPage = 1;
      loadFollowUpRecords();
    } else {
      const errorText = await response.text();
      console.error('保存反馈失败:', errorText);
      alert('保存失败，请稍后重试');
    }
  } catch (error) {
    console.error('Error saving feedback:', error);
    alert('网络错误，请稍后重试');
  }
}

// 刷新跟进记录列表
function refreshFollowUpRecords() {
  loadFollowUpRecords(true);
}

// 根据记录类型过滤跟进记录
function filterActivitiesByKind() {
  const filterSelect = document.getElementById('activity-kind-filter');
  const selectedKind = filterSelect.value;
  
  // 获取所有跟进记录元素
  const activityElements = document.querySelectorAll('.follow-up-item');
  
  activityElements.forEach(element => {
    const activityId = element.getAttribute('data-activity-id');
    const activity = window.activitiesModule.allActivities.find(a => a.id == activityId);
    
    if (!activity) {
      element.style.display = 'none';
      return;
    }
    
    // 如果选择"全部"或者记录类型匹配，则显示
    if (selectedKind === 'all' || activity.kind === selectedKind) {
      element.style.display = 'block';
    } else {
      element.style.display = 'none';
    }
  });
  
  // 更新显示状态
  updateFilteredDisplay();
}

// 更新过滤后的显示状态
function updateFilteredDisplay() {
  const container = document.getElementById('follow-up-list-container');
  const visibleItems = container.querySelectorAll('.follow-up-item[style*="block"], .follow-up-item:not([style*="none"])');
  
  // 如果没有可见的记录，显示提示信息
  let noResultsMsg = container.querySelector('.no-filtered-results');
  if (visibleItems.length === 0) {
    if (!noResultsMsg) {
      noResultsMsg = document.createElement('div');
      noResultsMsg.className = 'no-filtered-results';
      noResultsMsg.textContent = '没有符合条件的跟进记录';
      container.appendChild(noResultsMsg);
    }
    noResultsMsg.style.display = 'block';
  } else {
    if (noResultsMsg) {
      noResultsMsg.style.display = 'none';
    }
  }
}

// 导出函数到全局作用域
window.loadFollowUpRecords = loadFollowUpRecords;
window.loadMoreActivities = loadMoreActivities;
window.openAddFeedbackModal = openAddFeedbackModal;
window.openEditFeedbackModal = openEditFeedbackModal;
window.closeFeedbackModal = closeFeedbackModal;
window.saveFeedback = saveFeedback;
window.refreshFollowUpRecords = refreshFollowUpRecords;
window.filterActivitiesByKind = filterActivitiesByKind;