import random
import time
from multiprocessing import Pool


class MonteCarloSimulationPi:
    def __init__(self, number_of_processes):
        self.number_of_processes = number_of_processes
        self.parallel_flag = False
        self.experiment_flag = False

    # mcs stands for Monte Carlo Simulation
    # pi=3.1415926535
    def simulation_pi(self, number_of_simulations):
        if self.experiment_flag==True:
            inside = 0
            for _ in range(number_of_simulations):
                x = random.random()
                y = random.random()
                # The unit circle is the circle of radius 1 centered at the origin(0, 0)
                # in the Cartesia coordinate system in the Euclidean plane.
                if x * x + y * y < 1:
                    inside = inside + 1
            return inside
        else:
            if self.parallel_flag == False:
                # r before string converts normal string to raw string
                path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP\Pharo\PythonPiSerial.txt"
            else:
                # r before string converts normal string to raw string
                path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP\Pharo\PythonPiParallel.txt"
            out_file = open(path, "w")
            inside = 0
            for _ in range(number_of_simulations):
                x = random.random()
                y = random.random()
                # Pharo for Data Visualization. Circle of radius 250 centered at the point(250, 250).
                # To create a Rectangle in Pharo you must provide the top left and the bottom right points.
                out_file.write(str(int(x * 500)) + ' ' + str(int(y * 500)) + '\n')
                # The unit circle is the circle of radius 1 centered at the origin(0, 0)
                # in the Cartesia coordinate system in the Euclidean plane.
                if x * x + y * y < 1:
                    inside = inside + 1
            out_file.close()
            return inside


    def mcs_pi_serial(self, number_of_simulations):
        self.parallel_flag = False
        return 4 * self.simulation_pi(number_of_simulations) / number_of_simulations

    def mcs_pi_parallel(self, number_of_simulations):
        self.parallel_flag = True
        pool = Pool(processes=self.number_of_processes)
        number_of_simulations_per_process = int(number_of_simulations / self.number_of_processes)
        simulations_per_process = []
        # Append the same value multiple times to a list
        # To add v, n times, to l:
        # l += n * [v]
        simulations_per_process += self.number_of_processes * [number_of_simulations_per_process]
        inside_sum = pool.map(self.simulation_pi, simulations_per_process)
        return 4 * sum(inside_sum) / number_of_simulations




if __name__ == "__main__":
    monte_carlo_simulation_pi = MonteCarloSimulationPi(5)
    start = time.time()
    monte_carlo_simulation_pi.experiment_flag = True
    print("Approximation of Pi by using the Monte Carlo simulation serial version:" + str(
    monte_carlo_simulation_pi.mcs_pi_serial(100000000)))

    #print("Approximation of Pi by using the Monte Carlo simulation parallel version:" + str(
        #monte_carlo_simulation_pi.mcs_pi_parallel(100000000)))
    end = time.time()
    duration = end - start
    print(f"Duration {duration} seconds")
