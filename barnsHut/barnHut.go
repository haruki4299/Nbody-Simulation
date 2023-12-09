package barnsHut

import "math"

const (
	softening = 1e-9
	theta     = 0.5 // Determines distance to use internal node for calculations
)

// Calculate the force exerted on the particle from the node space
func calcForcesCOM(p *Particle, dt float64, node *OctTreeNode) {
	// Calculate the net Particle Force on the i'th Particle
	dx := p.X - node.centerOfMass[0]
	dy := p.Y - node.centerOfMass[1]
	dz := p.Z - node.centerOfMass[2]
	distSqr := dx*dx + dy*dy + dz*dz + softening
	invDist := 1.0 / math.Sqrt(distSqr)
	invDist3 := invDist * invDist * invDist

	fx := dx * invDist3
	fy := dy * invDist3
	fz := dz * invDist3

	p.VX += dt * fx
	p.VY += dt * fy
	p.VZ += dt * fz
}

func CallBarnsHut(p *Particle, dt float64, tree *OctTree) {
	barnsHutForceCalc(p, dt, tree.rootPtr.children[0])
}

/*
If the particle p and the center of mass of the particles within the space represented by node
is far enough, we will use the center of mass of the particle cluster to do force calculations.
Otherwise, we must look at each of the octants within the space.
*/
func barnsHutForceCalc(p *Particle, dt float64, node *OctTreeNode) {
	if node.totalMass == 1 {
		// If there is only one particle
		calcForcesCOM(p, dt, node)
	} else {
		dx := node.centerOfMass[0] - p.X
		dy := node.centerOfMass[1] - p.Y
		dz := node.centerOfMass[2] - p.Z
		r := math.Sqrt(dx*dx + dy*dy + dz*dz)
		d := node.totalMass
		if d/r < theta {
			// The center of mass if far enough
			calcForcesCOM(p, dt, node)
		} else {
			// Not far enough. Look into each child.
			for i := 0; i < N; i += 1 {
				if node.children[i] != nil {
					barnsHutForceCalc(p, dt, node.children[i])
				}
			}
		}
	}
}
