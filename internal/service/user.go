package service

import (
	"context"
	"fmt"
	"strconv"
	"vr-shope/internal/models"
	"vr-shope/internal/repository"
	"vr-shope/internal/utils"
	"vr-shope/internal/utils/uuids"
)

type UserService struct {
	repo *repository.UserStorage
}

func NewUserService(repo *repository.UserStorage) *UserService {
	return &UserService{repo}
}

func (s *UserService) CreateUser(ctx context.Context, userServ *models.User) error {
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

	userRepo := repository.User{
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

func (s *UserService) Get(ctx context.Context, id int) (*models.User, error) {
	exists, err := s.repo.ExistsByID(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	user, err := s.repo.GetByID(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		return nil, err
	}

	userServ := models.User{
		ID:              uuids.UUIDToInt(user.ID),
		Login:           user.Login,
		Name:            user.Name,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		WalletUSDT:      user.WalletUSDT,
		NumberPurchases: user.NumberPurchases,
	}

	return &userServ, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	usersRepo, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var usersServ []*models.User
	for _, user := range usersRepo {
		userServ := models.User{
			ID:              uuids.UUIDToInt(user.ID),
			Login:           user.Login,
			Name:            user.Name,
			LastName:        user.LastName,
			PhoneNumber:     user.PhoneNumber,
			Email:           user.Email,
			WalletUSDT:      user.WalletUSDT,
			NumberPurchases: user.NumberPurchases,
		}

		usersServ = append(usersServ, &userServ)
	}

	return usersServ, nil
}

func (s *UserService) Update(ctx context.Context, userServ *models.User) error {
	exists, err := s.repo.ExistsByID(ctx, uuids.IntToUUID(int64(userServ.ID)))
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

	email := userServ.Email
	em := utils.IsValidEmail(email)
	if !em {
		return fmt.Errorf("invalid email: %s", email)
	}

	userRepo := repository.User{
		ID:              uuids.IntToUUID(int64(userServ.ID)),
		Login:           userServ.Login,
		Name:            userServ.Name,
		LastName:        userServ.LastName,
		PhoneNumber:     userServ.PhoneNumber,
		Email:           userServ.Email,
		WalletUSDT:      userServ.WalletUSDT,
		NumberPurchases: userServ.NumberPurchases,
	}

	err = s.repo.Update(ctx, &userRepo)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	exists, err := s.repo.ExistsByID(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	err = s.repo.Delete(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
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

	userServ := models.User{
		ID:              uuids.UUIDToInt(user.ID),
		Login:           user.Login,
		Name:            user.Name,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		WalletUSDT:      user.WalletUSDT,
		NumberPurchases: user.NumberPurchases,
	}

	return &userServ, nil
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

func (s *UserService) GetUsersWithPagination(ctx context.Context, limit, offset string) ([]*models.User, error) {
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

	var users []*models.User
	for _, repoUser := range repoUsers {
		user := &models.User{
			ID:              uuids.UUIDToInt(repoUser.ID),
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
