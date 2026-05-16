package follow_module

import "gorm.io/gorm"

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) Follow(requesterID, targetID uint) (*FollowResponse, error) {
	follow, err := s.repo.Follow(requesterID, targetID)
	if err != nil {
		return nil, err
	}
	return &FollowResponse{
		FollowerID: follow.FollowerID,
		FolloweeID: follow.FolloweeID,
		Status:     string(follow.Status),
	}, nil
}

func (s *Service) Unfollow(requesterID, targetID uint) error {
	return s.repo.Unfollow(requesterID, targetID)
}

// AcceptFollow — requesterID คือ followee (เจ้าของ account ที่กด accept)
func (s *Service) AcceptFollow(requesterID, followerID uint) error {
	return s.repo.AcceptFollow(followerID, requesterID)
}

// RejectFollow — requesterID คือ followee (เจ้าของ account ที่กด reject)
func (s *Service) RejectFollow(requesterID, followerID uint) error {
	return s.repo.RejectFollow(followerID, requesterID)
}

func (s *Service) GetFollowers(userID uint) ([]UserSummary, error) {
	rows, err := s.repo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

func (s *Service) GetFollowing(userID uint) ([]UserSummary, error) {
	rows, err := s.repo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

// GetPendingRequests — เฉพาะเจ้าของ account เท่านั้นที่ดูได้
func (s *Service) GetPendingRequests(requesterID, userID uint) ([]UserSummary, error) {
	if requesterID != userID {
		return nil, ErrForbidden
	}
	rows, err := s.repo.GetPendingRequests(userID)
	if err != nil {
		return nil, err
	}
	return toSummaryList(rows), nil
}

func toSummaryList(rows []listRow) []UserSummary {
	result := make([]UserSummary, len(rows))
	for i, r := range rows {
		result[i] = UserSummary{
			ID:          r.ID,
			Username:    r.Username,
			DisplayName: r.DisplayName,
			AvatarURL:   r.AvatarURL,
			Status:      r.Status,
		}
	}
	return result
}
