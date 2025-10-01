# 2D Game Engine

一个使用Go编写的简单2D游戏引擎。

## 项目结构

```
.
├── cmd/                  # 主应用程序
│   └── main.go          # 入口点
├── config/              # 配置文件
│   └── config.yaml
├── internal/            # 私有应用代码
│   ├── engine/          # 核心游戏引擎
│   ├── entities/        # 游戏实体
│   ├── components/      # 实体组件
│   ├── systems/         # 游戏系统
│   ├── scenes/          # 游戏场景
│   ├── input/           # 输入管理
│   ├── graphics/        # 图形渲染
│   ├── audio/           # 音频管理
│   ├── managers/        # 各种管理器
│   └── utils/           # 工具函数
├── assets/              # 游戏资源
│   ├── images/          # 图像文件
│   ├── sounds/          # 音频文件
│   ├── fonts/          # 字体文件
│   ├── shaders/         # 着色器文件
│   └── config/          # 资源配置
├── docs/                # 文档
└── tests/               # 测试文件
```

