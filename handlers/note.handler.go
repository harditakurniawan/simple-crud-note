package handlers

import (
	"errors"
	"fmt"
	"math"
	databases "simple-crud-notes/databases/pgsql"
	"simple-crud-notes/databases/pgsql/entities"
	"simple-crud-notes/utils"
	"simple-crud-notes/utils/enum"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NoteResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   *string   `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const noteCondition = "id = ? AND user_id = ?"

func CreateNote(c *fiber.Ctx) error {
	validRequest, ok := c.Locals("validatedDTO").(*utils.CreateNoteDto)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			c.OriginalURL(),
			"Failed to get validated request",
		))
	}

	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				enum.USER_NOT_FOUND_IN_CONTEXT,
			),
		)
	}

	// convert string of userID to uint
	var userIDUint uint
	_, err := fmt.Sscanf(userInfo.UserID, "%d", &userIDUint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Failed to convert user ID to uint",
			),
		)
	}

	newNote := entities.Note{
		Title:   validRequest.Title,
		Content: &validRequest.Content,
		UserID:  userIDUint,
	}

	result := databases.DB.Create(&newNote)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				result.Error.Error(),
			),
		)
	}

	response := NoteResponse{
		ID:        newNote.ID,
		Title:     newNote.Title,
		Content:   newNote.Content,
		CreatedAt: newNote.CreatedAt,
		UpdatedAt: newNote.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(
		utils.SuccessResponse(
			fiber.StatusCreated,
			c.OriginalURL(),
			response,
		),
	)
}

func GetNotes(c *fiber.Ctx) error {
	order, ok := c.Locals("order").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Order not found in context",
			),
		)
	}

	offset, ok := c.Locals("offset").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Offset not found in context",
			),
		)
	}

	limit, ok := c.Locals("limit").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Limit not found in context",
			),
		)
	}

	pageNumber, ok := c.Locals("pageNumber").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"pageNumber not found in context",
			),
		)
	}

	search, ok := c.Locals("search").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Search not found in context",
			),
		)
	}

	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				enum.USER_NOT_FOUND_IN_CONTEXT,
			),
		)
	}

	var (
		notes        []entities.Note
		totalRecords int64
	)

	queryAllNotes := databases.DB.Where("user_id = ?", userInfo.UserID)
	queryTotalRecords := databases.DB.Model(&entities.Note{}).Where("user_id = ?", userInfo.UserID)

	if search != "" {
		queryTotalRecords = queryTotalRecords.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
		queryAllNotes = queryAllNotes.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	queryAllNotes.Order(order).Limit(limit).Offset(offset).Find(&notes)
	if queryAllNotes.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				queryAllNotes.Error.Error(),
			),
		)
	}

	queryTotalRecords.Count(&totalRecords)
	if queryTotalRecords.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				queryTotalRecords.Error.Error(),
			),
		)
	}

	response := make([]NoteResponse, len(notes))
	for i, note := range notes {
		response[i] = NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.PaginationResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			response,
			utils.PaginationMeta{
				LAST_PAGE: int(math.Ceil(float64(totalRecords) / float64(limit))),
				PER_PAGE:  limit,
				PAGE:      pageNumber,
				TOTAL:     totalRecords,
			},
		),
	)
}

func DetailNote(c *fiber.Ctx) error {
	noteId := c.Params("id")
	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				enum.USER_NOT_FOUND_IN_CONTEXT,
			),
		)
	}

	var detailNote entities.Note
	result := databases.DB.First(&detailNote, noteCondition, noteId, userInfo.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				utils.ErrorResponse(
					fiber.StatusNotFound,
					c.OriginalURL(),
					"Note not found",
				),
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				result.Error.Error(),
			),
		)
	}

	response := NoteResponse{
		ID:        detailNote.ID,
		Title:     detailNote.Title,
		Content:   detailNote.Content,
		CreatedAt: detailNote.CreatedAt,
		UpdatedAt: detailNote.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.SuccessResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			response,
		),
	)
}

func UpdateNote(c *fiber.Ctx) error {
	validRequest, ok := c.Locals("validatedDTO").(*utils.UpdateNoteDto)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			c.OriginalURL(),
			"Failed to get validated request",
		))
	}

	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				enum.USER_NOT_FOUND_IN_CONTEXT,
			),
		)
	}

	noteId, err := c.ParamsInt("id")
	if err != nil || noteId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			utils.ErrorResponse(
				fiber.StatusBadRequest,
				c.OriginalURL(),
				"Invalid note ID",
			),
		)
	}

	var note entities.Note
	result := databases.DB.Where(noteCondition, noteId, userInfo.UserID).First(&note)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				utils.ErrorResponse(
					fiber.StatusNotFound,
					c.OriginalURL(),
					"Note not found or does not belong to user",
				),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				result.Error.Error(),
			),
		)
	}

	if validRequest.Title != nil {
		note.Title = *validRequest.Title
	}
	if validRequest.Content != nil {
		note.Content = validRequest.Content
	}

	if result := databases.DB.Save(&note); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				result.Error.Error(),
			),
		)
	}

	response := NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.SuccessResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			response,
		),
	)
}

func DeleteNote(c *fiber.Ctx) error {
	noteId := c.Params("id")
	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				enum.USER_NOT_FOUND_IN_CONTEXT,
			),
		)
	}

	result := databases.DB.Delete(&entities.Note{}, noteCondition, noteId, userInfo.UserID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				result.Error.Error(),
			),
		)
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			utils.ErrorResponse(
				fiber.StatusNotFound,
				c.OriginalURL(),
				"Note not found or you do not have permission to delete this note",
			),
		)
	}

	response := map[string]interface{}{
		"message": "Note deleted successfully",
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.SuccessResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			response,
		),
	)
}
