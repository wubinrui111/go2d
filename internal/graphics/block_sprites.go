package graphics

// BlockSpriteMapping 定义方块类型与其在精灵表中索引的映射关系
var BlockSpriteMapping = map[string]int{
	// 玩家使用索引0
	"Player": 0,
	
	// 默认方块使用索引1
	"SmallBlock": 1,
	"stone":      1,
	"small_block": 1,
	
	// 红色方块使用索引2
	"RedBlock":  2,
	"red_block": 2,
	
	// 蓝色方块使用索引3
	"BlueBlock":  3,
	"blue_block": 3,
	
	// 绿色方块使用索引4
	"GreenBlock":  4,
	"green_block": 4,
	
	// 土方块使用索引5
	"dirt": 5,
	
	// 木方块使用索引6
	"wood": 6,
}

// GetSpriteIndex 根据方块名称获取对应的精灵索引
func GetSpriteIndex(blockName string) (int, bool) {
	index, exists := BlockSpriteMapping[blockName]
	return index, exists
}

// GetAllBlockTypes 返回所有已知的方块类型
func GetAllBlockTypes() []string {
	types := make([]string, 0, len(BlockSpriteMapping))
	
	// 排除玩家类型，只返回方块类型
	for blockType := range BlockSpriteMapping {
		if blockType != "Player" {
			types = append(types, blockType)
		}
	}
	
	return types
}