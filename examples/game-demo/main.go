package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dshills/goterm"
)

// Game states
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
	StateVictory
)

// Entity types
type EntityType int

const (
	EntityPlayer EntityType = iota
	EntityEnemy
	EntityItem
	EntityWall
	EntityFloor
)

// Direction for movement
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

// Entity represents any game object
type Entity struct {
	X, Y       int
	Ch         rune
	Color      goterm.Color
	Type       EntityType
	Health     int
	MaxHealth  int
	Damage     int
	Speed      float64
	MoveTimer  float64
	IsActive   bool
	AIState    string
	TargetX    int
	TargetY    int
}

// Game holds the game state
type Game struct {
	State       GameState
	Player      *Entity
	Enemies     []*Entity
	Items       []*Entity
	Map         [][]EntityType
	MapWidth    int
	MapHeight   int
	Score       int
	Time        float64
	DeltaTime   float64
	LastUpdate  time.Time
	GameAreaX   int
	GameAreaY   int
	GameAreaW   int
	GameAreaH   int
	Messages    []string
	MaxMessages int
}

func main() {
	// Initialize the terminal
	screen, err := goterm.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close screen: %v\n", err)
		}
	}()

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create game
	game := NewGame()

	// Game loop
	ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
	defer ticker.Stop()

	gameDuration := 45 * time.Second
	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			// Check if demo should end
			if time.Since(startTime) > gameDuration {
				return
			}

			// Update delta time
			now := time.Now()
			game.DeltaTime = now.Sub(game.LastUpdate).Seconds()
			game.LastUpdate = now

			// Clear screen
			screen.Clear()

			// Update and render based on state
			switch game.State {
			case StateMenu:
				game.UpdateMenu()
				game.RenderMenu(screen)
			case StatePlaying:
				game.Update()
				game.Render(screen)
			case StateGameOver:
				game.RenderGameOver(screen)
			case StateVictory:
				game.RenderVictory(screen)
			}

			// Show the frame
			if err := screen.Show(); err != nil {
				return
			}
		}
	}
}

func NewGame() *Game {
	g := &Game{
		State:       StateMenu,
		MapWidth:    40,
		MapHeight:   20,
		MaxMessages: 5,
		Messages:    make([]string, 0),
		LastUpdate:  time.Now(),
		GameAreaX:   2,
		GameAreaY:   2,
		GameAreaW:   44,
		GameAreaH:   24,
	}

	g.InitGame()
	return g
}

func (g *Game) InitGame() {
	// Create map
	g.Map = make([][]EntityType, g.MapHeight)
	for y := 0; y < g.MapHeight; y++ {
		g.Map[y] = make([]EntityType, g.MapWidth)
		for x := 0; x < g.MapWidth; x++ {
			// Create walls around edges
			if x == 0 || x == g.MapWidth-1 || y == 0 || y == g.MapHeight-1 {
				g.Map[y][x] = EntityWall
			} else {
				g.Map[y][x] = EntityFloor
			}
		}
	}

	// Add some random walls
	for i := 0; i < 15; i++ {
		x := rand.Intn(g.MapWidth-4) + 2
		y := rand.Intn(g.MapHeight-4) + 2
		for dx := 0; dx < 3; dx++ {
			for dy := 0; dy < 3; dy++ {
				if x+dx < g.MapWidth-1 && y+dy < g.MapHeight-1 {
					g.Map[y+dy][x+dx] = EntityWall
				}
			}
		}
	}

	// Create player
	g.Player = &Entity{
		X:         5,
		Y:         5,
		Ch:        '@',
		Color:     goterm.ColorRGB(100, 200, 255),
		Type:      EntityPlayer,
		Health:    100,
		MaxHealth: 100,
		Damage:    20,
		IsActive:  true,
	}

	// Create enemies
	g.Enemies = make([]*Entity, 0)
	for i := 0; i < 5; i++ {
		enemy := g.SpawnEnemy()
		if enemy != nil {
			g.Enemies = append(g.Enemies, enemy)
		}
	}

	// Create items
	g.Items = make([]*Entity, 0)
	for i := 0; i < 8; i++ {
		item := g.SpawnItem()
		if item != nil {
			g.Items = append(g.Items, item)
		}
	}

	g.Score = 0
	g.Time = 0
	g.AddMessage("Welcome to the Dungeon!")
	g.AddMessage("Collect items and defeat enemies!")
}

func (g *Game) SpawnEnemy() *Entity {
	// Find random empty spot
	for attempts := 0; attempts < 100; attempts++ {
		x := rand.Intn(g.MapWidth-2) + 1
		y := rand.Intn(g.MapHeight-2) + 1

		if g.Map[y][x] == EntityFloor && !g.IsOccupied(x, y) {
			// Different enemy types
			enemyType := rand.Intn(3)
			var ch rune
			var color goterm.Color
			var health, damage int
			var speed float64

			switch enemyType {
			case 0: // Goblin - fast, weak
				ch = 'g'
				color = goterm.ColorGreen
				health = 30
				damage = 10
				speed = 2.0
			case 1: // Orc - medium
				ch = 'o'
				color = goterm.ColorYellow
				health = 50
				damage = 15
				speed = 1.5
			case 2: // Troll - slow, strong
				ch = 'T'
				color = goterm.ColorRed
				health = 80
				damage = 25
				speed = 1.0
			}

			return &Entity{
				X:         x,
				Y:         y,
				Ch:        ch,
				Color:     color,
				Type:      EntityEnemy,
				Health:    health,
				MaxHealth: health,
				Damage:    damage,
				Speed:     speed,
				IsActive:  true,
				AIState:   "patrol",
			}
		}
	}
	return nil
}

func (g *Game) SpawnItem() *Entity {
	for attempts := 0; attempts < 100; attempts++ {
		x := rand.Intn(g.MapWidth-2) + 1
		y := rand.Intn(g.MapHeight-2) + 1

		if g.Map[y][x] == EntityFloor && !g.IsOccupied(x, y) {
			// Different item types
			itemType := rand.Intn(3)
			var ch rune
			var color goterm.Color

			switch itemType {
			case 0: // Health potion
				ch = '♥'
				color = goterm.ColorRed
			case 1: // Gold
				ch = '$'
				color = goterm.ColorYellow
			case 2: // Gem
				ch = '◆'
				color = goterm.ColorMagenta
			}

			return &Entity{
				X:        x,
				Y:        y,
				Ch:       ch,
				Color:    color,
				Type:     EntityItem,
				IsActive: true,
			}
		}
	}
	return nil
}

func (g *Game) IsOccupied(x, y int) bool {
	if g.Player.X == x && g.Player.Y == y {
		return true
	}
	for _, enemy := range g.Enemies {
		if enemy.IsActive && enemy.X == x && enemy.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) UpdateMenu() {
	// Auto-advance from menu after 2 seconds
	g.Time += g.DeltaTime
	if g.Time > 2.0 {
		g.State = StatePlaying
		g.Time = 0
	}
}

func (g *Game) Update() {
	g.Time += g.DeltaTime

	// Simulate player movement (AI-controlled for demo)
	g.Player.MoveTimer += g.DeltaTime
	if g.Player.MoveTimer > 0.3 {
		g.Player.MoveTimer = 0
		g.MovePlayerAI()
	}

	// Update enemies
	for _, enemy := range g.Enemies {
		if !enemy.IsActive {
			continue
		}

		enemy.MoveTimer += g.DeltaTime * enemy.Speed
		if enemy.MoveTimer > 1.0 {
			enemy.MoveTimer = 0
			g.MoveEnemy(enemy)
		}
	}

	// Check victory condition
	activeEnemies := 0
	for _, enemy := range g.Enemies {
		if enemy.IsActive {
			activeEnemies++
		}
	}

	activeItems := 0
	for _, item := range g.Items {
		if item.IsActive {
			activeItems++
		}
	}

	if activeEnemies == 0 && activeItems == 0 {
		g.State = StateVictory
	}

	// Check game over
	if g.Player.Health <= 0 {
		g.State = StateGameOver
	}
}

func (g *Game) MovePlayerAI() {
	// Simple AI: move towards nearest item or enemy
	var targetX, targetY int
	var found bool

	// First, look for nearby items
	minDist := 999999.0
	for _, item := range g.Items {
		if !item.IsActive {
			continue
		}
		dist := g.Distance(g.Player.X, g.Player.Y, item.X, item.Y)
		if dist < minDist {
			minDist = dist
			targetX = item.X
			targetY = item.Y
			found = true
		}
	}

	// If no items nearby, move towards enemies
	if !found || minDist > 10 {
		minDist = 999999.0
		for _, enemy := range g.Enemies {
			if !enemy.IsActive {
				continue
			}
			dist := g.Distance(g.Player.X, g.Player.Y, enemy.X, enemy.Y)
			if dist < minDist {
				minDist = dist
				targetX = enemy.X
				targetY = enemy.Y
				found = true
			}
		}
	}

	if !found {
		return
	}

	// Move towards target
	dx := 0
	dy := 0

	if targetX > g.Player.X {
		dx = 1
	} else if targetX < g.Player.X {
		dx = -1
	}

	if targetY > g.Player.Y {
		dy = 1
	} else if targetY < g.Player.Y {
		dy = -1
	}

	// Try to move
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	if g.CanMove(newX, newY) {
		g.Player.X = newX
		g.Player.Y = newY

		// Check for item collection
		g.CheckItemCollection()

		// Check for enemy collision (attack)
		g.CheckEnemyCollision()
	} else if dx != 0 && g.CanMove(g.Player.X+dx, g.Player.Y) {
		// Try moving only horizontally
		g.Player.X += dx
		g.CheckItemCollection()
		g.CheckEnemyCollision()
	} else if dy != 0 && g.CanMove(g.Player.X, g.Player.Y+dy) {
		// Try moving only vertically
		g.Player.Y += dy
		g.CheckItemCollection()
		g.CheckEnemyCollision()
	}
}

func (g *Game) MoveEnemy(enemy *Entity) {
	// Simple AI: chase player if close, otherwise patrol
	dist := g.Distance(enemy.X, enemy.Y, g.Player.X, g.Player.Y)

	if dist < 8 {
		// Chase player
		dx := 0
		dy := 0

		if g.Player.X > enemy.X {
			dx = 1
		} else if g.Player.X < enemy.X {
			dx = -1
		}

		if g.Player.Y > enemy.Y {
			dy = 1
		} else if g.Player.Y < enemy.Y {
			dy = -1
		}

		newX := enemy.X + dx
		newY := enemy.Y + dy

		if g.CanMoveEnemy(newX, newY, enemy) {
			enemy.X = newX
			enemy.Y = newY

			// Check if reached player (attack)
			if enemy.X == g.Player.X && enemy.Y == g.Player.Y {
				g.Player.Health -= enemy.Damage
				g.AddMessage(fmt.Sprintf("%c hit you for %d damage!", enemy.Ch, enemy.Damage))
			}
		}
	} else {
		// Random patrol
		dir := rand.Intn(4)
		dx, dy := 0, 0
		switch dir {
		case 0:
			dy = -1
		case 1:
			dy = 1
		case 2:
			dx = -1
		case 3:
			dx = 1
		}

		newX := enemy.X + dx
		newY := enemy.Y + dy
		if g.CanMoveEnemy(newX, newY, enemy) {
			enemy.X = newX
			enemy.Y = newY
		}
	}
}

func (g *Game) CanMove(x, y int) bool {
	if x < 0 || x >= g.MapWidth || y < 0 || y >= g.MapHeight {
		return false
	}
	if g.Map[y][x] == EntityWall {
		return false
	}
	return true
}

func (g *Game) CanMoveEnemy(x, y int, self *Entity) bool {
	if !g.CanMove(x, y) {
		return false
	}
	// Check if another enemy is there
	for _, enemy := range g.Enemies {
		if enemy != self && enemy.IsActive && enemy.X == x && enemy.Y == y {
			return false
		}
	}
	return true
}

func (g *Game) CheckItemCollection() {
	for _, item := range g.Items {
		if !item.IsActive {
			continue
		}
		if item.X == g.Player.X && item.Y == g.Player.Y {
			item.IsActive = false

			switch item.Ch {
			case '♥':
				heal := 30
				g.Player.Health += heal
				if g.Player.Health > g.Player.MaxHealth {
					g.Player.Health = g.Player.MaxHealth
				}
				g.AddMessage(fmt.Sprintf("Healed %d HP!", heal))
				g.Score += 10
			case '$':
				g.Score += 50
				g.AddMessage("Found gold! +50")
			case '◆':
				g.Score += 100
				g.AddMessage("Found rare gem! +100")
			}
		}
	}
}

func (g *Game) CheckEnemyCollision() {
	for _, enemy := range g.Enemies {
		if !enemy.IsActive {
			continue
		}
		if enemy.X == g.Player.X && enemy.Y == g.Player.Y {
			// Attack enemy
			enemy.Health -= g.Player.Damage
			g.AddMessage(fmt.Sprintf("Hit %c for %d damage!", enemy.Ch, g.Player.Damage))

			if enemy.Health <= 0 {
				enemy.IsActive = false
				points := 30
				switch enemy.Ch {
				case 'g':
					points = 30
				case 'o':
					points = 50
				case 'T':
					points = 100
				}
				g.Score += points
				g.AddMessage(fmt.Sprintf("Defeated %c! +%d", enemy.Ch, points))
			}
		}
	}
}

func (g *Game) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return dx*dx + dy*dy // Using squared distance (faster)
}

func (g *Game) AddMessage(msg string) {
	g.Messages = append(g.Messages, msg)
	if len(g.Messages) > g.MaxMessages {
		g.Messages = g.Messages[1:]
	}
}

func (g *Game) RenderMenu(screen *goterm.Screen) {
	w, h := screen.Size()

	// Title
	title := "DUNGEON CRAWLER"
	titleX := (w - len(title)) / 2
	screen.DrawText(titleX, h/2-5, title,
		goterm.ColorRGB(255, 100, 50),
		goterm.ColorDefault(),
		goterm.StyleBold)

	// Subtitle
	subtitle := "Terminal Game Demo"
	subtitleX := (w - len(subtitle)) / 2
	screen.DrawText(subtitleX, h/2-3, subtitle,
		goterm.ColorCyan,
		goterm.ColorDefault(),
		goterm.StyleItalic)

	// Instructions
	instructions := []string{
		"Collect items (♥ $ ◆)",
		"Defeat enemies (g o T)",
		"Avoid taking damage!",
		"",
		"Auto-playing demo...",
	}

	for i, line := range instructions {
		x := (w - len(line)) / 2
		screen.DrawText(x, h/2+i, line,
			goterm.ColorWhite,
			goterm.ColorDefault(),
			goterm.StyleNone)
	}

	// Legend
	legendY := h/2 + 7
	screen.DrawText((w-30)/2, legendY, "Legend:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleBold)
	legendY++

	legends := []struct {
		ch    string
		desc  string
		color goterm.Color
	}{
		{"@", "Player", goterm.ColorRGB(100, 200, 255)},
		{"g", "Goblin (fast)", goterm.ColorGreen},
		{"o", "Orc (medium)", goterm.ColorYellow},
		{"T", "Troll (strong)", goterm.ColorRed},
		{"♥", "Health potion", goterm.ColorRed},
		{"$", "Gold", goterm.ColorYellow},
		{"◆", "Gem", goterm.ColorMagenta},
	}

	for i, leg := range legends {
		x := (w - 30) / 2
		screen.DrawText(x, legendY+i, leg.ch, leg.color, goterm.ColorDefault(), goterm.StyleBold)
		screen.DrawText(x+3, legendY+i, leg.desc, goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	}
}

func (g *Game) Render(screen *goterm.Screen) {
	// Draw game border
	g.DrawBorder(screen, g.GameAreaX-1, g.GameAreaY-1, g.GameAreaW+2, g.GameAreaH+2,
		goterm.ColorCyan, "DUNGEON LEVEL 1")

	// Draw map
	for y := 0; y < g.MapHeight; y++ {
		for x := 0; x < g.MapWidth; x++ {
			screenX := g.GameAreaX + x
			screenY := g.GameAreaY + y

			var ch rune
			var color goterm.Color

			switch g.Map[y][x] {
			case EntityWall:
				ch = '█'
				color = goterm.ColorRGB(100, 100, 100)
			case EntityFloor:
				ch = '·'
				color = goterm.ColorRGB(60, 60, 60)
			}

			screen.SetCell(screenX, screenY,
				goterm.NewCell(ch, color, goterm.ColorDefault(), goterm.StyleNone))
		}
	}

	// Draw items
	for _, item := range g.Items {
		if !item.IsActive {
			continue
		}
		screenX := g.GameAreaX + item.X
		screenY := g.GameAreaY + item.Y
		screen.SetCell(screenX, screenY,
			goterm.NewCell(item.Ch, item.Color, goterm.ColorDefault(), goterm.StyleBold))
	}

	// Draw enemies
	for _, enemy := range g.Enemies {
		if !enemy.IsActive {
			continue
		}
		screenX := g.GameAreaX + enemy.X
		screenY := g.GameAreaY + enemy.Y

		// Blink effect for damaged enemies
		style := goterm.StyleBold
		if enemy.Health < enemy.MaxHealth/2 {
			if int(g.Time*4)%2 == 0 {
				style = goterm.StyleBold | goterm.StyleReverse
			}
		}

		screen.SetCell(screenX, screenY,
			goterm.NewCell(enemy.Ch, enemy.Color, goterm.ColorDefault(), style))
	}

	// Draw player (with glow effect)
	screenX := g.GameAreaX + g.Player.X
	screenY := g.GameAreaY + g.Player.Y
	screen.SetCell(screenX, screenY,
		goterm.NewCell(g.Player.Ch, g.Player.Color, goterm.ColorDefault(), goterm.StyleBold))

	// Draw UI panels
	g.DrawUI(screen)
}

func (g *Game) DrawUI(screen *goterm.Screen) {
	w, h := screen.Size()

	// Stats panel (right side)
	statsX := g.GameAreaX + g.GameAreaW + 3
	statsY := g.GameAreaY

	g.DrawBorder(screen, statsX-1, statsY-1, 32, 12, goterm.ColorYellow, "STATS")

	// Player health bar
	screen.DrawText(statsX, statsY, "Health:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	g.DrawBar(screen, statsX, statsY+1, 28, g.Player.Health, g.Player.MaxHealth,
		goterm.ColorGreen, goterm.ColorRed)

	// Score
	screen.DrawText(statsX, statsY+3, fmt.Sprintf("Score: %d", g.Score),
		goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

	// Time
	screen.DrawText(statsX, statsY+4, fmt.Sprintf("Time: %.1fs", g.Time),
		goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	// Enemy count
	activeEnemies := 0
	for _, enemy := range g.Enemies {
		if enemy.IsActive {
			activeEnemies++
		}
	}
	screen.DrawText(statsX, statsY+6, fmt.Sprintf("Enemies: %d", activeEnemies),
		goterm.ColorRed, goterm.ColorDefault(), goterm.StyleNone)

	// Item count
	activeItems := 0
	for _, item := range g.Items {
		if item.IsActive {
			activeItems++
		}
	}
	screen.DrawText(statsX, statsY+7, fmt.Sprintf("Items: %d", activeItems),
		goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleNone)

	// FPS
	fps := 1.0 / g.DeltaTime
	if g.DeltaTime == 0 {
		fps = 0
	}
	screen.DrawText(statsX, statsY+9, fmt.Sprintf("FPS: %.0f", fps),
		goterm.ColorRGB(150, 150, 150), goterm.ColorDefault(), goterm.StyleDim)

	// Message log (bottom)
	logY := h - 8
	g.DrawBorder(screen, 1, logY-1, w-2, 7, goterm.ColorGreen, "MESSAGE LOG")

	for i, msg := range g.Messages {
		y := logY + i
		color := goterm.ColorWhite
		if i == len(g.Messages)-1 {
			color = goterm.ColorYellow // Highlight latest message
		}
		screen.DrawText(3, y, msg, color, goterm.ColorDefault(), goterm.StyleNone)
	}
}

func (g *Game) DrawBar(screen *goterm.Screen, x, y, width, value, maxValue int, fullColor, emptyColor goterm.Color) {
	if maxValue == 0 {
		maxValue = 1
	}

	filled := (value * width) / maxValue
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	// Draw filled portion
	for i := 0; i < filled; i++ {
		screen.SetCell(x+i, y, goterm.NewCell('█', fullColor, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Draw empty portion
	for i := filled; i < width; i++ {
		screen.SetCell(x+i, y, goterm.NewCell('░', emptyColor, goterm.ColorDefault(), goterm.StyleDim))
	}

	// Draw value text
	text := fmt.Sprintf("%d/%d", value, maxValue)
	textX := x + (width-len(text))/2
	screen.DrawText(textX, y, text, goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
}

func (g *Game) DrawBorder(screen *goterm.Screen, x, y, width, height int, color goterm.Color, title string) {
	// Corners
	screen.SetCell(x, y, goterm.NewCell('┌', color, goterm.ColorDefault(), goterm.StyleNone))
	screen.SetCell(x+width-1, y, goterm.NewCell('┐', color, goterm.ColorDefault(), goterm.StyleNone))
	screen.SetCell(x, y+height-1, goterm.NewCell('└', color, goterm.ColorDefault(), goterm.StyleNone))
	screen.SetCell(x+width-1, y+height-1, goterm.NewCell('┘', color, goterm.ColorDefault(), goterm.StyleNone))

	// Horizontal lines
	for i := 1; i < width-1; i++ {
		screen.SetCell(x+i, y, goterm.NewCell('─', color, goterm.ColorDefault(), goterm.StyleNone))
		screen.SetCell(x+i, y+height-1, goterm.NewCell('─', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Vertical lines
	for i := 1; i < height-1; i++ {
		screen.SetCell(x, y+i, goterm.NewCell('│', color, goterm.ColorDefault(), goterm.StyleNone))
		screen.SetCell(x+width-1, y+i, goterm.NewCell('│', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Title
	if title != "" {
		titleText := fmt.Sprintf(" %s ", title)
		titleX := x + (width-len(titleText))/2
		screen.DrawText(titleX, y, titleText, color, goterm.ColorDefault(), goterm.StyleBold)
	}
}

func (g *Game) RenderGameOver(screen *goterm.Screen) {
	w, h := screen.Size()

	// Draw red border
	g.DrawBorder(screen, 2, 2, w-4, h-4, goterm.ColorRed, "GAME OVER")

	// Title
	title := "GAME OVER"
	screen.DrawText((w-len(title))/2, h/2-3, title,
		goterm.ColorRed, goterm.ColorDefault(), goterm.StyleBold|goterm.StyleReverse)

	// Final score
	scoreText := fmt.Sprintf("Final Score: %d", g.Score)
	screen.DrawText((w-len(scoreText))/2, h/2, scoreText,
		goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleBold)

	timeText := fmt.Sprintf("Survived: %.1f seconds", g.Time)
	screen.DrawText((w-len(timeText))/2, h/2+1, timeText,
		goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)
}

func (g *Game) RenderVictory(screen *goterm.Screen) {
	w, h := screen.Size()

	// Draw gold border
	g.DrawBorder(screen, 2, 2, w-4, h-4, goterm.ColorYellow, "VICTORY")

	// Title with animation
	title := "★ VICTORY! ★"
	titleColor := goterm.ColorYellow
	if int(g.Time*2)%2 == 0 {
		titleColor = goterm.ColorRGB(255, 215, 0)
	}
	screen.DrawText((w-len(title))/2, h/2-3, title,
		titleColor, goterm.ColorDefault(), goterm.StyleBold)

	// Messages
	messages := []string{
		fmt.Sprintf("Final Score: %d", g.Score),
		fmt.Sprintf("Time: %.1f seconds", g.Time),
		"",
		"You cleared the dungeon!",
	}

	for i, msg := range messages {
		screen.DrawText((w-len(msg))/2, h/2+i, msg,
			goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	}

	// Fireworks effect (simple)
	for i := 0; i < 5; i++ {
		x := rand.Intn(w-10) + 5
		y := rand.Intn(h-10) + 5
		colors := []goterm.Color{
			goterm.ColorRed, goterm.ColorYellow, goterm.ColorGreen,
			goterm.ColorCyan, goterm.ColorMagenta,
		}
		color := colors[rand.Intn(len(colors))]
		screen.SetCell(x, y, goterm.NewCell('*', color, goterm.ColorDefault(), goterm.StyleBold))
	}
}
