package game

type State [3][3]*int
type game struct {
	state State
}

var TicTacToe = game{state: [3][3]*int{{nil, nil, nil}, {nil, nil, nil}, {nil, nil, nil}}}

func (g game) GetState() [3][3]*int {
	return g.state
}

func (g *game) MakeMove(row, col int, role int) bool {
	if g.state[row][col] == nil {
		g.state[row][col] = &role
		return true
	}
	return false
}

func (g game) IsGameOver() bool {
	// Check rows
	for row := 0; row < 3; row++ {
		if g.state[row][0] != nil && g.state[row][0] == g.state[row][1] && g.state[row][0] == g.state[row][2] {
			return true
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if g.state[0][col] != nil && g.state[0][col] == g.state[1][col] && g.state[0][col] == g.state[2][col] {
			return true
		}
	}

	// Check diagonals
	if g.state[0][0] != nil && g.state[0][0] == g.state[1][1] && g.state[0][0] == g.state[2][2] {
		return true
	}
	if g.state[0][2] != nil && g.state[0][2] == g.state[1][1] && g.state[0][2] == g.state[2][0] {
		return true
	}

	// Check for a draw
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if g.state[row][col] == nil {
				return false
			}
		}
	}

	return true
}
