package particles

import (
	"image/color"
	"math/rand"

	"github.com/gremour/grue"
)

// GlitterEdge generates particles that glitter
// on the edge of rect.
type GlitterEdge struct {
	Rect        grue.Rect
	Placer      Placer
	Image       string
	Color       color.Color
	InitialSize float64
	MaxSize     float64
	// Life time of one particle.
	LifeTime float64
	// Interval between spawns. Not applicable when below MinParticles.
	SpawnTempo   float64
	MinParticles int
	MaxParticles int

	// Function to multiply particle Size with, e.g. math.Sin.
	SizeFunc func(x float64) float64
	// Arg of SizeFunc will gradually change from 0 to this over LifeTime.
	SizeFuncMaxArg float64

	timeToSpawn float64
	lastTime    float64
}

// Generate ...
func (ge *GlitterEdge) Generate(time float64, curNum int) Particle {
	defer func() {
		ge.lastTime = time
	}()

	if ge.LifeTime == 0 || ge.Placer == nil || curNum >= ge.MaxParticles {
		return Particle{}
	}

	// If particles above minimum.
	if curNum > ge.MinParticles {
		// Calculate time passed since last spawn.
		dt := time - ge.lastTime
		// Reduce time to spawn
		ge.timeToSpawn -= dt
		if ge.timeToSpawn < 0 {
			// Time has come to spawn
			ge.timeToSpawn += ge.SpawnTempo
		} else {
			// Not yet.
			return Particle{}
		}
	}

	p := Particle{
		Image:   ge.Image,
		Spawned: time,
		Expires: time + ge.LifeTime,
	}
	p.Initial.Pos = ge.Placer.Place(ge.Rect, time)
	p.Initial.Color = ge.Color
	p.Current = p.Initial
	p.Processor = func(p *Particle, time float64) {
		relt := (time - p.Spawned) / ge.LifeTime
		sz := ge.MaxSize
		if ge.SizeFunc != nil {
			ka := ge.SizeFuncMaxArg
			if ka == 0 {
				ka = 1
			}
			sz *= ge.SizeFunc(ka * relt)
		} else {
			sz *= relt
		}
		p.Current.Size = grue.V(sz, sz)
	}
	return p
}

// BorderPlacer generates postions on the rect border.
type BorderPlacer struct {
}

// Place ...
func (bp BorderPlacer) Place(r grue.Rect, time float64) grue.Vec {
	var pos grue.Vec
	w := r.W()
	h := r.H()
	dist := rand.Float64() * (w*2 + h*2)
	switch {
	case dist < w:
		pos = grue.V(dist, 0)
	case dist < w*2:
		pos = grue.V(dist-w, h)
	case dist < w*2+h:
		pos = grue.V(0, dist-w*2)
	default:
		pos = grue.V(w, dist-w*2-h)
	}
	return pos.Add(r.Min)
}
