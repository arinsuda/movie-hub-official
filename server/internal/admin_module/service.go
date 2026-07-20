package admin_module

import (
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
)

type Service struct {
	repo Repository
	hub  *notification_module.Hub
}

func NewService(repo Repository, hub *notification_module.Hub) *Service {
	return &Service{
		repo: repo,
		hub:  hub,
	}
}

func (s *Service) GetOverview() (*OverviewStats, error) {
	onlineCount := 0
	if s.hub != nil {
		onlineCount = s.hub.UniqueOnlineCount()
	}
	return s.repo.GetOverviewStats(onlineCount)
}

func (s *Service) GetGrowth() ([]GrowthPoint, error) {
	return s.repo.Get12MonthGrowthTrend()
}

func (s *Service) ListUsers(filter UserFilter) (PaginatedResponse[AdminUserRow], error) {
	items, total, err := s.repo.ListUsers(filter)
	if err != nil {
		return PaginatedResponse[AdminUserRow]{}, err
	}
	return NewPaginatedResponse(items, total, filter.Page, filter.Limit), nil
}

func (s *Service) ListReviews(filter ReviewFilter) (PaginatedResponse[AdminReviewRow], error) {
	items, total, err := s.repo.ListReviews(filter)
	if err != nil {
		return PaginatedResponse[AdminReviewRow]{}, err
	}
	return NewPaginatedResponse(items, total, filter.Page, filter.Limit), nil
}

func (s *Service) ListAuditLogs(filter AuditLogFilter) (PaginatedResponse[AdminAuditLogRow], error) {
	items, total, err := s.repo.ListAuditLogs(filter)
	if err != nil {
		return PaginatedResponse[AdminAuditLogRow]{}, err
	}
	return NewPaginatedResponse(items, total, filter.Page, filter.Limit), nil
}

func (s *Service) UpdateUserRole(adminID, targetUserID uint, req UpdateRoleRequest) error {
	if req.Role != "admin" && req.Role != "user" {
		return ErrRoleNotFound
	}
	return s.repo.UpdateUserRole(adminID, targetUserID, req.Role, req.Reason)
}

func (s *Service) UpdateUserStatus(adminID, targetUserID uint, req UpdateStatusRequest) error {
	return s.repo.UpdateUserStatus(adminID, targetUserID, req.IsActive, req.Reason)
}

func (s *Service) DeleteReview(adminID, reviewID uint, req DeleteReviewRequest) error {
	return s.repo.DeleteReview(adminID, reviewID, req.Reason)
}
