package particles

import (
	"image/color"

	"github.com/gremour/grue"
)

// ParticleData contains set of particle parameters.
type ParticleData struct {
	Pos   grue.Vec
	Size  grue.Vec
	Color color.Color
}

// Particle describes particle.
type Particle struct {
	Initial ParticleData
	Current ParticleData
	Image   string

	// Times of spawn and expiration.
	Spawned float64
	Expires float64

	// Processor updates current particle data based on
	// time passed since from Born.
	Processor func(p *Particle, time float64)
}

// Group of particles with one generator.
type Group struct {
	Particles        map[string][]Particle
	ParticlesHardCap int
	Generator        Generator
}

// Processor modifies current particle state.
type Processor interface {
	Process(p *Particle)
}

// Generator creates new particles.
type Generator interface {
	Generate(time float64, curNum int) Particle
}

// Placer generates position for new particle.
type Placer interface {
	Place(r grue.Rect, time float64) grue.Vec
}

// Process creates new particles, updates existing and removes expired.
func (g *Group) Process(time float64) {
	if g.Particles == nil {
		g.Clear()
	}
	num := 0
	// Remove expired particles.
	for k, v := range g.Particles {
		for i, p := range v {
			if p.Expires <= time {
				lv := len(v)
				if i < lv-1 {
					copy(v[i:lv-1], v[i+1:lv])
				}
				v = v[:lv-1]
			}
		}
		num += len(v)
		g.Particles[k] = v
	}
	cap := g.ParticlesHardCap
	if cap == 0 {
		cap = 128
	}
	// Generate new particles.
	for {
		p := g.Generator.Generate(time, num)
		if p.Image == "" || num >= cap {
			break
		}
		num++
		g.Particles[p.Image] = append(g.Particles[p.Image], p)
	}
	// Process existing particles.
	for _, v := range g.Particles {
		for i := range v {
			if v[i].Processor != nil {
				v[i].Processor(&v[i], time)
			}
		}
	}
}

// Draw particle on surface.
func (p Particle) Draw(s grue.Surface) {
	r := grue.Rect{Max: p.Current.Size}
	r = r.SetCenter(p.Current.Pos)
	s.DrawImageStretched(p.Image, r, p.Current.Color)
}

// Draw group of particles on surface.
func (g *Group) Draw(s grue.Surface) {
	for _, v := range g.Particles {
		for _, p := range v {
			p.Draw(s)
		}
	}
}

// Clear removes all particles. This is equivalent
// of creating new group with same generator.
func (g *Group) Clear() {
	g.Particles = make(map[string][]Particle, 32)
}
