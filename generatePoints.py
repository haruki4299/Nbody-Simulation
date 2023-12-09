import random

def generate_particle_file(filename, particle_count, cluster_size):
    with open(filename, 'w') as file:
        # Write the total particle count in the first line
        file.write(f"{particle_count}\n")

        # Generate and write particle data for the first cluster
        for _ in range(cluster_size):
            x, y, z = random.uniform(-50, -40), random.uniform(-50, -40), random.uniform(-50, -40)
            vx, vy, vz = random.uniform(-1, 1), random.uniform(-1, 1), random.uniform(-1, 1)
            file.write(f"{x} {y} {z} {vx} {vy} {vz}\n")

        # Generate and write particle data for the second cluster
        for _ in range(particle_count - cluster_size):
            x, y, z = random.uniform(40, 50), random.uniform(40, 50), random.uniform(40, 50)
            vx, vy, vz = random.uniform(-1, 1), random.uniform(-1, 1), random.uniform(-1, 1)
            file.write(f"{x} {y} {z} {vx} {vy} {vz}\n")

if __name__ == "__main__":
    output_filename = "generated_particle_data.txt"  # Change this to the desired output filename
    total_particle_count = 3000  # Change this to the desired number of particles
    first_cluster_size = 1500  # Change this to the desired size of the first cluster

    generate_particle_file(output_filename, total_particle_count, first_cluster_size)
    print(f"Generated {total_particle_count} particles in {output_filename}")
