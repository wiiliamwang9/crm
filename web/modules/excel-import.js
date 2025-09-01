// Excel导入模块
class ExcelImport {
    constructor() {
        this.isUploading = false;
        this.API_BASE = '/api/v1';
        this.init();
    }

    init() {
        this.bindEvents();
    }

    bindEvents() {
        // 文件选择处理
        const excelFile = document.getElementById('excelFile');
        if (excelFile) {
            excelFile.addEventListener('change', (e) => this.handleFileChange(e));
        }

        // Excel上传处理
        const uploadForm = document.getElementById('uploadForm');
        if (uploadForm) {
            uploadForm.addEventListener('submit', (e) => this.handleUpload(e));
        }

        // 拖拽上传功能
        this.setupDragAndDrop();
    }

    handleFileChange(e) {
        const file = e.target.files[0];
        const label = document.getElementById('fileLabel');
        const importBtn = document.getElementById('importBtn');
        
        if (file) {
            label.textContent = `📄 ${file.name}`;
            label.classList.add('file-selected');
            importBtn.disabled = false;
        } else {
            label.textContent = '📁 选择Excel文件';
            label.classList.remove('file-selected');
            importBtn.disabled = true;
        }
    }

    async handleUpload(e) {
        e.preventDefault();
        
        if (this.isUploading) return;
        
        const fileInput = document.getElementById('excelFile');
        const file = fileInput.files[0];
        const importBtn = document.getElementById('importBtn');
        const progressBar = document.getElementById('progressBar');
        const progressFill = document.getElementById('progressFill');
        
        if (!file) {
            this.showMessage('请选择Excel文件', 'error');
            return;
        }
        
        // 验证文件类型
        if (!this.validateFile(file)) {
            return;
        }
        
        // 检查文件大小（最大10MB）
        if (file.size > 10 * 1024 * 1024) {
            this.showMessage('文件大小不能超过10MB', 'error');
            return;
        }
        
        const formData = new FormData();
        formData.append('file', file);
        
        // 开始上传
        this.isUploading = true;
        importBtn.disabled = true;
        importBtn.innerHTML = '⏳ 正在上传...';
        progressBar.style.display = 'block';
        progressFill.style.width = '0%';
        
        // 模拟进度条
        const progressInterval = this.simulateProgress(progressFill);
        
        try {
            const response = await fetch(`${this.API_BASE}/upload-excel`, {
                method: 'POST',
                body: formData
            });
            
            clearInterval(progressInterval);
            progressFill.style.width = '100%';
            
            const result = await response.json();
            
            if (response.ok) {
                this.showMessage(`导入成功！${result.message}`, 'success');
                if (result.warnings && result.warnings.length > 0) {
                    console.warn('导入警告:', result.warnings);
                    this.showMessage(`导入完成，但有 ${result.warnings.length} 个警告，请查看控制台`, 'warning');
                }
                
                // 重置表单
                this.resetForm();
                
                // 刷新数据
                setTimeout(() => {
                    if (window.customerManager) {
                        window.customerManager.loadCustomers();
                    }
                }, 1000);
                
            } else {
                this.showMessage(`上传失败：${result.error || '未知错误'}`, 'error');
            }
            
        } catch (error) {
            clearInterval(progressInterval);
            console.error('上传错误:', error);
            this.showMessage(`上传错误：${error.message || '网络连接失败'}`, 'error');
        } finally {
            this.resetUploadState();
        }
    }

    validateFile(file) {
        const validTypes = ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'application/vnd.ms-excel'];
        const validExtensions = ['.xlsx', '.xls'];
        const fileName = file.name.toLowerCase();
        const hasValidExtension = validExtensions.some(ext => fileName.endsWith(ext));
        
        if (!validTypes.includes(file.type) && !hasValidExtension) {
            this.showMessage('请选择有效的Excel文件（.xlsx或.xls格式）', 'error');
            return false;
        }
        return true;
    }

    simulateProgress(progressFill) {
        let progress = 0;
        return setInterval(() => {
            if (progress < 90) {
                progress += Math.random() * 10;
                progressFill.style.width = Math.min(progress, 90) + '%';
            }
        }, 200);
    }

    resetForm() {
        const fileInput = document.getElementById('excelFile');
        const fileLabel = document.getElementById('fileLabel');
        
        fileInput.value = '';
        fileLabel.textContent = '📁 选择Excel文件';
        fileLabel.classList.remove('file-selected');
    }

    resetUploadState() {
        const importBtn = document.getElementById('importBtn');
        const progressBar = document.getElementById('progressBar');
        const progressFill = document.getElementById('progressFill');
        const fileInput = document.getElementById('excelFile');
        
        this.isUploading = false;
        importBtn.disabled = fileInput.files.length === 0;
        importBtn.innerHTML = '📤 上传并导入';
        
        // 隐藏进度条
        setTimeout(() => {
            progressBar.style.display = 'none';
            progressFill.style.width = '0%';
        }, 2000);
    }

    setupDragAndDrop() {
        const importSection = document.querySelector('.import-section');
        if (!importSection) return;
        
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            importSection.addEventListener(eventName, (e) => {
                e.preventDefault();
                e.stopPropagation();
            }, false);
        });

        ['dragenter', 'dragover'].forEach(eventName => {
            importSection.addEventListener(eventName, () => {
                importSection.style.transform = 'scale(1.02)';
                importSection.style.boxShadow = '0 8px 16px rgba(102, 126, 234, 0.3)';
            }, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            importSection.addEventListener(eventName, () => {
                importSection.style.transform = 'scale(1)';
                importSection.style.boxShadow = '0 4px 12px rgba(102, 126, 234, 0.2)';
            }, false);
        });

        importSection.addEventListener('drop', (e) => {
            const dt = e.dataTransfer;
            const files = dt.files;
            
            if (files.length > 0) {
                const fileInput = document.getElementById('excelFile');
                fileInput.files = files;
                fileInput.dispatchEvent(new Event('change', { bubbles: true }));
            }
        }, false);
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

export default ExcelImport;