-- 客户表系统标签字段升级脚本
-- 添加系统标签字段用于保存客户的标签信息

-- 为客户表添加系统标签字段（使用int4[]类型更适合存储标签ID数组）
ALTER TABLE customers 
ADD COLUMN IF NOT EXISTS system_tags int4[] DEFAULT '{}';

-- 创建系统标签字段的索引，支持数组查询
CREATE INDEX IF NOT EXISTS idx_customers_system_tags ON customers USING GIN (system_tags);

-- 添加备注说明
COMMENT ON COLUMN customers.system_tags IS '系统标签字段，存储整型数组格式的标签ID数组，例如: {1,2,3}';