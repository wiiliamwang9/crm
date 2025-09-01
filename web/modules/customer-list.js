// 客户列表模块
class CustomerList {
    constructor() {
        this.currentPage = 1;
        this.totalPages = 1;
        this.currentSearch = '';
        this.API_BASE = '/api/v1';
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadCustomers();
    }

    bindEvents() {
        // 搜索框回车事件
        const searchInput = document.getElementById('searchInput');
        if (searchInput) {
            searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    this.searchCustomers();
                }
            });
        }
        
        // 全局化方法供HTML调用
        window.searchCustomers = () => this.searchCustomers();
        window.loadCustomers = () => this.loadCustomers();
        window.previousPage = () => this.previousPage();
        window.nextPage = () => this.nextPage();
        window.customerManager = this;
    }

    async loadCustomers(page = 1, search = '') {
        try {
            document.getElementById('loading').style.display = 'block';
            document.getElementById('customerTable').style.display = 'none';
            document.getElementById('pagination').style.display = 'none';
            
            const url = `${this.API_BASE}/customers?page=${page}&limit=20${search ? '&search=' + encodeURIComponent(search) : ''}`;
            const response = await fetch(url);
            
            if (!response.ok) {
                if (response.status === 404) {
                    throw new Error('API服务未启动或地址错误');
                }
                throw new Error(`获取客户数据失败 (${response.status})`);
            }
            
            const result = await response.json();
            const data = result.data || result; // 兼容新旧API格式
            
            if (data.customers && Array.isArray(data.customers)) {
                this.renderCustomers(data.customers);
                this.updatePagination(data.page || page, data.total || 0, data.limit || 20);
                
                this.currentPage = data.page || page;
                this.currentSearch = search;
                
                document.getElementById('customerTable').style.display = 'table';
                document.getElementById('pagination').style.display = 'flex';
            } else {
                throw new Error('返回数据格式错误');
            }
            
        } catch (error) {
            console.error('加载客户数据错误:', error);
            this.showMessage(`加载数据失败：${error.message}`, 'error');
            
            // 显示空表格
            document.getElementById('customerTableBody').innerHTML = 
                '<tr><td colspan="11" style="text-align: center; color: #6c757d; padding: 40px;">暂无数据或服务连接失败</td></tr>';
            document.getElementById('customerTable').style.display = 'table';
            
        } finally {
            document.getElementById('loading').style.display = 'none';
        }
    }

    renderCustomers(customers) {
        const tbody = document.getElementById('customerTableBody');
        tbody.innerHTML = '';
        
        if (!customers || customers.length === 0) {
            tbody.innerHTML = '<tr><td colspan="11" style="text-align: center; color: #6c757d; padding: 40px;">暂无客户数据</td></tr>';
            return;
        }
        
        customers.forEach((customer, index) => {
            const row = document.createElement('tr');
            
            // 计算序号，每页都从1开始
            const serialNumber = index + 1;
            
            const levelText = {
                0: '未分级', 1: 'S级', 2: 'A级', 3: 'B级', 4: 'C级', 10: 'X级'
            }[customer.level] || '未知';
            
            const stateText = {
                0: '未知', 1: '未开发', 2: '开发中', 3: '已开发', 
                4: '已拉黑', 5: '已倒闭', 6: '同事', 7: '叛徒', 8: '同行'
            }[customer.state] || '未知';
            
            const addressInfo = [customer.province, customer.city, customer.district, customer.address]
                .filter(Boolean).join(' ');
            
            const phones = customer.phones && Array.isArray(customer.phones) 
                ? customer.phones.join(', ') : (customer.phones || '');
            
            const wechats = customer.wechats && Array.isArray(customer.wechats) 
                ? customer.wechats.join(', ') : (customer.wechats || '');
            
            // 优先显示销售员姓名，如果没有则显示销售员ID
            const sellerDisplay = customer.saller_name || 
                (customer.sellers && Array.isArray(customer.sellers) 
                    ? customer.sellers.join(', ') : (customer.sellers || ''));
            
            const createdAt = customer.created_at 
                ? new Date(customer.created_at).toLocaleDateString('zh-CN') 
                : '';
            
            row.innerHTML = `
                <td>${serialNumber}</td>
                <td>${customer.name || ''}</td>
                <td>${customer.contact_name || ''}</td>
                <td class="array-field">${phones}</td>
                <td class="array-field">${wechats}</td>
                <td class="array-field" title="${addressInfo}">${addressInfo}</td>
                <td>${stateText}</td>
                <td>${sellerDisplay}</td>
                <td>${customer.import_source || customer.source || ''}</td>
                <td title="${customer.remark || ''}">${(customer.remark || '').substring(0, 30)}${(customer.remark || '').length > 30 ? '...' : ''}</td>
                <td>${createdAt}</td>
            `;
            tbody.appendChild(row);
        });
    }

    updatePagination(page, total, limit) {
        this.totalPages = Math.ceil(total / limit);
        const pageInfo = document.getElementById('pageInfo');
        pageInfo.textContent = `第 ${page} 页 / 共 ${this.totalPages} 页 (总计 ${total} 条记录)`;
        
        const prevBtn = document.getElementById('prevBtn');
        const nextBtn = document.getElementById('nextBtn');
        
        prevBtn.disabled = page <= 1;
        nextBtn.disabled = page >= this.totalPages;
        
        prevBtn.style.opacity = page <= 1 ? '0.5' : '1';
        nextBtn.style.opacity = page >= this.totalPages ? '0.5' : '1';
    }

    searchCustomers() {
        const search = document.getElementById('searchInput').value.trim();
        this.loadCustomers(1, search);
    }

    previousPage() {
        if (this.currentPage > 1) {
            this.loadCustomers(this.currentPage - 1, this.currentSearch);
        }
    }

    nextPage() {
        if (this.currentPage < this.totalPages) {
            this.loadCustomers(this.currentPage + 1, this.currentSearch);
        }
    }

    showMessage(message, type = 'success') {
        const messageDiv = document.getElementById('message');
        if (!messageDiv) return;
        
        messageDiv.className = `message ${type}`;
        messageDiv.textContent = message;
        messageDiv.style.display = 'block';
        
        // 自动隐藏消息
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, type === 'error' ? 8000 : 5000);
    }
}

export default CustomerList;