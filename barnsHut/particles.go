package barnsHut

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

/*
Represents a particle in 3D space.
Includes XYZ space coordinates and velocity
*/
type Particle struct {
	X, Y, Z    float64 // Position
	VX, VY, VZ float64 // Velocity
}

/*
Contains the left, right, down, upper, front, back space bounds for
a 3D space. Used to patition particles in the Barns Hut Oct Tree.
*/
type Octant struct {
	lb, rb, db, ub, fb, bb float64
}

/*
Calculates the new position of the particle after time step dt with the velocities
*/
func (p *Particle) NewPos(dt float64) {
	p.X += p.VX * dt
	p.Y += p.VY * dt
	p.Z += p.VZ * dt
}

/*
Prints the octant bounds
*/
func (octant *Octant) PrintOctant() {
	fmt.Printf("lb: %f, rb: %f, db: %f, ub: %f, fb: %f, bb: %f\n", octant.lb, octant.rb, octant.db, octant.ub, octant.fb, octant.bb)
}

/*
Creates Particles in 3D space with random positions and velocities. Returns particles
as list of particle objects.

nParticles: Number of Particles to generate
*/
func CreateRandParticles(nParticles int) []Particle {
	// Create an array of particles
	particles := make([]Particle, 0, nParticles)
	for i := 0; i < nParticles; i += 1 {
		newParticle := Particle{}
		// generate random float:
		// https://stackoverflow.com/questions/49746992/generate-random-float64-numbers-in-specific-range-using-golang
		// Position range -5.0 and 5.0
		newParticle.X = rand.Float64()*20 - 10.0
		newParticle.Y = rand.Float64()*20 - 10.0
		newParticle.Z = rand.Float64()*20 - 10.0
		// Initial Velocity range -0.1 and 0.1
		newParticle.VX = rand.Float64()*50 - 25.0
		newParticle.VY = rand.Float64()*50 - 25.0
		newParticle.VZ = rand.Float64()*50 - 25.0
		// Append Particle to the list
		particles = append(particles, newParticle)
	}

	return particles
}

func ReadParticles(filename string) []Particle {
	file, err := os.Open(filename)
	if err != nil {
		panic("Could not open input file of custom points")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the first line to get the number of particles
	if !scanner.Scan() {
		panic("Empty File")
	}
	nParticles, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic("Failed to read first line of file. Must be an integer representing how many particles.")
	}
	particles := make([]Particle, 0, nParticles)

	// Read particle data
	for i := 0; i < nParticles && scanner.Scan(); i++ {
		line := scanner.Text()
		values := strings.Fields(line)
		newParticle := Particle{}

		// Ensure that there are six values per line
		if len(values) != 6 {
			panic("Need 6 float values per particle for position and velocity in 3D")
		}

		// Parse and store the values in the Particle struct
		for j := 0; j < 6; j++ {
			val, err := strconv.ParseFloat(values[j], 64)
			if err != nil {
				panic("Failed to convert particle value into float")
			}
			switch j {
			case 0:
				newParticle.X = val
			case 1:
				newParticle.Y = val
			case 2:
				newParticle.Z = val
			case 3:
				newParticle.VX = val
			case 4:
				newParticle.VY = val
			case 5:
				newParticle.VZ = val
			}
		}
		particles = append(particles, newParticle)
	}

	return particles
}

/* Finds the bounds of the particle space in 3D space based on the list of particles provided. */
func CreateBounds(particles []Particle) Octant {
	lb := particles[0].X
	rb := particles[0].X
	db := particles[0].Y
	ub := particles[0].Y
	fb := particles[0].Z
	bb := particles[0].Z
	for _, particle := range particles {
		if particle.X < lb {
			lb = particle.X
		}
		if particle.X > rb {
			rb = particle.X
		}
		if particle.Y > ub {
			ub = particle.Y
		}
		if particle.Y < db {
			db = particle.Y
		}
		if particle.Z < fb {
			fb = particle.Z
		}
		if particle.Z > bb {
			bb = particle.Z
		}
	}
	return Octant{lb: lb, rb: rb, ub: ub, db: db, fb: fb, bb: bb}
}
