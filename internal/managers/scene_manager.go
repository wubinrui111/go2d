package managers

import (
	"github.com/yourusername/2d-game/internal/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneManager struct {
	currentScene scenes.Scene
}

func NewSceneManager() *SceneManager {
	return &SceneManager{}
}

func (sm *SceneManager) SetScene(scene scenes.Scene) {
	sm.currentScene = scene
}

func (sm *SceneManager) Update() error {
	return sm.currentScene.Update()
}

func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.currentScene.Draw(screen)
}