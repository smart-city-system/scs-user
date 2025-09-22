package http

import (
	"scs-user/internal/dto"
	services "scs-user/internal/services"
	"scs-user/pkg/errors"
	"scs-user/pkg/validation"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Handler
type UserHandler struct {
	svc services.UserService
}

// NewHandler constructor
func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserDto true "User creation data"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		createUserDto := &dto.CreateUserDto{}
		if err := c.Bind(createUserDto); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}

		// Validate the DTO
		if err := validation.ValidateStruct(createUserDto); err != nil {
			return err
		}

		createdUser, err := h.svc.CreateUser(c.Request().Context(), createUserDto)
		if err != nil {
			return err
		}
		createdUser.Password = ""
		return c.JSON(201, createdUser)
	}
}

// GetUsers godoc
// @Summary Get paginated list of users
// @Description Retrieve a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dto.UserListResponse "Users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		page := c.QueryParam("page")
		limit := c.QueryParam("limit")
		if page == "" {
			page = "1"
		}
		if limit == "" {
			limit = "10"
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return errors.NewBadRequestError("Invalid page number")
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return errors.NewBadRequestError("Invalid limit")
		}

		users, err := h.svc.GetUsers(c.Request().Context(), pageInt, limitInt)
		if err != nil {
			return err
		}

		return c.JSON(200, users)
	}
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the profile information of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User "User profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("user_id").(string)
		user, err := h.svc.GetUserByID(c.Request().Context(), userId)
		if err != nil {
			return err
		}
		return c.JSON(200, user)
	}
}

// VerifyAccount godoc
// @Summary Verify user account
// @Description Verify a user account using a verification token
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.VerifyAccountRequest true "Verification token"
// @Success 200 {object} string "Account verified successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/verify [post]
func (h *UserHandler) VerifyAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		verifyAccountDto := &dto.VerifyAccountRequest{}
		if err := c.Bind(verifyAccountDto); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}
		// Validate the DTO
		if err := validation.ValidateStruct(verifyAccountDto); err != nil {
			return err
		}

		err := h.svc.VerifyAccount(c.Request().Context(), verifyAccountDto.Token)
		if err != nil {
			return err
		}
		return c.JSON(200, "success")
	}
}

// func (h *UserHandler) GetAssignments() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		userID := "72b194cd-3cb1-4653-b7d5-ed2fc032ed62"
// 		assignments, err := h.svc.GetAssignments(c.Request().Context(), userID)
// 		if err != nil {
// 			return err
// 		}
// 		return c.JSON(200, assignments)
// 	}
// }
// func (h *UserHandler) CompleteStep() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		assignmentId := c.Param("assignmentId")
// 		stepId := c.Param("stepId")
// 		err := h.svc.CompleteStep(c.Request().Context(), assignmentId, stepId)
// 		if err != nil {
// 			return err
// 		}
// 		return c.JSON(200, "success")
// 	}
// }
