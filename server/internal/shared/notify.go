package shared

import (
	"context"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
)

func TrackAndNotify(
	ctx context.Context,
	achieveSvc achievementsmodule.Service,
	notifSvc *notification_module.Service,
	userID uint,
	actionType string,
	count int,
) {
	if achieveSvc == nil {
		return
	}
	unlocked, _ := achieveSvc.Track(userID, actionType, count)
	if notifSvc == nil {
		return
	}
	for _, u := range unlocked {
		_ = notifSvc.PushAchievementUnlocked(ctx, userID, u.Achievement.ID, u.Achievement.Name, u.ExpGained)
	}
}
