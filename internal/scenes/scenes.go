// internal/scenes/scenes.go
package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update() error
	Draw(*ebiten.Image)
}
