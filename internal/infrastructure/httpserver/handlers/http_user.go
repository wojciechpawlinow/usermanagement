package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/wojciechpawlinow/usermanagement/internal/application/service"
	"github.com/wojciechpawlinow/usermanagement/internal/domain"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
)

type UserHTTPHandler struct {
	validator   *validator.Validate
	userService service.UserPort
}

type createUserRequest struct {
	Email       string                      `json:"email" binding:"required,email" validate:"email"`
	Password    string                      `json:"password" binding:"required" validate:"min=8"`
	FirstName   string                      `json:"first_name" binding:"required" validate:"required,min=1,max=50"`
	LastName    string                      `json:"last_name" binding:"required" validate:"required,min=1,max=50"`
	PhoneNumber string                      `json:"phone_number" binding:"required" validate:"omitempty,min=9,max=15,numeric"`
	Addresses   []*createUserAddressRequest `json:"addresses" binding:"required" validate:"required,min=1,dive"`
}

type createUserAddressRequest struct {
	Type       int    `json:"type" binding:"required" validate:"required,oneof=1 2 3"`
	Street     string `json:"street" binding:"required" validate:"required,min=1,max=255"`
	City       string `json:"city" binding:"required" validate:"required,min=1,max=100"`
	State      string `json:"state" binding:"required" validate:"omitempty,min=1,max=100"`
	PostalCode string `json:"postal_code" binding:"required" validate:"required,min=1,max=20,alphanum"`
	Country    string `json:"country" binding:"required" validate:"omitempty,min=1,max=100,alpha"`
}

type updateUserRequest struct {
	Password    *string                     `json:"password" binding:"omitempty" validate:"omitempty,min=8"`
	FirstName   *string                     `json:"first_name" binding:"omitempty" validate:"omitempty,min=1,max=50"`
	LastName    *string                     `json:"last_name" binding:"omitempty" validate:"omitempty,min=1,max=50"`
	PhoneNumber *string                     `json:"phone_number" binding:"omitempty" validate:"omitempty,min=9,max=15,numeric"`
	Addresses   []*updateUserAddressRequest `json:"addresses" binding:"omitempty" validate:"omitempty,dive"`
}

type updateUserAddressRequest struct {
	Type       int     `json:"type" binding:"required" validate:"required,oneof=1 2 3"`
	Street     *string `json:"street" binding:"omitempty" validate:"omitempty,min=1,max=255"`
	City       *string `json:"city" binding:"omitempty" validate:"omitempty,min=1,max=100"`
	State      *string `json:"state" binding:"omitempty" validate:"omitempty,min=1,max=100"`
	PostalCode *string `json:"postal_code" binding:"omitempty" validate:"omitempty,min=1,max=20,alphanum"`
	Country    *string `json:"country" binding:"omitempty" validate:"omitempty,min=1,max=100,alpha"`
}

func NewUserHTTPHandler(v *validator.Validate, userService service.UserPort) *UserHTTPHandler {
	return &UserHTTPHandler{
		validator:   v,
		userService: userService,
	}
}

func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed hashing password"})
		return
	}

	userID := domain.NewID()

	createUserDTO := &service.CreateUserDTO{
		ID:          userID,
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Addresses:   make([]*service.CreateUserAddress, 0, len(req.Addresses)),
	}

	typeIsPresent := make(map[user.AddressType]struct{})
	for _, addr := range req.Addresses {
		if _, ok := typeIsPresent[user.AddressType(addr.Type)]; !ok {
			createUserDTO.Addresses = append(createUserDTO.Addresses, &service.CreateUserAddress{
				Type:       addr.Type,
				Street:     addr.Street,
				City:       addr.City,
				State:      addr.State,
				PostalCode: addr.PostalCode,
				Country:    addr.Country,
			})
		}
		typeIsPresent[user.AddressType(addr.Type)] = struct{}{}
	}

	if err = h.userService.Create(c.Request.Context(), createUserDTO); err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		default:
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"}) // do not leak the actual error reason
		}
		return
	}

	c.JSON(http.StatusCreated, map[string]string{"uuid": userID.String()})
}

func (h *UserHTTPHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateUserDTO := &service.UpdateUserDTO{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Addresses:   make([]*service.UpdateUserAddress, 0, len(req.Addresses)),
	}

	if req.Password != nil {
		hashedPassword, err := hashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed hashing password"})
			return
		}

		updateUserDTO.Password = &hashedPassword
	}

	typeIsPresent := make(map[user.AddressType]struct{})
	for _, addr := range req.Addresses {
		if _, ok := typeIsPresent[user.AddressType(addr.Type)]; !ok {
			updateUserDTO.Addresses = append(updateUserDTO.Addresses, &service.UpdateUserAddress{
				Type:       addr.Type,
				Street:     addr.Street,
				City:       addr.City,
				State:      addr.State,
				PostalCode: addr.PostalCode,
				Country:    addr.Country,
			})
		}
		typeIsPresent[user.AddressType(addr.Type)] = struct{}{}
	}

	if err := h.userService.Update(c.Request.Context(), userID, updateUserDTO); err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"}) // do not leak the actual error reason
		}
		return
	}

	c.JSON(http.StatusOK, "ok")
}

func (h *UserHTTPHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userService.Delete(c.Request.Context(), userID); err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"}) // do not leak the actual error reason
		}
		return
	}

	c.JSON(http.StatusOK, "ok")
}

func (h *UserHTTPHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	domainUser, err := h.userService.GetByUUID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"}) // do not leak the actual error reason
		}
		return
	}

	c.JSON(http.StatusOK, domainUser)
}

func (h *UserHTTPHandler) Get(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	iPage, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page param"})
		return
	}
	size := c.DefaultQuery("size", "5")
	iSize, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size param"})
		return
	}

	domainUsers, err := h.userService.Get(c.Request.Context(), iPage, iSize)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"}) // do not leak the actual error reason
		return
	}

	c.JSON(http.StatusOK, domainUsers)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
