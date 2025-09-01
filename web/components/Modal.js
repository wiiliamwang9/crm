// 模态框组件
class Modal {
    constructor(options = {}) {
        this.options = {
            title: '提示',
            content: '',
            width: '500px',
            height: 'auto',
            closable: true,
            maskClosable: true,
            className: '',
            onOpen: null,
            onClose: null,
            ...options
        };
        
        this.isVisible = false;
        this.overlay = null;
        this.modal = null;
    }
    
    create() {
        // 创建遮罩层
        this.overlay = document.createElement('div');
        this.overlay.className = `modal-overlay ${this.options.className}`;
        this.overlay.style.cssText = `
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.5);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 10000;
            opacity: 0;
            transition: opacity 0.3s ease;
        `;
        
        // 创建模态框
        this.modal = document.createElement('div');
        this.modal.className = 'modal-dialog';
        this.modal.style.cssText = `
            background: white;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
            max-height: 90vh;
            overflow: hidden;
            width: ${this.options.width};
            height: ${this.options.height};
            transform: scale(0.7);
            transition: transform 0.3s ease;
            display: flex;
            flex-direction: column;
        `;
        
        // 创建头部
        if (this.options.title) {
            const header = document.createElement('div');
            header.className = 'modal-header';
            header.style.cssText = `
                padding: 16px 20px;
                border-bottom: 1px solid #f0f0f0;
                display: flex;
                justify-content: space-between;
                align-items: center;
                font-size: 16px;
                font-weight: 500;
            `;
            
            const title = document.createElement('div');
            title.textContent = this.options.title;
            header.appendChild(title);
            
            if (this.options.closable) {
                const closeBtn = document.createElement('button');
                closeBtn.innerHTML = '×';
                closeBtn.style.cssText = `
                    background: none;
                    border: none;
                    font-size: 24px;
                    color: #999;
                    cursor: pointer;
                    padding: 0;
                    width: 24px;
                    height: 24px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                `;
                closeBtn.addEventListener('click', () => this.hide());
                header.appendChild(closeBtn);
            }
            
            this.modal.appendChild(header);
        }
        
        // 创建内容区
        const body = document.createElement('div');
        body.className = 'modal-body';
        body.style.cssText = `
            padding: 20px;
            flex: 1;
            overflow-y: auto;
        `;
        
        if (typeof this.options.content === 'string') {
            body.innerHTML = this.options.content;
        } else if (this.options.content instanceof HTMLElement) {
            body.appendChild(this.options.content);
        }
        
        this.modal.appendChild(body);
        this.overlay.appendChild(this.modal);
        
        // 绑定事件
        if (this.options.maskClosable) {
            this.overlay.addEventListener('click', (e) => {
                if (e.target === this.overlay) {
                    this.hide();
                }
            });
        }
        
        // ESC键关闭
        this.handleKeydown = (e) => {
            if (e.key === 'Escape' && this.isVisible) {
                this.hide();
            }
        };
        
        return this.modal;
    }
    
    show() {
        if (this.isVisible) return;
        
        if (!this.overlay) {
            this.create();
        }
        
        document.body.appendChild(this.overlay);
        document.addEventListener('keydown', this.handleKeydown);
        
        // 禁止body滚动
        document.body.style.overflow = 'hidden';
        
        // 显示动画
        requestAnimationFrame(() => {
            this.overlay.style.opacity = '1';
            this.modal.style.transform = 'scale(1)';
        });
        
        this.isVisible = true;
        
        if (this.options.onOpen) {
            this.options.onOpen(this);
        }
        
        return this;
    }
    
    hide() {
        if (!this.isVisible) return;
        
        // 隐藏动画
        this.overlay.style.opacity = '0';
        this.modal.style.transform = 'scale(0.7)';
        
        setTimeout(() => {
            if (this.overlay && this.overlay.parentNode) {
                this.overlay.parentNode.removeChild(this.overlay);
            }
            
            // 恢复body滚动
            document.body.style.overflow = '';
            
            document.removeEventListener('keydown', this.handleKeydown);
            
            this.isVisible = false;
            
            if (this.options.onClose) {
                this.options.onClose(this);
            }
        }, 300);
        
        return this;
    }
    
    setContent(content) {
        if (!this.modal) return;
        
        const body = this.modal.querySelector('.modal-body');
        if (body) {
            if (typeof content === 'string') {
                body.innerHTML = content;
            } else if (content instanceof HTMLElement) {
                body.innerHTML = '';
                body.appendChild(content);
            }
        }
        
        return this;
    }
    
    setTitle(title) {
        if (!this.modal) return;
        
        const header = this.modal.querySelector('.modal-header');
        if (header) {
            const titleEl = header.querySelector('div');
            if (titleEl) {
                titleEl.textContent = title;
            }
        }
        
        return this;
    }
    
    destroy() {
        this.hide();
        
        setTimeout(() => {
            if (this.overlay && this.overlay.parentNode) {
                this.overlay.parentNode.removeChild(this.overlay);
            }
            this.overlay = null;
            this.modal = null;
        }, 300);
    }
}

// 静态方法
Modal.alert = function(message, title = '提示', onClose) {
    const modal = new Modal({
        title,
        content: `
            <div style="padding: 20px 0; text-align: center; font-size: 16px;">
                ${message}
            </div>
            <div style="text-align: center; padding-top: 20px;">
                <button class="btn-primary" style="padding: 8px 24px; background: #1890ff; color: white; border: none; border-radius: 4px; cursor: pointer;">确定</button>
            </div>
        `,
        closable: false,
        maskClosable: false,
        onOpen: function(modalInstance) {
            const btn = modalInstance.modal.querySelector('.btn-primary');
            btn.addEventListener('click', () => {
                modalInstance.hide();
                if (onClose) onClose();
            });
        }
    });
    
    modal.show();
    return modal;
};

Modal.confirm = function(message, title = '确认', onConfirm, onCancel) {
    const modal = new Modal({
        title,
        content: `
            <div style="padding: 20px 0; font-size: 16px;">
                ${message}
            </div>
            <div style="text-align: right; padding-top: 20px;">
                <button class="btn-cancel" style="padding: 8px 16px; border: 1px solid #d9d9d9; background: white; color: #333; border-radius: 4px; cursor: pointer; margin-right: 12px;">取消</button>
                <button class="btn-confirm" style="padding: 8px 16px; border: none; background: #1890ff; color: white; border-radius: 4px; cursor: pointer;">确定</button>
            </div>
        `,
        closable: false,
        maskClosable: false,
        onOpen: function(modalInstance) {
            const btnCancel = modalInstance.modal.querySelector('.btn-cancel');
            const btnConfirm = modalInstance.modal.querySelector('.btn-confirm');
            
            btnCancel.addEventListener('click', () => {
                modalInstance.hide();
                if (onCancel) onCancel();
            });
            
            btnConfirm.addEventListener('click', () => {
                modalInstance.hide();
                if (onConfirm) onConfirm();
            });
        }
    });
    
    modal.show();
    return modal;
};

export default Modal;