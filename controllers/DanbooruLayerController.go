package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"topaz-workspace/config"

	"github.com/gin-gonic/gin"
)

type DanbooruLayerController struct {
	Client *http.Client
}

func NewDanbooruLayerController(client *http.Client) *DanbooruLayerController {
	return &DanbooruLayerController{
		Client: client,
	}
}

func (c *DanbooruLayerController) rewriteURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	return fmt.Sprintf("http://localhost"+config.PORT+"/api/danbooru/media?url=%s", url.QueryEscape(rawURL))
}

func (c *DanbooruLayerController) proxyDanbooru(ctx *gin.Context, targetURL string) {
	req, err := http.NewRequestWithContext(ctx.Request.Context(), "GET", targetURL, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	req.Header.Set("User-Agent", "SvelteDanbooruApp/1.0 (Go-Backend)")

	resp, err := c.Client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("DoH request failed: %v", err)})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Danbooru response"})
		return
	}

	ctx.Data(resp.StatusCode, "application/json", body)
}

func (c *DanbooruLayerController) GetAutocomplete(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusOK, []interface{}{})
		return
	}

	limit := ctx.DefaultQuery("limit", "20")
	targetURL := fmt.Sprintf(
		"https://danbooru.donmai.us/autocomplete.json?search[query]=%s&search[type]=tag_query&version=3&limit=%s",
		url.QueryEscape(query),
		limit,
	)

	c.proxyDanbooru(ctx, targetURL)
}

func (c *DanbooruLayerController) GetMedia(ctx *gin.Context) {
	mediaURL := ctx.Query("url")
	if mediaURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing url parameter"})
		return
	}

	parsed, err := url.Parse(mediaURL)
	if err != nil || (!strings.HasSuffix(parsed.Host, "donmai.us")) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Only donmai.us media can be proxied"})
		return
	}

	req, err := http.NewRequestWithContext(ctx.Request.Context(), "GET", mediaURL, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("User-Agent", "SvelteDanbooruApp/1.0 (Go-Backend)")

	resp, err := c.Client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch image via DoH"})
		return
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		ctx.Header("Content-Type", contentType)
	}

	ctx.Header("Cache-Control", "public, max-age=86400")
	ctx.Status(resp.StatusCode)

	io.Copy(ctx.Writer, resp.Body)
}

func (c *DanbooruLayerController) GetPosts(ctx *gin.Context) {
	tags := ctx.DefaultQuery("tags", "")
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "20")

	targetURL := fmt.Sprintf(
		"https://danbooru.donmai.us/posts.json?tags=%s&page=%s&limit=%s",
		url.QueryEscape(tags), page, limit,
	)

	req, _ := http.NewRequestWithContext(ctx.Request.Context(), "GET", targetURL, nil)
	req.Header.Set("User-Agent", "SvelteDanbooruApp/1.0 (Go-Backend)")

	resp, err := c.Client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "DoH fetch failed"})
		return
	}
	defer resp.Body.Close()

	var posts []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Danbooru response"})
		return
	}

	for i := range posts {
		isBanned, _ := posts[i]["is_banned"].(bool)
		isDeleted, _ := posts[i]["is_deleted"].(bool)

		if isBanned || isDeleted {
			posts[i]["preview_file_url"] = nil
			posts[i]["large_file_url"] = nil
			posts[i]["file_url"] = nil
			continue
		}

		// Rewrite top-level legacy URLs
		if preview, ok := posts[i]["preview_file_url"].(string); ok && preview != "" {
			posts[i]["preview_file_url"] = c.rewriteURL(preview)
		}
		if large, ok := posts[i]["large_file_url"].(string); ok && large != "" {
			posts[i]["large_file_url"] = c.rewriteURL(large)
		}
		if file, ok := posts[i]["file_url"].(string); ok && file != "" {
			posts[i]["file_url"] = c.rewriteURL(file)
		}

		// Rewrite nested media_asset.variants URLs
		if mediaAsset, ok := posts[i]["media_asset"].(map[string]interface{}); ok {
			if variants, ok := mediaAsset["variants"].([]interface{}); ok {
				for j := range variants {
					if variant, ok := variants[j].(map[string]interface{}); ok {
						if vURL, ok := variant["url"].(string); ok && vURL != "" {
							variant["url"] = c.rewriteURL(vURL)
						}
					}
				}
			}
		}
	}

	ctx.JSON(resp.StatusCode, posts)
}
