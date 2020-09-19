import random
import time
import numpy as np
import math
from multiprocessing import Pool


def calculate_execution_time(function):
    def calculate_duration(*args, **kwargs):
        start_time = time.time()
        executing_function = function(*args, **kwargs)
        end_time = time.time()
        execution_time = round(end_time - start_time, 7)
        return executing_function, execution_time

    return calculate_duration


class MonteCarloSimulationIntegration:
    def __init__(self, number_of_processes):
        self.number_of_processes = number_of_processes
        self.parallel_flag = False
        self.experiment_flag = False
        # Upper and Lower Bounds of Integral.
        self.LOWER_BOUND = 1
        self.UPPER_BOUND = 2
        # The area under the graph of a function can be found by adding slices that approach zero in width.
        self.SLICE_SIZE = 0.01

    # The function f(x) to be integrated is called the integrand.
    # The function we are integrating must be non-negative continuous function between lower bound and upper bound
    # Non-negative function: is a function when it attain non negative values only. A function would be called a
    # positive function if its values are positive for all arguments of its domain, or a non-negative function
    # if all of its values are non-negative. The function graph sits above or on the x-axis.
    # Continuous function: is a function with no holes, jumps or vertical asymptotes
    # (where the function heads up/down towards infinity). A vertical asymptote between lower bound and
    # upper bound affects the definite integral.
    def function(self, x):
        return 2*x

    def simulation_integration(self, number_of_simulations):
        if self.experiment_flag == True:
            # Points under the graph of a function.
            below = 0
            lower_bound_interval = self.LOWER_BOUND
            upper_bound_interval = self.UPPER_BOUND
            # Define the interval between the lower and upper bound.
            x = []
            # Function Values
            y = []
            # Maximum of the function f(x) on the interval[lower_bound, upper_bound]
            f_max = self.function(self.LOWER_BOUND)

            while lower_bound_interval < upper_bound_interval:
                x.append(lower_bound_interval)
                t = self.function(lower_bound_interval)
                y.append(t)
                if t > f_max:
                    f_max = t
                lower_bound_interval += self.SLICE_SIZE

            for _ in range(number_of_simulations):
                x_rand = self.LOWER_BOUND + (self.UPPER_BOUND - self.LOWER_BOUND) * random.random()
                y_rand = 0 + f_max * random.random()
                if y_rand < self.function(x_rand):
                    below = below + 1
            # Rectangle area that surrounds the area under the graph of a function.
            a = self.UPPER_BOUND - self.LOWER_BOUND
            b = f_max - 0
            rectangle_area = a * b
            # bellow = Points under the graph of a function.
            # number_of_simulations = Total number of points = Points inside rectangle
            proportion = below / number_of_simulations
            integral = proportion * rectangle_area
            return integral
        else:
            if self.parallel_flag == False:
                # r before string converts normal string to raw string
                path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
                       r"\Execution Results\Integration\PythonIntegrationSerial.txt"
            else:
                # r before string converts normal string to raw string
                path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
                       r"\Execution Results\Integration\PythonIntegrationParallel.txt"
            out_file = open(path, "w")

            # Points under the graph of a function.
            below = 0
            lower_bound_interval = self.LOWER_BOUND
            upper_bound_interval = self.UPPER_BOUND
            # Define the interval between the lower and upper bound.
            x = []
            # Function Values
            y = []
            # Maximum of the function f(x) on the interval[lower_bound, upper_bound]
            f_max = self.function(self.LOWER_BOUND)

            while lower_bound_interval < upper_bound_interval:
                x.append(lower_bound_interval)
                t = self.function(lower_bound_interval)
                y.append(t)
                if t > f_max:
                    f_max = t
                lower_bound_interval += self.SLICE_SIZE

            for _ in range(number_of_simulations):
                x_rand = self.LOWER_BOUND + (self.UPPER_BOUND - self.LOWER_BOUND) * random.random()
                y_rand = 0 + f_max * random.random()
                out_file.write(str(round(x_rand, 2)) + ' ' + str(round(y_rand, 2)) + '\n')
                if y_rand < self.function(x_rand):
                    below = below + 1
            # Rectangle area that surrounds the area under the graph of a function.
            a = self.UPPER_BOUND - self.LOWER_BOUND
            b = f_max - 0
            rectangle_area = a * b
            # bellow = Points under the graph of a function.
            # number_of_simulations = Total number of points = Points inside rectangle
            proportion = below / number_of_simulations
            integral = proportion * rectangle_area
            out_file.close()
            return integral


    @calculate_execution_time
    def mcs_integration_serial(self, number_of_simulations):
        self.parallel_flag = False
        integral = self.simulation_integration(number_of_simulations)
        return integral

    @calculate_execution_time
    def mcs_integration_parallel(self, number_of_simulations):
        self.parallel_flag = True
        pool = Pool(processes=self.number_of_processes)
        number_of_simulations_per_process = int(number_of_simulations / self.number_of_processes)
        simulations_per_process = []
        # Append the same value multiple times to a list
        # To add v, n times, to l:
        # l += n * [v]
        # Mapping a function with multiple arguments to a multiprocessing pool will distribute
        # the input data across processes to be run with the referenced function.
        simulations_per_process += self.number_of_processes * [number_of_simulations_per_process]
        # list of partial result per process
        list_of_integral_per_process = pool.map(self.simulation_integration, simulations_per_process)
        # cumulative result, aggregating partial results
        integral_per_processes = sum(list_of_integral_per_process)
        integral = integral_per_processes / self.number_of_processes
        return integral


if __name__ == "__main__":

    number_of_simulations_serial = 100
    number_of_processes_serial = 1
    monte_carlo_simulation_integration_serial = MonteCarloSimulationIntegration(number_of_processes_serial)
    monte_carlo_simulation_integration_serial.experiment_flag = False
    print("Integral Approximation by using the Monte Carlo simulation serial version")
    serial_integration, serial_execution_time = monte_carlo_simulation_integration_serial.mcs_integration_serial(
        number_of_simulations_serial)
    print("Integral(n = {}, p = {}) = {}".format(number_of_simulations_serial, number_of_processes_serial,
                                                 serial_integration))
    print("Execution time (duration): {} seconds".format(serial_execution_time))

    number_of_simulations_parallel = 100
    number_of_processes_parallel = 4
    monte_carlo_simulation_integration_parallel = MonteCarloSimulationIntegration(number_of_processes_parallel)
    monte_carlo_simulation_integration_parallel.experiment_flag = False
    print("Integral Approximation by using the Monte Carlo simulation parallel version")
    parallel_integration, parallel_execution_time = monte_carlo_simulation_integration_parallel.mcs_integration_parallel(
        number_of_simulations_parallel)
    print("Integral(n = {}, p = {}) = {}".format(number_of_simulations_parallel, number_of_processes_parallel,
                                                 parallel_integration))
    print("Execution time (duration): {} seconds".format(parallel_execution_time))
