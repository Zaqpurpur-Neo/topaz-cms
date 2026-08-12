package controllers

import (
	"archive/zip"
	"art-cms/config"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type KritaArtworkController struct {
	RoutePath string
}

// [struct(XML)]
type XMLDoc struct {
	XMLName xml.Name `xml:"DOC"`
	Image   XMLImage `xml:"IMAGE"`
}

// [struct(XML)]
type XMLImage struct {
	Name       string    `xml:"name,attr"`
	Width      int       `xml:"width,attr"`
	Height     int       `xml:"height,attr"`
	XRes       float64   `xml:"x-res,attr"`
	ColorSpace string    `xml:"colorspacename,attr"`
	Layers     XMLLayers `xml:"layers"`
}

// [struct(XML)]
type XMLLayers struct {
	LayerList []XMLLayer `xml:"layer"`
}

// [struct(XML)]
type XMLLayer struct {
	Name        string        `xml:"name,attr"`
	NodeType    string        `xml:"nodetype,attr"`
	Visible     string        `xml:"visible,attr"`
	Opacity     float64       `xml:"opacity,attr"`
	Filename    string        `xml:"filename,attr"`
	X           float64       `xml:"x,attr"`
	Y           float64       `xml:"y,attr"`
	CompositeOp string        `xml:"compositeop,attr"`
	RefImage    []XMLRefImage `xml:"referenceimage"`
	ChildLayers *XMLLayers    `xml:"layers"`
}

// [struct(XML)]
type XMLRefImage struct {
	Src       string  `xml:"src,attr"`
	Width     float64 `xml:"width,attr"`
	Height    float64 `xml:"height,attr"`
	Opacity   float64 `xml:"opacity,attr"`
	Transform string  `xml:"transform,attr"`
}

// [struct(JSON)]
type CanvasItem struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"` // "main_artwork" | "reference_image"
	URL      string  `json:"url"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Opacity  float64 `json:"opacity"`
	Filename string  `json:"filename"`
}

// [struct(JSON)]
type LayerNode struct {
	Name        string      `json:"name"`
	NodeType    string      `json:"nodetype"`
	Visible     bool        `json:"visible"`
	Opacity     float64     `json:"opacity"`
	CompositeOp string      `json:"compositeop"`
	Filename    string      `json:"filename"`
	Children    []LayerNode `json:"children"`
}

// [struct(JSON)]
type BoardResponse struct {
	ArtworkName string       `json:"artworkName"`
	Filename    string       `json:"filename"`
	Canvas      interface{}  `json:"canvas"`
	Items       []CanvasItem `json:"items"`
	Layers      []LayerNode  `json:"layers"`
}

func NewKritaArtworkController(routePath string, apiPath string) *KritaArtworkController {
	return &KritaArtworkController{
		RoutePath: apiPath + routePath,
	}
}

// private

func (c *KritaArtworkController) parseTranslate(transform string) (float64, float64) {
	if strings.Contains(transform, "translate(") {
		start := strings.Index(transform, "translate(") + 10
		end := strings.Index(transform[start:], ")") + start
		parts := strings.FieldsFunc(transform[start:end], func(r rune) bool {
			return r == ',' || r == ' '
		})
		if len(parts) >= 2 {
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			return x, y
		}
	}
	return 0, 0
}

func (c *KritaArtworkController) parseLayerTree(layers []XMLLayer) []LayerNode {
	var nodes []LayerNode
	for _, layer := range layers {
		if layer.NodeType == "referenceimages" {
			continue
		}

		node := LayerNode{
			Name:        layer.Name,
			NodeType:    layer.NodeType,
			Visible:     layer.Visible == "1",
			Opacity:     layer.Opacity / 255.0,
			CompositeOp: layer.CompositeOp,
			Filename:    layer.Filename,
		}

		if layer.ChildLayers != nil && len(layer.ChildLayers.LayerList) > 0 {
			node.Children = c.parseLayerTree(layer.ChildLayers.LayerList)
		}

		nodes = append(nodes, node)
	}
	return nodes
}

func (c *KritaArtworkController) getKritaImage(fileNameFullPath string, isPreview bool) (io.ReadCloser, int64, error) {
	zipReader, err := zip.OpenReader(fileNameFullPath)
	fileTarget := "preview.png"
	if !isPreview {
		fileTarget = "mergedimage.png"
	}
	if err != nil {
		return nil, 0, err
	}

	for _, file := range zipReader.File {
		if file.Name == fileTarget {
			rc, err := file.Open()
			if err != nil {
				zipReader.Close()
				return nil, 0, err
			}

			wrappedReader := &zipFileReadCloser{
				ReadCloser:      rc,
				fileCloser:      rc,
				zipReaderCloser: zipReader,
			}
			return wrappedReader, file.FileInfo().Size(), nil
		}
	}

	zipReader.Close()
	return nil, 0, fmt.Errorf("preview.png not found")
}

// public

func (c *KritaArtworkController) GetArtworkPreview(ctx *gin.Context) {
	workspace := ctx.Param("workspace")
	if workspace == "" {
		ctx.JSON(400, gin.H{"error": "workspace is required"})
		return
	}

	artwork := ctx.Param("artwork")
	if artwork == "" {
		ctx.JSON(400, gin.H{"error": "artwork is required"})
		return
	}

	folderPath := config.KRITA_WORKSPACE_PATH + "/" + workspace
	filePath := folderPath + "/" + artwork + config.KRITA_FILE_EXTENSION

	imgStream, size, err := c.getKritaImage(filePath, true)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer imgStream.Close()

	ctx.DataFromReader(http.StatusOK, size, "image/png", imgStream, nil)
}

func (c *KritaArtworkController) GetArtworkOriginal(ctx *gin.Context) {
	workspace := ctx.Param("workspace")
	if workspace == "" {
		ctx.JSON(400, gin.H{"error": "workspace is required"})
		return
	}

	artwork := ctx.Param("artwork")
	if artwork == "" {
		ctx.JSON(400, gin.H{"error": "artwork is required"})
		return
	}

	folderPath := config.KRITA_WORKSPACE_PATH + "/" + workspace
	filePath := folderPath + "/" + artwork + config.KRITA_FILE_EXTENSION

	imgStream, size, err := c.getKritaImage(filePath, false)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer imgStream.Close()

	ctx.DataFromReader(http.StatusOK, size, "image/png", imgStream, nil)
}

func (c *KritaArtworkController) GetArtworkReferences(ctx *gin.Context) {
	workspace := ctx.Param("workspace")
	artwork := ctx.Param("artwork")
	refName := ctx.Param("refname")
	refIndicator := "reference_images/"

	folderPath := config.KRITA_WORKSPACE_PATH + "/" + workspace
	filePath := folderPath + "/" + artwork + config.KRITA_FILE_EXTENSION

	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if strings.Contains(file.Name, refIndicator) && strings.HasSuffix(file.Name, refName) {
			rc, _ := file.Open()
			defer rc.Close()
			ctx.Header("Content-Type", "image/png")
			io.Copy(ctx.Writer, rc)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "reference not found"})
}

func (c *KritaArtworkController) GetBoard(ctx *gin.Context) {
	workspace := ctx.Param("workspace")
	artwork := ctx.Param("artwork")

	filePathFull := config.KRITA_WORKSPACE_PATH + "/" + workspace + "/" + artwork + config.KRITA_FILE_EXTENSION

	zipReader, err := zip.OpenReader(filePathFull)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer zipReader.Close()

	var dox XMLDoc

	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, "maindoc.xml") {
			rc, _ := file.Open()
			bytes, _ := io.ReadAll(rc)
			xml.Unmarshal(bytes, &dox)
			rc.Close()
			break
		}
	}

	dpi := dox.Image.XRes
	if dpi <= 0 {
		dpi = 72.0 // Fallback default
	}
	ptToPx := dpi / 72.0

	items := []CanvasItem{
		{
			ID:       "main-canvas",
			Type:     "main_artwork",
			URL:      fmt.Sprintf("%s/%s/original/%s", c.RoutePath, workspace, artwork),
			X:        0,
			Y:        0,
			Width:    float64(dox.Image.Width),
			Height:   float64(dox.Image.Height),
			Opacity:  1.0,
			Filename: "mergedimage.png",
		},
	}

	for _, layer := range dox.Image.Layers.LayerList {
		if layer.NodeType == "referenceimages" {
			for idx, ref := range layer.RefImage {
				refX, refY := c.parseTranslate(ref.Transform)
				refName := filepath.Base(ref.Src)

				items = append(items, CanvasItem{
					ID:       fmt.Sprintf("ref-%d", idx),
					Type:     "reference_image",
					URL:      fmt.Sprintf("%s/%s/references/%s/%s", c.RoutePath, workspace, artwork, refName),
					X:        refX * ptToPx,
					Y:        refY * ptToPx,
					Width:    ref.Width * ptToPx,
					Height:   ref.Height * ptToPx,
					Opacity:  ref.Opacity,
					Filename: refName,
				})
			}
		}
	}

	layerTree := c.parseLayerTree(dox.Image.Layers.LayerList)
	ctx.JSON(http.StatusOK, BoardResponse{
		ArtworkName: dox.Image.Name,
		Filename:    artwork,
		Canvas: gin.H{
			"width":      dox.Image.Width,
			"height":     dox.Image.Height,
			"colorSpace": dox.Image.ColorSpace,
			"name":       dox.Image.Name,
			"dpi":        dpi,
		},
		Items:  items,
		Layers: layerTree,
	})
}
