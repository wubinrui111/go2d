// internal/input/input_manager.go
package input

import (
	"github.com/yourusername/2d-game/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct{}

func (im *InputManager) Update(velocity *components.Velocity, onGround bool) {
	const speed = 100.0
	const jumpForce = 400.0

	// 根据按键设置水平速度 (支持方向键和WASD)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		velocity.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		velocity.X += speed
	}

	// 处理跳跃 (支持空格键和W键)
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if onGround {
			velocity.Y = -jumpForce
		}
	}

	// 支持S键向下移动（在某些游戏中可能有用）
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		velocity.Y += speed
	}
}
