package stat

import (
	"go/url-shortening/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) { // Статистика по кликам
	// Если нет статистики за сегодня по ссылке - создаем
	// Если есть - увеличиваем на 1

	var stat Stat
	currentDate := datatypes.Date(time.Now())

	repo.Db.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)
	if stat.ID == 0 { // Делаем запись за сегодня
		repo.Db.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks += 1 // Добовляем единицу
		repo.Db.Save(&stat)
	}

}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {

	// Делаем запрос:

	// SELECT to_char(date, 'YYYY-MM') as period, sum(clicks) FROM stats
	// WHERE date BETWEEN '01/01/2024' and '05/01/2025'
	// GROUP BY period
	// ORDER BY period

	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	repo.DB.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats

}
