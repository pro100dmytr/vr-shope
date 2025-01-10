package service

import (
	"context"
	"fmt"
	"strconv"
	"vr-shope/internal/models/repositories"
	"vr-shope/internal/models/services"
	"vr-shope/internal/repository"
	"vr-shope/internal/utils"
	"vr-shope/internal/utils/uuids"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) CreateUser(ctx context.Context, userServ *services.User) error {
	err := utils.ValidateUser(userServ)
	if err != nil {
		return err
	}

	var email string = userServ.Email
	em := utils.IsValidEmail(email)
	if !em {
		return fmt.Errorf("invalid email: %s", email)
	}

	exists, err := s.repo.ExistsByEmail(ctx, userServ.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user with this email already exists")
	}

	hashedPassword, salt, err := utils.HashPassword(userServ.Password)
	if err != nil {
		return err
	}

	userServ.Password = hashedPassword

	userRepo := repositories.User{
		ID:          uuids.IntToUUID(int64(userServ.ID)),
		Login:       userServ.Login,
		Name:        userServ.Name,
		LastName:    userServ.LastName,
		PhoneNumber: userServ.PhoneNumber,
		Password:    userServ.Password,
		Email:       userServ.Email,
		Salt:        salt,
	}

	err = s.repo.Create(ctx, &userRepo)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Get(ctx context.Context, id int) (*services.User, error) {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*services.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) Update(ctx context.Context, userServ *services.User) error {
	exists, err := s.repo.ExistsByID(ctx, userServ.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	err = utils.ValidateUser(userServ)
	if err != nil {
		return err
	}

	var email string = userServ.Email
	em := utils.IsValidEmail(email)
	if !em {
		return fmt.Errorf("invalid email: %s", email)
	}

	err = s.repo.Update(ctx, userServ)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*services.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is empty")
	}

	em := utils.IsValidEmail(email)
	if !em {
		return nil, fmt.Errorf("invalid email: %s", email)
	}

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("email not found")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetToken(ctx context.Context, login string, password string) (string, error) {
	if login == "" || password == "" {
		return "", fmt.Errorf("invalid login or password")
	}

	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	isValidPassword, err := utils.CheckPassword(password, user.Password, user.Salt)
	if err != nil || !isValidPassword {
		return "", fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateToken(int(uuids.UUIDToInt(user.ID)))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUsersWithPagination(ctx context.Context, limit, offset string) ([]*services.User, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		return nil, err
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		return nil, err
	}

	repoUsers, err := s.repo.GetUsers(ctx, offsetInt, limitInt)
	if err != nil {
		return nil, err
	}

	var users []*services.User
	for _, repoUser := range repoUsers {
		user := &services.User{
			ID:              int(uuids.UUIDToInt(repoUser.ID)),
			Login:           repoUser.Login,
			Name:            repoUser.Name,
			LastName:        repoUser.LastName,
			PhoneNumber:     repoUser.PhoneNumber,
			Email:           repoUser.Email,
			WalletUSDT:      repoUser.WalletUSDT,
			NumberPurchases: repoUser.NumberPurchases,
		}
		users = append(users, user)
	}

	return users, nil
}
