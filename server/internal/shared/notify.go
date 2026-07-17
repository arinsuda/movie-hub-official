package shared

import (
	"context"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
)

type NeutralUnlocked struct {
	ID   uint
	Name string
}

func TrackAndNotify(
	ctx context.Context,
	achieveSvc achievementsmodule.Service,
	notifSvc *notification_module.Service,
	userID uint,
	actionType string,
	count int,
) []NeutralUnlocked {
	if achieveSvc == nil {
		return nil
	}
	unlocked, _ := achieveSvc.Track(userID, actionType, count)
	
	var result []NeutralUnlocked
	for _, u := range unlocked {
		result = append(result, NeutralUnlocked{
			ID:   u.Achievement.ID,
			Name: u.Achievement.Name,
		})
		if notifSvc != nil {
			_ = notifSvc.PushAchievementUnlocked(ctx, userID, u.Achievement.ID, u.Achievement.Name, u.ExpGained)
		}
	}
	return result
}
