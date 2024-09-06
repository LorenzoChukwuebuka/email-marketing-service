package adminrepository

type AdminUsersRepository struct {
}

func NewAdminUserRepository() *AdminUsersRepository {
	return &AdminUsersRepository{}
}

func (r *AdminUsersRepository) GetAllUsers() error {
	return nil
}

func (r *AdminUsersRepository) BlockUser(userId string) error {
	return nil
}

func (r *AdminUsersRepository) GetSingleUser() error {
	return nil
}
