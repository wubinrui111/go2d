// internal/scenes/main_scene.go
package scenes

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wubinrui111/2d-game/internal/entities"
	"github.com/wubinrui111/2d-game/internal/input"
	"github.com/wubinrui111/2d-game/internal/components"
	"github.com/wubinrui111/2d-game/internal/graphics"
	graphicsSystem "github.com/wubinrui111/2d-game/internal/systems"
)

const (
	// Game constants
	GroundLevel = 550.0 // Y position of the ground surface
	GridSize    = 32.0  // Size of the grid for alignment
)

type MainScene struct {
	player    *entities.Player
	inputMgr  *input.InputManager
	blocks    []*entities.SmallBlock
	itemDrops []*entities.ItemDrop // 掉落物列表
	cameraX   float64  // 添加摄像机X坐标
	cameraY   float64  // 添加摄像机Y坐标
	selectedBlock *entities.SmallBlock // 添加选中的方块
	// 添加鼠标点击状态跟踪，避免重复处理同一点击
	leftClickProcessed   bool
	rightClickProcessed  bool
	// 添加网格显示控制字段
	showGrid  bool
	// 添加帧率计算相关字段
	fps       float64
	frameCount int
	lastFpsUpdate time.Time
	// 添加F3按键状态跟踪
	f3Pressed bool
	// 添加物品系统相关字段
	itemTypes []string     // 可用物品类型列表
	currentItemIndex int   // 当前选中的物品索引
	
	// 添加物品栏系统相关字段
	inventory       *components.Inventory
	inventorySystem *graphicsSystem.InventorySystem
	
	// 添加鼠标跟随方块相关字段
	draggedBlockType string // 被拖拽的方块类型
	draggedBlockColor color.RGBA // 被拖拽方块的颜色
	showDraggedBlock bool   // 是否显示被拖拽的方块
	
	// 添加精灵相关字段
	playerSprite *ebiten.Image
	blockSprites map[string]*ebiten.Image
}

// NewMainScene creates a new main scene
func NewMainScene() *MainScene {
	// Create the scene
	scene := &MainScene{
		player: entities.NewPlayer(320, 160), // Start player at (320, 160)
		inputMgr: &input.InputManager{},
		blocks: []*entities.SmallBlock{
			entities.NewSmallBlock(192, 192),  // 对齐到32像素网格 (32*6, 32*6)
			entities.NewSmallBlock(384, 288),  // 对齐到32像素网格 (32*12, 32*9)
			entities.NewSmallBlock(96, 384),   // 对齐到32像素网格 (32*3, 32*12)
			entities.NewSmallBlock(288, 128),  // 对齐到32像素网格 (32*9, 32*4)
			entities.NewSmallBlock(480, 224),  // 对齐到32像素网格 (32*15, 32*7)
			
			entities.NewSmallBlock(192, 576),  // 对齐到32像素网格 (32*6, 32*18)
			entities.NewSmallBlock(224, 640),  // 对齐到32像素网格 (32*7, 32*20)
			entities.NewSmallBlock(288, 672),  // 对齐到32像素网格 (32*9, 32*21)
			
			entities.NewSmallBlock(384, 96),   // 空中方块 (32*12, 32*3)
			entities.NewSmallBlock(480, 192),  // 空中方块 (32*15, 32*6)
			entities.NewSmallBlock(576, 288),  // 空中方块 (32*18, 32*9)
			
			entities.NewSmallBlock(128, 672),  // 深地下方块 (32*4, 32*21)
			entities.NewSmallBlock(288, 768),  // 深地下方块 (32*9, 32*24)
			entities.NewSmallBlock(416, 864),  // 深地下方块 (32*13, 32*27)
			
			entities.NewSmallBlock(96, 448),   // 地面附近方块 (32*3, 32*14)
			entities.NewSmallBlock(192, 480),  // 地面附近方块 (32*6, 32*15)
			entities.NewSmallBlock(320, 512),  // 地面附近方块 (32*10, 32*16)
		},
		itemDrops: []*entities.ItemDrop{}, // 初始化空的掉落物列表
		cameraX:   0,
		cameraY:      0,
		selectedBlock: nil,
		leftClickProcessed:   false,
		rightClickProcessed:  false,
		showGrid:  false,
		fps:       0,
		frameCount: 0,
		lastFpsUpdate: time.Now(),
		f3Pressed: false,
		itemTypes: []string{"SmallBlock", "RedBlock", "BlueBlock", "GreenBlock"},
		currentItemIndex: 0,
		inventory:       components.NewInventory(27, 9),
		inventorySystem: graphicsSystem.NewInventorySystem(),
		draggedBlockType: "",
		draggedBlockColor: color.RGBA{0, 0, 0, 0},
		showDraggedBlock: false,
		blockSprites: make(map[string]*ebiten.Image), // 初始化方块精灵映射
	}
	
	// 初始化摄像机位置跟随玩家
	scene.cameraX = scene.player.Position.X - 400  // 400是屏幕宽度的一半
	scene.cameraY = scene.player.Position.Y - 300  // 300是屏幕高度的一半
	
	// 尝试加载精灵表
	spriteSheet, err := graphics.NewSpriteSheet("./image/test.png", 32, 32)
	if err == nil {
		// 获取精灵映射
		spriteMap := spriteSheet.GetSpriteMap()
		
		// 为玩家获取精灵（使用索引0）
		if playerSprite, exists := spriteMap[0]; exists {
			scene.playerSprite = playerSprite
		}
		
		// 为不同类型方块建立精灵映射
		for blockType, index := range graphics.BlockSpriteMapping {
			if blockSprite, exists := spriteMap[index]; exists {
				// 为玩家和方块分别处理
				if blockType == "Player" {
					scene.playerSprite = blockSprite
				} else {
					scene.blockSprites[blockType] = blockSprite
				}
			}
		}
		
		// 将方块精灵映射传递给物品栏系统
		scene.inventorySystem.SetBlockSprites(scene.blockSprites)
	} else {
		// 如果加载失败，打印错误信息但继续运行（使用默认颜色渲染）
		fmt.Printf("Failed to load sprite sheet: %v\n", err)
		scene.playerSprite = nil
		scene.blockSprites = nil
	}
	
	// 添加一些初始物品到物品栏
	scene.initializeInventory()
	
	return scene
}

// updateFps 更新帧率计算
func (ms *MainScene) updateFps() {
	ms.frameCount++
	
	// 每秒更新一次FPS
	if time.Since(ms.lastFpsUpdate) >= time.Second {
		ms.fps = float64(ms.frameCount) / time.Since(ms.lastFpsUpdate).Seconds()
		ms.frameCount = 0
		ms.lastFpsUpdate = time.Now()
	}
}

// Update updates the scene state
func (ms *MainScene) Update() error {
	// Update inventory system (handles key presses for inventory, etc.)
	ms.inventorySystem.Update(ms.inventory)
	
	// Update player with gravity and input
	ms.inputMgr.Update(&ms.player.Velocity, &ms.player.Acceleration, ms.player.OnGround)
	
	// Apply gravity if enabled
	if ms.player.Gravity.Enabled {
		ms.player.Velocity.Y += ms.player.Gravity.Force / 60.0 // 60 FPS
	}
	
	// Apply friction/resistance based on whether player is on ground
	if ms.player.OnGround {
		ms.player.Velocity.X *= ms.player.Acceleration.GroundFriction
	} else {
		ms.player.Velocity.X *= ms.player.Acceleration.AirResistance
	}
	
	// Store previous Y velocity for fall damage calculation
	prevVelocityY := ms.player.Velocity.Y
	
	// Update player position
	ms.player.Position.X += ms.player.Velocity.X / 60.0 // 60 FPS
	ms.player.Position.Y += ms.player.Velocity.Y / 60.0 // 60 FPS
	
	// Reset onGround status
	ms.player.OnGround = false
	
	// Update player box position to match player position
	ms.player.UpdateBoxPosition()
	
	// Check ground collision (no longer used)
	// ms.checkGroundCollision()
	
	// Check block collisions
	ms.resolveCollisions()
	
	// Update FPS counter
	ms.updateFps()
	
	// Handle player fall damage
	ms.handleFallDamage(prevVelocityY)
	
	// Example: Take damage when colliding with certain blocks
	// Only take damage in survival mode
	if ms.inventorySystem.GameMode == 0 {
		for _, block := range ms.blocks {
			// Check if block is a "dangerous" block (example implementation)
			if block.Name == "LavaBlock" && ms.player.Box.Intersects(&block.Box) {
				ms.player.TakeDamage(5) // Take 5 damage per frame
			}
		}
	}
	
	// Check if player is dead
	if !ms.player.IsAlive() {
		// In creative mode, player cannot die
		if ms.inventorySystem.GameMode == 1 {
			ms.player.Heal(ms.player.Max) // Restore full health
		} else {
			// Respawn player at initial position
			ms.player.Position.X = 320
			ms.player.Position.Y = 160
			ms.player.Velocity.X = 0
			ms.player.Velocity.Y = 0
			ms.player.Heal(ms.player.Max) // Restore full health
			
			// Reset player state to prevent getting stuck
			ms.player.OnGround = false
		}
	}
	
	// 平滑跟随摄像机实现
	// 计算摄像机目标位置（玩家位置居中）
	targetX := ms.player.Position.X - 400
	targetY := ms.player.Position.Y - 300
	
	// 使用线性插值实现平滑跟随
	smoothFactor := 0.1
	ms.cameraX += (targetX - ms.cameraX) * smoothFactor
	ms.cameraY += (targetY - ms.cameraY) * smoothFactor
	
	// 处理鼠标点击事件
	ms.handleMouseInput()
	
	// 处理物品切换
	ms.handleItemSwitching()
	
	// 处理F3按键切换网格显示
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		ms.showGrid = !ms.showGrid
	}
	
	// Update item drops
	for i := len(ms.itemDrops) - 1; i >= 0; i-- {
		itemDrop := ms.itemDrops[i]
		
		// Convert blocks to BoxHolder interface
		boxHolders := make([]components.BoxHolder, len(ms.blocks))
		for j, block := range ms.blocks {
			boxHolders[j] = block
		}
		
		// Update item drop behavior with block collision
		itemDrop.Update(ms.player.Position, 1.0/60.0, boxHolders)
		
		// Check if item should disappear (lifetime exceeded)
		if itemDrop.ShouldDisappear() {
			// Remove item drop from scene
			ms.itemDrops = append(ms.itemDrops[:i], ms.itemDrops[i+1:]...)
			continue
		}
		
		// Check if item should be picked up
		if itemDrop.ShouldPickup(ms.player.Position) {
			// Create item to add to inventory
			itemToAdd := components.Item{
				ID:       itemDrop.GetItem().ID,
				Name:     itemDrop.GetItem().Name,
				Count:    1,
				MaxStack: itemDrop.GetItem().MaxStack,
				Color:    itemDrop.GetItem().Color,
			}
			
			// Add item to inventory
			ms.inventory.AddItem(itemToAdd)
			
			// Remove item drop from scene
			ms.itemDrops = append(ms.itemDrops[:i], ms.itemDrops[i+1:]...)
		}
	}
	
	return nil
}

// handleFallDamage calculates and applies fall damage to the player
func (ms *MainScene) handleFallDamage(prevVelocityY float64) {
	// 在创造模式下不受到摔落伤害
	if ms.inventorySystem.GameMode == 1 {
		return
	}
	
	// Only apply fall damage when landing (current velocity is positive (falling) 
	// and previous velocity was also positive (falling), but now on ground)
	if ms.player.OnGround && prevVelocityY > 0 && ms.player.Velocity.Y == 0 {
		// Calculate fall damage based on the speed at which the player hit the ground
		// The faster the fall, the more damage
		fallSpeed := prevVelocityY
		
		// Define thresholds for fall damage
		const (
			// Minimum speed to start taking fall damage
			minFallSpeed = 400.0
			
			// Speed at which maximum damage is dealt
			maxFallSpeed = 800.0
			
			// Maximum fall damage that can be dealt
			maxFallDamage = 20
		)
		
		// Only apply damage if falling fast enough
		if fallSpeed > minFallSpeed {
			// Calculate damage based on fall speed
			var damage int
			if fallSpeed >= maxFallSpeed {
				// Maximum damage
				damage = maxFallDamage
			} else {
				// Scale damage based on fall speed
				damage = int(maxFallDamage * (fallSpeed - minFallSpeed) / (maxFallSpeed - minFallSpeed))
			}
			
			// Apply the damage to the player
			ms.player.TakeDamage(damage)
		}
	}
}

// handleItemSwitching 处理物品切换输入
func (ms *MainScene) handleItemSwitching() {
	// 处理滚轮切换物品（原始方向）
	_, wheelY := ebiten.Wheel()
	if wheelY > 0 {
		// 向上滚动，切换到下一个物品
		ms.currentItemIndex = (ms.currentItemIndex + 1) % len(ms.itemTypes)
	} else if wheelY < 0 {
		// 向下滚动，切换到上一个物品
		ms.currentItemIndex = (ms.currentItemIndex - 1 + len(ms.itemTypes)) % len(ms.itemTypes)
	}
	
	// 处理数字键切换物品 (1-9)
	for i := 1; i <= 9 && i <= len(ms.itemTypes); i++ {
		key := ebiten.Key(int(ebiten.Key1) + i - 1)
		if ebiten.IsKeyPressed(key) {
			ms.currentItemIndex = i - 1
			break
		}
	}
}

// handleMouseInput 处理鼠标输入事件
func (ms *MainScene) handleMouseInput() {
	// 获取鼠标位置并转换为世界坐标
	mouseX, mouseY := ebiten.CursorPosition()
	worldX := float64(mouseX) + ms.cameraX
	worldY := float64(mouseY) + ms.cameraY
	
	// 更新鼠标跟随方块的位置
	ms.updateDraggedBlock(worldX, worldY)
	
	// 处理左键点击（破坏方块）
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// 连续破坏方块
		ms.removeBlockAt(worldX, worldY)
		// 注意：我们不设置leftClickProcessed为true，这样可以实现连续破坏
	} else {
		ms.leftClickProcessed = false
	}
	
	// 处理右键点击（放置方块）
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		// 连续放置方块
		ms.placeBlockAt(worldX, worldY)
		// 注意：我们不设置rightClickProcessed为true，这样可以实现连续放置
	} else {
		ms.rightClickProcessed = false
	}
	
	// 处理中键点击（拾取方块）
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		ms.pickBlockAt(worldX, worldY)
	}
}

// removeBlockAt 在指定位置移除方块
func (ms *MainScene) removeBlockAt(x, y float64) {
	// 计算方块应该所在的网格位置（强制对齐到GridSize像素网格）
	gridX := math.Floor(x/GridSize) * GridSize
	gridY := math.Floor(y/GridSize) * GridSize
	
	for i, block := range ms.blocks {
		// 检查方块是否在计算出的网格位置上
		if block.Position.X == gridX && block.Position.Y == gridY {
			// 创建掉落物
			var item *components.Item
			
			// 根据方块类型创建相应的物品
			switch block.Name {
			case "RedBlock":
				item = &components.Item{
					ID:       "red_block",
					Name:     "Red Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{200, 50, 50, 255},
				}
			case "BlueBlock":
				item = &components.Item{
					ID:       "blue_block",
					Name:     "Blue Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{50, 50, 200, 255},
				}
			case "GreenBlock":
				item = &components.Item{
					ID:       "green_block",
					Name:     "Green Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{50, 200, 50, 255},
				}
			case "SmallBlock":
				fallthrough
			default:
				// 根据颜色判断方块类型
				blockColor := block.GetColor()
				if blockColor.R == 128 && blockColor.G == 128 && blockColor.B == 128 {
					// Stone block
					item = &components.Item{
						ID:       "stone",
						Name:     "Stone",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{128, 128, 128, 255},
					}
				} else if blockColor.R == 100 && blockColor.G == 50 && blockColor.B == 0 {
					// Dirt block
					item = &components.Item{
						ID:       "dirt",
						Name:     "Dirt",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{100, 50, 0, 255},
					}
				} else if blockColor.R == 100 && blockColor.G == 70 && blockColor.B == 30 {
					// Wood block
					item = &components.Item{
						ID:       "wood",
						Name:     "Wood",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{100, 70, 30, 255},
					}
				} else {
					// Default small block
					item = &components.Item{
						ID:       "small_block",
						Name:     "Small Block",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{200, 200, 50, 255},
					}
				}
			}
			
			// 在方块位置创建掉落物，稍微偏移一点位置以避免重叠
			itemDrop := entities.NewItemDrop(gridX+8, gridY+8, item)
			ms.itemDrops = append(ms.itemDrops, itemDrop)
			
			// 从切片中移除方块
			ms.blocks = append(ms.blocks[:i], ms.blocks[i+1:]...)
			
			// 成功移除方块后立即返回
			// 注意：由于我们移除了一个元素，索引i现在指向下一个元素
			// 但由于我们直接返回，所以不需要处理这种情况
			return
		}
	}
	// 如果没有找到相交的方块，什么也不做
}

// placeBlockAt 在指定位置放置新方块
func (ms *MainScene) placeBlockAt(x, y float64) {
	// 计算方块应该放置的网格位置（强制对齐到GridSize像素网格）
	gridX := math.Floor(x/GridSize) * GridSize
	gridY := math.Floor(y/GridSize) * GridSize
	
	// 检查该位置是否已经有方块
	for _, block := range ms.blocks {
		if block.Position.X == gridX && block.Position.Y == gridY {
			// 该位置已有方块，不放置新方块
			return
		}
	}
	
	// 检查是否试图在玩家位置放置方块
	playerGridX := math.Floor(ms.player.Position.X/GridSize) * GridSize
	playerGridY := math.Floor(ms.player.Position.Y/GridSize) * GridSize
	if gridX == playerGridX && gridY == playerGridY {
		// 不允许在玩家位置放置方块
		return
	}
	
	// 检查是否有选中的物品
	selectedItem := ms.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Count <= 0 {
		// 没有选中物品或物品数量不足
		return
	}
	
	// 根据当前选中的物品类型创建相应的方块
	var newBlock *entities.SmallBlock
	switch selectedItem.Item.ID {
	case "red_block":
		newBlock = entities.NewRedBlock(gridX, gridY)
	case "blue_block":
		newBlock = entities.NewBlueBlock(gridX, gridY)
	case "green_block":
		newBlock = entities.NewGreenBlock(gridX, gridY)
	case "stone":
		newBlock = entities.NewSmallBlockWithColor(gridX, gridY, color.RGBA{128, 128, 128, 255})
		newBlock.SetName("stone")
	case "dirt":
		newBlock = entities.NewSmallBlockWithColor(gridX, gridY, color.RGBA{100, 50, 0, 255})
		newBlock.SetName("dirt")
	case "wood":
		newBlock = entities.NewSmallBlockWithColor(gridX, gridY, color.RGBA{100, 70, 30, 255})
		newBlock.SetName("wood")
	default: // "small_block" 或其他情况
		newBlock = entities.NewSmallBlock(gridX, gridY)
		newBlock.SetName("small_block")
	}
	
	// 减少物品数量（创造模式下不减少物品数量）
	if ms.inventorySystem.GameMode == 0 { // 生存模式才减少物品
		if !ms.inventory.RemoveItem(selectedItem.Item.ID, 1) {
			// 移除物品失败
			return
		}
	}
	
	ms.blocks = append(ms.blocks, newBlock)
}

// checkGroundCollision checks if the player has hit the ground
func (ms *MainScene) checkGroundCollision() {
	// 不再检查地面碰撞，允许玩家无限向下移动
	// 玩家现在可以穿过地面表面进入地下世界
	ms.player.OnGround = false
}

// resolveCollisions checks and resolves collisions between the player and blocks
func (ms *MainScene) resolveCollisions() {
	// Check collision with all blocks including ground surface blocks
	for _, block := range ms.blocks {
		if ms.player.Box.Intersects(&block.Box) {
			// Calculate intersection depth
			xDepth, yDepth := ms.player.Box.GetIntersectionDepth(&block.Box)
			
			// Determine the minimum translation vector
			if math.Abs(xDepth) < math.Abs(yDepth) {
				// Horizontal collision - push along the x-axis
				ms.player.Position.X += xDepth
				ms.player.Velocity.X = 0
			} else {
				// Vertical collision - push along the y-axis
				ms.player.Position.Y += yDepth
				
				// If moving down and hit the top of a block, set on ground
				if yDepth < 0 && ms.player.Velocity.Y > 0 {
					ms.player.Velocity.Y = 0
					ms.player.OnGround = true
				}
				// If moving up and hit the bottom of a block, stop upward movement
				if yDepth > 0 && ms.player.Velocity.Y < 0 {
					ms.player.Velocity.Y = 0
				}
			}
			
			// Update the player's collision box position
			ms.player.UpdateBoxPosition()
		}
	}
}

// Draw renders the scene
func (ms *MainScene) Draw(screen *ebiten.Image) {
	// 获取鼠标位置并应用摄像机偏移
	mouseX, mouseY := ebiten.CursorPosition()
	mouseXFloat := float64(mouseX) + ms.cameraX
	mouseYFloat := float64(mouseY) + ms.cameraY
	
	// 计算鼠标所在的网格位置
	mouseGridX := math.Floor(mouseXFloat/GridSize) * GridSize
	mouseGridY := math.Floor(mouseYFloat/GridSize) * GridSize
	
	// 绘制背景
	ebitenutil.DrawRect(screen, 0, 0, 800, 600, color.RGBA{135, 206, 235, 255}) // Sky blue background
	
	// 绘制坐标系网格（如果启用）
	if ms.showGrid {
		ms.drawCoordinateSystem(screen)
	}
	
	// Draw the player
	if ms.playerSprite != nil {
		// 使用精灵渲染玩家
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(ms.player.Position.X-ms.cameraX, ms.player.Position.Y-ms.cameraY)
		screen.DrawImage(ms.playerSprite, opts)
	} else {
		// 回退到纯色矩形渲染
		playerColor := ms.player.GetColor()
		ms.drawBoxWithBorder(screen, ms.player.Position.X, ms.player.Position.Y, ms.player.Box.Width, ms.player.Box.Height, playerColor, color.RGBA{0, 0, 0, 255})
	}

	// Draw blocks
	for _, block := range ms.blocks {
		if ms.blockSprites != nil {
			// 使用精灵渲染方块
			var blockSprite *ebiten.Image
			
			// 根据方块名称选择对应的精灵
			if sprite, exists := ms.blockSprites[block.Name]; exists {
				blockSprite = sprite
			} else {
				// 如果找不到对应名称的精灵，使用默认精灵
				blockSprite = ms.blockSprites["SmallBlock"]
			}
			
			if blockSprite != nil {
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(block.Position.X-ms.cameraX, block.Position.Y-ms.cameraY)
				screen.DrawImage(blockSprite, opts)
			} else {
				// 回退到纯色渲染
				blockColor := block.GetColor()
				
				// 检查鼠标是否悬停在方块上
				if block.IsMouseOver(mouseXFloat, mouseYFloat) {
					// 如果鼠标悬停，绘制高亮边框
					ms.drawBoxWithHighlight(screen, block.Position.X, block.Position.Y, block.Box.Width, block.Box.Height, blockColor)
				} else {
					// 否则绘制普通边框
					ms.drawBoxWithBorder(screen, block.Position.X, block.Position.Y, block.Box.Width, block.Box.Height, blockColor, color.RGBA{0, 0, 0, 255})
				}
			}
		} else {
			// 回退到纯色矩形渲染
			blockColor := block.GetColor()
			
			// 检查鼠标是否悬停在方块上
			if block.IsMouseOver(mouseXFloat, mouseYFloat) {
				// 如果鼠标悬停，绘制高亮边框
				ms.drawBoxWithHighlight(screen, block.Position.X, block.Position.Y, block.Box.Width, block.Box.Height, blockColor)
			} else {
				// 否则绘制普通边框
				ms.drawBoxWithBorder(screen, block.Position.X, block.Position.Y, block.Box.Width, block.Box.Height, blockColor, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	
	// 绘制鼠标跟随方块（如果启用）
	if ms.showDraggedBlock {
		if ms.blockSprites != nil {
			// 使用精灵渲染鼠标跟随方块
			var blockSprite *ebiten.Image
			
			// 根据方块类型选择对应的精灵
			switch ms.draggedBlockType {
			case "red_block":
				blockSprite = ms.blockSprites["red_block"]
			case "blue_block":
				blockSprite = ms.blockSprites["blue_block"]
			case "green_block":
				blockSprite = ms.blockSprites["green_block"]
			case "dirt":
				blockSprite = ms.blockSprites["dirt"]
			case "wood":
				blockSprite = ms.blockSprites["wood"]
			case "stone":
				blockSprite = ms.blockSprites["stone"]
			default:
				blockSprite = ms.blockSprites["small_block"]
			}
			
			if blockSprite != nil {
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(mouseGridX-ms.cameraX, mouseGridY-ms.cameraY)
				screen.DrawImage(blockSprite, opts)
			} else {
				// 回退到纯色矩形渲染
				ms.drawBoxWithBorder(screen, mouseGridX, mouseGridY, GridSize, GridSize, ms.draggedBlockColor, color.RGBA{255, 255, 255, 255})
			}
		} else {
			// 回退到纯色矩形渲染
			ms.drawBoxWithBorder(screen, mouseGridX, mouseGridY, GridSize, GridSize, ms.draggedBlockColor, color.RGBA{255, 255, 255, 255})
		}
	}
	
	// 绘制鼠标所在的网格位置指示器（半透明红色方框）
	ms.drawBoxWithBorder(screen, mouseGridX, mouseGridY, GridSize, GridSize, color.RGBA{0, 0, 0, 0}, color.RGBA{255, 0, 0, 100})
	
	// 绘制帧率和坐标信息
	ms.drawDebugInfo(screen, mouseXFloat, mouseYFloat)
	
	// Draw item drops
	for _, itemDrop := range ms.itemDrops {
		// Apply camera offset
		x := itemDrop.Position.X - ms.cameraX
		y := itemDrop.Position.Y - ms.cameraY
		
		// Get current size (considering shrink effect)
		width, height := itemDrop.GetCurrentSize()
		
		// Center the item based on its current size
		x += (entities.ItemDropSize - width) / 2
		y += (entities.ItemDropSize - height) / 2
		
		// Only draw if on screen
		if x >= -entities.ItemDropSize && x <= 800+entities.ItemDropSize && y >= -entities.ItemDropSize && y <= 600+entities.ItemDropSize {
			if ms.blockSprites != nil {
				// Use sprite for item drop based on item type
				var blockSprite *ebiten.Image
				
				// 根据物品ID选择对应的精灵
				switch itemDrop.GetItem().ID {
				case "red_block":
					blockSprite = ms.blockSprites["RedBlock"]
				case "blue_block":
					blockSprite = ms.blockSprites["BlueBlock"]
				case "green_block":
					blockSprite = ms.blockSprites["GreenBlock"]
				case "dirt":
					blockSprite = ms.blockSprites["dirt"]
				case "wood":
					blockSprite = ms.blockSprites["wood"]
				default:
					blockSprite = ms.blockSprites["SmallBlock"]
				}
				
				if blockSprite != nil {
					// Create a scaled version of the sprite
					opts := &ebiten.DrawImageOptions{}
					
					// Calculate scale factors
					scaleX := width / 32.0
					scaleY := height / 32.0
					
					// Apply scaling
					opts.GeoM.Scale(scaleX, scaleY)
					opts.GeoM.Translate(x, y)
					screen.DrawImage(blockSprite, opts)
				} else {
					// Fallback to colored rectangle
					ebitenutil.DrawRect(screen, x, y, width, height, itemDrop.GetItem().Color)
				}
			} else {
				// Fallback to colored rectangle
				ebitenutil.DrawRect(screen, x, y, width, height, itemDrop.GetItem().Color)
			}
			
			// Only draw border if item is not too small
			if width > 4 && height > 4 {
				// Draw a subtle border
				ebitenutil.DrawRect(screen, x, y, width, 1, color.RGBA{0, 0, 0, 255}) // Top
				ebitenutil.DrawRect(screen, x, y, 1, height, color.RGBA{0, 0, 0, 255}) // Left
				ebitenutil.DrawRect(screen, x+width-1, y, 1, height, color.RGBA{0, 0, 0, 255}) // Right
				ebitenutil.DrawRect(screen, x, y+height-1, width, 1, color.RGBA{0, 0, 0, 255}) // Bottom
			}
		}
	}
	
	// 绘制物品栏
	ms.inventorySystem.Draw(screen, ms.inventory)
	
	// Draw player health bar
	ms.drawHealthBar(screen)
	
	// Draw FPS counter
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ms.fps), 10, 10)
	
}

// drawCoordinateSystem 绘制坐标系网格
func (ms *MainScene) drawCoordinateSystem(screen *ebiten.Image) {
	// 获取屏幕边界（考虑摄像机偏移）
	screenMinX := ms.cameraX
	screenMaxX := ms.cameraX + 800  // 屏幕宽度
	screenMinY := ms.cameraY
	screenMaxY := ms.cameraY + 600  // 屏幕高度
	
	// 计算需要绘制的网格线范围
	startX := math.Floor(screenMinX/GridSize) * GridSize
	startY := math.Floor(screenMinY/GridSize) * GridSize
	endX := math.Ceil(screenMaxX/GridSize) * GridSize
	endY := math.Ceil(screenMaxY/GridSize) * GridSize
	
	// 绘制垂直网格线
	for x := startX; x <= endX; x += GridSize {
		// 应用摄像机偏移
		screenX := x - ms.cameraX
		
		// 绘制半透明的网格线
		ebitenutil.DrawRect(screen, screenX, 0, 1, 600, color.RGBA{200, 200, 200, 50})
		
		// 每128像素绘制坐标标签（避免过于密集）
		if math.Mod(x, GridSize*4) == 0 {
			// 绘制X坐标标签
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", x), int(screenX)+2, 2)
		}
	}
	
	// 绘制水平网格线
	for y := startY; y <= endY; y += GridSize {
		// 应用摄像机偏移
		screenY := y - ms.cameraY
		
		// 绘制半透明的网格线
		ebitenutil.DrawRect(screen, 0, screenY, 800, 1, color.RGBA{200, 200, 200, 50})
		
		// 每128像素绘制坐标标签（避免过于密集）
		if math.Mod(y, GridSize*4) == 0 {
			// 绘制Y坐标标签
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", y), 2, int(screenY)+2)
		}
	}
}

// drawDebugInfo 绘制调试信息（帧率和坐标）
func (ms *MainScene) drawDebugInfo(screen *ebiten.Image, mouseX, mouseY float64) {
	// 绘制玩家坐标信息
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Player: (%.0f, %.0f)", 
		ms.player.Position.X, ms.player.Position.Y), 10, 30)
	
	// 绘制鼠标坐标信息
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: (%.0f, %.0f)", 
		mouseX, mouseY), 10, 50)
	
	// 绘制网格状态信息
	gridStatus := "ON"
	if !ms.showGrid {
		gridStatus = "OFF"
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Grid: %s (Press F3 to toggle)", gridStatus), 10, 70)
	
	// 绘制当前选中的物品信息
	if len(ms.itemTypes) > 0 {
		// 显示所有物品类型，当前选中的加粗显示
		itemsText := "Items: "
		for i, itemType := range ms.itemTypes {
			if i == ms.currentItemIndex {
				itemsText += fmt.Sprintf("[%s] ", itemType) // 用方括号标记当前选中
			} else {
				itemsText += fmt.Sprintf("%s ", itemType)
			}
		}
		ebitenutil.DebugPrintAt(screen, itemsText, 10, 90)
		
		// 显示切换提示
		ebitenutil.DebugPrintAt(screen, "Use mouse wheel or 1-4 keys to switch items", 10, 110)
	}
}

// 添加绘制带高亮边框矩形的辅助方法
func (ms *MainScene) drawBoxWithHighlight(screen *ebiten.Image, x, y, width, height float64, fillColor color.Color) {
	// 应用摄像机偏移
	x -= ms.cameraX
	y -= ms.cameraY
	
	// 绘制填充矩形
	ebitenutil.DrawRect(screen, x, y, width, height, fillColor)
	
	// 绘制普通黑色边框
	// 上边框
	ebitenutil.DrawRect(screen, x, y, width, 1, color.RGBA{0, 0, 0, 255})
	// 下边框
	ebitenutil.DrawRect(screen, x, y + height - 1, width, 1, color.RGBA{0, 0, 0, 255})
	// 左边框
	ebitenutil.DrawRect(screen, x, y, 1, height, color.RGBA{0, 0, 0, 255})
	// 右边框
	ebitenutil.DrawRect(screen, x + width - 1, y, 1, height, color.RGBA{0, 0, 0, 255})
	
	// 绘制黄色高亮边框（更粗）
	highlightThickness := 2.0
	// 上边框高亮
	ebitenutil.DrawRect(screen, x-highlightThickness, y-highlightThickness, width+2*highlightThickness, highlightThickness, color.RGBA{255, 255, 0, 255})
	// 下边框高亮
	ebitenutil.DrawRect(screen, x-highlightThickness, y + height, width+2*highlightThickness, highlightThickness, color.RGBA{255, 255, 0, 255})
	// 左边框高亮
	ebitenutil.DrawRect(screen, x-highlightThickness, y, highlightThickness, height, color.RGBA{255, 255, 0, 255})
	// 右边框高亮
	ebitenutil.DrawRect(screen, x + width, y, highlightThickness, height, color.RGBA{255, 255, 0, 255})
}

// 添加绘制带边框矩形的辅助方法
func (ms *MainScene) drawBoxWithBorder(screen *ebiten.Image, x, y, width, height float64, fillColor, borderColor color.Color) {
	// 应用摄像机偏移
	x -= ms.cameraX
	y -= ms.cameraY

	// 绘制填充矩形
	ebitenutil.DrawRect(screen, x, y, width, height, fillColor)

	// 绘制黑色边框（1像素宽）
	// 上边框
	ebitenutil.DrawRect(screen, x, y, width, 1, borderColor)
	// 下边框
	ebitenutil.DrawRect(screen, x, y + height - 1, width, 1, borderColor)
	// 左边框
	ebitenutil.DrawRect(screen, x, y, 1, height, borderColor)
	// 右边框
	ebitenutil.DrawRect(screen, x + width - 1, y, 1, height, borderColor)
}

// updateDraggedBlock 更新鼠标跟随方块的位置和显示状态
func (ms *MainScene) updateDraggedBlock(mouseX, mouseY float64) {
	// 获取当前选中的物品
	selectedItem := ms.inventory.GetSelectedItem()
	
	// 如果有选中物品，显示鼠标跟随方块
	if selectedItem != nil && selectedItem.Count > 0 {
		ms.showDraggedBlock = true
		ms.draggedBlockColor = selectedItem.Item.Color
		
		// 根据物品ID设置方块类型
		ms.draggedBlockType = selectedItem.Item.ID
	} else {
		ms.showDraggedBlock = false
	}
}

// pickBlockAt 在指定位置拾取方块
func (ms *MainScene) pickBlockAt(x, y float64) {
	// 计算方块应该所在的网格位置（强制对齐到GridSize像素网格）
	gridX := math.Floor(x/GridSize) * GridSize
	gridY := math.Floor(y/GridSize) * GridSize
	
	// 查找该位置的方块
	for _, block := range ms.blocks {
		if block.Position.X == gridX && block.Position.Y == gridY {
			// 根据方块颜色确定方块类型
			blockColor := block.GetColor()
			
			// 查找匹配的物品类型
			var targetItemID string
			switch blockColor {
			case color.RGBA{200, 50, 50, 255}:
				targetItemID = "red_block"
			case color.RGBA{50, 50, 200, 255}:
				targetItemID = "blue_block"
			case color.RGBA{50, 200, 50, 255}:
				targetItemID = "green_block"
			default:
				targetItemID = "small_block"
			}
			
			// 查找匹配的物品槽位并选中它
			for i, slot := range ms.inventory.Slots {
				if slot.Item != nil && slot.Item.ID == targetItemID {
					ms.inventory.SelectSlot(i)
					return
				}
			}
			return
		}
	}
}

// breakBlock breaks a block at the given position and creates an item drop
func (ms *MainScene) breakBlock(x, y float64, blockType string) {
	// Remove the block from the blocks slice
	for i, block := range ms.blocks {
		if block.Position.X == x && block.Position.Y == y {
			// Remove the block
			ms.blocks = append(ms.blocks[:i], ms.blocks[i+1:]...)
			
			// Create an item based on the block type
			var item components.Item
			switch blockType {
			case "RedBlock":
				item = components.Item{
					ID:       "red_block",
					Name:     "Red Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{255, 0, 0, 255},
				}
			case "BlueBlock":
				item = components.Item{
					ID:       "blue_block",
					Name:     "Blue Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{0, 0, 255, 255},
				}
			case "GreenBlock":
				item = components.Item{
					ID:       "green_block",
					Name:     "Green Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{0, 255, 0, 255},
				}
			case "dirt":
				item = components.Item{
					ID:       "dirt",
					Name:     "Dirt Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{139, 69, 19, 255},
				}
			case "wood":
				item = components.Item{
					ID:       "wood",
					Name:     "Wood Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{160, 120, 40, 255},
				}
			default:
				item = components.Item{
					ID:       "small_block",
					Name:     "Small Block",
					Count:    1,
					MaxStack: 64,
					Color:    color.RGBA{128, 128, 128, 255},
				}
			}
			
			// Create an item drop at the block's position with random initial velocity
			itemDrop := entities.NewItemDrop(x, y, &item)
			
			// Add some random initial velocity
			itemDrop.Velocity.X = (0.5 - rand.Float64()) * 100 // Random X velocity
			itemDrop.Velocity.Y = -100 - rand.Float64()*50     // Upward velocity with some randomness
			
			// Add the item drop to the scene
			ms.itemDrops = append(ms.itemDrops, itemDrop)
			
			break
		}
	}
}

// initializeInventory adds some initial items to the inventory
func (ms *MainScene) initializeInventory() {
	// 添加一些示例物品到物品栏
	items := []components.Item{
		{
			ID:       "stone",
			Name:     "Stone",
			Count:    64,
			MaxStack: 64,
			Color:    color.RGBA{128, 128, 128, 255},
		},
		{
			ID:       "dirt",
			Name:     "Dirt",
			Count:    32,
			MaxStack: 64,
			Color:    color.RGBA{100, 50, 0, 255},
		},
		{
			ID:       "wood",
			Name:     "Wood",
			Count:    16,
			MaxStack: 64,
			Color:    color.RGBA{100, 70, 30, 255},
		},
		{
			ID:       "small_block",
			Name:     "Small Block",
			Count:    10,
			MaxStack: 64,
			Color:    color.RGBA{200, 200, 50, 255},
		},
		{
			ID:       "red_block",
			Name:     "Red Block",
			Count:    10,
			MaxStack: 64,
			Color:    color.RGBA{200, 50, 50, 255},
		},
		{
			ID:       "blue_block",
			Name:     "Blue Block",
			Count:    10,
			MaxStack: 64,
			Color:    color.RGBA{50, 50, 200, 255},
		},
		{
			ID:       "green_block",
			Name:     "Green Block",
			Count:    10,
			MaxStack: 64,
			Color:    color.RGBA{50, 200, 50, 255},
		},
	}
	
	// 添加物品到物品栏
	for _, item := range items {
		// 创建物品副本以避免引用问题
		itemCopy := item
		ms.inventory.AddItem(itemCopy)
	}
}

// drawHealthBar draws the player's health bar on the screen
func (ms *MainScene) drawHealthBar(screen *ebiten.Image) {
	// Health bar position and size
	barX := 10
	barY := 30
	barWidth := 200
	barHeight := 20
	
	// Calculate health percentage
	healthPercentage := ms.player.Health.GetHealthPercentage()
	
	// Draw background (red)
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(barWidth), float64(barHeight), color.RGBA{100, 0, 0, 200})
	
	// Draw health (green to yellow based on health percentage)
	healthWidth := int(float64(barWidth) * healthPercentage / 100)
	var healthColor color.RGBA
	if healthPercentage > 50 {
		// Green to yellow transition (high health)
		greenValue := uint8(255)
		redValue := uint8(255 * (100 - healthPercentage) / 50)
		healthColor = color.RGBA{redValue, greenValue, 0, 255}
	} else {
		// Yellow to red transition (low health)
		redValue := uint8(255)
		greenValue := uint8(255 * healthPercentage / 50)
		healthColor = color.RGBA{redValue, greenValue, 0, 255}
	}
	
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(healthWidth), float64(barHeight), healthColor)
	
	// Draw border
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(barWidth), 1, color.RGBA{255, 255, 255, 255}) // Top
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), 1, float64(barHeight), color.RGBA{255, 255, 255, 255}) // Left
	ebitenutil.DrawRect(screen, float64(barX+barWidth-1), float64(barY), 1, float64(barHeight), color.RGBA{255, 255, 255, 255}) // Right
	ebitenutil.DrawRect(screen, float64(barX), float64(barY+barHeight-1), float64(barWidth), 1, color.RGBA{255, 255, 255, 255}) // Bottom
	
	// Draw health text
	healthText := fmt.Sprintf("Health: %d/%d (%.0f%%)", ms.player.Health.Current, ms.player.Health.Max, healthPercentage)
	ebitenutil.DebugPrintAt(screen, healthText, barX, barY+barHeight+5)
}
