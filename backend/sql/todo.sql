-- 待办表结构设计
-- 删除已存在的枚举类型（如果存在）
DROP TYPE IF EXISTS todo_status CASCADE;
DROP TYPE IF EXISTS reminder_type CASCADE;
DROP TYPE IF EXISTS todo_priority CASCADE;
DROP TYPE IF EXISTS todo_action CASCADE;

-- 创建枚举类型
CREATE TYPE todo_status AS ENUM ('pending', 'completed', 'overdue', 'cancelled');
CREATE TYPE reminder_type AS ENUM ('wechat', 'enterprise_wechat', 'both');
CREATE TYPE todo_priority AS ENUM ('low', 'medium', 'high', 'urgent');
CREATE TYPE creation_type AS ENUM ('manual', 'auto');
CREATE TYPE alert_status AS ENUM ('none', 'warning', 'overdue');

-- 删除已存在的表（如果存在）
DROP TABLE IF EXISTS todo_logs CASCADE;
DROP TABLE IF EXISTS todos CASCADE;

-- 待办表
CREATE TABLE todos (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    creator_id BIGINT NOT NULL,
    executor_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    status todo_status DEFAULT 'pending',
    planned_time TIMESTAMP NOT NULL,
    completed_time TIMESTAMP NULL,
    creation_type creation_type DEFAULT 'manual',
    warning_time TIMESTAMP NULL,
    alert_status alert_status DEFAULT 'none',
    last_alert_time TIMESTAMP NULL,
    is_reminder BOOLEAN DEFAULT FALSE,
    reminder_type reminder_type NULL,
    reminder_user_id BIGINT NULL,
    reminder_time TIMESTAMP NULL,
    priority todo_priority DEFAULT 'medium',
    tags JSONB NULL,
    attachments JSONB NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    is_deleted BOOLEAN DEFAULT FALSE
);

-- 添加表注释
COMMENT ON TABLE todos IS '待办事项表';
COMMENT ON COLUMN todos.id IS '待办ID';
COMMENT ON COLUMN todos.customer_id IS '关联客户ID';
COMMENT ON COLUMN todos.creator_id IS '创建人ID';
COMMENT ON COLUMN todos.executor_id IS '执行人ID';
COMMENT ON COLUMN todos.title IS '待办标题';
COMMENT ON COLUMN todos.content IS '待办内容详情';
COMMENT ON COLUMN todos.status IS '待办状态：未完成、已完成、延期、取消';
COMMENT ON COLUMN todos.planned_time IS '计划执行时间';
COMMENT ON COLUMN todos.completed_time IS '完成时间';
COMMENT ON COLUMN todos.creation_type IS '创建方式：manual-手动创建，auto-自动创建';
COMMENT ON COLUMN todos.warning_time IS '预警时间（一般为计划时间前12小时）';
COMMENT ON COLUMN todos.alert_status IS '告警状态：none-无告警，warning-预警，overdue-逾期告警';
COMMENT ON COLUMN todos.last_alert_time IS '最后告警时间';
COMMENT ON COLUMN todos.is_reminder IS '是否提醒：false-否，true-是';
COMMENT ON COLUMN todos.reminder_type IS '提醒方式：微信、企业微信、两者';
COMMENT ON COLUMN todos.reminder_user_id IS '提醒人ID';
COMMENT ON COLUMN todos.reminder_time IS '提醒时间';
COMMENT ON COLUMN todos.priority IS '优先级';
COMMENT ON COLUMN todos.tags IS '标签（JSONB格式）';
COMMENT ON COLUMN todos.attachments IS '附件信息（JSONB格式）';
COMMENT ON COLUMN todos.created_at IS '创建时间';
COMMENT ON COLUMN todos.updated_at IS '更新时间';
COMMENT ON COLUMN todos.deleted_at IS '删除时间（软删除）';
COMMENT ON COLUMN todos.is_deleted IS '是否删除：false-否，true-是';

-- 创建索引
CREATE INDEX idx_todos_customer_id ON todos(customer_id);
CREATE INDEX idx_todos_creator_id ON todos(creator_id);
CREATE INDEX idx_todos_executor_id ON todos(executor_id);
CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_planned_time ON todos(planned_time);
CREATE INDEX idx_todos_created_at ON todos(created_at);
CREATE INDEX idx_todos_deleted ON todos(is_deleted);
CREATE INDEX idx_todos_creation_type ON todos(creation_type);
CREATE INDEX idx_todos_warning_time ON todos(warning_time);
CREATE INDEX idx_todos_alert_status ON todos(alert_status);
CREATE INDEX idx_todos_last_alert_time ON todos(last_alert_time);
-- 复合索引用于监控查询
CREATE INDEX idx_todos_monitor ON todos(status, alert_status, warning_time, planned_time) WHERE is_deleted = FALSE;

-- 外键约束已移除

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建自动设置预警时间的触发器函数
CREATE OR REPLACE FUNCTION set_warning_time()
RETURNS TRIGGER AS $$
BEGIN
    -- 如果预警时间为空，自动设置为计划时间前12小时
    IF NEW.warning_time IS NULL AND NEW.planned_time IS NOT NULL THEN
        NEW.warning_time = NEW.planned_time - INTERVAL '12 hours';
    END IF;
    
    -- 如果状态变为已完成或取消，重置告警状态
    IF NEW.status IN ('completed', 'cancelled') THEN
        NEW.alert_status = 'none';
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 删除已存在的触发器（如果存在）
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
DROP TRIGGER IF EXISTS set_todos_warning_time ON todos;

-- 创建触发器
CREATE TRIGGER update_todos_updated_at
    BEFORE UPDATE ON todos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_todos_warning_time
    BEFORE INSERT OR UPDATE ON todos
    FOR EACH ROW
    EXECUTE FUNCTION set_warning_time();

-- 待办操作日志表（记录待办的变更历史）
-- 创建操作类型枚举
CREATE TYPE todo_action AS ENUM ('create', 'update', 'delete', 'complete', 'cancel');

-- 创建待办操作日志表
CREATE TABLE todo_logs (
    id BIGSERIAL PRIMARY KEY,
    todo_id BIGINT NOT NULL,
    operator_id BIGINT NOT NULL,
    action todo_action NOT NULL,
    old_data JSONB NULL,
    new_data JSONB NULL,
    remark VARCHAR(500) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 添加表注释
COMMENT ON TABLE todo_logs IS '待办操作日志表';
COMMENT ON COLUMN todo_logs.id IS '日志ID';
COMMENT ON COLUMN todo_logs.todo_id IS '待办ID';
COMMENT ON COLUMN todo_logs.operator_id IS '操作人ID';
COMMENT ON COLUMN todo_logs.action IS '操作类型';
COMMENT ON COLUMN todo_logs.old_data IS '变更前数据';
COMMENT ON COLUMN todo_logs.new_data IS '变更后数据';
COMMENT ON COLUMN todo_logs.remark IS '操作备注';
COMMENT ON COLUMN todo_logs.created_at IS '操作时间';

-- 创建索引
CREATE INDEX idx_todo_logs_todo_id ON todo_logs(todo_id);
CREATE INDEX idx_todo_logs_operator_id ON todo_logs(operator_id);
CREATE INDEX idx_todo_logs_action ON todo_logs(action);
CREATE INDEX idx_todo_logs_created_at ON todo_logs(created_at);

-- 外键约束已移除

-- 插入示例数据
INSERT INTO todos (customer_id, creator_id, executor_id, title, content, status, planned_time, creation_type, alert_status, is_reminder, reminder_type, priority) VALUES
(1, 1, 1, '跟进客户需求', '联系客户了解最新采购需求，确认订单意向', 'pending', '2024-01-15 10:00:00', 'manual', 'none', TRUE, 'enterprise_wechat', 'high'),
(1, 1, 1, '发送产品样品', '向客户发送新款茶叶样品，收集反馈意见', 'completed', '2024-01-10 14:00:00', 'manual', 'none', FALSE, NULL, 'medium'),
(2, 1, 2, '客户回访', '电话回访客户使用情况，维护客户关系', 'overdue', '2024-01-12 09:00:00', 'auto', 'overdue', TRUE, 'wechat', 'medium'),
(3, 1, 1, '系统自动提醒', '系统根据客户行为自动生成的跟进提醒', 'pending', CURRENT_TIMESTAMP + INTERVAL '2 hours', 'auto', 'warning', TRUE, 'enterprise_wechat', 'medium');

-- 插入示例日志数据
INSERT INTO todo_logs (todo_id, operator_id, action, new_data, remark) VALUES
(1, 1, 'create', '{"title": "跟进客户需求", "status": "pending", "creation_type": "manual"}', '创建待办事项'),
(2, 1, 'complete', '{"status": "completed", "completed_time": "2024-01-10 15:30:00"}', '完成待办事项'),
(3, 1, 'create', '{"title": "客户回访", "status": "pending", "creation_type": "auto"}', '自动创建待办事项'),
(4, 1, 'create', '{"title": "系统自动提醒", "status": "pending", "creation_type": "auto"}', '系统自动创建待办事项');

-- 监控查询示例
-- 1. 查询需要预警的待办（当前时间已超过预警时间且状态为pending）
/*
SELECT id, title, executor_id, planned_time, warning_time, alert_status
FROM todos 
WHERE status = 'pending' 
  AND warning_time <= CURRENT_TIMESTAMP 
  AND alert_status = 'none'
  AND is_deleted = FALSE;
*/

-- 2. 查询已逾期的待办（当前时间已超过计划时间且状态为pending）
/*
SELECT id, title, executor_id, planned_time, alert_status
FROM todos 
WHERE status = 'pending' 
  AND planned_time <= CURRENT_TIMESTAMP 
  AND is_deleted = FALSE;
*/

-- 3. 更新告警状态的SQL（用于定时任务）
/*
-- 更新预警状态
UPDATE todos 
SET alert_status = 'warning', last_alert_time = CURRENT_TIMESTAMP
WHERE status = 'pending' 
  AND warning_time <= CURRENT_TIMESTAMP 
  AND alert_status = 'none'
  AND is_deleted = FALSE;

-- 更新逾期状态
UPDATE todos 
SET status = 'overdue', alert_status = 'overdue', last_alert_time = CURRENT_TIMESTAMP
WHERE status = 'pending' 
  AND planned_time <= CURRENT_TIMESTAMP 
  AND is_deleted = FALSE;
*/