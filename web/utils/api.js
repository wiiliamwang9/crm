// API工具类
class ApiClient {
    constructor() {
        this.baseURL = '/api';
    }

    async request(url, options = {}) {
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        };

        try {
            const response = await fetch(this.baseURL + url, config);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    async get(url, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const fullUrl = queryString ? `${url}?${queryString}` : url;
        return this.request(fullUrl, { method: 'GET' });
    }

    async post(url, data = {}) {
        return this.request(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async put(url, data = {}) {
        return this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    async delete(url) {
        return this.request(url, { method: 'DELETE' });
    }
}

// 客户API
class CustomerAPI extends ApiClient {
    constructor() {
        super();
        this.basePath = '/customers';
    }

    async getCustomers(page = 1, limit = 20, search = '') {
        return this.get(this.basePath, { page, limit, search });
    }

    async getCustomer(id) {
        return this.get(`${this.basePath}/${id}`);
    }

    async createCustomer(customerData) {
        return this.post(this.basePath, customerData);
    }

    async updateCustomer(id, customerData) {
        return this.put(`${this.basePath}/${id}`, customerData);
    }

    async deleteCustomer(id) {
        return this.delete(`${this.basePath}/${id}`);
    }

    async searchCustomers(searchData) {
        return this.post(`${this.basePath}/search`, searchData);
    }
}

// 待办API
class TodoAPI extends ApiClient {
    constructor() {
        super();
        this.basePath = '/todos';
    }

    async getTodos(params = {}) {
        return this.get(this.basePath, params);
    }

    async getTodo(id) {
        return this.get(`${this.basePath}/${id}`);
    }

    async createTodo(todoData) {
        return this.post(this.basePath, todoData);
    }

    async updateTodo(id, todoData) {
        return this.put(`${this.basePath}/${id}`, todoData);
    }

    async deleteTodo(id) {
        return this.delete(`${this.basePath}/${id}`);
    }

    async completeTodo(id) {
        return this.post(`${this.basePath}/${id}/complete`);
    }

    async cancelTodo(id) {
        return this.post(`${this.basePath}/${id}/cancel`);
    }

    async getTodoStats(customerID, executorID) {
        const params = {};
        if (customerID) params.customer_id = customerID;
        if (executorID) params.executor_id = executorID;
        return this.get(`${this.basePath}/stats`, params);
    }
}

// 导出API实例
const customerAPI = new CustomerAPI();
const todoAPI = new TodoAPI();

export { customerAPI, todoAPI, ApiClient };