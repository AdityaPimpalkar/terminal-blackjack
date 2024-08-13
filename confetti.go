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

type frameMsg time.Time

func AnimateCelebrate() tea.Cmd {
	return tea.Tick(time.Second/framesPerSecond, func(t time.Time) tea.Msg {
		return frameMsg(t)
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
