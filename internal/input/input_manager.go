// internal/input/input_manager.go
package input

import (
	"github.com/wubinrui111/2d-game/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct{}

func (im *InputManager) Update(velocity *components.Velocity, acceleration *components.Acceleration, onGround bool) {
	var speed float64
	if onGround {
		speed = acceleration.GroundSpeed
	} else {
		speed = acceleration.AirSpeed
	}

	// 根据按键设置水平速度 (支持方向键和WASD)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		velocity.X -= speed / 60.0 // Apply acceleration over time
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		velocity.X += speed / 60.0 // Apply acceleration over time
	}

	// 处理跳跃 (支持空格键和W键)
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if onGround {
			velocity.Y = -acceleration.JumpForce
		}
	}

	// 支持S键向下移动（在某些游戏中可能有用）
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		downSpeed := speed
		if onGround {
			downSpeed *= 0.5 // 向下移动速度较慢
		}
		velocity.Y += downSpeed / 60.0 // Apply downward acceleration over time
	}
}