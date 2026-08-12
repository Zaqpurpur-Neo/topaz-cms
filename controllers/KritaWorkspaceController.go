package controllers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"topaz-workspace/config"

	"github.com/gin-gonic/gin"
)

type KritaWorkspaceController struct {
	RoutePath string
}

type KritaArtworkModel struct {
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	RoutePath   string `json:"routePath"`
	PreviewPath string `json:"previewPath"`
	RealPath    string `json:"realPath"`
	CreatedAt   string `json:"createdAt"`
}

type KritaWorkspaceModel struct {
	Name      string              `json:"name"`
	RoutePath string              `json:"routePath"`
	RealPath  string              `json:"realPath"`
	FileCount int                 `json:"fileCount"`
	Artworks  []KritaArtworkModel `json:"artworks"`
}

type KritaWorkspaceResponse struct {
	WorkspacePath string                `json:"workspacePath"`
	Workspaces    []KritaWorkspaceModel `json:"workspaces"`
}

type zipFileReadCloser struct {
	io.ReadCloser
	fileCloser      io.Closer
	zipReaderCloser io.Closer
}

func (z *zipFileReadCloser) Close() error {
	if z.zipReaderCloser != nil {
		z.zipReaderCloser.Close()
	}
	if z.fileCloser != nil {
		z.fileCloser.Close()
	}
	return nil
}

func NewKritaWorkspaceController(routePath string, apiPath string) *KritaWorkspaceController {
	return &KritaWorkspaceController{RoutePath: apiPath + routePath}
}

func (c *KritaWorkspaceController) isWorkspaceExist(folderPath string) bool {
	_, err := os.Stat(folderPath)
	return !os.IsNotExist(err)
}

func (c *KritaWorkspaceController) getFolderAsWorkspace(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var folderList []string
	for _, entry := range entries {
		if entry.IsDir() {
			fullPath := path + "/" + entry.Name()
			folderList = append(folderList, fullPath)
		}
	}
	return folderList, nil
}

func (c *KritaWorkspaceController) getListKritaFile(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), config.KRITA_FILE_EXTENSION) {
			fileList = append(fileList, entry.Name())
		}
	}

	return fileList, nil
}

func (c *KritaWorkspaceController) makeWorkspaces(folder string) (KritaWorkspaceModel, error) {
	fileList, err := c.getListKritaFile(folder)
	if err != nil {
		return KritaWorkspaceModel{}, err
	}

	folderName := filepath.Base(folder)
	artwork := make([]KritaArtworkModel, 0, len(fileList))
	for _, file := range fileList {
		fileInfo, err := os.Stat(folder + "/" + file)
		fileName := strings.TrimSuffix(file, config.KRITA_FILE_EXTENSION)
		if err != nil {
			continue
		}
		artworkModel := KritaArtworkModel{
			Name:        fileName,
			Filename:    file,
			RealPath:    folder,
			RoutePath:   c.RoutePath + "/" + folderName + "/board/" + fileName,
			PreviewPath: c.RoutePath + "/" + folderName + "/preview/" + fileName,
			CreatedAt:   fileInfo.ModTime().Format(time.RFC3339),
		}
		artwork = append(artwork, artworkModel)
	}

	return KritaWorkspaceModel{
		Name:      folderName,
		RealPath:  folder,
		FileCount: len(fileList),
		Artworks:  artwork,
		RoutePath: c.RoutePath + "/" + folderName,
	}, nil
}

func (c *KritaWorkspaceController) GetKritaWorkspace(ctx *gin.Context) {
	var response KritaWorkspaceResponse
	response.WorkspacePath = config.KRITA_WORKSPACE_PATH

	folderList, err := c.getFolderAsWorkspace(config.KRITA_WORKSPACE_PATH)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, folder := range folderList {
		workspace, err := c.makeWorkspaces(folder)
		if err != nil {
			continue
		}
		response.Workspaces = append(response.Workspaces, workspace)
	}

	ctx.JSON(200, response)
}

func (c *KritaWorkspaceController) GetKritaWorkspaceByName(ctx *gin.Context) {
	workspace := ctx.Param("workspace")
	if workspace == "" {
		ctx.JSON(400, gin.H{"error": "workspace is required"})
		return
	}

	folderPath := config.KRITA_WORKSPACE_PATH + "/" + workspace

	if !c.isWorkspaceExist(folderPath) {
		ctx.JSON(404, gin.H{"error": "workspace not found"})
		return
	}

	workspaceRes, err := c.makeWorkspaces(folderPath)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, workspaceRes)
}
