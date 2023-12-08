# Barns-Hut 3D NBody Problem

## Project Description

The Parallelized N-Body Simulation project aims to efficiently model the gravitational interactions between a large number of celestial bodies by implementing the Barnes-Hut algorithm. The Barnes-Hut algorithm is a hierarchical algorithm that reduces the computational complexity of the traditional N-body simulation, making it well-suited for simulating systems with a large number of particles.

In the context of n-body simulations, determining the positions of particles in the next time step involves calculating the forces exerted by each particle on every other particle. The brute-force approach entails individually computing forces for each particle pair and then aggregating these individual calculations to obtain the final force. This process can be computationally intensive, especially with a large number of particles.

The Barnes-Hut algorithm introduces a significant improvement by treating groups of particles that are sufficiently distant as a single combined particle. This strategic abstraction minimizes the computational workload during each iteration, leading to more efficient simulations.

# Key Features

- Barns-Hut Algorithm Implementations

- Parallel Computing (golang)

- Work Stealing / Pipelining using channels

- Visualization using python

- Ability to read custom input files

- Performance evaluation (using randomly generated particles)

See write up pdf file for more details on the implementation and the performance evaluation.

## How To Use

- Recreate the experiment by running the `benchmark.sh` script in the directory. This will run the program multiple times and produce a speed up graph. (Note: Go into nbody.go and comment out the human readable print statement to the print statement that only prints the float value of the time.)

- May run individual cases of the nbody problem by running `go run nbody/nbody.go <p/s> <num of particles> <num of threads> <num of iterations>` in your command line. mode "p" or "s" represents running the program in parallel mode or sequential mode.

- By default, the program is made to run with randomly generated particles. However, we have also added a functionality to read in custom particle data text file. The directory contains a sample input file `generated_particle_data.txt`. The file should start with a line containing a single integer indicating the number of particles. It is followed by that many lines of particle data. Each line should consist of 6 float values (three for XYZ coordinates and three for XYZ velocities). In order to use the custom input file, run the program as `go run nbody/nbody.go <p/s> <num of particles> <num of threads> <num of iterations> <file_name>` (eg. `go run nbody/nbody.go p 3000 4 1000 generated_particle_data.txt`).

- After running the program, you should see a points.txt file generated in the directory. This is a file that tracks all of the coordinates of the points during the nbody simulation. You can run `python3 plot3D.py --nParticles <number of particles>` to see a visual representation of the nbody problem we simulated.

- You also may go into particles.go and fix the coefficients for the initialization of the random particles to play with.

- There is a python script `generatePoints.py` which I used to generate a large particle input file. Feel free to play around with it to generate some interesting initial point distributions.

## References

- https://en.wikipedia.org/wiki/Barnes%E2%80%93Hut_simulation

- https://jheer.github.io/barnes-hut/
