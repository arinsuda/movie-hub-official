package achievementsmodule

import (
	"context"
	"errors"
	"time"

	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	"gorm.io/gorm"
)

type Service interface {
	ListAllAchievements(filter PaginationQuery) (ListAchievementsResponse, error)
	ListUserAchievements(ctx context.Context, userID, requesterID uint, filter UserAchievementFilter) (ListUserAchievementsResponse, error)

	Track(userID uint, actionType string, newCount int) ([]NewlyUnlocked, error)
}

type service struct {
	repo   Repository
	policy privacy_policy.UserAccessPolicy
}

func newService(repo Repository, policy privacy_policy.UserAccessPolicy) Service {
	return &service{repo: repo, policy: policy}
}

func (s *service) ListAllAchievements(filter PaginationQuery) (ListAchievementsResponse, error) {
	filter.Normalize()
	achievements, total, err := s.repo.ListAllAchievements(filter)
	if err != nil {
		return ListAchievementsResponse{}, err
	}

	data := make([]AchievementResponse, len(achievements))
	for i, a := range achievements {
		data[i] = toAchievementResponse(a)
	}

	return ListAchievementsResponse{
		Data:       data,
		Pagination: newPaginationMeta(filter.Page, filter.Limit, total),
	}, nil
}

func (s *service) ListUserAchievements(ctx context.Context, userID, requesterID uint, filter UserAchievementFilter) (ListUserAchievementsResponse, error) {
	if s.policy != nil {
		canView, err := s.policy.CanViewProfileSection(ctx, requesterID, userID, privacy_policy.SectionAchievements)
		if err != nil {
			return ListUserAchievementsResponse{}, err
		}
		if !canView {
			return ListUserAchievementsResponse{}, privacy_policy.ErrForbidden
		}
	}

	filter.Normalize()
	list, total, err := s.repo.ListUserAchievements(userID, filter)
	if err != nil {
		return ListUserAchievementsResponse{}, err
	}

	data := make([]UserAchievementResponse, len(list))
	for i, ua := range list {
		data[i] = toUserAchievementResponse(ua)
	}

	return ListUserAchievementsResponse{
		Data:       data,
		Pagination: newPaginationMeta(filter.Page, filter.Limit, total),
	}, nil
}

func (s *service) Track(userID uint, actionType string, newCount int) ([]NewlyUnlocked, error) {
	achievements, err := s.repo.FindAchievementsByActionType(actionType)
	if err != nil {
		return nil, err
	}

	var newlyUnlocked []NewlyUnlocked

	for _, a := range achievements {
		ua, err := s.repo.FindUserAchievement(userID, a.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
				ua = &UserAchievement{
					UserID:        userID,
					AchievementID: a.ID,
					CurrentCount:  0,
					IsUnlocked:    false,
				}
			} else {
				return nil, err
			}
		}

		if ua.IsUnlocked {
			continue
		}

		ua.CurrentCount = newCount

		if newCount >= a.TargetCount {
			now := time.Now().UTC()
			ua.IsUnlocked = true
			ua.UnlockedAt = &now

			newlyUnlocked = append(newlyUnlocked, NewlyUnlocked{
				Achievement: toAchievementResponse(a),
				ExpGained:   a.ExpReward,
			})
		}

		if err := s.repo.UpsertUserAchievement(ua); err != nil {
			return nil, err
		}
	}

	return newlyUnlocked, nil
}

// ── Response Mapping ──────────────────────────────────────────────

func toAchievementResponse(a Achievement) AchievementResponse {
	return AchievementResponse{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		ExpReward:   a.ExpReward,
		ActionType:  a.ActionType,
		TargetCount: a.TargetCount,
	}
}

func toUserAchievementResponse(ua UserAchievement) UserAchievementResponse {
	pct := 0.0
	if ua.Achievement.TargetCount > 0 {
		pct = float64(ua.CurrentCount) / float64(ua.Achievement.TargetCount) * 100
		if pct > 100 {
			pct = 100
		}
	}
	return UserAchievementResponse{
		AchievementID: ua.AchievementID,
		Achievement:   toAchievementResponse(ua.Achievement),
		CurrentCount:  ua.CurrentCount,
		IsUnlocked:    ua.IsUnlocked,
		UnlockedAt:    ua.UnlockedAt,
		ProgressPct:   pct,
	}
}
