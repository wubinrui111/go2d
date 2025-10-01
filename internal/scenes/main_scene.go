// internal/scenes/main_scene.go
package scenes

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yourusername/2d-game/internal/entities"
	"github.com/yourusername/2d-game/internal/input"
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
}

func NewMainScene() *MainScene {
	// Create some small blocks for testing
	// 使用与玩家大小相仿的小方块，并确保对齐到网格
	blocks := []*entities.SmallBlock{
		entities.NewSmallBlock(192, 192),  // 对齐到32像素网格 (32*6, 32*6)
		entities.NewSmallBlock(384, 288),  // 对齐到32像素网格 (32*12, 32*9)
		entities.NewSmallBlock(96, 384),   // 对齐到32像素网格 (32*3, 32*12)
		entities.NewSmallBlock(288, 128),  // 对齐到32像素网格 (32*9, 32*4)
		entities.NewSmallBlock(480, 224),  // 对齐到32像素网格 (32*15, 32*7)
		// 添加一些地下方块
		entities.NewSmallBlock(192, 576),  // 对齐到32像素网格 (32*6, 32*18)
		entities.NewSmallBlock(224, 640),  // 对齐到32像素网格 (32*7, 32*20)
		entities.NewSmallBlock(288, 672),  // 对齐到32像素网格 (32*9, 32*21)
		// 添加更多方块（确保对齐到网格）
		entities.NewSmallBlock(384, 96),   // 空中方块 (32*12, 32*3)
		entities.NewSmallBlock(480, 192),  // 空中方块 (32*15, 32*6)
		entities.NewSmallBlock(576, 288),  // 空中方块 (32*18, 32*9)
		entities.NewSmallBlock(128, 672),  // 深地下方块 (32*4, 32*21)
		entities.NewSmallBlock(288, 768),  // 深地下方块 (32*9, 32*24)
		entities.NewSmallBlock(416, 864),  // 深地下方块 (32*13, 32*27)
		entities.NewSmallBlock(96, 448),   // 地面附近方块 (32*3, 32*14)
		entities.NewSmallBlock(192, 480),  // 地面附近方块 (32*6, 32*15)
		entities.NewSmallBlock(320, 512),  // 地面附近方块 (32*10, 32*16)
		// 注意：不要在地面表面放置方块，避免测试失败
	}
	
	// 添加地面表面方块（每隔32像素一个方块）
	for x := 0.0; x < 800; x += GridSize {
		groundBlock := entities.NewSmallBlock(x, 544) // 544 = 32*17, 确保对齐到网格
		// 设置地面方块的颜色
		groundBlock.SetColor(color.RGBA{0, 180, 0, 255}) // Darker green color for ground surface
		blocks = append(blocks, groundBlock)
	}

	scene := &MainScene{
		player:    entities.NewPlayer(400, 300),
		inputMgr:  &input.InputManager{},
		blocks:    blocks,
		cameraX:   0,  // 初始化摄像机位置
		cameraY:   0,
		showGrid:  true, // 默认显示网格
		lastFpsUpdate: time.Now(), // 初始化FPS更新时间
		f3Pressed: false, // 初始化F3按键状态
		// 初始化物品系统
		itemTypes: []string{"SmallBlock", "RedBlock", "BlueBlock", "GreenBlock"},
		currentItemIndex: 0, // 默认选择第一个物品
	}
	
	// 初始化摄像机位置跟随玩家
	scene.cameraX = scene.player.Position.X - 400  // 400是屏幕宽度的一半
	scene.cameraY = scene.player.Position.Y - 300  // 300是屏幕高度的一半
	
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

func (ms *MainScene) Update() error {
	// 更新帧率计算
	ms.updateFps()
	
	// Handle player input
	ms.inputMgr.Update(&ms.player.Velocity, ms.player.OnGround)
	
	// Handle keyboard input for toggling grid visibility
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		if !ms.f3Pressed {
			ms.showGrid = !ms.showGrid
			ms.f3Pressed = true
		}
	} else {
		ms.f3Pressed = false
	}
	
	// 处理物品切换输入
	ms.handleItemSwitching()
	
	// Apply gravity to the player
	ms.player.Velocity.Y += 500 * 1.0/60.0 // 500 pixels/sec^2 gravity

	// Move the player based on velocity
	ms.player.Position.X += ms.player.Velocity.X * 1.0/60.0
	ms.player.Position.Y += ms.player.Velocity.Y * 1.0/60.0

	// Update player's collision box
	ms.player.UpdateBoxPosition()

	// Check for ground collision
	ms.checkGroundCollision()

	// Check for collisions with blocks and resolve them
	ms.resolveCollisions()

	// Reset horizontal velocity to prevent sliding
	ms.player.Velocity.X *= 0.8

	// Limit falling speed
	if ms.player.Velocity.Y > 500 {
		ms.player.Velocity.Y = 500
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

	return nil
}

// handleItemSwitching 处理物品切换输入
func (ms *MainScene) handleItemSwitching() {
	// 处理滚轮切换物品
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
}

// removeBlockAt 在指定位置移除方块
func (ms *MainScene) removeBlockAt(x, y float64) {
	// 计算方块应该所在的网格位置（强制对齐到GridSize像素网格）
	gridX := math.Floor(x/GridSize) * GridSize
	gridY := math.Floor(y/GridSize) * GridSize
	
	for i, block := range ms.blocks {
		// 检查方块是否在计算出的网格位置上
		if block.Position.X == gridX && block.Position.Y == gridY {
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
	
	// 根据当前选中的物品类型创建相应的方块
	var newBlock *entities.SmallBlock
	switch ms.itemTypes[ms.currentItemIndex] {
	case "RedBlock":
		newBlock = entities.NewRedBlock(gridX, gridY)
	case "BlueBlock":
		newBlock = entities.NewBlueBlock(gridX, gridY)
	case "GreenBlock":
		newBlock = entities.NewGreenBlock(gridX, gridY)
	default: // "SmallBlock" 或其他情况
		newBlock = entities.NewSmallBlock(gridX, gridY)
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

func (ms *MainScene) Draw(screen *ebiten.Image) {
	// 获取鼠标位置并应用摄像机偏移
	mouseX, mouseY := ebiten.CursorPosition()
	mouseXFloat := float64(mouseX) + ms.cameraX
	mouseYFloat := float64(mouseY) + ms.cameraY
	
	// 计算鼠标所在的网格位置
	mouseGridX := math.Floor(mouseXFloat/GridSize) * GridSize
	mouseGridY := math.Floor(mouseYFloat/GridSize) * GridSize
	
	// 绘制坐标系网格（如果启用）
	if ms.showGrid {
		ms.drawCoordinateSystem(screen)
	}
	
	// Draw the player with border
	playerColor := ms.player.GetColor()
	ms.drawBoxWithBorder(screen, ms.player.Position.X, ms.player.Position.Y, ms.player.Box.Width, ms.player.Box.Height, playerColor, color.RGBA{0, 0, 0, 255})

	// Draw blocks with border
	for _, block := range ms.blocks {
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
	
	// 绘制鼠标所在的网格位置指示器（半透明红色方框）
	ms.drawBoxWithBorder(screen, mouseGridX, mouseGridY, GridSize, GridSize, color.RGBA{0, 0, 0, 0}, color.RGBA{255, 0, 0, 100})
	
	// 绘制帧率和坐标信息
	ms.drawDebugInfo(screen, mouseXFloat, mouseYFloat)
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
	// 绘制帧率信息
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ms.fps), 10, 10)
	
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
