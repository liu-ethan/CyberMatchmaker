-- 异步状态更新（对应步骤2）：
UPDATE fortune_records SET status = 'in_process', updated_at = NOW() WHERE id = ?;
-- LLM 调用完成，存入结果：
UPDATE fortune_records
SET status = 'completed',
    best_city = ?,
    recent_fortune = ?,
    partner_persona = ?,
    updated_at = NOW()
WHERE id = ?;
-- 匹配检索逻辑（对应步骤3）：
-- 当用户 A（假设是男性，寻找女性，要求年龄在 20-30 岁，在“北京”）发起匹配时，你可以使用 pgvector 提供的 <=> 操作符进行余弦距离排序。
SELECT
    user_id,
    wechat_id,
    -- 1 减去余弦距离即为余弦相似度
    1 - (partner_embedding <=> '[0.12, 0.45, ...]') AS similarity
FROM
    match_profiles
WHERE
    deleted_at IS NULL        -- 排除已退出的用户
  AND gender = 'female'     -- 异性限制
  AND city = '北京'          -- 城市过滤
  AND birth_date BETWEEN '1996-01-01' AND '2006-01-01' -- 岁数区间
ORDER BY
    partner_embedding <=> '[0.12, 0.45, ...]' -- 按余弦距离从小到大排序（最相似的排前面）
LIMIT 10;
-- 退出匹配（对应步骤4）：
UPDATE match_profiles SET deleted_at = NOW() WHERE user_id = ?;