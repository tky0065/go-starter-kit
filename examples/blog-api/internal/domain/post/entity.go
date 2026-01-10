package post

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Post représente un article de blog
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Contenu
	Title   string `gorm:"not null;size:255" json:"title" validate:"required,max=255"`
	Slug    string `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Content string `gorm:"type:text;not null" json:"content" validate:"required"`

	// Métadonnées
	Tags      string `gorm:"size:500" json:"tags"`
	Published bool   `gorm:"default:false" json:"published"`

	// Relations
	AuthorID uint `gorm:"not null" json:"author_id"`
}

// BeforeCreate génère automatiquement un slug unique avant l'insertion
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		p.Slug = slugify(p.Title)
	}
	return nil
}

// slugify convertit un titre en slug URL-friendly
// Exemple: "Mon Super Article!" -> "mon-super-article"
func slugify(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Supprimer les caractères spéciaux
	replacer := strings.NewReplacer(
		"!", "", "?", "", ".", "", ",", "",
		"'", "", "\"", "", ":", "", ";", "",
		"(", "", ")", "", "[", "", "]", "",
	)
	slug = replacer.Replace(slug)

	// Supprimer les tirets multiples
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Supprimer les tirets en début/fin
	slug = strings.Trim(slug, "-")

	return slug
}
