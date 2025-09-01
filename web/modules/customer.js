// 客户模块
import { customerAPI } from '../utils/api.js';
import { formatDate, showMessage, showConfirm, showLoading } from '../utils/utils.js';
import Modal from '../components/Modal.js';

class CustomerModule {
    constructor() {
        this.currentPage = 1;
        this.pageSize = 20;
        this.searchKeyword = '';
        this.customers = [];
        this.total = 0;
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.loadCustomers();
    }
    
    bindEvents() {
        // 搜索框事件
        const searchInput = document.getElementById('searchInput');
        if (searchInput) {
            searchInput.addEventListener('input', this.debounce((e) => {
                this.searchKeyword = e.target.value;
                this.currentPage = 1;
                this.loadCustomers();
            }, 500));
        }
        
        // 添加客户按钮
        const addCustomerBtn = document.getElementById('addCustomerBtn');
        if (addCustomerBtn) {
            addCustomerBtn.addEventListener('click', () => this.showCustomerModal());
        }
        
        // 导出按钮
        const exportBtn = document.getElementById('exportBtn');
        if (exportBtn) {
            exportBtn.addEventListener('click', () => this.exportCustomers());
        }
        
        // 导入按钮
        const importBtn = document.getElementById('importBtn');
        if (importBtn) {
            importBtn.addEventListener('click', () => this.importCustomers());
        }
    }
    
    async loadCustomers() {
        const container = document.getElementById('customerList');
        if (!container) return;
        
        const loading = showLoading(container, '加载客户数据...');
        
        try {
            const response = await customerAPI.getCustomers(this.currentPage, this.pageSize, this.searchKeyword);
            this.customers = response.customers || [];
            this.total = response.total || 0;
            
            this.renderCustomerList();
            this.renderPagination();
            
        } catch (error) {
            console.error('加载客户失败:', error);
            showMessage('加载客户失败，请重试', 'error');
        } finally {
            loading.hide();
        }
    }
    
    renderCustomerList() {
        const container = document.getElementById('customerList');
        if (!container) return;
        
        if (this.customers.length === 0) {
            container.innerHTML = `
                <div class="empty-state" style="text-align: center; padding: 60px 20px; color: #999;">
                    <div style="font-size: 48px; margin-bottom: 16px;">📋</div>
                    <div style="font-size: 16px; margin-bottom: 8px;">暂无客户数据</div>
                    <div style="font-size: 14px;">点击"添加客户"开始管理您的客户</div>
                </div>
            `;
            return;
        }
        
        const customerCards = this.customers.map(customer => this.renderCustomerCard(customer)).join('');
        container.innerHTML = customerCards;
        
        // 绑定卡片事件
        this.bindCustomerCardEvents();
    }
    
    renderCustomerCard(customer) {
        const mainPhone = customer.phones && customer.phones.length > 0 ? customer.phones[0] : '';
        const mainWechat = customer.wechats && customer.wechats.length > 0 ? customer.wechats[0] : '';
        
        return `
            <div class="customer-card" data-id="${customer.id}" style="
                border: 1px solid #e8e8e8;
                border-radius: 8px;
                padding: 20px;
                margin-bottom: 16px;
                background: white;
                transition: all 0.3s ease;
                cursor: pointer;
            ">
                <div class="customer-header" style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px;">
                    <div>
                        <h3 style="margin: 0 0 4px 0; font-size: 18px; color: #333;">${customer.name || '未命名客户'}</h3>
                        <p style="margin: 0; color: #666; font-size: 14px;">联系人: ${customer.contact_name || '未知'}</p>
                    </div>
                    <div class="customer-actions">
                        <button class="btn-edit" data-id="${customer.id}" style="
                            background: #1890ff;
                            color: white;
                            border: none;
                            padding: 6px 12px;
                            border-radius: 4px;
                            cursor: pointer;
                            margin-right: 8px;
                            font-size: 12px;
                        ">编辑</button>
                        <button class="btn-delete" data-id="${customer.id}" style="
                            background: #ff4d4f;
                            color: white;
                            border: none;
                            padding: 6px 12px;
                            border-radius: 4px;
                            cursor: pointer;
                            font-size: 12px;
                        ">删除</button>
                    </div>
                </div>
                
                <div class="customer-info" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 12px;">
                    ${mainPhone ? `
                        <div style="display: flex; align-items: center;">
                            <span style="color: #666; margin-right: 8px;">📞</span>
                            <span>${mainPhone}</span>
                        </div>
                    ` : ''}
                    
                    ${mainWechat ? `
                        <div style="display: flex; align-items: center;">
                            <span style="color: #666; margin-right: 8px;">💬</span>
                            <span>${mainWechat}</span>
                        </div>
                    ` : ''}
                    
                    <div style="display: flex; align-items: center;">
                        <span style="color: #666; margin-right: 8px;">📅</span>
                        <span>${formatDate(customer.created_at)}</span>
                    </div>
                </div>
                
                ${customer.remark ? `
                    <div class="customer-remark" style="margin-top: 12px; padding-top: 12px; border-top: 1px solid #f0f0f0;">
                        <span style="color: #999; font-size: 12px;">备注: </span>
                        <span style="color: #666; font-size: 14px;">${customer.remark}</span>
                    </div>
                ` : ''}
            </div>
        `;
    }
    
    bindCustomerCardEvents() {
        // 编辑按钮
        document.querySelectorAll('.btn-edit').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const customerId = parseInt(btn.dataset.id);
                this.showCustomerModal(customerId);
            });
        });
        
        // 删除按钮
        document.querySelectorAll('.btn-delete').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const customerId = parseInt(btn.dataset.id);
                this.deleteCustomer(customerId);
            });
        });
        
        // 卡片点击查看详情
        document.querySelectorAll('.customer-card').forEach(card => {
            card.addEventListener('click', () => {
                const customerId = parseInt(card.dataset.id);
                this.showCustomerDetail(customerId);
            });
        });
    }
    
    showCustomerModal(customerId = null) {
        const isEdit = customerId !== null;
        const title = isEdit ? '编辑客户' : '添加客户';
        
        const modalContent = document.createElement('div');
        modalContent.innerHTML = `
            <form id="customerForm" style="display: grid; gap: 16px;">
                <div>
                    <label style="display: block; margin-bottom: 4px; font-weight: 500;">客户名称 *</label>
                    <input type="text" id="customerName" required style="
                        width: 100%;
                        padding: 8px 12px;
                        border: 1px solid #d9d9d9;
                        border-radius: 4px;
                        box-sizing: border-box;
                    ">
                </div>
                
                <div>
                    <label style="display: block; margin-bottom: 4px; font-weight: 500;">联系人</label>
                    <input type="text" id="contactName" style="
                        width: 100%;
                        padding: 8px 12px;
                        border: 1px solid #d9d9d9;
                        border-radius: 4px;
                        box-sizing: border-box;
                    ">
                </div>
                
                <div>
                    <label style="display: block; margin-bottom: 4px; font-weight: 500;">手机号码</label>
                    <input type="text" id="phones" placeholder="多个号码用逗号分隔" style="
                        width: 100%;
                        padding: 8px 12px;
                        border: 1px solid #d9d9d9;
                        border-radius: 4px;
                        box-sizing: border-box;
                    ">
                </div>
                
                <div>
                    <label style="display: block; margin-bottom: 4px; font-weight: 500;">微信号</label>
                    <input type="text" id="wechats" placeholder="多个微信号用逗号分隔" style="
                        width: 100%;
                        padding: 8px 12px;
                        border: 1px solid #d9d9d9;
                        border-radius: 4px;
                        box-sizing: border-box;
                    ">
                </div>
                
                <div>
                    <label style="display: block; margin-bottom: 4px; font-weight: 500;">备注</label>
                    <textarea id="remark" rows="3" style="
                        width: 100%;
                        padding: 8px 12px;
                        border: 1px solid #d9d9d9;
                        border-radius: 4px;
                        box-sizing: border-box;
                        resize: vertical;
                    "></textarea>
                </div>
                
                <div style="text-align: right; padding-top: 16px; border-top: 1px solid #f0f0f0;">
                    <button type="button" id="cancelBtn" style="
                        padding: 8px 16px;
                        border: 1px solid #d9d9d9;
                        background: white;
                        color: #333;
                        border-radius: 4px;
                        cursor: pointer;
                        margin-right: 12px;
                    ">取消</button>
                    <button type="submit" style="
                        padding: 8px 16px;
                        border: none;
                        background: #1890ff;
                        color: white;
                        border-radius: 4px;
                        cursor: pointer;
                    ">${isEdit ? '更新' : '添加'}</button>
                </div>
            </form>
        `;
        
        const modal = new Modal({
            title,
            content: modalContent,
            width: '500px',
            closable: true
        });
        
        modal.show();
        
        // 如果是编辑模式，加载客户数据
        if (isEdit) {
            this.loadCustomerData(customerId, modal);
        }
        
        // 绑定表单事件
        this.bindCustomerFormEvents(modal, customerId);
    }
    
    async loadCustomerData(customerId, modal) {
        try {
            const customer = await customerAPI.getCustomer(customerId);
            
            document.getElementById('customerName').value = customer.name || '';
            document.getElementById('contactName').value = customer.contact_name || '';
            document.getElementById('phones').value = (customer.phones || []).join(', ');
            document.getElementById('wechats').value = (customer.wechats || []).join(', ');
            document.getElementById('remark').value = customer.remark || '';
            
        } catch (error) {
            console.error('加载客户数据失败:', error);
            showMessage('加载客户数据失败', 'error');
            modal.hide();
        }
    }
    
    bindCustomerFormEvents(modal, customerId) {
        const form = document.getElementById('customerForm');
        const cancelBtn = document.getElementById('cancelBtn');
        
        cancelBtn.addEventListener('click', () => modal.hide());
        
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = {
                name: document.getElementById('customerName').value.trim(),
                contact_name: document.getElementById('contactName').value.trim(),
                phones: document.getElementById('phones').value.split(',').map(p => p.trim()).filter(p => p),
                wechats: document.getElementById('wechats').value.split(',').map(w => w.trim()).filter(w => w),
                remark: document.getElementById('remark').value.trim()
            };
            
            if (!formData.name) {
                showMessage('请输入客户名称', 'error');
                return;
            }
            
            try {
                if (customerId) {
                    await customerAPI.updateCustomer(customerId, formData);
                    showMessage('客户更新成功', 'success');
                } else {
                    await customerAPI.createCustomer(formData);
                    showMessage('客户添加成功', 'success');
                }
                
                modal.hide();
                this.loadCustomers();
                
            } catch (error) {
                console.error('保存客户失败:', error);
                showMessage('保存失败，请重试', 'error');
            }
        });
    }
    
    async deleteCustomer(customerId) {
        const customer = this.customers.find(c => c.id === customerId);
        if (!customer) return;
        
        showConfirm(
            `确定要删除客户 "${customer.name}" 吗？此操作不可恢复。`,
            async () => {
                try {
                    await customerAPI.deleteCustomer(customerId);
                    showMessage('客户删除成功', 'success');
                    this.loadCustomers();
                } catch (error) {
                    console.error('删除客户失败:', error);
                    showMessage('删除失败，请重试', 'error');
                }
            }
        );
    }
    
    showCustomerDetail(customerId) {
        // 跳转到客户详情页面
        window.location.href = `/customer-detail.html?id=${customerId}`;
    }
    
    renderPagination() {
        const container = document.getElementById('pagination');
        if (!container) return;
        
        const totalPages = Math.ceil(this.total / this.pageSize);
        if (totalPages <= 1) {
            container.innerHTML = '';
            return;
        }
        
        let pagination = '<div class="pagination-wrapper" style="display: flex; justify-content: center; align-items: center; gap: 8px; margin-top: 20px;">';
        
        // 上一页
        if (this.currentPage > 1) {
            pagination += `<button class="page-btn" data-page="${this.currentPage - 1}" style="padding: 8px 12px; border: 1px solid #d9d9d9; background: white; cursor: pointer; border-radius: 4px;">上一页</button>`;
        }
        
        // 页码
        const startPage = Math.max(1, this.currentPage - 2);
        const endPage = Math.min(totalPages, this.currentPage + 2);
        
        for (let i = startPage; i <= endPage; i++) {
            const isActive = i === this.currentPage;
            pagination += `
                <button class="page-btn" data-page="${i}" style="
                    padding: 8px 12px;
                    border: 1px solid ${isActive ? '#1890ff' : '#d9d9d9'};
                    background: ${isActive ? '#1890ff' : 'white'};
                    color: ${isActive ? 'white' : '#333'};
                    cursor: pointer;
                    border-radius: 4px;
                ">${i}</button>
            `;
        }
        
        // 下一页
        if (this.currentPage < totalPages) {
            pagination += `<button class="page-btn" data-page="${this.currentPage + 1}" style="padding: 8px 12px; border: 1px solid #d9d9d9; background: white; cursor: pointer; border-radius: 4px;">下一页</button>`;
        }
        
        pagination += `<span style="margin-left: 16px; color: #666; font-size: 14px;">共 ${this.total} 条</span>`;
        pagination += '</div>';
        
        container.innerHTML = pagination;
        
        // 绑定分页事件
        container.querySelectorAll('.page-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                this.currentPage = parseInt(btn.dataset.page);
                this.loadCustomers();
            });
        });
    }
    
    debounce(func, wait) {
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
    
    async exportCustomers() {
        try {
            const response = await fetch('/api/export-excel', {
                method: 'GET',
            });
            
            if (!response.ok) {
                throw new Error('导出失败');
            }
            
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = `客户数据_${formatDate(new Date(), 'YYYY-MM-DD')}.xlsx`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
            
            showMessage('导出成功', 'success');
            
        } catch (error) {
            console.error('导出失败:', error);
            showMessage('导出失败，请重试', 'error');
        }
    }
    
    importCustomers() {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.xlsx,.xls';
        input.style.display = 'none';
        
        input.addEventListener('change', async (e) => {
            const file = e.target.files[0];
            if (!file) return;
            
            const formData = new FormData();
            formData.append('file', file);
            
            try {
                const response = await fetch('/api/upload-excel', {
                    method: 'POST',
                    body: formData
                });
                
                if (!response.ok) {
                    throw new Error('导入失败');
                }
                
                const result = await response.json();
                showMessage(`导入成功，共处理 ${result.count || 0} 条数据`, 'success');
                this.loadCustomers();
                
            } catch (error) {
                console.error('导入失败:', error);
                showMessage('导入失败，请重试', 'error');
            }
        });
        
        document.body.appendChild(input);
        input.click();
        document.body.removeChild(input);
    }
}

export default CustomerModule;