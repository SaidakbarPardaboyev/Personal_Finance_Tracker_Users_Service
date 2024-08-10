package handler

import (
	"auth_service/models"
	"auth_service/pkg/logger"
	"auth_service/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	log     logger.ILogger
	storage storage.IStorage
}

func NewHandler(log logger.ILogger, storage storage.IStorage) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func handleResponse(ctx *gin.Context, log logger.ILogger, msg string, statusCode int, data interface{}) {

	resp := models.Response{}

	switch code := statusCode; {
	case code == 200:
		resp.Description = "OK"
		log.Info("~~~> OK", logger.String("msg", msg), logger.Any("status", statusCode))
	case code == 401:
		resp.Description = "Unauthorized"
		log.Info("???? Unauthorized", logger.String("msg", msg), logger.Any("status", statusCode))
	case code == 500:
		resp.Description = "Bad Request"
		log.Info("!!!! BAD REQUEST", logger.String("msg", msg), logger.Any("status", statusCode), logger.Any("Error", data))
	default:
		resp.Description = "Internal Server Error"
		log.Info("!!!! INTERNAL SERVER ERROR", logger.String("msg", msg), logger.Any("status", statusCode), logger.Any("Error", data))
	}

	resp.StatusCode = statusCode
	resp.Data = data

	ctx.JSON(resp.StatusCode, resp)
}
