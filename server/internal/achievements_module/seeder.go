package achievementsmodule

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedFromFile(db *gorm.DB, jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	type jsonItem struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ExpReward   int    `json:"exp_reward"`
		ActionType  string `json:"action_type"`
		TargetCount int    `json:"target_count"`
	}

	var items []jsonItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	achievements := make([]Achievement, len(items))
	for i, it := range items {
		achievements[i] = Achievement{
			Name:        it.Name,
			Description: it.Description,
			ExpReward:   it.ExpReward,
			ActionType:  it.ActionType,
			TargetCount: it.TargetCount,
		}
	}

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"description",
			"exp_reward",
			"action_type",
			"target_count",
			"updated_at",
		}),
	}).CreateInBatches(&achievements, 50)

	if result.Error != nil {
		return result.Error
	}

	log.Printf("✅ Achievements seeded: %d records", result.RowsAffected)
	return nil
}
