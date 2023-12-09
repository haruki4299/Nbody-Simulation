import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
from matplotlib.animation import FuncAnimation
import argparse

# Parse command line arguments
parser = argparse.ArgumentParser(description='Particle Animation')
parser.add_argument('--nParticles', type=int,
                    help='Number of particles per iteration')
args = parser.parse_args()

# Read particle positions from the file
with open('points.txt', 'r') as file:
    lines = file.readlines()

# Parse coordinates
coordinates = [list(map(float, line.strip().split(','))) for line in lines]

# Set up the 3D plot
fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')
sc = ax.scatter([], [], [], c='r', marker='o')

# Set labels
ax.set_xlabel('X Label')
ax.set_ylabel('Y Label')
ax.set_zlabel('Z Label')
ax.set_title('Particle Animation')

# Number of particles per iteration
particles_per_iteration = args.nParticles

# Set initial axis limits
x_min, x_max = float('inf'), float('-inf')
y_min, y_max = float('inf'), float('-inf')
z_min, z_max = float('inf'), float('-inf')

for coord in coordinates:
    x_min = min(x_min, coord[0])
    x_max = max(x_max, coord[0])
    y_min = min(y_min, coord[1])
    y_max = max(y_max, coord[1])
    z_min = min(z_min, coord[2])
    z_max = max(z_max, coord[2])

ax.set_xlim(x_min, x_max)
ax.set_ylim(y_min, y_max)
ax.set_zlim(z_min, z_max)

# Animation update function
def update(iteration):
    start_index = iteration * particles_per_iteration
    end_index = (iteration + 1) * particles_per_iteration
    x_data = [coord[0] for coord in coordinates[start_index:end_index]]
    y_data = [coord[1] for coord in coordinates[start_index:end_index]]
    z_data = [coord[2] for coord in coordinates[start_index:end_index]]
    sc._offsets3d = (x_data, y_data, z_data)
    ax.set_title(f'Particle Animation - Iteration {iteration}')

# Calculate the number of iterations
num_iterations = len(coordinates) // particles_per_iteration

# Create the animation
animation = FuncAnimation(fig, update, frames=num_iterations, interval=1, repeat=False)

# Show the animation
plt.show()
