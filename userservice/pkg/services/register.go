package services

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"userservice/pkg/clients/logger"
	"userservice/pkg/data/entity"
	"userservice/pkg/data/repository"
	"userservice/pkg/utils"
)

type registerRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,CheckPassword"`
	FirstName   string `json:"first_name" validate:"required"`
	Lastname    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// Register service
func Register(c *fiber.Ctx) error {
	request := c.Request()

	var req registerRequest

	err := json.Unmarshal(request.Body(), &req)
	if err != nil {

		resp := ErrorResponse{
			Status:       "failed",
			Error:        err,
			ErrorMessage: "invalid request body",
		}

		logger.SendLog(logger.Log{
			Source:   utils.CurrentTrace(),
			Request:  request.String(),
			Response: resp,
		})

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	validate := validator.New()
	_ = validate.RegisterValidation("CheckPassword", utils.CheckPassword)
	err = validate.Struct(req)
	if err != nil {

		resp := ErrorResponse{
			Status:       "failed",
			Error:        err,
			ErrorMessage: "invalid request body",
		}

		logger.SendLog(logger.Log{
			Source:   utils.CurrentTrace(),
			Request:  request.String(),
			Response: resp,
		})

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	pass, _ := utils.HashPassword(req.Password)

	user := entity.User{
		Username:    req.Username,
		Password:    pass,
		FirstName:   req.FirstName,
		LastName:    req.Lastname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	userId, err := repository.InsertUser(user)
	if err != nil {

		resp := ErrorResponse{
			Status:       "failed",
			Error:        err,
			ErrorMessage: "error while inserting user",
		}

		logger.SendErrLog(logger.ErrLog{
			Level:   "[Critical]",
			Source:  utils.CurrentTrace(),
			Message: resp.ErrorMessage,
			Error:   err.Error(),
		})

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	log.Default().Println("register new user : ", userId)
	user.Id = userId

	err = SendMailCode(user)
	if err != nil {

		resp := ErrorResponse{
			Status:       "failed",
			Error:        err,
			ErrorMessage: "error while sending verification code",
		}

		logger.SendErrLog(logger.ErrLog{
			Level:   "[Critical]",
			Source:  utils.CurrentTrace(),
			Message: resp.ErrorMessage,
			Error:   err.Error(),
		})

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	respUser := UserResponse{
		UserId:            userId,
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		PhoneNumber:       user.PhoneNumber,
		IsBlocked:         user.IsBlocked,
		LoginAttemptCount: user.LoginAttemptCount,
		IsVerified:        user.IsVerified,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}

	resp := Response{
		Status: "success",
		Data:   make(map[string]interface{}),
	}
	resp.Data["user"] = respUser

	logger.SendLog(logger.Log{
		Source:   utils.CurrentTrace(),
		Request:  request.String(),
		Response: resp,
	})

	return c.Status(fiber.StatusOK).JSON(resp)
}
