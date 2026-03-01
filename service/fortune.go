/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"errors"
	"fmt"
	"time" // 别忘了引入 time 包

	"CyberMatchmaker/mapper"
	"CyberMatchmaker/model"
	"CyberMatchmaker/mq"
)

// SubmitFortune 提交算命请求（异步）
func SubmitFortune(userID int64, realName, gender, birthDateStr, birthTime, currentCity string) (int64, error) {
	// 1. 将前端传入的字符串 "YYYY-MM-DD" 解析为 time.Time 对象
	// 注意：Go 语言规定格式化模板必须使用固定的诞生时间 "2006-01-02"
	parsedDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return 0, errors.New("出生日期格式错误，请使用 YYYY-MM-DD 格式")
	}

	// 2. 先在数据库里占个坑位，状态设为 pending
	record := &model.FortuneRecord{
		UserID:      userID,
		RealName:    realName,
		Gender:      gender,
		BirthDate:   parsedDate, // 这里传入解析后的 time.Time 对象
		BirthTime:   birthTime,
		CurrentCity: currentCity,
		Status:      "pending",
	}

	if err := mapper.CreateFortuneRecord(record); err != nil {
		return 0, err
	}

	// 3. 拼接给大模型的精细提示词（Prompt）
	prompt := fmt.Sprintf(
		"我叫%s，性别%s，出生于%s %s，目前在%s。请帮我算一下八字排盘、五行属性、生肖，以及最适合发展的方位和近期运势。请根据你的专业知识给出一段详细的命理描述。",
		realName, gender, birthDateStr, birthTime, currentCity,
	)

	// 4. 构造任务载荷并推送到 MQ
	task := mq.FortuneTask{
		RecordID: record.ID,
		UserID:   userID,
		Prompt:   prompt,
	}

	if err := mq.PublishFortuneTask(task); err != nil {
		// 补偿机制：如果进队列失败，直接把 DB 状态标为 failed
		_ = mapper.UpdateFortuneRecord(record.ID, map[string]interface{}{"status": "failed"})
		return 0, errors.New("算命任务投递失败，请稍后重试")
	}

	// 5. 返回入库的记录 ID，供前端轮询
	return record.ID, nil
}

// GetFortuneResult 查询算命结果
// 参数增加 userID，用于配合 mapper 进行越权校验
func GetFortuneResult(recordID int64, userID int64) (interface{}, error) {
	// 1. 调用更新后的 mapper 层获取记录
	// 传入 userID 确保用户只能查询属于自己的记录
	record, err := mapper.GetFortuneRecordByID(recordID, userID)
	if err != nil {
		// 如果查不到或 userID 不匹配，统一返回“记录未找到”以保护隐私
		return nil, errors.New("算命记录不存在或无权访问")
	}

	// 2. 根据业务逻辑判断状态
	// 如果任务未完成 (pending/failed)，按照接口文档约定仅返回 status 字段
	if record.Status != "completed" {
		return map[string]interface{}{
			"status": record.Status,
		}, nil
	}

	// 3. 如果状态是 "completed"，返回完整的命理结果
	// 这里的 Map Key 必须严格遵守接口文档定义的下划线命名规范
	result := map[string]interface{}{
		"status":         record.Status,
		"bazi":           record.Bazi,          // 八字
		"five_elements":  record.FiveElements,  // 五行
		"zodiac_sign":    record.ZodiacSign,    // 生肖
		"best_city":      record.BestCity,      // 适合发展的城市
		"recent_fortune": record.RecentFortune, // 近期运势
		"description":    record.Description,   // 详细命理描述
	}

	return result, nil
}
