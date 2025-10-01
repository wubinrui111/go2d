// internal/engine/game.go
package game

import (
	"github.com/yourusername/2d-game/internal/scenes"
	"github.com/yourusername/2d-game/internal/managers"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneManager *managers.SceneManager
}

func (g *Game) Update() error {
	return g.sceneManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func Run() error {
	game := &Game{
		sceneManager: managers.NewSceneManager(),
	}

	game.sceneManager.SetScene(scenes.NewMainScene())
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("2D Game Engine")

	if err := ebiten.RunGame(game); err != nil {
		return err
	}

	return nil
}
