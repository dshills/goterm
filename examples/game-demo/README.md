# Terminal Game Development Demo

A complete dungeon crawler game demonstration showcasing everything a terminal game developer needs when building games with goterm.

## Overview

This demo simulates a **roguelike dungeon crawler** with auto-play AI, demonstrating all the essential components of terminal game development.

## Running the Demo

```bash
go run ./examples/game-demo
```

Or build and run:

```bash
go build -o dungeon-crawler ./examples/game-demo
./dungeon-crawler
```

## Game Features Demonstrated

### Core Game Systems

#### 1. **Game Loop** (~30 FPS)
- Fixed timestep with delta time
- State management (Menu, Playing, GameOver, Victory)
- Update/Render separation
- Frame timing control

```go
ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
for {
    // Update game state
    game.Update()
    // Render to screen
    game.Render(screen)
    // Show frame
    screen.Show()
}
```

#### 2. **Entity System**
- Player entity with health, damage, position
- Enemy entities with different types:
  - **Goblins** (g) - Fast, weak (30 HP, 10 DMG)
  - **Orcs** (o) - Medium (50 HP, 15 DMG)
  - **Trolls** (T) - Slow, strong (80 HP, 25 DMG)
- Item entities:
  - **Health Potions** (♥) - Restore 30 HP
  - **Gold** ($) - 50 points
  - **Gems** (◆) - 100 points

#### 3. **Tile-based Map System**
- 40x20 tile map
- Wall and floor tiles
- Procedural wall generation
- Collision detection
- Pathfinding-ready structure

#### 4. **AI Systems**

**Player AI** (for auto-demo):
- Seeks nearest items
- Engages enemies when close
- Pathfinding around obstacles

**Enemy AI**:
- **Chase behavior** - Pursues player when within range (8 tiles)
- **Patrol behavior** - Random movement when player far away
- **Different speeds** - Based on enemy type

#### 5. **Combat System**
- Collision-based combat
- Damage calculation
- Health tracking
- Death/defeat handling
- Score rewards

#### 6. **UI System**

**Game Border**:
- Colored border with title
- Box drawing characters
- Reusable border function

**Status Panel**:
- Health bar with gradient colors
- Score display
- Timer display
- Enemy counter
- Item counter
- FPS meter

**Message Log**:
- Scrolling message history
- Recent message highlighting
- Event notifications

#### 7. **Visual Effects**

**Color Coding**:
- Player: Blue (RGB)
- Enemies: Type-specific colors
- Items: Meaningful colors (red = health, gold = money)
- Map: Gray tones for walls/floors

**Animations**:
- Blinking effect for damaged enemies
- Victory screen particle effects
- Title screen color cycling
- Smooth movement with delta time

#### 8. **State Management**
- Menu state
- Playing state
- Game Over state
- Victory state
- Smooth state transitions

## What Game Developers Will Learn

### Map Rendering
```go
// Tile-based map rendering
for y := 0; y < mapHeight; y++ {
    for x := 0; x < mapWidth; x++ {
        screenX := gameAreaX + x
        screenY := gameAreaY + y

        switch mapTile {
        case Wall:
            screen.SetCell(screenX, screenY,
                goterm.NewCell('█', wallColor, bg, style))
        case Floor:
            screen.SetCell(screenX, screenY,
                goterm.NewCell('·', floorColor, bg, style))
        }
    }
}
```

### Entity Rendering
```go
// Render all entities in z-order
// 1. Items (bottom layer)
for _, item := range items {
    screen.SetCell(item.X, item.Y,
        goterm.NewCell(item.Ch, item.Color, bg, goterm.StyleBold))
}

// 2. Enemies
for _, enemy := range enemies {
    screen.SetCell(enemy.X, enemy.Y,
        goterm.NewCell(enemy.Ch, enemy.Color, bg, goterm.StyleBold))
}

// 3. Player (top layer)
screen.SetCell(player.X, player.Y,
    goterm.NewCell('@', playerColor, bg, goterm.StyleBold))
```

### Health Bars
```go
func DrawBar(screen *goterm.Screen, x, y, width, value, maxValue int) {
    filled := (value * width) / maxValue

    // Draw filled portion (green)
    for i := 0; i < filled; i++ {
        screen.SetCell(x+i, y, goterm.NewCell('█', green, bg, style))
    }

    // Draw empty portion (red)
    for i := filled; i < width; i++ {
        screen.SetCell(x+i, y, goterm.NewCell('░', red, bg, style))
    }

    // Overlay text
    text := fmt.Sprintf("%d/%d", value, maxValue)
    screen.DrawText(x, y, text, white, bg, goterm.StyleBold)
}
```

### Collision Detection
```go
func CheckCollision(x, y int) bool {
    // Check map boundaries
    if x < 0 || x >= mapWidth || y < 0 || y >= mapHeight {
        return true
    }

    // Check walls
    if gameMap[y][x] == Wall {
        return true
    }

    // Check other entities
    for _, entity := range entities {
        if entity.X == x && entity.Y == y {
            return true
        }
    }

    return false
}
```

### Delta Time Movement
```go
type Entity struct {
    X, Y        int
    Speed       float64  // Tiles per second
    MoveTimer   float64  // Accumulator
}

func Update(entity *Entity, deltaTime float64) {
    entity.MoveTimer += deltaTime * entity.Speed

    // Move when timer exceeds threshold
    if entity.MoveTimer >= 1.0 {
        entity.MoveTimer -= 1.0
        // Move entity
        entity.X += direction.X
        entity.Y += direction.Y
    }
}
```

### Message System
```go
type Game struct {
    Messages    []string
    MaxMessages int
}

func AddMessage(msg string) {
    g.Messages = append(g.Messages, msg)
    if len(g.Messages) > g.MaxMessages {
        g.Messages = g.Messages[1:] // Remove oldest
    }
}

func RenderMessages(screen *goterm.Screen) {
    for i, msg := range g.Messages {
        color := goterm.ColorWhite
        if i == len(g.Messages)-1 {
            color = goterm.ColorYellow // Highlight latest
        }
        screen.DrawText(x, y+i, msg, color, bg, style)
    }
}
```

### State Machine
```go
type GameState int

const (
    StateMenu GameState = iota
    StatePlaying
    StateGameOver
    StateVictory
)

func Update() {
    switch game.State {
    case StateMenu:
        UpdateMenu()
    case StatePlaying:
        UpdateGame()
    case StateGameOver:
        // No update needed
    case StateVictory:
        UpdateVictoryEffects()
    }
}
```

## Game Architecture Patterns

### Entity Component System (Simplified)
```go
type Entity struct {
    // Position
    X, Y int

    // Rendering
    Ch    rune
    Color goterm.Color

    // Gameplay
    Type      EntityType
    Health    int
    MaxHealth int
    Damage    int

    // AI/Movement
    Speed      float64
    MoveTimer  float64
    AIState    string

    // State
    IsActive bool
}
```

### Game State Pattern
Separate update/render logic per state:
- Menu: Simple timer, then transition
- Playing: Full game logic
- GameOver: Static display
- Victory: Particle effects

### Spatial Partitioning
Using a 2D grid for fast entity lookups and collision detection.

## Performance Considerations

### Frame Rate Control
```go
// Target 30 FPS for responsive gameplay
ticker := time.NewTicker(33 * time.Millisecond)

// Calculate delta time for smooth movement
deltaTime := now.Sub(lastUpdate).Seconds()
```

### Efficient Rendering
```go
// Only draw active entities
for _, enemy := range enemies {
    if !enemy.IsActive {
        continue
    }
    // Render enemy
}

// Use screen buffer to batch updates
// Call Show() once per frame, not per entity
```

### Minimizing Allocations
```go
// Reuse color constants
var (
    playerColor = goterm.ColorRGB(100, 200, 255)
    wallColor   = goterm.ColorRGB(100, 100, 100)
)

// Reuse slices, mark entities inactive instead of removing
entity.IsActive = false
```

## Game Development Checklist

Based on this demo, here's what you need for a terminal game:

- [x] **Game loop** with fixed timestep
- [x] **Delta time** for smooth movement
- [x] **State management** (menu, playing, game over)
- [x] **Entity system** for game objects
- [x] **Map/Grid system** for level layout
- [x] **Collision detection**
- [x] **AI system** for NPCs
- [x] **Combat/interaction** system
- [x] **UI panels** (health, score, messages)
- [x] **Visual feedback** (colors, styles, effects)
- [x] **Message log** for events
- [x] **Victory/defeat conditions**
- [x] **Score/progression tracking**

## Extending This Demo

### Add Keyboard Input
```go
// Replace AI-controlled player with keyboard input
// You'll need to implement keyboard event handling
// (not yet in goterm, but can use syscall or termbox)
```

### More Enemy Types
```go
// Add ranged enemies
type Enemy struct {
    // ... existing fields
    AttackRange int
    ProjectileSpeed float64
}

// Add boss enemies with multiple phases
// Add spawners that create other enemies
```

### Power-ups and Equipment
```go
type Item struct {
    // ... existing fields
    Effect     string  // "health", "damage", "speed"
    Value      int
    Duration   float64 // For temporary effects
}
```

### Advanced AI
```go
// A* pathfinding
// Flocking behavior for groups
// State machines with multiple states
// Line of sight checking
```

### Procedural Generation
```go
// BSP dungeon generation
// Cellular automata for caves
// Random room placement
// Corridor generation
```

### Save/Load System
```go
// Serialize game state to JSON
// Save player progress
// High score table
```

## Technical Details

**Frame Rate**: ~30 FPS (33ms per frame)
**Map Size**: 40x20 tiles
**Game Area**: 44x24 characters (including border)
**UI Panels**: Right side stats, bottom message log
**Entity Types**: Player (1), Enemies (5), Items (8)
**Auto-play Duration**: 45 seconds
**States**: 4 (Menu, Playing, GameOver, Victory)

## Code Structure

```
game-demo/
├── main.go           # Complete game implementation
│   ├── Entity system
│   ├── Game loop
│   ├── AI logic
│   ├── Rendering
│   ├── UI system
│   └── State management
└── README.md         # This file
```

## For Game Developers

This demo provides a **production-ready template** for terminal game development:

✅ **Copy this structure** for your own games
✅ **Extend the entity system** for your game objects
✅ **Modify the AI** for your gameplay
✅ **Customize the UI** for your needs
✅ **Add keyboard input** for player control
✅ **Expand the map system** for your levels

The code is well-commented and organized to serve as a learning resource and starting point for your terminal game projects.

## Requirements

- Terminal with color support
- Go 1.25.3 or later
- Minimum 80x30 terminal size recommended

## Demo Flow

1. **Menu Screen** (2 seconds) - Shows title and legend
2. **Gameplay** (Auto-play) - AI-controlled demonstration
3. **Victory Screen** - If all enemies defeated and items collected
4. **Game Over Screen** - If player health reaches zero

The demo runs for approximately **45 seconds** total, showcasing all game systems in action.

---

**Perfect for learning terminal game development with Go!**

*Built with goterm - Making terminal games beautiful and fun.*
