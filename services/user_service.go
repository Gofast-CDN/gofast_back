package services

import (
	"errors"

	"gofast/models"
	"gofast/utils/auth"
	"gofast/utils/validator"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Token   string `json:"token"`
	UserID  string `json:"userId"`
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

func (s *UserService) Register(req *RegisterRequest) error {
	if err := s.emailValidator.Validate(req.Email); err != nil {
		return err
	}

	existingUser := &models.User{}
	err := s.collection.First(bson.M{"email": req.Email}, existingUser)
	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "USER",
	}

	return s.collection.Create(user)
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user := &models.User{}
	err := s.collection.First(bson.M{"email": req.Email}, user)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Message: "Login successful",
		Email:   user.Email,
		Role:    user.Role,
		Token:   token,
		UserID:  user.ID.Hex(),
	}, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user := &models.User{}
	err = s.collection.FindByID(objectID, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
