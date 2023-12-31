package managers

//go:generate mockgen -destination=../../../../mocks/mock_BaseUsersManager.go -package=mocks github.com/erupshis/zero_agency_test/internal/auth/users/managers BaseUsersManager
type BaseUsersManager interface {
	AddUser(login string, password string) (int64, error)
	GetUserId(login string) (int64, error)
	GetUserRole(userID int64) (int, error)
	ValidateUser(login string, password string) (bool, error)
}
