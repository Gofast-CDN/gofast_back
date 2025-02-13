package services

import (
	"context"
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

type RegisterResponse struct {
	UserID string `json:"userId"`
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

func (s *UserService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	if err := s.emailValidator.Validate(req.Email); err != nil {
		return nil, err
	}

	existingUser := &models.User{}
	err := s.collection.First(bson.M{"email": req.Email}, existingUser)
	if err == nil {
		return nil, errors.New("user already exists")
	}

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

	return &RegisterResponse{UserID: user.ID.Hex()}, nil
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

func (s *UserService) DeleteUserByID(userID string) error {
	// Convert userID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Perform the delete operation
	// Pass context.TODO() or an actual context, and the filter for deletion
	result, err := s.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return errors.New("error deleting user: " + err.Error())
	}

	// Check if a user was actually deleted
	if result.DeletedCount == 0 {
		return errors.New("no user found with the given ID")
	}

	return nil
}
