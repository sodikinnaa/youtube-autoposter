package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"youtube-autoposter/internal/cli"
	"youtube-autoposter/internal/domain"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/usecase"
)

// JSON-RPC 2.0 & MCP Structs
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema ToolSchema  `json:"inputSchema"`
}

type ToolSchema struct {
	Type       string              `json:"type"`
	Properties map[string]PropDesc `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type PropDesc struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolCallResult struct {
	Content []TextContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// StartMCPServer runs the stdio JSON-RPC 2.0 MCP Server
func StartMCPServer(ctx context.Context, secretFile string) {
	reader := bufio.NewReader(os.Stdin)
	oauthProvider := oauth.NewGoogleOAuthProvider()
	getChannelInfoUseCase := usecase.NewGetChannelInfoUseCase(oauthProvider)
	uploadUseCase := usecase.NewUploadVideoUseCase(oauthProvider)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(nil, -32700, "Parse error")
			continue
		}

		handleRequest(ctx, req, secretFile, oauthProvider, getChannelInfoUseCase, uploadUseCase)
	}
}

func handleRequest(ctx context.Context, req JSONRPCRequest, secretFile string, oauthProvider *oauth.GoogleOAuthProvider, getChannelUseCase *usecase.GetChannelInfoUseCase, uploadUseCase *usecase.UploadVideoUseCase) {
	switch req.Method {
	case "initialize":
		sendResponse(req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]string{
				"name":    "youtube-autoposter-mcp",
				"version": "1.0.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		})

	case "notifications/initialized":
		// Notification - no response needed

	case "tools/list":
		tools := []MCPTool{
			{
				Name:        "list_profiles",
				Description: "List all saved YouTube account profiles/tokens.",
				InputSchema: ToolSchema{
					Type:       "object",
					Properties: map[string]PropDesc{},
				},
			},
			{
				Name:        "list_channels",
				Description: "List all YouTube channels connected to a profile.",
				InputSchema: ToolSchema{
					Type: "object",
					Properties: map[string]PropDesc{
						"profile": {Type: "string", Description: "Profile name (e.g. 'akun_dua' or empty for default)"},
					},
				},
			},
			{
				Name:        "list_videos",
				Description: "Scan local project directory recursively for video files (.mp4, .mkv, .mov).",
				InputSchema: ToolSchema{
					Type: "object",
					Properties: map[string]PropDesc{
						"directory": {Type: "string", Description: "Directory path to scan (default '.')"},
					},
				},
			},
			{
				Name:        "upload_video",
				Description: "Upload a video file to YouTube with title, description, tags, thumbnail, and privacy status.",
				InputSchema: ToolSchema{
					Type: "object",
					Properties: map[string]PropDesc{
						"file":        {Type: "string", Description: "Path to video file (e.g. ./video.mp4)"},
						"profile":     {Type: "string", Description: "Profile name (optional)"},
						"title":       {Type: "string", Description: "Video title"},
						"description": {Type: "string", Description: "Video description"},
						"tags":        {Type: "string", Description: "Comma-separated tags (e.g. 'coding,golang')"},
						"thumbnail":   {Type: "string", Description: "Path to custom thumbnail image"},
						"privacy":     {Type: "string", Description: "Privacy status: 'public', 'private', or 'unlisted'"},
						"publish_at":  {Type: "string", Description: "RFC3339 timestamp for scheduled publish"},
					},
					Required: []string{"file"},
				},
			},
		}
		sendResponse(req.ID, map[string]interface{}{"tools": tools})

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			sendError(req.ID, -32602, "Invalid params")
			return
		}

		res := executeTool(ctx, params.Name, params.Arguments, secretFile, getChannelUseCase, uploadUseCase)
		sendResponse(req.ID, res)

	default:
		sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func executeTool(ctx context.Context, toolName string, args map[string]interface{}, secretFile string, getChannelUseCase *usecase.GetChannelInfoUseCase, uploadUseCase *usecase.UploadVideoUseCase) ToolCallResult {
	switch toolName {
	case "list_profiles":
		profiles, err := oauth.ListProfiles()
		if err != nil {
			return formatResultError(err)
		}
		b, _ := json.MarshalIndent(profiles, "", "  ")
		return formatResultText(string(b))

	case "list_channels":
		profile, _ := args["profile"].(string)
		tokenFile := oauth.GetTokenFileForProfile(profile)
		channels, err := getChannelUseCase.ExecuteList(ctx, secretFile, tokenFile)
		if err != nil {
			return formatResultError(err)
		}
		b, _ := json.MarshalIndent(channels, "", "  ")
		return formatResultText(string(b))

	case "list_videos":
		dir, _ := args["directory"].(string)
		if dir == "" {
			dir = "."
		}
		foundVideos, err := cli.ScanVideoFiles(dir)
		if err != nil {
			return formatResultError(err)
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
		b, _ := json.MarshalIndent(scanned, "", "  ")
		return formatResultText(string(b))

	case "upload_video":
		file, _ := args["file"].(string)
		if file == "" {
			return formatResultError(fmt.Errorf("parameter 'file' wajib diisi"))
		}

		profile, _ := args["profile"].(string)
		tokenFile := oauth.GetTokenFileForProfile(profile)

		title, _ := args["title"].(string)
		description, _ := args["description"].(string)
		tagsStr, _ := args["tags"].(string)
		thumbnail, _ := args["thumbnail"].(string)
		privacy, _ := args["privacy"].(string)
		publishAt, _ := args["publish_at"].(string)

		if privacy == "" {
			privacy = "private"
		}

		var tags []string
		if tagsStr != "" {
			for _, t := range strings.Split(tagsStr, ",") {
				if trimmed := strings.TrimSpace(t); trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		input := usecase.UploadVideoInput{
			FilePath:      file,
			ThumbnailPath: thumbnail,
			Title:         title,
			Description:   description,
			Tags:          tags,
			CategoryID:    "22",
			PrivacyStatus: privacy,
			PublishAt:     publishAt,
			SecretFile:    secretFile,
			TokenFile:     tokenFile,
		}

		res, err := uploadUseCase.Execute(ctx, input)
		if err != nil {
			return formatResultError(err)
		}

		b, _ := json.MarshalIndent(res, "", "  ")
		return formatResultText(string(b))

	default:
		return formatResultError(fmt.Errorf("tool '%s' tidak ditemukan", toolName))
	}
}

func formatResultText(text string) ToolCallResult {
	return ToolCallResult{
		Content: []TextContent{
			{Type: "text", Text: text},
		},
	}
}

func formatResultError(err error) ToolCallResult {
	return ToolCallResult{
		Content: []TextContent{
			{Type: "text", Text: fmt.Sprintf("Error: %v", err)},
		},
		IsError: true,
	}
}

func sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}

func sendError(id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}
