package usecase

import (
	"context"
	"errors"
	"time"

	"stream-rec-core/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)

	CreateItem(ctx context.Context, item *domain.Item) error
	GetItems(ctx context.Context) ([]domain.Item, error)
	GetItemByID(ctx context.Context, id primitive.ObjectID) (*domain.Item, error)
	UpdateItem(ctx context.Context, id primitive.ObjectID, item *domain.Item) error
	DeleteItem(ctx context.Context, id primitive.ObjectID) error
	SearchItems(ctx context.Context, filter domain.ItemSearchFilter) ([]domain.Item, error)
	GetRecommendations(ctx context.Context, userID primitive.ObjectID, limit int) ([]domain.Item, error)

	RecordInteraction(ctx context.Context, inter *domain.Interaction) error
	GetUserInteractions(ctx context.Context, userID primitive.ObjectID) ([]domain.Interaction, error)
}

type PlatformService struct {
	repo      RepositoryInterface
	jwtSecret []byte
}

func NewPlatformService(repo RepositoryInterface, jwtSecret string) *PlatformService {
	return &PlatformService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *PlatformService) generateJWT(u *domain.User) (string, error) {
	claims := domain.JWTClaims{
		UserID: u.ID.Hex(),
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // Токен на 3 суток
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *PlatformService) ValidateToken(tokenStr string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *PlatformService) Register(ctx context.Context, u *domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	u.Password = ""
	return s.repo.CreateUser(ctx, u)
}

func (s *PlatformService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := s.generateJWT(u)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *PlatformService) GetProfile(ctx context.Context, idHex string) (*domain.User, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, id)
}

func (s *PlatformService) CreateItem(ctx context.Context, item *domain.Item) error {
	return s.repo.CreateItem(ctx, item)
}

func (s *PlatformService) GetItems(ctx context.Context) ([]domain.Item, error) {
	return s.repo.GetItems(ctx)
}

func (s *PlatformService) GetItemByID(ctx context.Context, idHex string) (*domain.Item, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	return s.repo.GetItemByID(ctx, id)
}

func (s *PlatformService) UpdateItem(ctx context.Context, idHex string, item *domain.Item) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	return s.repo.UpdateItem(ctx, id, item)
}

func (s *PlatformService) DeleteItem(ctx context.Context, idHex string) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	return s.repo.DeleteItem(ctx, id)
}

func (s *PlatformService) SearchItems(ctx context.Context, filter domain.ItemSearchFilter) ([]domain.Item, error) {
	return s.repo.SearchItems(ctx, filter)
}

func (s *PlatformService) RecordInteraction(ctx context.Context, inter *domain.Interaction) error {
	return s.repo.RecordInteraction(ctx, inter)
}

func (s *PlatformService) GetUserHistory(ctx context.Context, userIDHex string) ([]domain.Interaction, error) {
	id, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserInteractions(ctx, id)
}

func (s *PlatformService) GetRecommendations() {}
