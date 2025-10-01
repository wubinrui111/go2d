// internal/graphics/render.go
package graphics

import (
	"github.com/wubinrui111/2d-game/internal/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// RenderPlayer renders the player to the screen
func RenderPlayer(player *entities.Player, screen *ebiten.Image, sprite *ebiten.Image) {
	if sprite != nil {
		// Render with sprite image
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(player.Position.X, player.Position.Y)
		screen.DrawImage(sprite, opts)
	} else {
		// Fallback to simple rectangle rendering
		ebitenutil.DrawRect(screen, player.Position.X, player.Position.Y, 32, 32, player.GetColor())
	}
}

// RenderBlock renders a block to the screen
func RenderBlock(block *entities.SmallBlock, screen *ebiten.Image, sprite *ebiten.Image) {
	if sprite != nil {
		// Render with sprite image
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(block.Position.X, block.Position.Y)
		screen.DrawImage(sprite, opts)
	} else {
		// Fallback to simple rectangle rendering
		ebitenutil.DrawRect(screen, block.Position.X, block.Position.Y, 32, 32, block.GetColor())
	}
}