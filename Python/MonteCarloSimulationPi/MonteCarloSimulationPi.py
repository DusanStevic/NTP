import random
import time
from multiprocessing import Pool


def calculate_execution_time(function):
    def calculate_duration(*args, **kwargs):
        start_time = time.time()
        executing_function = function(*args, **kwargs)
        end_time = time.time()

        execution_time = round(end_time - start_time, 7)
        return executing_function, execution_time

    return calculate_duration


class MonteCarloSimulationPi:
    def __init__(self, number_of_processes):
        self.number_of_processes = number_of_processes
        self.parallel_flag = False
        self.experiment_flag = False

    # mcs stands for Monte Carlo Simulation
    # pi=3.1415926535
    def simulation_pi(self, number_of_simulations):
        if self.experiment_flag == True:
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

    @calculate_execution_time
    def mcs_pi_serial(self, number_of_simulations):
        self.parallel_flag = False
        pi = 4 * self.simulation_pi(number_of_simulations) / number_of_simulations


        return pi

    @calculate_execution_time
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
        pi = 4 * sum(inside_sum) / number_of_simulations

        print("Pi(n=" + str(number_of_simulations) + ") = " + str(pi))
        return pi


if __name__ == "__main__":
    number_of_simulations_n = 2500000
    number_of_processes_p = 4
    monte_carlo_simulation_pi = MonteCarloSimulationPi(number_of_processes_p)
    monte_carlo_simulation_pi.experiment_flag = False
    
    print("Approximation of Pi by using the Monte Carlo simulation serial version")
    pi, execution_time = monte_carlo_simulation_pi.mcs_pi_serial(number_of_simulations_n)
    print("Pi(n = {}, p = {}) = {}".format(number_of_simulations_n,number_of_processes_p,pi))
    print("Execution time (duration): {} seconds".format(execution_time))

    print("Approximation of Pi by using the Monte Carlo simulation parallel version")
    pi, execution_time = monte_carlo_simulation_pi.mcs_pi_parallel(number_of_simulations_n)
    print("Pi(n = {}, p = {}) = {}".format(number_of_simulations_n,number_of_processes_p,pi))
    print("Execution time (duration): {} seconds".format(execution_time))

