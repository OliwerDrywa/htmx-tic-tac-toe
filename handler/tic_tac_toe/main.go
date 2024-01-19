package tic_tac_toe

type Game struct {
	State      [3][3]int
	NowPlaying int // 1 or 2
}

func New() *Game {
	return &Game{
		State:      [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		NowPlaying: 1,
	}
}

func (g *Game) MakeMove(row, col int, role int) (err error) {
	if g.State[row][col] != 0 {
		return err
	}
	g.NowPlaying = 3 - g.NowPlaying
	g.State[row][col] = role
	return nil
}

func (g *Game) IsGameOver() int {
	// Check rows
	for row := 0; row < 3; row++ {
		if g.State[row][0] != 0 && g.State[row][0] == g.State[row][1] && g.State[row][0] == g.State[row][2] {
			return g.State[row][0]
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if g.State[0][col] != 0 && g.State[0][col] == g.State[1][col] && g.State[0][col] == g.State[2][col] {
			return g.State[0][col]
		}
	}

	// Check diagonals
	if g.State[0][0] != 0 && g.State[0][0] == g.State[1][1] && g.State[0][0] == g.State[2][2] {
		return g.State[0][0]
	}
	if g.State[0][2] != 0 && g.State[0][2] == g.State[1][1] && g.State[0][2] == g.State[2][0] {
		return g.State[0][2]
	}

	// Check for a draw
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if g.State[row][col] == 0 {
				return 0
			}
		}
	}

	return 0
}
