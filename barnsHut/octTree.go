package barnsHut

import (
	"fmt"
)

const (
	D = 3 // Dimensions
	N = 8 // Number of children in octTree. 2^D = 8
)

/*
Node of the OctTree Structure

Each subtree represents an octant in the space. There can be more than one particle in this space.
*/
type OctTreeNode struct {
	particle     *Particle  // Pointer to the particle
	whichChild   int        // Which octant the particle belongs in
	centerOfMass [D]float64 // Center of mass of all particles within the octant (subtree)
	totalMass    float64    // Total mass of the subtree
	octant       Octant     // x, y, z space boundaries

	parent   *OctTreeNode    // Pointer to the parent node
	children [N]*OctTreeNode /* Children will be indexed 0-7 in the order of DLF, DLB, DRF, DRB, TLF, TLB, TRF, TRB */
}

/*
OctTree stuct. nPart represents the number of particles.

rootPtr: pointer to the root of the octTree. The rootPtr itself is an internal node that is empty with child 0
pointing to the root node of the tree.
*/
type OctTree struct {
	rootPtr *OctTreeNode
}

/*
Initialize the octTree strucutre and return a pointer to the octTree
*/
func InitOctTree(octant Octant) *OctTree {
	quadTree := &OctTree{}
	// The root pointer will be an internal node with the child 0 pointing to the root of the tree
	quadTree.rootPtr = createInternalNode(octant, nil, 0)
	// The root of the tree will start from an empty leaf node
	quadTree.rootPtr.children[0] = createLeafNode(octant, nil, quadTree.rootPtr, 0)
	return quadTree
}

/*
Add the particle to the octTree.
*/
func (tree *OctTree) AddParticle(p *Particle) {
	// The rootPtr points to the actual root
	tree.rootPtr.children[0].insertParticle(p)
}

/*
Inserts the particle into the subtree with root.
*/
func (root *OctTreeNode) insertParticle(p *Particle) {
	// Different nodes: Internal nodes that represents space with multiple particles,
	// Leaf Nodes with Particles, Empty Leaf Nodes
	if root.totalMass > 1 {
		// Is an internal node
		// Simply find the child where the particle belongs
		root.totalMass += 1
		for i := 0; i < N; i += 1 {
			if withinSpace(p, root.children[i].octant) {
				root.children[i].insertParticle(p)
				break
			}
		}
	} else if root.totalMass == 1 { // is a leaf node
		// Replace with internal node and reinsert two particles
		whichChild := root.whichChild
		octant := root.octant
		curP := root.particle
		newInternalNode := createInternalNode(octant, root.parent, whichChild)
		root.parent.children[whichChild] = newInternalNode
		newInternalNode.totalMass = 2
		for i := 0; i < N; i += 1 {
			if withinSpace(curP, newInternalNode.children[i].octant) {
				newInternalNode.children[i].insertParticle(curP)
				break
			}
		}
		for i := 0; i < N; i += 1 {
			if withinSpace(p, newInternalNode.children[i].octant) {
				newInternalNode.children[i].insertParticle(p)
				break
			}
		}
	} else { // Empty leaf node
		root.fillEmptyLeaf(p)
	}
}

/*
Fills an empty leaf node with the particle
*/
func (node *OctTreeNode) fillEmptyLeaf(p *Particle) {
	node.particle = p
	node.centerOfMass = [D]float64{p.X, p.Y, p.Z}
	node.totalMass = 1
}

/*
Print the tree with leaf nodes and internal nodes
*/
func (tree *OctTree) PrintTree() {
	fmt.Println("Printing Tree")
	tree.rootPtr.children[0].printNode()
}

/*
Helper function for PrintTree()
*/
func (node *OctTreeNode) printNode() {
	if node != nil {
		if node.totalMass > 1 {
			fmt.Printf("Internal Node total Mass: %f COM: %f, %f, %f\n", node.totalMass, node.centerOfMass[0], node.centerOfMass[1], node.centerOfMass[2])
			for i := 0; i < N; i += 1 {
				node.children[i].printNode()
			}
		} else if node.totalMass == 1 {
			fmt.Printf("Particle: (x,y,z) = (%f, %f, %f) Velocity: (%f, %f, %f)\n", node.centerOfMass[0], node.centerOfMass[1], node.centerOfMass[2], node.particle.VX, node.particle.VY, node.particle.VZ)
		} else {
			fmt.Printf("Empty Node\n")
		}
	}
}

/*
Creates a leaf node that represents a particle. If p == nil, creates an empty leaf node.
*/
func createLeafNode(octant Octant, p *Particle, parent *OctTreeNode, whichChild int) *OctTreeNode {
	newLeaf := &OctTreeNode{
		particle:   p,
		whichChild: whichChild,
		octant:     octant,
		parent:     parent,
	}
	// if p == nil, create an empty leaf node
	if p != nil {
		newLeaf.centerOfMass = [D]float64{p.X, p.Y, p.Z}
		newLeaf.totalMass = 1
	} else {
		newLeaf.centerOfMass = [D]float64{0, 0, 0}
		newLeaf.totalMass = 0
	}

	return newLeaf
}

/*
Creates and internal node. This node represents a space or octant where there is more than
two particles within the space.
*/
func createInternalNode(octant Octant, parent *OctTreeNode, whichChild int) *OctTreeNode {
	newInternal := &OctTreeNode{
		particle:     nil,
		whichChild:   whichChild,
		octant:       octant,
		parent:       parent,
		centerOfMass: [D]float64{},
	}
	newInternal.children = [N]*OctTreeNode{}
	for i := 0; N > i; i += 1 {
		newOctant := newInternal.findOctant(i)
		newInternal.children[i] = createLeafNode(newOctant, nil, newInternal, i)
	}
	return newInternal
}

/*
After inserting all particles call this function on the tree.
This will calculate the center of mass for the internal nodes and prune the empty leaf nodes
from the tree.
*/
func (tree *OctTree) CalculateCOMandPrune() {
	tree.rootPtr.children[0].recCOMPrune()
}

// Helper function for CalculateCOMandPrune()
// Returns true if the node needs to be pruned
func (node *OctTreeNode) recCOMPrune() bool {
	if node.totalMass == 0 {
		// Return true since this must be removed
		return true
	}
	if node.totalMass > 1 {
		totalX := 0.0
		totalY := 0.0
		totalZ := 0.0
		for i := 0; i < N; i += 1 {
			if node.children[i].recCOMPrune() {
				// Remove the node
				node.children[i] = nil
			} else {
				// Calculate center of mass
				totalX += node.children[i].centerOfMass[0] * node.children[i].totalMass
				totalY += node.children[i].centerOfMass[1] * node.children[i].totalMass
				totalZ += node.children[i].centerOfMass[2] * node.children[i].totalMass
			}
		}
		node.centerOfMass[0] = totalX / node.totalMass
		node.centerOfMass[1] = totalY / node.totalMass
		node.centerOfMass[2] = totalZ / node.totalMass
		return false
	}
	return false
}

// Returns true if the particle is within the given octant
func withinSpace(p *Particle, octant Octant) bool {
	if octant.lb <= p.X && p.X <= octant.rb && octant.db <= p.Y && p.Y <= octant.ub && octant.fb <= p.Z && p.Z <= octant.bb {
		return true
	} else {
		return false
	}
}

/*
Returns the correct bounds for the octant of the given space
See OctTreeNode struct for octant numbering
*/
func (root *OctTreeNode) findOctant(num int) Octant {
	lb, rb, db, ub, fb, bb := root.octant.lb, root.octant.rb, root.octant.db, root.octant.ub, root.octant.fb, root.octant.bb
	dx := (rb - lb) / 2
	dy := (ub - db) / 2
	dz := (bb - fb) / 2
	switch num {
	case 0:
		return Octant{lb: lb, rb: lb + dx, db: db, ub: db + dy, fb: fb, bb: fb + dz}
	case 1:
		return Octant{lb: lb, rb: lb + dx, db: db, ub: db + dy, fb: fb + dz, bb: bb}
	case 2:
		return Octant{lb: lb + dx, rb: rb, db: db, ub: db + dy, fb: fb, bb: fb + dz}
	case 3:
		return Octant{lb: lb + dx, rb: rb, db: db, ub: db + dy, fb: fb + dz, bb: bb}
	case 4:
		return Octant{lb: lb, rb: lb + dx, db: db + dy, ub: ub, fb: fb, bb: fb + dz}
	case 5:
		return Octant{lb: lb, rb: lb + dx, db: db + dy, ub: ub, fb: fb + dz, bb: bb}
	case 6:
		return Octant{lb: lb + dx, rb: rb, db: db + dy, ub: ub, fb: fb, bb: fb + dz}
	case 7:
		return Octant{lb: lb + dx, rb: rb, db: db + dy, ub: ub, fb: fb + dz, bb: bb}
	default:
		panic("Invalid integer value for createBounds (must be 0-7)")
	}
}
