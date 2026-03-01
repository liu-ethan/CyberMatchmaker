-- ==========================================
-- 0. 环境准备与旧表清理
-- ==========================================
-- 启用 pgvector 向量扩展（如果已经启用过则无视）
CREATE EXTENSION IF NOT EXISTS vector;

-- 按照依赖关系的倒序删除表，防止外键冲突。CASCADE 会连带删除相关的索引
DROP TABLE IF EXISTS match_profile CASCADE;
DROP TABLE IF EXISTS fortune_record CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;


-- ==========================================
-- 1. 用户表 ("user")
-- 核心作用：管理用户的登录账号和密码
-- 注意：'user' 在 PGSQL 中是保留字，表名必须带双引号
-- ==========================================
CREATE TABLE "user" (
                        id BIGSERIAL PRIMARY KEY,                             -- 唯一主键，用户的内部 ID
                        username VARCHAR(50) UNIQUE NOT NULL,                 -- 登录账号，不可重复
                        password VARCHAR(255) NOT NULL,                       -- 登录密码（按要求存明文，方便小项目开发调试）
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),    -- 账号注册时间
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),    -- 账号信息最后一次修改的时间
                        is_deleted SMALLINT DEFAULT 0                         -- 逻辑删除标记 (0: 正常使用中, 1: 账号已注销删除)
);

-- 为逻辑删除字段建立索引，加快过滤已注销用户的查询速度
CREATE INDEX idx_user_is_deleted ON "user"(is_deleted);


-- ==========================================
-- 2. 算命记录与对外展示表 (fortune_record)
-- 核心作用：记录用户提交的八字参数、异步大模型的计算状态，以及算命的详细结果
-- ==========================================
CREATE TABLE fortune_record (
                                id BIGSERIAL PRIMARY KEY,                             -- 唯一主键，每次算命都会生成一条新记录
                                user_id BIGINT NOT NULL REFERENCES "user"(id),        -- 外键：这条算命记录属于哪个用户

    -- 【前端传入的原始参数】
                                real_name VARCHAR(50) NOT NULL,                       -- 真实姓名
                                gender VARCHAR(10) NOT NULL,                          -- 性别（用于后续匹配异性时的硬性过滤）
                                birth_date DATE NOT NULL,                             -- 出生年月日（例如：1998-05-20）
                                birth_time VARCHAR(20) NOT NULL,                      -- 出生时辰（例如：子时 / 23:00）
                                current_city VARCHAR(100) NOT NULL,                   -- 所在城市（用于后续同城匹配的硬性过滤）

    -- 【大模型计算出的命理专业字段】
                                bazi VARCHAR(50),                                     -- 八字排盘结果（例如：甲子 乙丑 丙寅 丁卯）
                                five_elements VARCHAR(50),                            -- 五行属性分析（例如：水旺缺火）
                                zodiac_sign VARCHAR(10),                              -- 生肖（例如：龙）

    -- 【给用户自己看的隐私算命结果】
                                best_city VARCHAR(100),                               -- 最适合发展的城市或方位
                                recent_fortune TEXT,                                  -- 近期运势吉凶分析

    -- 【给匹配对象看的公开展示信息】
                                description TEXT,                                     -- 大模型生成的“玄学自我介绍”（例如：该命主五行属水，性格温润...）

    -- 【系统与消息队列状态控制】
                                status VARCHAR(20) DEFAULT 'pending',                 -- 异步任务状态 (pending:排队中, in_process:大模型计算中, completed:完成, failed:失败)
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),    -- 提交算命的时间
                                updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),    -- 状态最后更新的时间（用于前端轮询查进度）
                                is_deleted SMALLINT DEFAULT 0                         -- 逻辑删除标记 (0: 正常, 1: 用户删除了这条算命记录)
);

-- 为高频查询的 status 和 user_id 建立索引，极大地提高前端轮询查询进度的速度
CREATE INDEX idx_fortune_record_status ON fortune_record(status);
CREATE INDEX idx_fortune_record_user_id ON fortune_record(user_id);


-- ==========================================
-- 3. 匹配交友广场与向量表 (match_profile)
-- 核心作用：存储用于相亲匹配的数据，利用 pgvector 插件进行相似度检索
-- ==========================================
CREATE TABLE match_profile (
                               id BIGSERIAL PRIMARY KEY,                             -- 唯一主键
                               user_id BIGINT UNIQUE NOT NULL REFERENCES "user"(id), -- 外键：参与匹配的用户 ID（加了 UNIQUE 保证一人在广场只有一个坑位）

    -- 外键：强绑定某一次算命记录。匹配成功后，通过这个 ID 去 fortune_record 表拿 description 给对方看
                               fortune_record_id BIGINT NOT NULL REFERENCES fortune_record(id),

                               wechat_id VARCHAR(100) NOT NULL,                      -- 用户的微信号（匹配成功后发放给对方的“奖品”）

    -- 【用于 SQL 预过滤的标量字段】（为了性能，避免每次检索都去关联查询）
                               gender VARCHAR(10) NOT NULL,                          -- 冗余的性别字段（用于过滤，找男的还是找女的）
                               city VARCHAR(100) NOT NULL,                           -- 冗余的城市字段（用于过滤，是否要求同城）

    -- 【核心黑科技：大模型生成的向量】
    -- 存储“伴侣画像”的特征向量。假设你用 OpenAI 的 text-embedding-3-small，固定为 1536 维
                               partner_embedding vector(1536),

                               created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),    -- 加入匹配广场的时间
                               updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),    -- 匹配资料最后修改的时间
                               is_deleted SMALLINT DEFAULT 0                         -- 逻辑删除标记 (0: 在匹配池中, 1: 退出匹配广场，不再被别人搜到)
);

-- 普通索引：加速过滤“已退出广场的人”
CREATE INDEX idx_match_profile_is_deleted ON match_profile(is_deleted);
-- 复合索引：加速“性别+城市”的前置条件筛选（只对未退出广场的人建立索引，节省空间）
CREATE INDEX idx_match_profile_filters ON match_profile(gender, city) WHERE is_deleted = 0;

-- 向量索引 (HNSW 算法)：配合 vector_cosine_ops 使用余弦相似度计算，极速查出向量距离最近的异性
CREATE INDEX idx_match_profile_embedding ON match_profile USING hnsw (partner_embedding vector_cosine_ops);