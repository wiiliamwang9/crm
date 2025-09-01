-- 添加 sms 到 reminder_type 枚举
-- 注意：PostgreSQL 不支持直接修改枚举，需要先创建新类型，然后替换

-- 创建新的枚举类型
CREATE TYPE reminder_type_new AS ENUM ('wechat', 'enterprise_wechat', 'both', 'sms');

-- 更新 todos 表
ALTER TABLE todos 
ALTER COLUMN reminder_type TYPE reminder_type_new 
USING reminder_type::text::reminder_type_new;

-- 更新 reminders 表（如果存在）
-- ALTER TABLE reminders 
-- ALTER COLUMN type TYPE reminder_type_new 
-- USING type::text::reminder_type_new;

-- 更新 reminder_templates 表（如果存在）
-- ALTER TABLE reminder_templates 
-- ALTER COLUMN type TYPE reminder_type_new 
-- USING type::text::reminder_type_new;

-- 删除旧的枚举类型
DROP TYPE IF EXISTS reminder_type;

-- 重命名新类型
ALTER TYPE reminder_type_new RENAME TO reminder_type;