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
