package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type uploadHandler struct {
	uploadDir string
	log       *slog.Logger
}

func (h *uploadHandler) uploadIcon(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20) // 2 MB

	file, header, err := c.Request.FormFile("icon")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "файл не найден"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "допустимы только jpg, png, webp, gif"})
		return
	}

	// Read first 512 bytes to detect content type
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ошибка чтения файла"})
		return
	}
	ct := http.DetectContentType(buf[:n])
	if !strings.HasPrefix(ct, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "файл не является изображением"})
		return
	}

	rb := make([]byte, 12)
	rand.Read(rb)
	filename := hex.EncodeToString(rb) + ext

	dst := filepath.Join(h.uploadDir, filename)
	out, err := os.Create(dst)
	if err != nil {
		h.log.Error("create upload file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ошибка сохранения файла"})
		return
	}
	defer out.Close()

	if _, err := out.Write(buf[:n]); err != nil {
		h.log.Error("write upload header", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ошибка сохранения файла"})
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		h.log.Error("write upload body", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ошибка сохранения файла"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("/api/uploads/icons/%s", filename)})
}
