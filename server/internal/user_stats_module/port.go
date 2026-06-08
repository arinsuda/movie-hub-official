package user_stats_module

// ExpAdder is the port that other modules depend on.
// Any module that awards EXP imports only this interface — not the concrete Service.
type ExpAdder interface {
	AddExperience(userID uint, exp int) error
}
