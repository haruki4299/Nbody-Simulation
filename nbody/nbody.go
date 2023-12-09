package main

import (
	"Nbody-Simulation/barnsHut"
	"Nbody-Simulation/workStealing"
	"fmt"
	"os"
	"strconv"
	"time"
)

/*
Runs the sequential version of the Barn Hut Nbody Simulation

nParticles: total Number of Particles to Simulate

nIters: number of iterations

dt: time step (per iterations)
*/
func seqNBody(nParticles int, nIters int, dt float64, input_file_name string) {
	var particles []barnsHut.Particle
	if input_file_name == "" {
		// Generate random particles
		particles = barnsHut.CreateRandParticles(nParticles)
	} else {
		// Custom Particles
		particles = barnsHut.ReadParticles(input_file_name)
	}

	// Create write file
	file, err := os.Create("points.txt")
	if err != nil {
		panic("Error creating file")
	}
	defer file.Close()

	// Write initial positions of particles
	for i := 0; i < nParticles; i += 1 {
		fmt.Fprintf(file, "%f,%f,%f\n", particles[i].X, particles[i].Y, particles[i].Z)
	}

	// Nbody Simulation
	for i := 0; i < nIters; i += 1 {
		startTime := time.Now()
		// Find the bounds of the particles
		octant := barnsHut.CreateBounds(particles)

		// Create and Insert into Tree
		tree := barnsHut.InitOctTree(octant)

		for i := range particles {
			newPart := particles[i]
			tree.AddParticle(&newPart)
		}

		tree.CalculateCOMandPrune()

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		elapsedTime := duration.Seconds()
		if i == 0 {
			fmt.Printf("Initialization of Tree: %f seconds\n", elapsedTime)
		}

		startTime2 := time.Now()
		// Calculate the forces and find the velocity of the particle
		for i := 0; i < nParticles; i += 1 {
			barnsHut.CallBarnsHut(&particles[i], dt, tree)
		}

		endTime2 := time.Now()
		duration2 := endTime2.Sub(startTime2)
		elapsedTime2 := duration2.Seconds()
		if i == 0 || i == nIters-1 {
			fmt.Printf("Barns Hut Force Calculation (%dth Iteration): %f seconds\n", i+1, elapsedTime2)
		}

		// Calculate the new positions based on their velocity
		startTime3 := time.Now()
		for i := 0; i < nParticles; i += 1 {
			particles[i].NewPos(dt)
		}
		endTime3 := time.Now()
		duration3 := endTime3.Sub(startTime3)
		elapsedTime3 := duration3.Seconds()
		if i == 0 {
			fmt.Printf("New Position Calculation: %f seconds\n", elapsedTime3)
		}

		// Write new positions to file
		startTime4 := time.Now()
		for i := 0; i < nParticles; i += 1 {
			fmt.Fprintf(file, "%f,%f,%f\n", particles[i].X, particles[i].Y, particles[i].Z)
		}
		endTime4 := time.Now()
		duration4 := endTime4.Sub(startTime4)
		elapsedTime4 := duration4.Seconds()
		if i == 0 {
			fmt.Printf("Writing to points.txt: %f seconds\n", elapsedTime4)
		}
	}
}

/*
Runs the parallel version of the Barn Hut Nbody Simulation

nParticles: total Number of Particles to Simulate

nThread: number of threads to use

nIters: number of iterations

dt: time step (per iterations)
*/
func parallelNBody(nParticles int, nThreads int, nIters int, dt float64, input_file_name string) {
	var particles []barnsHut.Particle
	if input_file_name == "" {
		// Generate random particles
		particles = barnsHut.CreateRandParticles(nParticles)
	} else {
		// Custom Particles
		particles = barnsHut.ReadParticles(input_file_name)
	}

	// Create write file
	file, err := os.Create("points.txt")
	if err != nil {
		panic("Error creating file")
	}
	defer file.Close()

	for i := 0; i < nParticles; i += 1 {
		fmt.Fprintf(file, "%f,%f,%f\n", particles[i].X, particles[i].Y, particles[i].Z)
	}

	// Nbody Simulation

	// Create a channel to receive completion signals from threads
	completionChannel := make(chan int64)

	for i := 0; i < nIters; i += 1 {
		// Find the bounds of the particles
		octant := barnsHut.CreateBounds(particles)

		// Create and Insert into Tree
		tree := barnsHut.InitOctTree(octant)

		for i := range particles {
			newPart := particles[i]
			tree.AddParticle(&newPart)
		}

		tree.CalculateCOMandPrune()

		// Global Variable indicating if all threads work is done
		allWorkDone := new(bool)
		*allWorkDone = false

		// Create a list of threads
		threadList := make([]workStealing.WorkStealingThread, nThreads)
		for i := 0; i < nThreads; i += 1 {
			threadData := workStealing.CreateThreadData(i, nThreads, allWorkDone)
			threadList[i] = *threadData
		}

		// Create goroutines that calculate the forces on each particle
		for i := 0; i < nThreads; i += 1 {
			// Simple division to determine work load
			start := (nParticles / nThreads) * i
			end := ((nParticles / nThreads) * (i + 1)) - 1
			if i == nThreads-1 {
				end = nParticles - 1
			}
			go barnsHutWorker(start, end, dt, tree, threadList, &threadList[i], completionChannel, particles)
		}

		// Once each particle is processed, this information will be relayed via the channel
		// Then we can start calculating the new positions for the particle
		done := 0
		for done != nParticles {
			num := <-completionChannel
			particles[num].NewPos(dt)
			done += 1
		}
		// All Work Done signals to the threads that they can be done once
		// they are done with the current work
		*allWorkDone = true

		// Write to file the new Particle positions
		for i := 0; i < nParticles; i += 1 {
			fmt.Fprintf(file, "%f,%f,%f\n", particles[i].X, particles[i].Y, particles[i].Z)
		}
	}

}

/*
Each goroutine will first enqueue the tasks they are assigned. (Particle numbers)
Then they will start working on their tasks. If they complete all of their original tasks,
They will steal tasks from other threads if they exist.
*/
func barnsHutWorker(start int, end int, dt float64, tree *barnsHut.OctTree, threadList []workStealing.WorkStealingThread, threadData *workStealing.WorkStealingThread, completionChannel chan<- int64, particles []barnsHut.Particle) {
	for i := start; i <= end; i += 1 {
		threadData.Queue.PushBottom(int64(i))
	}
	threadData.WorkLeft = true

	for !(*threadData.AllWorkDone) {
		if threadData.WorkLeft {
			// Work on own tasks
			num := threadData.Queue.PopBottom()
			if num == -1 {
				threadData.WorkLeft = false
			} else {
				barnsHut.CallBarnsHut(&particles[num], dt, tree)
				completionChannel <- num // Signal to the main thread that this particle is done being processed
			}
		} else {
			// Find another thread to steal work from
			victim := threadData.Rand.Intn(threadData.TotalThreads)
			if victim != threadData.ThreadID {
				num := threadList[victim].Queue.PopTop()
				if num != -1 {
					barnsHut.CallBarnsHut(&particles[num], dt, tree)
					completionChannel <- num // Signal to the main thread that this particle is done being processed
				}
			}
		}
	}
}

func main() {
	// Read input
	if len(os.Args) < 5 {
		panic("Not enough arguments go run <path>/nbody.go <mode> <nParticles> <nThreads> <nIterations>")
	}

	mode := os.Args[1]
	if mode != "s" && mode != "p" {
		panic("Invalid mode (Usage: go run <path>/nbody.go <mode> <nParticles> <nThreads> <nIterations>)")
	}

	dt := 0.01
	nParticles, err1 := strconv.Atoi(os.Args[2])
	nThreads, err2 := strconv.Atoi(os.Args[3])
	nIters, err3 := strconv.Atoi(os.Args[4])
	input_file := ""
	if len(os.Args) == 6 {
		// Custom Points
		input_file = os.Args[5]
	}

	if err1 != nil || err2 != nil || err3 != nil {
		panic("Error Converting Arguments into Integers")
	}

	// Nbody simulation
	startTime := time.Now()
	if mode == "s" {
		seqNBody(nParticles, nIters, dt, input_file)
	} else {
		parallelNBody(nParticles, nThreads, nIters, dt, input_file)
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	elapsedTime := duration.Seconds()
	fmt.Printf("Elapsed time: %f seconds\n", elapsedTime) // Human Readable Output
	//fmt.Printf("%f\n", elapsedTime) // Experimental use
}
