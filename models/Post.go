package models

type Post struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	SlugName      string   `json:"slug_name"`
	InstagramUrl  string   `json:"instagram_url"`
	ThumbnailPath string   `json:"thumbnail_path"`
	Contents      []string `json:"contents"`
	PostDate      string   `json:"post_date"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
