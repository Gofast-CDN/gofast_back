package services

import (
	"errors"

	"gofast/models"
	"gofast/utils/auth"
	"gofast/utils/validator"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

type UserService struct {
	collection     *mgm.Collection
	emailValidator *validator.EmailValidator
}

func NewUserService() *UserService {
	return &UserService{
		collection:     mgm.Coll(&models.User{}),
		emailValidator: validator.NewEmailValidator(),
	}
}

func (s *UserService) Create(req *CreateUserRequest) (*UserResponse, error) {
	if err := s.emailValidator.Validate(req.Email); err != nil {
		return nil, err
	}

	existingUser := &models.User{}
	err := s.collection.First(bson.M{"email": req.Email}, existingUser)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "USER",
	}

	if err := s.collection.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	}, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	err := s.collection.FindByID(id, user)
	return user, err
}
