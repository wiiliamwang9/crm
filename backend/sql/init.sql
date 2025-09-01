CREATE TABLE "public"."activities"(
                                      "id" int8 NOT NULL DEFAULT nextval('activities_id_seq'::regclass),
                                      "customer_id" int8 NOT NULL,
                                      "user_id" int8 NOT NULL,
                                      "kind" "public"."activity_kind" DEFAULT 'other'::activity_kind,
                                      "title" varchar(255),
                                      "data" jsonb,
                                      "remark" text,
                                      "duration" int4,
                                      "location" varchar(255),
                                      "next_follow_time" timestamp(6),
                                      "attachments" jsonb,
                                      "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                      "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                      "deleted_at" timestamp(6),
                                      "is_deleted" bool DEFAULT false,
                                      CONSTRAINT "activities_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."activities" OWNER TO "postgres";
CREATE INDEX "idx_activities_customer_id" ON "public"."activities" USING btree (customer_id);
CREATE INDEX "idx_activities_user_id" ON "public"."activities" USING btree (user_id);
CREATE INDEX "idx_activities_kind" ON "public"."activities" USING btree (kind);
CREATE INDEX "idx_activities_created_at" ON "public"."activities" USING btree (created_at);
CREATE INDEX "idx_activities_next_follow_time" ON "public"."activities" USING btree (next_follow_time);
CREATE INDEX "idx_activities_deleted" ON "public"."activities" USING btree (is_deleted);
CREATE INDEX "idx_activities_customer_user" ON "public"."activities" USING btree (customer_id, user_id) WHERE (is_deleted = false);
CREATE INDEX "idx_activities_user_time" ON "public"."activities" USING btree (user_id, created_at) WHERE (is_deleted = false);
COMMENT ON COLUMN "public"."activities"."id" IS '记录ID';
COMMENT ON COLUMN "public"."activities"."customer_id" IS '客户ID';
COMMENT ON COLUMN "public"."activities"."user_id" IS '用户ID';
COMMENT ON COLUMN "public"."activities"."kind" IS '跟进类型：电话、拜访、邮件、微信、会议、其他';
COMMENT ON COLUMN "public"."activities"."title" IS '跟进标题';
COMMENT ON COLUMN "public"."activities"."data" IS '沟通数据，格式：{"content": "沟通内容", "result": "沟通结果"}';
COMMENT ON COLUMN "public"."activities"."remark" IS '备注';
COMMENT ON COLUMN "public"."activities"."duration" IS '持续时长（分钟）';
COMMENT ON COLUMN "public"."activities"."location" IS '地点';
COMMENT ON COLUMN "public"."activities"."next_follow_time" IS '下次跟进时间';
COMMENT ON COLUMN "public"."activities"."attachments" IS '附件信息（JSONB格式）';
COMMENT ON COLUMN "public"."activities"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."activities"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."activities"."deleted_at" IS '删除时间（软删除）';
COMMENT ON COLUMN "public"."activities"."is_deleted" IS '是否删除：false-否，true-是';
COMMENT ON TABLE "public"."activities" IS '跟进记录表';

CREATE TABLE "public"."customers"(
                                     "id" int4 NOT NULL DEFAULT nextval('customers_id_seq'::regclass),
                                     "name" varchar(256),
                                     "contact_name" varchar(256),
                                     "gender" int8,
                                     "avatar" varchar(2048),
                                     "photos" varchar[],
                                     "remark" text,
                                     "source" varchar(256),
                                     "created_at" timestamptz(6) DEFAULT now(),
                                     "created_by" int8,
                                     "updated_at" timestamptz(6) DEFAULT now(),
                                     "updated_by" int8,
                                     "phones" varchar[] DEFAULT NULL::character varying[],
                                     "wechats" varchar[],
                                     "douyins" varchar[],
                                     "kwais" varchar[],
                                     "redbooks" varchar[],
                                     "wework_openids" varchar[],
                                     "work_phone" varchar[],
                                     "work_wechat" varchar[],
                                     "province" varchar(256),
                                     "city" varchar(256),
                                     "district" varchar(256),
                                     "district_id" int8,
                                     "street" varchar(256),
                                     "address" varchar(2048),
                                     "lat" numeric,
                                     "lon" numeric,
                                     "category" varchar(256),
                                     "tags" varchar[],
                                     "level" int8,
                                     "state" int8,
                                     "kind" int8,
                                     "added_wechat" bool,
                                     "credit_sale" numeric DEFAULT 0,
                                     "sellers" int4[],
                                     "last_visited" timestamptz(6),
                                     "last_called" timestamptz(6),
                                     "last_order_date" timestamptz(6),
                                     "order_count" int4 DEFAULT 0,
                                     "avg_order_value" numeric,
                                     "group_id" int4[],
                                     "birth_place" varchar(256),
                                     "birth_year" int8,
                                     "birth_month" int8,
                                     "birth_date" int8,
                                     "favors" jsonb[],
                                     "products" varchar[],
                                     "annual_turnover" varchar(512),
                                     "shipping_infos" jsonb[],
                                     "preferred_delivery_method" varchar(128),
                                     "extra_info" jsonb,
                                     "flags" int8,
                                     "original_customer_id" varchar(256),
                                     "import_source" varchar(256),
                                     "saller_name" varchar(256),
                                     "system_tags" int4[] DEFAULT '{}'::integer[],
                                     CONSTRAINT "customers_gender_check" CHECK (((gender >= 0) AND (gender <= 2))),
                                     CONSTRAINT "customers_kind_check" CHECK (((kind >= 0) AND (kind <= 4))),
                                     CONSTRAINT "customers_level_check" CHECK ((level = ANY (ARRAY[0, 1, 2, 3, 4, 10]))),
                                     CONSTRAINT "customers_pkey" PRIMARY KEY (id),
                                     CONSTRAINT "customers_state_check" CHECK (((state >= 0) AND (state <= 8)))
);
ALTER TABLE "public"."customers" OWNER TO "postgres";
CREATE INDEX "idx_status" ON "public"."customers" USING btree (state);
CREATE INDEX "idx_phone" ON "public"."customers" USING btree (phones);
CREATE INDEX "idx_wechats" ON "public"."customers" USING btree (wechats);
CREATE INDEX "idx_customers_system_tags" ON "public"."customers" USING gin (system_tags);
CREATE INDEX "idx_customers_last_order_date" ON "public"."customers" USING btree (last_order_date);
CREATE TRIGGER "trg_customers_updated"
    BEFORE UPDATE ON "public"."customers"
    FOR EACH ROW
    EXECUTE FUNCTION "public"."set_customers_updated_at"() ;

COMMENT ON COLUMN "public"."customers"."id" IS '主键，自增';
COMMENT ON COLUMN "public"."customers"."name" IS '客户名称（门店招牌或常用昵称）';
COMMENT ON COLUMN "public"."customers"."contact_name" IS '联系人姓名（如“李雨亮 李老板”）';
COMMENT ON COLUMN "public"."customers"."gender" IS '性别：0=未知 1=男 2=女';
COMMENT ON COLUMN "public"."customers"."avatar" IS '头像 URL';
COMMENT ON COLUMN "public"."customers"."photos" IS '门头或店内照片 URL 数组';
COMMENT ON COLUMN "public"."customers"."remark" IS '人工备注（富文本，可存特殊要求、忌讳等）';
COMMENT ON COLUMN "public"."customers"."source" IS '客户来源描述（如“京东快递客户”“抖音私信”）';
COMMENT ON COLUMN "public"."customers"."created_at" IS '首次建档时间';
COMMENT ON COLUMN "public"."customers"."created_by" IS '建档人（后台用户 ID）';
COMMENT ON COLUMN "public"."customers"."updated_at" IS '最后更新时间（触发器自动维护）';
COMMENT ON COLUMN "public"."customers"."updated_by" IS '最后更新人（后台用户 ID）';
COMMENT ON COLUMN "public"."customers"."phones" IS '客户手机号数组（去重，支持多个）';
COMMENT ON COLUMN "public"."customers"."wechats" IS '微信号数组（去重，支持多个）';
COMMENT ON COLUMN "public"."customers"."douyins" IS '客户抖音号数组';
COMMENT ON COLUMN "public"."customers"."kwais" IS '客户快手号数组';
COMMENT ON COLUMN "public"."customers"."redbooks" IS '客户小红书号数组';
COMMENT ON COLUMN "public"."customers"."wework_openids" IS '企业微信客户 OpenID 数组';
COMMENT ON COLUMN "public"."customers"."work_phone" IS '工作/座机号码数组';
COMMENT ON COLUMN "public"."customers"."work_wechat" IS '工作微信（企业微信）';
COMMENT ON COLUMN "public"."customers"."province" IS '省（如湖北）';
COMMENT ON COLUMN "public"."customers"."city" IS '市（如孝感）';
COMMENT ON COLUMN "public"."customers"."district" IS '县/区（如孝南区）';
COMMENT ON COLUMN "public"."customers"."district_id" IS '区县行政代码（高德/民政部标准）';
COMMENT ON COLUMN "public"."customers"."street" IS '街道/乡镇';
COMMENT ON COLUMN "public"."customers"."address" IS '完整门牌地址';
COMMENT ON COLUMN "public"."customers"."lat" IS '纬度（WGS84）';
COMMENT ON COLUMN "public"."customers"."lon" IS '经度（WGS84）';
COMMENT ON COLUMN "public"."customers"."category" IS '客户分类（如“茶叶批发”“茶楼”）';
COMMENT ON COLUMN "public"."customers"."tags" IS '标签数组（如“大客户”“月结”“爱讲价”）';
COMMENT ON COLUMN "public"."customers"."level" IS '客户分级：0=未分级 1=S 2=A 3=B 4=C 10=X';
COMMENT ON COLUMN "public"."customers"."state" IS '客户状态：0=未知 1=未开发 2=开发中 3=已开发 4=已拉黑 5=已倒闭 6=同事 7=叛徒 8=同行';
COMMENT ON COLUMN "public"."customers"."kind" IS '店铺类型：0=未知 1=个体夫妻店 2=加盟连锁店 3=工厂直营店 4=其他';
COMMENT ON COLUMN "public"."customers"."added_wechat" IS '是否已添加微信好友';
COMMENT ON COLUMN "public"."customers"."credit_sale" IS '允许赊账额度（元）';
COMMENT ON COLUMN "public"."customers"."sellers" IS '所属销售员 ID 数组（支持多人跟进）';
COMMENT ON COLUMN "public"."customers"."last_visited" IS '最后线下拜访时间';
COMMENT ON COLUMN "public"."customers"."last_called" IS '最后线上联系时间';
COMMENT ON COLUMN "public"."customers"."last_order_date" IS '最后下单时间';
COMMENT ON COLUMN "public"."customers"."order_count" IS '累计下单次数（由销售单自动累计）';
COMMENT ON COLUMN "public"."customers"."avg_order_value" IS '平均订单金额（元）';
COMMENT ON COLUMN "public"."customers"."group_id" IS '所属群组 ID 数组（连锁、亲戚、战友等）';
COMMENT ON COLUMN "public"."customers"."birth_place" IS '出生地/口音描述（如“信阳固始口音”）';
COMMENT ON COLUMN "public"."customers"."birth_year" IS '出生年份';
COMMENT ON COLUMN "public"."customers"."birth_month" IS '出生月份';
COMMENT ON COLUMN "public"."customers"."birth_date" IS '出生日';
COMMENT ON COLUMN "public"."customers"."favors" IS '买货偏好 JSONB 数组：[{"product":"毛尖","avgPrice":158,"avgQuantity":10,"unit":"斤"}]';
COMMENT ON COLUMN "public"."customers"."products" IS '主营产品列表（去重后的商品名称数组）';
COMMENT ON COLUMN "public"."customers"."annual_turnover" IS '年营业额描述（如“50-100万”“100-300万”）';
COMMENT ON COLUMN "public"."customers"."shipping_infos" IS '收货信息 JSONB 数组：[{"receiver":"张三","phone":"18812345678","address":"湖北省孝感市...","isDefault":true}]';
COMMENT ON COLUMN "public"."customers"."preferred_delivery_method" IS '偏好发货方式（如“京东快递”“自提”）';
COMMENT ON COLUMN "public"."customers"."extra_info" IS '其他非结构化扩展信息（审核备注、促销记录等）';
COMMENT ON COLUMN "public"."customers"."saller_name" IS '销售员姓名';
COMMENT ON COLUMN "public"."customers"."system_tags" IS '系统标签字段，存储整型数组格式的标签ID数组，例如: {1,2,3}';
COMMENT ON TABLE "public"."customers" IS '客户主档案表：包含客户基础信息、联系方式、地址、交易画像、收货信息、赊账额度、分组关系等全量数据';
COMMENT ON INDEX "public"."idx_status" IS '状态';
COMMENT ON INDEX "public"."idx_phone" IS '电话';
COMMENT ON INDEX "public"."idx_wechats" IS '微信';

CREATE TABLE "public"."follow_up_records"(
                                             "id" int8 NOT NULL DEFAULT nextval('follow_up_records_id_seq'::regclass),
                                             "customer_id" int8 NOT NULL,
                                             "user_id" int8 NOT NULL,
                                             "kind" "public"."activity_kind_extended" DEFAULT 'other'::activity_kind_extended,
                                             "title" varchar(255),
                                             "content" text,
                                             "data" jsonb,
                                             "remark" text,
                                             "duration" int4,
                                             "location" varchar(255),
                                             "amount" numeric(15,2),
                                             "cost" numeric(15,2),
                                             "related_todo_id" int8,
                                             "parent_record_id" int8,
                                             "next_follow_time" timestamp(6),
                                             "next_follow_content" varchar(500),
                                             "attachments" jsonb,
                                             "photos" jsonb,
                                             "customer_satisfaction" int4,
                                             "customer_feedback" text,
                                             "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                             "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                             "deleted_at" timestamp(6),
                                             "is_deleted" bool DEFAULT false,
                                             CONSTRAINT "follow_up_records_customer_satisfaction_check" CHECK (((customer_satisfaction >= 1) AND (customer_satisfaction <= 5))),
                                             CONSTRAINT "follow_up_records_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."follow_up_records" OWNER TO "postgres";
CREATE INDEX "idx_follow_up_customer_id" ON "public"."follow_up_records" USING btree (customer_id);
CREATE INDEX "idx_follow_up_user_id" ON "public"."follow_up_records" USING btree (user_id);
CREATE INDEX "idx_follow_up_kind" ON "public"."follow_up_records" USING btree (kind);
CREATE INDEX "idx_follow_up_created_at" ON "public"."follow_up_records" USING btree (created_at);
CREATE INDEX "idx_follow_up_next_follow_time" ON "public"."follow_up_records" USING btree (next_follow_time);
CREATE INDEX "idx_follow_up_deleted" ON "public"."follow_up_records" USING btree (is_deleted);
CREATE INDEX "idx_follow_up_related_todo_id" ON "public"."follow_up_records" USING btree (related_todo_id);
CREATE INDEX "idx_follow_up_parent_record_id" ON "public"."follow_up_records" USING btree (parent_record_id);
CREATE INDEX "idx_follow_up_customer_time" ON "public"."follow_up_records" USING btree (customer_id, created_at DESC) WHERE (is_deleted = false);
CREATE INDEX "idx_follow_up_user_time" ON "public"."follow_up_records" USING btree (user_id, created_at DESC) WHERE (is_deleted = false);
CREATE INDEX "idx_follow_up_customer_kind" ON "public"."follow_up_records" USING btree (customer_id, kind) WHERE (is_deleted = false);
CREATE TRIGGER "update_follow_up_records_updated_at"
    BEFORE UPDATE ON "public"."follow_up_records"
    FOR EACH ROW
    EXECUTE FUNCTION "public"."update_follow_up_updated_at"() ;

COMMENT ON COLUMN "public"."follow_up_records"."id" IS '记录ID';
COMMENT ON COLUMN "public"."follow_up_records"."customer_id" IS '客户ID';
COMMENT ON COLUMN "public"."follow_up_records"."user_id" IS '创建用户ID';
COMMENT ON COLUMN "public"."follow_up_records"."kind" IS '跟进类型';
COMMENT ON COLUMN "public"."follow_up_records"."title" IS '跟进标题';
COMMENT ON COLUMN "public"."follow_up_records"."content" IS '跟进内容详情';
COMMENT ON COLUMN "public"."follow_up_records"."data" IS '结构化数据，JSON格式存储不同类型记录的特定字段';
COMMENT ON COLUMN "public"."follow_up_records"."remark" IS '备注';
COMMENT ON COLUMN "public"."follow_up_records"."duration" IS '持续时长（分钟）';
COMMENT ON COLUMN "public"."follow_up_records"."location" IS '地点';
COMMENT ON COLUMN "public"."follow_up_records"."amount" IS '金额（营业额等）';
COMMENT ON COLUMN "public"."follow_up_records"."cost" IS '成本';
COMMENT ON COLUMN "public"."follow_up_records"."related_todo_id" IS '关联的待办事项ID';
COMMENT ON COLUMN "public"."follow_up_records"."parent_record_id" IS '父记录ID（用于记录关联和回复）';
COMMENT ON COLUMN "public"."follow_up_records"."next_follow_time" IS '下次跟进时间';
COMMENT ON COLUMN "public"."follow_up_records"."next_follow_content" IS '下次跟进内容';
COMMENT ON COLUMN "public"."follow_up_records"."attachments" IS '附件信息（JSONB格式）';
COMMENT ON COLUMN "public"."follow_up_records"."photos" IS '照片信息（JSONB格式）';
COMMENT ON COLUMN "public"."follow_up_records"."customer_satisfaction" IS '客户满意度（1-5分）';
COMMENT ON COLUMN "public"."follow_up_records"."customer_feedback" IS '客户反馈';
COMMENT ON COLUMN "public"."follow_up_records"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."follow_up_records"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."follow_up_records"."deleted_at" IS '删除时间（软删除）';
COMMENT ON COLUMN "public"."follow_up_records"."is_deleted" IS '是否删除';
COMMENT ON TABLE "public"."follow_up_records" IS '跟进记录表（扩展版）';

CREATE TABLE "public"."groups"(
                                  "id" int8 NOT NULL DEFAULT nextval('groups_id_seq'::regclass),
                                  "name" varchar(256) NOT NULL,
                                  "roles" jsonb,
                                  "description" text,
                                  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                  "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                  "deleted_at" timestamp(6),
                                  "is_deleted" bool DEFAULT false,
                                  CONSTRAINT "groups_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."groups" OWNER TO "postgres";
CREATE INDEX "idx_groups_name" ON "public"."groups" USING btree (name);
CREATE INDEX "idx_groups_created_at" ON "public"."groups" USING btree (created_at);
CREATE INDEX "idx_groups_deleted" ON "public"."groups" USING btree (is_deleted);
COMMENT ON COLUMN "public"."groups"."id" IS '组ID';
COMMENT ON COLUMN "public"."groups"."name" IS '组名';
COMMENT ON COLUMN "public"."groups"."roles" IS '组成人员，格式：{"朋友": [2, 3], "同事": [4, 5]}';
COMMENT ON COLUMN "public"."groups"."description" IS '组描述';
COMMENT ON COLUMN "public"."groups"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."groups"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."groups"."deleted_at" IS '删除时间（软删除）';
COMMENT ON COLUMN "public"."groups"."is_deleted" IS '是否删除：false-否，true-是';
COMMENT ON TABLE "public"."groups" IS '客户组表';

CREATE TABLE "public"."tag_dimensions"(
                                          "id" int8 NOT NULL DEFAULT nextval('tag_dimensions_id_seq'::regclass),
                                          "name" varchar(128) NOT NULL,
                                          "description" text,
                                          "sort_order" int4 DEFAULT 0,
                                          "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                          "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                          "deleted_at" timestamp(6),
                                          "is_deleted" bool DEFAULT false,
                                          CONSTRAINT "tag_dimensions_name_key" UNIQUE (name),
                                          CONSTRAINT "tag_dimensions_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."tag_dimensions" OWNER TO "postgres";
CREATE INDEX "idx_tag_dimensions_name" ON "public"."tag_dimensions" USING btree (name);
CREATE INDEX "idx_tag_dimensions_sort_order" ON "public"."tag_dimensions" USING btree (sort_order);
CREATE INDEX "idx_tag_dimensions_deleted" ON "public"."tag_dimensions" USING btree (is_deleted);
COMMENT ON COLUMN "public"."tag_dimensions"."id" IS '维度ID';
COMMENT ON COLUMN "public"."tag_dimensions"."name" IS '维度名称';
COMMENT ON COLUMN "public"."tag_dimensions"."description" IS '维度描述';
COMMENT ON COLUMN "public"."tag_dimensions"."sort_order" IS '排序顺序';
COMMENT ON COLUMN "public"."tag_dimensions"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."tag_dimensions"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."tag_dimensions"."deleted_at" IS '删除时间（软删除）';
COMMENT ON COLUMN "public"."tag_dimensions"."is_deleted" IS '是否删除：false-否，true-是';
COMMENT ON TABLE "public"."tag_dimensions" IS '标签维度表';

CREATE TABLE "public"."tags"(
                                "id" int8 NOT NULL DEFAULT nextval('tags_id_seq'::regclass),
                                "dimension_id" int8 NOT NULL,
                                "name" varchar(128) NOT NULL,
                                "color" varchar(32),
                                "description" text,
                                "sort_order" int4 DEFAULT 0,
                                "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                "deleted_at" timestamp(6),
                                "is_deleted" bool DEFAULT false,
                                CONSTRAINT "tags_dimension_id_fkey" FOREIGN KEY (dimension_id) REFERENCES tag_dimensions(id),
                                CONSTRAINT "tags_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."tags" OWNER TO "postgres";
CREATE INDEX "idx_tags_dimension_id" ON "public"."tags" USING btree (dimension_id);
CREATE INDEX "idx_tags_name" ON "public"."tags" USING btree (name);
CREATE INDEX "idx_tags_sort_order" ON "public"."tags" USING btree (sort_order);
CREATE INDEX "idx_tags_deleted" ON "public"."tags" USING btree (is_deleted);
CREATE INDEX "idx_tags_dimension_name" ON "public"."tags" USING btree (dimension_id, name) WHERE (is_deleted = false);
COMMENT ON COLUMN "public"."tags"."id" IS '标签ID';
COMMENT ON COLUMN "public"."tags"."dimension_id" IS '维度ID';
COMMENT ON COLUMN "public"."tags"."name" IS '标签名称';
COMMENT ON COLUMN "public"."tags"."color" IS '标签颜色（十六进制）';
COMMENT ON COLUMN "public"."tags"."description" IS '标签描述';
COMMENT ON COLUMN "public"."tags"."sort_order" IS '排序顺序';
COMMENT ON COLUMN "public"."tags"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."tags"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."tags"."deleted_at" IS '删除时间（软删除）';
COMMENT ON COLUMN "public"."tags"."is_deleted" IS '是否删除：false-否，true-是';
COMMENT ON TABLE "public"."tags" IS '标签表';

CREATE TABLE "public"."todo_logs"(
                                     "id" int8 NOT NULL DEFAULT nextval('todo_logs_id_seq'::regclass),
                                     "todo_id" int8 NOT NULL,
                                     "operator_id" int8 NOT NULL,
                                     "action" "public"."todo_action" NOT NULL,
                                     "old_data" jsonb,
                                     "new_data" jsonb,
                                     "remark" varchar(500),
                                     "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                     CONSTRAINT "todo_logs_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."todo_logs" OWNER TO "postgres";
CREATE INDEX "idx_todo_logs_todo_id" ON "public"."todo_logs" USING btree (todo_id);
CREATE INDEX "idx_todo_logs_operator_id" ON "public"."todo_logs" USING btree (operator_id);
CREATE INDEX "idx_todo_logs_action" ON "public"."todo_logs" USING btree (action);
CREATE INDEX "idx_todo_logs_created_at" ON "public"."todo_logs" USING btree (created_at);
COMMENT ON COLUMN "public"."todo_logs"."id" IS '日志ID';
COMMENT ON COLUMN "public"."todo_logs"."todo_id" IS '待办ID';
COMMENT ON COLUMN "public"."todo_logs"."operator_id" IS '操作人ID';
COMMENT ON COLUMN "public"."todo_logs"."action" IS '操作类型';
COMMENT ON COLUMN "public"."todo_logs"."old_data" IS '变更前数据';
COMMENT ON COLUMN "public"."todo_logs"."new_data" IS '变更后数据';
COMMENT ON COLUMN "public"."todo_logs"."remark" IS '操作备注';
COMMENT ON COLUMN "public"."todo_logs"."created_at" IS '操作时间';
COMMENT ON TABLE "public"."todo_logs" IS '待办操作日志表';

CREATE TABLE "public"."todos"(
                                 "id" int8 NOT NULL DEFAULT nextval('todos_id_seq'::regclass),
                                 "customer_id" int8 NOT NULL,
                                 "creator_id" int8 NOT NULL,
                                 "executor_id" int8 NOT NULL,
                                 "title" varchar(255) NOT NULL,
                                 "content" text,
                                 "status" "public"."todo_status" DEFAULT 'pending'::todo_status,
                                 "planned_time" timestamp(6) NOT NULL,
                                 "completed_time" timestamp(6),
                                 "is_reminder" bool DEFAULT false,
                                 "reminder_type" "public"."reminder_type_new",
                                 "reminder_user_id" int8,
                                 "reminder_time" timestamp(6),
                                 "priority" "public"."todo_priority" DEFAULT 'medium'::todo_priority,
                                 "tags" jsonb,
                                 "attachments" jsonb,
                                 "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                 "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                 "deleted_at" timestamp(6),
                                 "is_deleted" bool DEFAULT false,
                                 CONSTRAINT "todos_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."todos" OWNER TO "postgres";
CREATE INDEX "idx_todos_customer_id" ON "public"."todos" USING btree (customer_id);
CREATE INDEX "idx_todos_creator_id" ON "public"."todos" USING btree (creator_id);
CREATE INDEX "idx_todos_executor_id" ON "public"."todos" USING btree (executor_id);
CREATE INDEX "idx_todos_status" ON "public"."todos" USING btree (status);
CREATE INDEX "idx_todos_planned_time" ON "public"."todos" USING btree (planned_time);
CREATE INDEX "idx_todos_created_at" ON "public"."todos" USING btree (created_at);
CREATE INDEX "idx_todos_deleted" ON "public"."todos" USING btree (is_deleted);
CREATE TRIGGER "update_todos_updated_at"
    BEFORE UPDATE ON "public"."todos"
    FOR EACH ROW
    EXECUTE FUNCTION "public"."update_updated_at_column"() ;

COMMENT ON COLUMN "public"."todos"."id" IS '待办ID';
COMMENT ON COLUMN "public"."todos"."customer_id" IS '关联客户ID';
COMMENT ON COLUMN "public"."todos"."creator_id" IS '创建人ID';
COMMENT ON COLUMN "public"."todos"."executor_id" IS '执行人ID';
COMMENT ON COLUMN "public"."todos"."title" IS '待办标题';
COMMENT ON COLUMN "public"."todos"."content" IS '待办内容详情';
COMMENT ON COLUMN "public"."todos"."status" IS '待办状态：未完成、已完成、延期、取消';
COMMENT ON COLUMN "public"."todos"."planned_time" IS '计划执行时间';
COMMENT ON COLUMN "public"."todos"."completed_time" IS '完成时间';
COMMENT ON COLUMN "public"."todos"."is_reminder" IS '是否提醒：false-否，true-是';
COMMENT ON COLUMN "public"."todos"."reminder_type" IS '提醒方式：微信、企业微信、两者';
COMMENT ON COLUMN "public"."todos"."reminder_user_id" IS '提醒人ID';
COMMENT ON COLUMN "public"."todos"."reminder_time" IS '提醒时间';
COMMENT ON COLUMN "public"."todos"."priority" IS '优先级';
COMMENT ON COLUMN "public"."todos"."tags" IS '标签（JSONB格式）';
COMMENT ON COLUMN "public"."todos"."attachments" IS '附件信息（JSONB格式）';
COMMENT ON COLUMN "public"."todos"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."todos"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."todos"."deleted_at" IS '删除时间（软删除）';
COMMENT ON COLUMN "public"."todos"."is_deleted" IS '是否删除：false-否，true-是';
COMMENT ON TABLE "public"."todos" IS '待办事项表';

CREATE TABLE "public"."users"(
                                 "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
                                 "name" varchar(256) NOT NULL,
                                 "manager_id" int8,
                                 "email" varchar(256),
                                 "phone" varchar(32),
                                 "department" varchar(128),
                                 "department_leader_id" int8,
                                 "position" varchar(128),
                                 "wechat_work_id" varchar(128),
                                 "wechat_id" varchar(128),
                                 "status" varchar(32) DEFAULT 'active'::character varying,
                                 "avatar_url" varchar(512),
                                 "last_login_at" timestamp(6),
                                 "created_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
                                 "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
                                 "deleted_at" timestamptz(6),
                                 "is_deleted" bool DEFAULT false,
                                 "username" text,
                                 CONSTRAINT "users_pkey" PRIMARY KEY (id)
);
ALTER TABLE "public"."users" OWNER TO "postgres";
CREATE INDEX "idx_users_deleted" ON "public"."users" USING btree (is_deleted);
CREATE INDEX "idx_users_manager_id" ON "public"."users" USING btree (manager_id);
CREATE INDEX "idx_users_phone" ON "public"."users" USING btree (phone);
CREATE INDEX "idx_users_department" ON "public"."users" USING btree (department);
CREATE INDEX "idx_users_department_leader_id" ON "public"."users" USING btree (department_leader_id);
CREATE INDEX "idx_users_wechat_work_id" ON "public"."users" USING btree (wechat_work_id);
CREATE INDEX "idx_users_wechat_id" ON "public"."users" USING btree (wechat_id);
CREATE INDEX "idx_users_name" ON "public"."users" USING btree (name);
CREATE INDEX "idx_users_email" ON "public"."users" USING btree (email);
CREATE INDEX "idx_users_status" ON "public"."users" USING btree (status);
CREATE INDEX "idx_users_created_at" ON "public"."users" USING btree (created_at);
COMMENT ON COLUMN "public"."users"."id" IS '用户ID';
COMMENT ON COLUMN "public"."users"."name" IS '用户名';
COMMENT ON COLUMN "public"."users"."manager_id" IS '主管ID';
COMMENT ON COLUMN "public"."users"."email" IS '邮箱';
COMMENT ON COLUMN "public"."users"."phone" IS '手机号';
COMMENT ON COLUMN "public"."users"."department" IS '部门';
COMMENT ON COLUMN "public"."users"."department_leader_id" IS '部门领导ID';
COMMENT ON COLUMN "public"."users"."position" IS '职位';
COMMENT ON COLUMN "public"."users"."wechat_work_id" IS '企业微信ID';
COMMENT ON COLUMN "public"."users"."wechat_id" IS '微信ID';
COMMENT ON COLUMN "public"."users"."status" IS '状态';
COMMENT ON COLUMN "public"."users"."avatar_url" IS '头像URL';
COMMENT ON COLUMN "public"."users"."last_login_at" IS '最后登录时间';
COMMENT ON COLUMN "public"."users"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."users"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."users"."deleted_at" IS '删除时间';
COMMENT ON COLUMN "public"."users"."is_deleted" IS '是否删除';
COMMENT ON TABLE "public"."users" IS '用户表';

