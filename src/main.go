package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	grid     [4][4]int
	score    int
	gameOver bool
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

func (m *model) SpawnNum() {
	var freeRowCol [][2]int
	for i, row := range m.grid {
		for j := range row {
			if m.grid[i][j] == 0 {
				freeRowCol = append(freeRowCol, [2]int{i, j})
			}
		}
	}

	if len(freeRowCol) == 0 {
		return
	}
	randomGrid := freeRowCol[rand.Intn(len(freeRowCol))]
	m.grid[randomGrid[0]][randomGrid[1]] = 2
}

func (m *model) MoveGrid(direction string) {
	if direction == "up" || direction == "down" {
		var checkMax int
		var inc int
		if direction == "up" {
			checkMax = 0
			inc = -1
		} else if direction == "down" {
			checkMax = 3
			inc = 1
		}

		for i, row := range m.grid {
			for j := range row {
				counter := i
				for counter != checkMax {
					if m.grid[counter][j] == m.grid[counter+inc][j] {
						m.grid[counter+inc][j] = m.grid[counter][j] * 2
						m.score += m.grid[i][counter+inc]
						m.grid[counter][j] = 0
					} else if m.grid[counter+inc][j] == 0 {
						m.grid[counter+inc][j] = m.grid[counter][j]
						m.grid[counter][j] = 0
					}

					if direction == "up" {
						counter -= 1
					} else {
						counter += 1
					}
				}
			}
		}

	} else if direction == "left" || direction == "right" {
		var checkMax int
		var inc int
		if direction == "left" {
			checkMax = 0
			inc = -1
		} else if direction == "right" {
			checkMax = 3
			inc = 1
		}

		for i, row := range m.grid {
			for j := range row {
				counter := j
				for counter != checkMax {
					if m.grid[i][counter] == m.grid[i][counter+inc] {
						m.grid[i][counter+inc] = m.grid[i][counter+inc] * 2
						m.score += m.grid[i][counter+inc]
						m.grid[i][counter] = 0
					} else if m.grid[i][counter+inc] == 0 {
						m.grid[i][counter+inc] = m.grid[i][counter]
						m.grid[i][counter] = 0
					}

					if direction == "left" {
						counter -= 1
					} else {
						counter += 1
					}
				}
			}
		}
	}
}

func (m model) IsGameOver() bool {
	var freeRowCol [][2]int
	noMove := true
	for i, row := range m.grid {
		for j := range row {
			if m.grid[i][j] == 0 {
				freeRowCol = append(freeRowCol, [2]int{i, j})
			}

			// check x and y
			for ii := 0; ii < 4; ii++ {
				if i == ii || ii > i+1 || ii < i-1 {
					continue
				}

				if m.grid[ii][j] == m.grid[i][j] {
					noMove = false
				}
			}
			for jj := 0; jj < 4; jj++ {
				if j == jj || jj > j+1 || jj < j-1 {
					continue
				}
				if m.grid[i][jj] == m.grid[i][j] {
					noMove = false
				}
			}
		}
	}

	if len(freeRowCol) == 0 && noMove {
		return true
	}

	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "w":
			m.MoveGrid("up")
			m.SpawnNum()

		case "down", "s":
			m.MoveGrid("down")
			m.SpawnNum()

		case "left", "a":
			m.MoveGrid("left")
			m.SpawnNum()

		case "right", "d":
			m.MoveGrid("right")
			m.SpawnNum()

		}
	}

	m.gameOver = m.IsGameOver()
	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Width(5)
	// The header
	s := "2048\n\n"

	if m.gameOver {
		s += "Game Over"
	}

	s += fmt.Sprintf("Score: %d\n", m.score)

	for _, row := range m.grid {
		for _, col := range row {
			s += style.Render(fmt.Sprintf("%d", col))
		}
		s += fmt.Sprint("\n")
	}

	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	rand.Seed(time.Now().Unix())
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
