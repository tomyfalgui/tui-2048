package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	grid  [4][4]int
	score int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func randomRowCol(seedRow, seedCol int) (int, int) {
	seed := rand.NewSource(time.Now().UnixNano())
	rand1 := rand.New(seed)

	row := rand1.Intn(4)
	col := rand1.Intn(4)
	for row == seedRow && col == seedCol {
		row = rand1.Intn(4)
		col = rand1.Intn(4)
	}

	return row, col
}

func initialModel() model {
	row1, col1 := randomRowCol(0, 0)
	row2, col2 := randomRowCol(row1, col1)

	grid := [4][4]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	grid[row1][col1] = 2
	grid[row2][col2] = 2

	return model{
		grid:  grid,
		score: 0,
	}
}

func (m *model) MoveGrid(direction string) {
	switch direction {
	case "up":
		for i, row := range m.grid {
			for j := range row {
				if i == 0 {
					continue
				}
				counter := i
				for counter != 0 {
					if m.grid[i][j] == m.grid[i-1][j] {
						m.grid[i-1][j] = m.grid[i][j] * 2
						m.grid[i][j] = 0
					} else if m.grid[i-1][j] == 0 {

						m.grid[i-1][j] = m.grid[i][j]
						m.grid[i][j] = 0
					}
					counter -= 1
				}
			}
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "w":
			m.MoveGrid("up")

		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := "2048\n\n"

	s += fmt.Sprintf("Score: %d\n", m.score)

	for _, row := range m.grid {
		for _, col := range row {
			s += fmt.Sprintf(" %d ", col)
		}
		s += fmt.Sprint("\n")
	}

	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
