package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// WuHistoryWeek 存储每周武勋历史数据
type WuHistoryWeek struct {
	Id          int       `json:"id" gorm:"column:id;primaryKey"`
	GroupId     string    `json:"group_id" gorm:"column:group_id"`                                     // 分组ID
	GroupName   string    `json:"group_name" gorm:"column:group_name;index:idx_group_record,unique"`   // 分组名称，与记录日期组成唯一索引
	MemberCount int       `json:"member_count" gorm:"column:member_count"`                             // 成员数量
	TotalWu     int       `json:"total_wu" gorm:"column:total_wu"`                                     // 总武勋
	AverageWu   int       `json:"average_wu" gorm:"column:average_wu"`                                 // 平均武勋
	ZeroWuCount int       `json:"zero_wu_count" gorm:"column:zero_wu_count"`                           // 0武勋人数
	RecordDate  time.Time `json:"record_date" gorm:"column:record_date;index:idx_group_record,unique"` // 记录日期，与分组名称组成唯一索引
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`                                 // 创建时间
}

func (WuHistoryWeek) TableName() string {
	return "wu_history_week"
}

// SaveGroupWuHistory 保存分组武勋历史数据
func SaveGroupWuHistory() error {
	// 查询当前各分组的统计数据
	type GroupWuStats struct {
		Group       string `json:"group"`
		MemberCount int    `json:"member_count"`
		TotalWu     int    `json:"total_wu"`
		AverageWu   int    `json:"average_wu"`
		ZeroWuCount int    `json:"zero_wu_count"`
	}

	var stats []GroupWuStats

	subQuery := Conn.Model(&TeamUser{}).
		Select("`group`, COUNT(*) as zero_wu_count").
		Where("wu = 0").
		Group("`group`")

	err := Conn.Model(&TeamUser{}).
		Select("`team_user`.`group`, COUNT(*) as member_count, SUM(wu) as total_wu, ROUND(AVG(wu)) as average_wu, IFNULL(sub.zero_wu_count, 0) as zero_wu_count").
		Joins("LEFT JOIN (?) as sub ON sub.`group` = `team_user`.`group`", subQuery).
		Group("`team_user`.`group`").
		Scan(&stats).Error

	if err != nil {
		return err
	}

	currentTime := time.Now()
	// 只保留日期部分，确保时间部分为00:00:00
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	// 使用 Upsert 确保数据一致性，防止重复插入
	return Conn.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			// 构造要插入或更新的记录
			history := WuHistoryWeek{
				GroupName:   stat.Group,
				GroupId:     stat.Group,
				MemberCount: stat.MemberCount,
				TotalWu:     stat.TotalWu,
				AverageWu:   stat.AverageWu,
				ZeroWuCount: stat.ZeroWuCount,
				RecordDate:  currentDate, // 使用规范化后的日期
				CreatedAt:   currentTime,
			}
			// 使用 GORM 的 OnConflict 来处理冲突情况
			// 当基于 group_name 和 record_date 的组合键冲突时，更新现有记录
			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "group_name"}, {Name: "record_date"}},
				DoUpdates: clause.AssignmentColumns([]string{"member_count", "total_wu", "average_wu", "zero_wu_count", "created_at"}),
			}).Create(&history)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}
