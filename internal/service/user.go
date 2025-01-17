package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
	"vr-shope/internal/models"
	"vr-shope/internal/repository"
	"vr-shope/internal/uuids"

	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	repo *repository.UserStorage
}

func NewUserService(repo *repository.UserStorage) *UserService {
	return &UserService{repo}
}

var secretKey = []byte("sfbwm37c7gd7c")

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func HashPassword(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)

	return hex.EncodeToString(hashedPassword), hex.EncodeToString(salt), nil
}

func CheckPassword(password, storedHash, storedSalt string) (bool, error) {
	storedHashBytes, err := hex.DecodeString(storedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode stored hash: %w", err)
	}

	storedSaltBytes, err := hex.DecodeString(storedSalt)
	if err != nil {
		return false, fmt.Errorf("failed to decode stored salt: %w", err)
	}

	hash := sha256.New()
	hash.Write(storedSaltBytes)
	hash.Write([]byte(password))
	computedHash := hash.Sum(nil)

	return bytes.Equal(computedHash, storedHashBytes), nil
}

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateUser(user *models.User) error {
	if user.Login == "" {
		return errors.New("login is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *UserService) CreateUser(ctx context.Context, userServ *models.User) error {
	err := ValidateUser(userServ)
	if err != nil {
		return err
	}

	var email string = userServ.Email
	em := IsValidEmail(email)
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

	hashedPassword, salt, err := HashPassword(userServ.Password)
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

	err = ValidateUser(userServ)
	if err != nil {
		return err
	}

	email := userServ.Email
	em := IsValidEmail(email)
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

	em := IsValidEmail(email)
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

	isValidPassword, err := CheckPassword(password, user.Password, user.Salt)
	if err != nil || !isValidPassword {
		return "", fmt.Errorf("invalid password")
	}

	token, err := GenerateToken(int(uuids.UUIDToInt(user.ID)))
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
