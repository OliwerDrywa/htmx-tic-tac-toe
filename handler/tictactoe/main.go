package tictactoe

type game struct {
	state [3][3]int
}

func New() *game {
	return &game{state: [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}}
}

func (g *game) GetState() [3][3]int {
	return g.state
}

func (g *game) MakeMove(row, col int, role int) bool {
	if g.state[row][col] == 0 {
		g.state[row][col] = role
		return true
	}
	return false
}

func (g *game) IsGameOver() bool {
	// Check rows
	for row := 0; row < 3; row++ {
		if g.state[row][0] != 0 && g.state[row][0] == g.state[row][1] && g.state[row][0] == g.state[row][2] {
			return true
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if g.state[0][col] != 0 && g.state[0][col] == g.state[1][col] && g.state[0][col] == g.state[2][col] {
			return true
		}
	}

	// Check diagonals
	if g.state[0][0] != 0 && g.state[0][0] == g.state[1][1] && g.state[0][0] == g.state[2][2] {
		return true
	}
	if g.state[0][2] != 0 && g.state[0][2] == g.state[1][1] && g.state[0][2] == g.state[2][0] {
		return true
	}

	// Check for a draw
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if g.state[row][col] == 0 {
				return false
			}
		}
	}

	return true
}
