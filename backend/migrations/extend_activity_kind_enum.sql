-- 扩展 activity_kind 枚举类型
-- 添加缺失的跟进记录类型

-- 检查并添加新的枚举值到现有的 activity_kind 类型
-- 使用 DO 块来避免重复添加已存在的值

DO $$
BEGIN
    -- 添加 'order' 如果不存在
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'order' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind')) THEN
        ALTER TYPE activity_kind ADD VALUE 'order';
    END IF;
    
    -- 添加 'sample' 如果不存在
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'sample' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind')) THEN
        ALTER TYPE activity_kind ADD VALUE 'sample';
    END IF;
    
    -- 添加 'feedback' 如果不存在
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'feedback' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind')) THEN
        ALTER TYPE activity_kind ADD VALUE 'feedback';
    END IF;
    
    -- 添加 'complaint' 如果不存在
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'complaint' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind')) THEN
        ALTER TYPE activity_kind ADD VALUE 'complaint';
    END IF;
    
    -- 添加 'payment' 如果不存在
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'payment' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind')) THEN
        ALTER TYPE activity_kind ADD VALUE 'payment';
    END IF;
    
    -- 添加 'system' 如果不存在
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'system' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind')) THEN
        ALTER TYPE activity_kind ADD VALUE 'system';
    END IF;
END $$;

-- 验证枚举值已添加
SELECT enumlabel FROM pg_enum WHERE enumtypid = (SELECT oid FROM pg_type WHERE typname = 'activity_kind') ORDER BY enumsortorder;