package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/elghazx/go-vault/internal/core/services"
	"github.com/elghazx/go-vault/templates"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	fileService *services.FileService
	authService *services.AuthService
}

func NewFileHandler(fileService *services.FileService, authService *services.AuthService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		authService: authService,
	}
}

func (h *FileHandler) Upload(c echo.Context) error {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		return templates.ErrorMessage("Please login first").Render(c.Request().Context(), c.Response().Writer)
	}

	if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
		return templates.ErrorMessage("File too large (max 32MB)").Render(c.Request().Context(), c.Response().Writer)
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return templates.ErrorMessage("No file selected").Render(c.Request().Context(), c.Response().Writer)
	}
	defer file.Close()

	isOneTime := c.FormValue("onetime") == "true"

	uploadedFile, err := h.fileService.UploadFile(
		c.Request().Context(),
		file,
		header.Filename,
		header.Header.Get("Content-Type"),
		userID,
		isOneTime,
	)
	if err != nil {
		return templates.ErrorMessage("Upload failed").Render(c.Request().Context(), c.Response().Writer)
	}

	return templates.UploadSuccess(uploadedFile).Render(c.Request().Context(), c.Response().Writer)
}

func (h *FileHandler) Download(c echo.Context) error {
	uuid := strings.TrimPrefix(c.Request().URL.Path, "/d/")
	if uuid == "" {
		return c.String(http.StatusBadRequest, "Invalid file ID")
	}

	reader, file, err := h.fileService.DownloadFile(c.Request().Context(), uuid)
	if err != nil {
		return c.String(http.StatusNotFound, "File not found")
	}
	defer reader.Close()

	c.Response().Header().Set("Content-Type", file.ContentType)
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\""+file.FileName+"\"")

	if file.FileSize > 0 {
		c.Response().Header().Set("Content-Length", strconv.FormatInt(file.FileSize, 10))
	}

	_, err = io.Copy(c.Response().Writer, reader)
	return err
}

func (h *FileHandler) GetMyFiles(c echo.Context) error {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		return templates.ErrorMessage("Please login first").Render(c.Request().Context(), c.Response().Writer)
	}

	files, err := h.fileService.GetUserFiles(c.Request().Context(), userID)
	if err != nil {
		return templates.ErrorMessage("Failed to load files").Render(c.Request().Context(), c.Response().Writer)
	}

	return templates.FileList(files).Render(c.Request().Context(), c.Response().Writer)
}

func (h *FileHandler) Preview(c echo.Context) error {
	uuid := strings.TrimPrefix(c.Request().URL.Path, "/f/")
	if uuid == "" {
		return c.String(http.StatusBadRequest, "Invalid file ID")
	}

	file, err := h.fileService.GetFileMetadata(c.Request().Context(), uuid)
	if err != nil {
		return h.renderNotFound(c)
	}

	if file.IsExpired() {
		return h.renderExpired(c)
	}

	return templates.FilePreview(file, c.Request().Host).Render(c.Request().Context(), c.Response().Writer)
}

func (h *FileHandler) getUserIDFromContext(c echo.Context) int64 {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return 0
	}

	userID, err := h.authService.ValidateToken(cookie.Value)
	if err != nil {
		return 0
	}

	return userID
}

func (h *FileHandler) renderNotFound(c echo.Context) error {
	return templates.NotFound().Render(c.Request().Context(), c.Response().Writer)
}

func (h *FileHandler) renderExpired(c echo.Context) error {
	return templates.FileExpired().Render(c.Request().Context(), c.Response().Writer)
}
