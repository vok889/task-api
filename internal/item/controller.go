package item

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB) Controller {
	return Controller{
		Service: NewService(db),
	}
}

type ApiError struct {
	Field  string
	Reason string
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "จำเป็นต้องกรอกข้อมูลนี้"
	case "email":
		return "Invalid email"
	case "gt":
		return fmt.Sprintf("Number must greater than %v", param)
	case "gte":
		return fmt.Sprintf("Number must greater than or equal %v", param)
	}
	return ""
}

func getValidationErrors(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Param())}
		}
		return out
	}
	return nil
}

func (controller Controller) CreateItem(ctx *gin.Context) {
	// Bind
	var request model.RequestCreateItem

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": getValidationErrors(err),
		})
		return
	}

	// get owner_id from context
	//ownerIdFloat := ctx.MustGet("uid").(float64)
	//ownerId := int(ownerIdFloat)
	ownerId := 1

	// Create item
	item, err := controller.Service.Create(request, ownerId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"data": item,
	})
}

func (controller Controller) FindItems(ctx *gin.Context) {
	// Bind query parameters
	var (
		request model.RequestFindItem
	)

	if err := ctx.BindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Find
	items, err := controller.Service.Find(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (controller Controller) FindItemByID(ctx *gin.Context) {
	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Find item
	item, err := controller.Service.FindByID(uint(id))
	if err != nil {
		// Check if the error is because the record was not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Item not found",
			})
			return
		}

		// Other types of errors (e.g., database issues)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error finding item: " + err.Error(),
		})
		return
	}

	// Successfully found item
	ctx.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func (controller Controller) UpdateItem(ctx *gin.Context) {
	// Bind
	var (
		request model.RequestUpdateItem
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Update item
	item, err := controller.Service.UpdateItem(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func (controller Controller) UpdateItemStatus(ctx *gin.Context) {
	// Bind
	var (
		request model.RequestPatchItemStatus
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Update status
	item, err := controller.Service.UpdateStatus(uint(id), request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func (controller Controller) DeleteItem(ctx *gin.Context) {
	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Delete
	if err := controller.Service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}
