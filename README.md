# Tetris Game (Terminal Version)

## Introduction
This is a simple Tetris game implemented in Go, designed to run in your terminal. The game features classic Tetris gameplay, colorful tetrominoes, and real-time controls. It is ideal for learning Go, terminal UI programming, or just enjoying a nostalgic game in your terminal.

## Technologies Used
- **Go (Golang):** The main programming language for the project.
- **Tcell:** A Go package for building modern, colorful terminal-based UIs.
- **Air:** A live reloading tool for Go, useful during development for instant feedback.

## How to Run

### 1. Clone the Repository
```
git clone https://github.com/harshan-dev-test/tetris-game-go-.git
cd tetris-game-go-
```

### 2. Checkout to the Master Branch
```
git checkout master
```

### 3. Install Dependencies
Make sure you have Go installed (https://golang.org/dl/).
Install dependencies using:
```
go mod tidy
```

### 4. Run the Game
You can run the game in two ways:

#### a) Developer Mode (with live reloading)
If you have [Air](https://github.com/cosmtrek/air) installed:
```
air
```

#### b) Standard Go Run
```
go run .
```

#### c) Directly run exe
```
run tetris-game.exe
```

## Game Controls
- **Q or q:** Quit the game
- **S or s:** Restart the game (when game over)
- **Up Arrow:** Rotate the tetromino
- **Down Arrow:** Move the tetromino down faster
- **Left Arrow:** Move the tetromino left
- **Right Arrow:** Move the tetromino right

Enjoy playing Tetris in your terminal! 