package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/maaslalani/confetty/array"
	"github.com/maaslalani/confetty/simulation"
)

const (
	framesPerSecond = 30.0
	numParticles    = 75
)

var (
	colors     = []string{"#a864fd", "#29cdff", "#78ff44", "#ff718d", "#fdff6a"}
	characters = []string{"█", "▓", "▒", "░", "▄", "▀"}
)

type Celebrate struct {
	confetti *simulation.System
}

type CelebrationMsg time.Time
type WinnerMsg string

func AnimateCelebration() tea.Cmd {
	return tea.Tick(time.Second/framesPerSecond, func(t time.Time) tea.Msg {
		return CelebrationMsg(t)
	})
}

func Winner() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return WinnerMsg("")
	})
}

func InitCelebration() Celebrate {
	return Celebrate{confetti: &simulation.System{
		Particles: []*simulation.Particle{},
		Frame:     simulation.Frame{},
	}}
}

func Spawn(width, height int) []*simulation.Particle {
	particles := []*simulation.Particle{}
	for i := 0; i < numParticles; i++ {
		x := float64(width / 2)
		y := float64(0)

		p := simulation.Particle{
			Physics: harmonica.NewProjectile(
				harmonica.FPS(framesPerSecond),
				harmonica.Point{X: x + (float64(width/4) * (rand.Float64() - 0.5)), Y: y, Z: 0},
				harmonica.Vector{X: (rand.Float64() - 0.5) * 100, Y: rand.Float64() * 50, Z: 0},
				harmonica.TerminalGravity,
			),
			Char: lipgloss.NewStyle().
				Foreground(lipgloss.Color(array.Sample(colors))).
				Render(array.Sample(characters)),
		}

		particles = append(particles, &p)
	}
	return particles
}

func CelebrationView(m MainModel) string {
	background := lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		m.celebrate.confetti.Render(),
	)
	return background
}
