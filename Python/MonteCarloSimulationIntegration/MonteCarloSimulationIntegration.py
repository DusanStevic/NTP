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
    # The area under the graph of a function can be found by adding slices that approach zero in width.
    # The function f(x) to be integrated is called the integrand.

    def simulation_integration(self, number_of_simulations, lower_bound, upper_bound, slice_size, function):
        if self.experiment_flag == True:
            # Points under the graph of a function.
            below = 0
            lower_bound_interval = lower_bound
            upper_bound_interval = upper_bound
            # Define the interval between the lower and upper bound.
            x = []
            # Function Values
            y = []
            # Maximum of the function f(x) on the interval[lower_bound, upper_bound]
            f_max = function(lower_bound)

            while lower_bound_interval < upper_bound_interval:
                x.append(lower_bound_interval)
                t = function(lower_bound_interval)
                y.append(t)
                if t > f_max:
                    f_max = t
                lower_bound_interval += slice_size

            for _ in range(number_of_simulations):
                x_rand = lower_bound + (upper_bound - lower_bound) * random.random()
                y_rand = 0 + f_max * random.random()
                if y_rand < function(x_rand):
                    below = below + 1
            # Rectangle area that surrounds the area under the graph of a function.
            a = upper_bound - lower_bound
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
            below = 0
            for _ in range(number_of_simulations):
                x = random.random()
                y = random.random()
                # Pharo for Data Visualization. Circle of radius 250 centered at the point(250, 250).
                # To create a Rectangle in Pharo you must provide the top left and the bottom right points.
                out_file.write(str(int(x * 500)) + ' ' + str(int(y * 500)) + '\n')
                # The unit circle is the circle of radius 1 centered at the origin(0, 0)
                # in the Cartesia coordinate system in the Euclidean plane.
                if x * x + y * y < 1:
                    below = below + 1
            out_file.close()
            return below

    @calculate_execution_time
    def mcs_integration_serial(self, number_of_simulations, lower_bound, upper_bound, slice_size, function):
        self.parallel_flag = False
        integral = self.simulation_integration(number_of_simulations, lower_bound, upper_bound, slice_size, function)
        return integral


if __name__ == "__main__":
    lower_bound = 1
    upper_bound = 2
    slice_size = 0.01


    # The function f(x) to be integrated is called the integrand.
    # The function we are integrating must be non-negative continuous function between lower bound and upper bound
    # Non-negative function: is a function when it attain non negative values only. A function would be called a
    # positive function if its values are positive for all arguments of its domain, or a non-negative function
    # if all of its values are non-negative. The function graph sits above or on the x-axis.
    # Continuous function: is a function with no holes, jumps or vertical asymptotes
    # (where the function heads up/down towards infinity). A vertical asymptote between lower bound and
    # upper bound affects the definite integral.

    def f(x):
        return 2 * x


    number_of_simulations_serial = 10000
    number_of_processes_serial = 1
    monte_carlo_simulation_integration_serial = MonteCarloSimulationIntegration(number_of_processes_serial)
    monte_carlo_simulation_integration_serial.experiment_flag = True
    print("Integral Approximation by using the Monte Carlo simulation serial version")
    serial_integration, serial_execution_time = monte_carlo_simulation_integration_serial.mcs_integration_serial(
        number_of_simulations_serial, lower_bound, upper_bound, slice_size, f)
    print("Integral(n = {}, p = {}) = {}".format(number_of_simulations_serial, number_of_processes_serial,
                                                 serial_integration))
    print("Execution time (duration): {} seconds".format(serial_execution_time))
