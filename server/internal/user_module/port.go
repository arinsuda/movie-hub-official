package user_module

type StatsProvider interface {
	GetLevel(userID uint) int
}

type EmailVerificationSender interface {
	SendVerification(userID uint, email string) error
}
