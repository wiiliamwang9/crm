-- 提醒记录表
CREATE TABLE reminders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '提醒ID',
    todo_id BIGINT NOT NULL COMMENT '关联待办ID',
    user_id BIGINT NOT NULL COMMENT '提醒用户ID',
    type ENUM('wechat', 'enterprise_wechat', 'both') NOT NULL COMMENT '提醒方式',
    title VARCHAR(255) NOT NULL COMMENT '提醒标题',
    content TEXT COMMENT '提醒内容',
    status ENUM('pending', 'sent', 'failed', 'cancelled') DEFAULT 'pending' COMMENT '提醒状态',
    frequency ENUM('once', 'daily', 'weekly', 'monthly') DEFAULT 'once' COMMENT '提醒频率',
    schedule_time DATETIME NOT NULL COMMENT '计划提醒时间',
    sent_time DATETIME NULL COMMENT '实际发送时间',
    fail_reason VARCHAR(500) NULL COMMENT '失败原因',
    retry_count INT DEFAULT 0 COMMENT '重试次数',
    max_retries INT DEFAULT 3 COMMENT '最大重试次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_todo_id (todo_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_schedule_time (schedule_time),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提醒记录表';

-- 提醒模板表
CREATE TABLE reminder_templates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '模板ID',
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    type ENUM('wechat', 'enterprise_wechat', 'both') NOT NULL COMMENT '适用的提醒方式',
    title VARCHAR(255) NOT NULL COMMENT '标题模板',
    content TEXT NOT NULL COMMENT '内容模板',
    variables JSON NULL COMMENT '可用变量说明',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    is_default TINYINT(1) DEFAULT 0 COMMENT '是否为默认模板',
    created_by BIGINT NOT NULL COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_type (type),
    INDEX idx_active (is_active),
    INDEX idx_default (is_default),
    
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提醒模板表';

-- 用户提醒配置表
CREATE TABLE reminder_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    user_id BIGINT NOT NULL UNIQUE COMMENT '用户ID',
    enable_wechat TINYINT(1) DEFAULT 1 COMMENT '启用微信提醒',
    enable_enterprise_wechat TINYINT(1) DEFAULT 1 COMMENT '启用企业微信提醒',
    wechat_user_id VARCHAR(100) NULL COMMENT '微信用户ID',
    enterprise_wechat_user_id VARCHAR(100) NULL COMMENT '企业微信用户ID',
    default_advance_minutes INT DEFAULT 30 COMMENT '默认提前提醒分钟数',
    quiet_start_time VARCHAR(5) DEFAULT '22:00' COMMENT '免打扰开始时间',
    quiet_end_time VARCHAR(5) DEFAULT '08:00' COMMENT '免打扰结束时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户提醒配置表';

-- 插入默认提醒模板
INSERT INTO reminder_templates (name, type, title, content, variables, is_default, created_by) VALUES
('默认待办提醒', 'both', '待办提醒：{{.Title}}', '您有一个待办事项需要处理：\n\n标题：{{.Title}}\n内容：{{.Content}}\n客户：{{.CustomerName}}\n计划时间：{{.PlannedTime}}\n优先级：{{.Priority}}\n\n请及时处理！', '{"Title":"待办标题","Content":"待办内容","CustomerName":"客户名称","PlannedTime":"计划时间","Priority":"优先级"}', 1, 1),
('延期提醒', 'both', '延期提醒：{{.Title}}', '您的待办事项已延期：\n\n标题：{{.Title}}\n客户：{{.CustomerName}}\n原计划时间：{{.PlannedTime}}\n已延期：{{.OverdueDays}}天\n\n请尽快处理！', '{"Title":"待办标题","CustomerName":"客户名称","PlannedTime":"计划时间","OverdueDays":"延期天数"}', 0, 1),
('完成提醒', 'both', '任务完成：{{.Title}}', '恭喜！您已完成待办事项：\n\n标题：{{.Title}}\n客户：{{.CustomerName}}\n完成时间：{{.CompletedTime}}', '{"Title":"待办标题","CustomerName":"客户名称","CompletedTime":"完成时间"}', 0, 1);

-- 创建默认用户配置（假设用户ID为1）
INSERT INTO reminder_configs (user_id, enable_wechat, enable_enterprise_wechat, default_advance_minutes, quiet_start_time, quiet_end_time) VALUES
(1, 1, 1, 30, '22:00', '08:00');

-- 创建示例提醒记录
INSERT INTO reminders (todo_id, user_id, type, title, content, status, schedule_time) VALUES
(1, 1, 'enterprise_wechat', '待办提醒：跟进客户需求', '您有一个待办事项需要处理：\n\n标题：跟进客户需求\n客户：示例客户\n计划时间：2024-01-15 10:00:00\n优先级：高\n\n请及时处理！', 'pending', '2024-01-15 09:30:00'),
(3, 1, 'wechat', '延期提醒：客户回访', '您的待办事项已延期：\n\n标题：客户回访\n客户：测试客户\n原计划时间：2024-01-12 09:00:00\n已延期：3天\n\n请尽快处理！', 'sent', '2024-01-15 09:00:00');