package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"youtube-autoposter/internal/cli"
	"youtube-autoposter/internal/domain"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/usecase"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Server struct {
	getChannelUseCase *usecase.GetChannelInfoUseCase
	uploadUseCase     *usecase.UploadVideoUseCase
	secretFile        string
	defaultTokenFile  string
	port              int
}

func NewServer(getChannelUC *usecase.GetChannelInfoUseCase, uploadUC *usecase.UploadVideoUseCase, secretFile, tokenFile string, port int) *Server {
	if port <= 0 {
		port = 8080
	}
	return &Server{
		getChannelUseCase: getChannelUC,
		uploadUseCase:     uploadUC,
		secretFile:        secretFile,
		defaultTokenFile:  tokenFile,
		port:              port,
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, APIResponse{
		Status:  "error",
		Message: message,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{
		Status: "ok",
		Data: map[string]interface{}{
			"service":   "YouTube Auto-Poster Web UI & REST API",
			"version":   "1.1.0",
			"timestamp": time.Now().Format(time.RFC3339),
			"port":      s.port,
		},
	})
}

func (s *Server) handleProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := oauth.ListProfiles()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal mengambil daftar profile: %v", err))
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{
		Status: "success",
		Data:   profiles,
	})
}

func (s *Server) handleChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile := r.URL.Query().Get("profile")
	tokenFile := s.defaultTokenFile
	if profile != "" {
		tokenFile = oauth.GetTokenFileForProfile(profile)
	}

	channels, err := s.getChannelUseCase.ExecuteList(ctx, s.secretFile, tokenFile)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Gagal mendapatkan informasi channel (%s): %v", tokenFile, err))
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "success",
		Data:   channels,
	})
}

func (s *Server) handleVideos(w http.ResponseWriter, r *http.Request) {
	foundVideos, err := cli.ScanVideoFiles(".")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal me-scan video: %v", err))
		return
	}

	var scanned []domain.ScannedVideo
	for _, v := range foundVideos {
		scanned = append(scanned, domain.ScannedVideo{
			Path:          v.Path,
			RelPath:       v.RelPath,
			SizeBytes:     v.SizeBytes,
			SizeFormatted: cli.FormatFileSize(v.SizeBytes),
		})
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "success",
		Data:   scanned,
	})
}

type UploadPayload struct {
	FilePath      string   `json:"file_path"`
	ThumbnailPath string   `json:"thumbnail_path"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	CategoryID    string   `json:"category_id"`
	PrivacyStatus string   `json:"privacy_status"`
	PublishAt     string   `json:"publish_at"`
	Profile       string   `json:"profile"`
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method harus POST")
		return
	}

	contentType := r.Header.Get("Content-Type")

	// Direct Multipart File Upload from Browser
	if strings.HasPrefix(contentType, "multipart/form-data") {
		s.handleMultipartUpload(w, r)
		return
	}

	// JSON Payload for Server-Local Files
	var payload UploadPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON body: %v", err))
		return
	}

	if payload.FilePath == "" {
		writeError(w, http.StatusBadRequest, "file_path tidak boleh kosong")
		return
	}

	tokenFile := s.defaultTokenFile
	if payload.Profile != "" {
		tokenFile = oauth.GetTokenFileForProfile(payload.Profile)
	}

	categoryID := payload.CategoryID
	if categoryID == "" {
		categoryID = "22"
	}

	privacyStatus := payload.PrivacyStatus
	if privacyStatus == "" {
		privacyStatus = "private"
	}

	input := usecase.UploadVideoInput{
		FilePath:      payload.FilePath,
		ThumbnailPath: payload.ThumbnailPath,
		Title:         payload.Title,
		Description:   payload.Description,
		Tags:          payload.Tags,
		CategoryID:    categoryID,
		PrivacyStatus: privacyStatus,
		PublishAt:     payload.PublishAt,
		SecretFile:    s.secretFile,
		TokenFile:     tokenFile,
	}

	result, err := s.uploadUseCase.Execute(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal mengunggah video: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Video berhasil diunggah!",
		Data:    result,
	})
}

func (s *Server) handleMultipartUpload(w http.ResponseWriter, r *http.Request) {
	// Max 2GB file upload limit in memory/temp disk parse
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Gagal memproses upload form: %v", err))
		return
	}

	file, header, err := r.FormFile("video_file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Form file 'video_file' tidak ditemukan")
		return
	}
	defer file.Close()

	// Ensure temp directory
	tempDir := filepath.Join(".", "temp_uploads")
	_ = os.MkdirAll(tempDir, 0755)

	destPath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename))
	destFile, err := os.Create(destPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menyimpan video sementara: %v", err))
		return
	}
	defer func() {
		destFile.Close()
		// Clean up temporary uploaded file after request completion
		os.Remove(destPath)
	}()

	if _, err := io.Copy(destFile, file); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menyalin file video: %v", err))
		return
	}

	// Optional Thumbnail file
	var thumbDestPath string
	thumbFile, thumbHeader, thumbErr := r.FormFile("thumbnail_file")
	if thumbErr == nil && thumbHeader != nil {
		defer thumbFile.Close()
		thumbDestPath = filepath.Join(tempDir, fmt.Sprintf("thumb_%d_%s", time.Now().Unix(), thumbHeader.Filename))
		tDestFile, err := os.Create(thumbDestPath)
		if err == nil {
			io.Copy(tDestFile, thumbFile)
			tDestFile.Close()
			defer os.Remove(thumbDestPath)
		}
	}

	// Form values
	profile := r.FormValue("profile")
	tokenFile := s.defaultTokenFile
	if profile != "" {
		tokenFile = oauth.GetTokenFileForProfile(profile)
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryID := r.FormValue("category_id")
	if categoryID == "" {
		categoryID = "22"
	}
	privacyStatus := r.FormValue("privacy_status")
	if privacyStatus == "" {
		privacyStatus = "private"
	}
	publishAt := r.FormValue("publish_at")

	var tags []string
	if tagsRaw := r.FormValue("tags"); tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			tr := strings.TrimSpace(t)
			if tr != "" {
				tags = append(tags, tr)
			}
		}
	}

	input := usecase.UploadVideoInput{
		FilePath:      destPath,
		ThumbnailPath: thumbDestPath,
		Title:         title,
		Description:   description,
		Tags:          tags,
		CategoryID:    categoryID,
		PrivacyStatus: privacyStatus,
		PublishAt:     publishAt,
		SecretFile:    s.secretFile,
		TokenFile:     tokenFile,
	}

	result, err := s.uploadUseCase.Execute(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal mengunggah video: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Video file berhasil diunggah ke YouTube!",
		Data:    result,
	})
}

func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	hasSecret := oauth.HasSecretFile(s.secretFile)
	profiles, _ := oauth.ListProfiles()

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "success",
		Data: map[string]interface{}{
			"has_secret": hasSecret,
			"secret_file": s.secretFile,
			"profiles":   profiles,
		},
	})
}

type SaveSecretPayload struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	JSONContent  string `json:"json_content"`
}

func (s *Server) handleSaveSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method harus POST")
		return
	}

	var payload SaveSecretPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	if payload.JSONContent != "" {
		if err := os.WriteFile(s.secretFile, []byte(payload.JSONContent), 0600); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menyimpan client_secret.json: %v", err))
			return
		}
	} else if payload.ClientID != "" && payload.ClientSecret != "" {
		if err := oauth.SaveClientSecretInputs(payload.ClientID, payload.ClientSecret, s.secretFile); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal menyimpan client_secret.json: %v", err))
			return
		}
	} else {
		writeError(w, http.StatusBadRequest, "Client ID & Client Secret atau JSON Content harus diisi")
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Client Secret credentials berhasil disimpan!",
	})
}

func (s *Server) handleGetAuthURL(w http.ResponseWriter, r *http.Request) {
	if !oauth.HasSecretFile(s.secretFile) {
		writeError(w, http.StatusBadRequest, "File client_secret.json belum disiapkan")
		return
	}

	redirectURL := fmt.Sprintf("http://localhost:%d/callback", s.port)
	authURL, err := oauth.GenerateAuthURL(s.secretFile, redirectURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal membuat OAuth URL: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status: "success",
		Data: map[string]interface{}{
			"auth_url":     authURL,
			"redirect_url": redirectURL,
		},
	})
}

type ExchangeCodePayload struct {
	Code    string `json:"code"`
	Profile string `json:"profile"`
}

func (s *Server) handleExchangeCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method harus POST")
		return
	}

	var payload ExchangeCodePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	if payload.Code == "" {
		writeError(w, http.StatusBadRequest, "Kode otorisasi (code) tidak boleh kosong")
		return
	}

	tokenFile := s.defaultTokenFile
	if payload.Profile != "" {
		tokenFile = oauth.GetTokenFileForProfile(payload.Profile)
	}

	redirectURL := fmt.Sprintf("http://localhost:%d/callback", s.port)
	err := oauth.ExchangeCodeForToken(r.Context(), s.secretFile, tokenFile, payload.Code, redirectURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal otorisasi token: %v", err))
		return
	}

	// Fetch channel info to confirm success
	channels, channelErr := s.getChannelUseCase.ExecuteList(r.Context(), s.secretFile, tokenFile)
	if channelErr != nil {
		writeJSON(w, http.StatusOK, APIResponse{
			Status:  "success",
			Message: "Token berhasil disimpan, tetapi gagal mengambil detail channel",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Otorisasi akun YouTube berhasil!",
		Data:    channels,
	})
}
