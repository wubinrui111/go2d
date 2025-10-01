package graphics

import (
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet represents a sprite sheet image
type SpriteSheet struct {
	image *ebiten.Image
	width int
	height int
	spriteWidth int
	spriteHeight int
}

// NewSpriteSheet loads a sprite sheet from a file
func NewSpriteSheet(filePath string, spriteWidth, spriteHeight int) (*SpriteSheet, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Convert to ebiten image
	ebitenImg := ebiten.NewImageFromImage(img)

	// Create and return the sprite sheet
	return &SpriteSheet{
		image: ebitenImg,
		width: img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		spriteWidth: spriteWidth,
		spriteHeight: spriteHeight,
	}, nil
}

// GetSprite returns a sub-image from the sprite sheet at the specified grid position
func (ss *SpriteSheet) GetSprite(x, y int) *ebiten.Image {
	// Calculate the position in pixels
	px := x * ss.spriteWidth
	py := y * ss.spriteHeight

	// Create options for drawing the sub-image
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(-px), float64(-py))

	// Create a new image for the sprite
	sprite := ebiten.NewImage(ss.spriteWidth, ss.spriteHeight)

	// Draw the sub-image onto the sprite
	sprite.DrawImage(ss.image, opts)

	return sprite
}

// GetSpriteByIndex returns a sub-image from the sprite sheet at the specified index
// Index goes from left to right, top to bottom
func (ss *SpriteSheet) GetSpriteByIndex(index int) *ebiten.Image {
	// Calculate grid position from index
	spritesPerRow := ss.width / ss.spriteWidth
	x := index % spritesPerRow
	y := index / spritesPerRow

	return ss.GetSprite(x, y)
}

// GetSpriteMap returns a map of sprites from the sprite sheet
// The map keys correspond to the indices in the sprite sheet
func (ss *SpriteSheet) GetSpriteMap() map[int]*ebiten.Image {
	spriteMap := make(map[int]*ebiten.Image)
	
	spritesPerRow := ss.width / ss.spriteWidth
	spritesPerColumn := ss.height / ss.spriteHeight
	totalSprites := spritesPerRow * spritesPerColumn
	
	for i := 0; i < totalSprites; i++ {
		spriteMap[i] = ss.GetSpriteByIndex(i)
	}
	
	return spriteMap
}