from MonteCarloSimulationIntegration import MonteCarloSimulationIntegration

# https://www.kth.se/blogs/pdc/2018/11/scalability-strong-and-weak-scaling/

# s + p = 1
# complementary values

# s is the proportion of execution time spent on the serial part
# s is part of the program which cannot be parallelized
# s is ratio of the serial part
SERIAL_PART_s = 0

# p is the proportion of execution time spent on the part that can be parallelized
# p is part of the program which can be parallelized
# p is ratio of the parallel part
PARALLEL_PART_p = 1


# Amdahl’s law and strong scaling
# Amdahl’s law can be formulated as follows speedup = 1 / (s + p / N) where
# s = SERIAL_PART_s is the proportion of execution time spent on the serial part,
# p = PARALLEL_PART_p is the proportion of execution time spent on the part that can be parallelized,
# and N = number_of_processes is the number of processors.
def calculate_amdahl_speedup(number_of_processes):
    return 1.0 / (SERIAL_PART_s + PARALLEL_PART_p / number_of_processes)


# Gustafson’s law and weak scaling
# Gustafson’s law can be formulated as follows speedup = s + p × N where
# s = SERIAL_PART_s is the proportion of execution time spent on the serial part,
# p = PARALLEL_PART_p is the proportion of execution time spent on the part that can be parallelized,
# and N = number_of_processes is the number of processors.
def calculate_gustafson_speedup(number_of_processes):
    return SERIAL_PART_s + PARALLEL_PART_p * number_of_processes


def strong_scaling():
    print("=======================")
    print("Start strong scaling:")
    print("=======================\n")
    # r before string converts normal string to raw string
    path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
           r"\Scaling Results\Integration\PythonIntegrationStrongScaling.csv"
    out_file = open(path, "w")
    out_file.write("number_of_processes,achieved_speedup,theoretical_maximum_speedup\n")
    number_of_simulations_n = 1500000
    number_of_processes_serial = 1
    monte_carlo_simulation_integration = MonteCarloSimulationIntegration(number_of_processes_serial)
    monte_carlo_simulation_integration.experiment_flag = True
    print("Integral approximation by using the Monte Carlo simulation serial version")
    serial_integration, serial_execution_time = monte_carlo_simulation_integration.mcs_integration_serial(number_of_simulations_n)
    print("Integral(n = {}, p = {}) = {}".format(number_of_simulations_n, number_of_processes_serial, serial_integration))
    print("Execution time (duration): {} seconds\n".format(serial_execution_time))
    for number_of_processes_parallel in range(2, 14):
        monte_carlo_simulation_integration = MonteCarloSimulationIntegration(number_of_processes_parallel)
        monte_carlo_simulation_integration.experiment_flag = True
        print("Integral approximation by using the Monte Carlo simulation parallel version")
        parallel_integration, parallel_execution_time = monte_carlo_simulation_integration.mcs_integration_parallel(number_of_simulations_n)
        print("Integral(n = {}, p = {}) = {}".format(number_of_simulations_n, number_of_processes_parallel, parallel_integration))
        print("Execution time (duration): {} seconds".format(parallel_execution_time))
        achieved_speedup = serial_execution_time / parallel_execution_time
        theoretical_maximum_speedup = calculate_amdahl_speedup(number_of_processes_parallel)
        print("Achieved speedup is: {} times.".format(achieved_speedup))
        print(
            "Theoretical maximum speedup according to Amdahl’s law is: {} times.\n".format(theoretical_maximum_speedup))
        out_file.write("{},{},{}\n".format(number_of_processes_parallel, achieved_speedup, theoretical_maximum_speedup))
    out_file.close()
    print("End strong scaling.")


def weak_scaling():
    print("=======================")
    print("Start weak scaling:")
    print("=======================\n")
    # r before string converts normal string to raw string
    path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
           r"\Scaling Results\Integration\PythonIntegrationWeakScaling.csv"
    out_file = open(path, "w")
    out_file.write("number_of_processes,achieved_speedup,theoretical_maximum_speedup\n")
    number_of_simulations_n = 1500000
    for number_of_processes_p in range(2, 14):
        increased_number_of_simulations = number_of_simulations_n * number_of_processes_p
        monte_carlo_simulation_integration = MonteCarloSimulationIntegration(number_of_processes_p)
        monte_carlo_simulation_integration.experiment_flag = True
        print("Integral approximation by using the Monte Carlo simulation serial version")
        serial_integration, serial_execution_time = monte_carlo_simulation_integration.mcs_integration_serial(
            increased_number_of_simulations)
        print("Integration(n = {}, p = {}) = {}".format(increased_number_of_simulations, 1, serial_integration))
        print("Execution time (duration): {} seconds".format(serial_execution_time))
        print("Integral approximation by using the Monte Carlo simulation parallel version")
        parallel_integration, parallel_execution_time = monte_carlo_simulation_integration.mcs_integration_parallel(
            increased_number_of_simulations)
        print("Integral(n = {}, p = {}) = {}".format(increased_number_of_simulations, number_of_processes_p, parallel_integration))
        print("Execution time (duration): {} seconds".format(parallel_execution_time))
        achieved_speedup = serial_execution_time / parallel_execution_time
        theoretical_maximum_speedup = calculate_gustafson_speedup(number_of_processes_p)
        print("Achieved speedup is: {} times.".format(achieved_speedup))
        print("Theoretical maximum speedup according to Gustafson’s law is: {} times.\n".format(
            theoretical_maximum_speedup))
        out_file.write("{},{},{}\n".format(number_of_processes_p, achieved_speedup, theoretical_maximum_speedup))
    out_file.close()
    print("End weak scaling.")


if __name__ == '__main__':
    strong_scaling()
    weak_scaling()