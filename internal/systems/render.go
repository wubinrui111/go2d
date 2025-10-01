// internal/graphics/render.go
package graphics

import (
	"github.com/yourusername/2d-game/internal/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func RenderPlayer(player *entities.Player, screen *ebiten.Image) {
	// 简单绘制一个矩形代表玩家
	ebitenutil.DrawRect(screen, player.Position.X, player.Position.Y, 32, 32, nil)
}