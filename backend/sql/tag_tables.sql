-- 标签系统表结构设计
-- 删除已存在的表（如果存在）
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS tag_dimensions CASCADE;

-- 标签维度表
CREATE TABLE tag_dimensions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    description TEXT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    is_deleted BOOLEAN DEFAULT FALSE
);

-- 添加表注释
COMMENT ON TABLE tag_dimensions IS '标签维度表';
COMMENT ON COLUMN tag_dimensions.id IS '维度ID';
COMMENT ON COLUMN tag_dimensions.name IS '维度名称';
COMMENT ON COLUMN tag_dimensions.description IS '维度描述';
COMMENT ON COLUMN tag_dimensions.sort_order IS '排序顺序';
COMMENT ON COLUMN tag_dimensions.created_at IS '创建时间';
COMMENT ON COLUMN tag_dimensions.updated_at IS '更新时间';
COMMENT ON COLUMN tag_dimensions.deleted_at IS '删除时间（软删除）';
COMMENT ON COLUMN tag_dimensions.is_deleted IS '是否删除：false-否，true-是';

-- 标签表
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    dimension_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    color VARCHAR(32) NULL,
    description TEXT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (dimension_id) REFERENCES tag_dimensions(id)
);

-- 添加表注释
COMMENT ON TABLE tags IS '标签表';
COMMENT ON COLUMN tags.id IS '标签ID';
COMMENT ON COLUMN tags.dimension_id IS '维度ID';
COMMENT ON COLUMN tags.name IS '标签名称';
COMMENT ON COLUMN tags.color IS '标签颜色（十六进制）';
COMMENT ON COLUMN tags.description IS '标签描述';
COMMENT ON COLUMN tags.sort_order IS '排序顺序';
COMMENT ON COLUMN tags.created_at IS '创建时间';
COMMENT ON COLUMN tags.updated_at IS '更新时间';
COMMENT ON COLUMN tags.deleted_at IS '删除时间（软删除）';
COMMENT ON COLUMN tags.is_deleted IS '是否删除：false-否，true-是';

-- 创建索引
CREATE INDEX idx_tag_dimensions_name ON tag_dimensions(name);
CREATE INDEX idx_tag_dimensions_sort_order ON tag_dimensions(sort_order);
CREATE INDEX idx_tag_dimensions_deleted ON tag_dimensions(is_deleted);

CREATE INDEX idx_tags_dimension_id ON tags(dimension_id);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_sort_order ON tags(sort_order);
CREATE INDEX idx_tags_deleted ON tags(is_deleted);
CREATE INDEX idx_tags_dimension_name ON tags(dimension_id, name) WHERE is_deleted = FALSE;

-- 插入维度数据
INSERT INTO tag_dimensions (name, description, sort_order) VALUES
('基本', '客户基本信息维度', 1),
('社会', '客户社会属性维度', 2),
('个性', '客户个性特征维度', 3),
('行为', '客户行为特征维度', 4),
('预测', '客户预测分析维度', 5);

-- 插入标签数据
-- 基本维度标签
INSERT INTO tags (dimension_id, name, color, description, sort_order) VALUES
((SELECT id FROM tag_dimensions WHERE name = '基本'), '汉族', '#2196F3', '汉族客户', 1),
((SELECT id FROM tag_dimensions WHERE name = '基本'), '少数民族', '#4CAF50', '少数民族客户', 2),
((SELECT id FROM tag_dimensions WHERE name = '基本'), '回族', '#FF9800', '回族客户', 3),
((SELECT id FROM tag_dimensions WHERE name = '基本'), '其他', '#9E9E9E', '其他民族客户', 4);

-- 社会维度标签
INSERT INTO tags (dimension_id, name, color, description, sort_order) VALUES
((SELECT id FROM tag_dimensions WHERE name = '社会'), '已婚', '#E91E63', '已婚客户', 1),
((SELECT id FROM tag_dimensions WHERE name = '社会'), '高收入', '#FFC107', '高收入客户', 2),
((SELECT id FROM tag_dimensions WHERE name = '社会'), '博士', '#9C27B0', '博士学历客户', 3),
((SELECT id FROM tag_dimensions WHERE name = '社会'), '茶行业', '#4CAF50', '茶行业相关客户', 4),
((SELECT id FROM tag_dimensions WHERE name = '社会'), '销售', '#FF5722', '销售行业客户', 5),
((SELECT id FROM tag_dimensions WHERE name = '社会'), '老板', '#795548', '企业老板客户', 6),
((SELECT id FROM tag_dimensions WHERE name = '社会'), '有子女', '#607D8B', '有子女的客户', 7);

-- 个性维度标签
INSERT INTO tags (dimension_id, name, color, description, sort_order) VALUES
((SELECT id FROM tag_dimensions WHERE name = '个性'), '爱足疗', '#E91E63', '喜欢足疗的客户', 1),
((SELECT id FROM tag_dimensions WHERE name = '个性'), '不喝酒', '#2196F3', '不饮酒的客户', 2),
((SELECT id FROM tag_dimensions WHERE name = '个性'), '不抽烟', '#4CAF50', '不吸烟的客户', 3),
((SELECT id FROM tag_dimensions WHERE name = '个性'), '四川辣', '#FF5722', '喜欢四川辣味的客户', 4),
((SELECT id FROM tag_dimensions WHERE name = '个性'), '豪爽', '#FF9800', '性格豪爽的客户', 5);

-- 行为维度标签
INSERT INTO tags (dimension_id, name, color, description, sort_order) VALUES
((SELECT id FROM tag_dimensions WHERE name = '行为'), '多次下单', '#4CAF50', '多次下单的客户', 1),
((SELECT id FROM tag_dimensions WHERE name = '行为'), '新客户', '#2196F3', '新注册的客户', 2),
((SELECT id FROM tag_dimensions WHERE name = '行为'), '从未下单', '#9E9E9E', '从未下单的客户', 3);

-- 预测维度标签
INSERT INTO tags (dimension_id, name, color, description, sort_order) VALUES
((SELECT id FROM tag_dimensions WHERE name = '预测'), '疑似流失', '#F44336', '有流失风险的客户', 1);

-- 常用查询示例
-- 1. 查询所有维度及其标签
/*
SELECT 
    td.name as dimension_name,
    td.description as dimension_desc,
    t.name as tag_name,
    t.color as tag_color,
    t.description as tag_desc
FROM tag_dimensions td
LEFT JOIN tags t ON td.id = t.dimension_id
WHERE td.is_deleted = FALSE 
  AND (t.is_deleted = FALSE OR t.id IS NULL)
ORDER BY td.sort_order, t.sort_order;
*/

-- 2. 按维度分组查询标签
/*
SELECT 
    td.name as dimension_name,
    COUNT(t.id) as tag_count,
    array_agg(t.name ORDER BY t.sort_order) as tag_names
FROM tag_dimensions td
LEFT JOIN tags t ON td.id = t.dimension_id AND t.is_deleted = FALSE
WHERE td.is_deleted = FALSE
GROUP BY td.id, td.name, td.sort_order
ORDER BY td.sort_order;
*/

-- 3. 查询特定维度的所有标签
/*
SELECT t.* 
FROM tags t
JOIN tag_dimensions td ON t.dimension_id = td.id
WHERE td.name = '基本' 
  AND t.is_deleted = FALSE 
  AND td.is_deleted = FALSE
ORDER BY t.sort_order;
*/