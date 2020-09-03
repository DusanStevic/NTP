from MonteCarloSimulationPi import MonteCarloSimulationPi
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
    # r before string converts normal string to raw string
    path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
           r"\Scaling Results\Pi\PythonPiStrongScaling.csv"
    out_file = open(path, "w")
    out_file.write("parallel_tasks_num,achieved_speedup,max_theoretical_speedup\n")
    number_of_simulations_n = 500000
    number_of_processes_p = 1
    monte_carlo_simulation_pi = MonteCarloSimulationPi(number_of_processes_p)
    monte_carlo_simulation_pi.experiment_flag = True
    print("Approximation of Pi by using the Monte Carlo simulation serial version")
    serial_pi, serial_execution_time = monte_carlo_simulation_pi.mcs_pi_serial(number_of_simulations_n)
    print("Pi(n = {}, p = {}) = {}".format(number_of_simulations_n,number_of_processes_p,serial_pi))
    print("Execution time (duration): {} seconds\n".format(serial_execution_time))
    for number_of_processes_p in range(2, 14):
        monte_carlo_simulation_pi = MonteCarloSimulationPi(number_of_processes_p)
        monte_carlo_simulation_pi.experiment_flag = True
        print("Approximation of Pi by using the Monte Carlo simulation parallel version")
        parallel_pi, parallel_execution_time = monte_carlo_simulation_pi.mcs_pi_parallel(number_of_simulations_n)
        print("Pi(n = {}, p = {}) = {}".format(number_of_simulations_n, number_of_processes_p, parallel_pi))
        print("Execution time (duration): {} seconds".format(parallel_execution_time))
        achieved_speedup = serial_execution_time / parallel_execution_time
        max_speedup = calculate_amdahl_speedup(number_of_processes_p)
        print("Achieved speedup is: {} times.".format(achieved_speedup))
        print("Maximum speedup according to Amdahl’s law is: {} times.\n".format(max_speedup))
        out_file.write("{},{},{}\n".format(number_of_processes_p,achieved_speedup,max_speedup))
    out_file.close()



def weak_scaling():
    # r before string converts normal string to raw string
    path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
           r"\Scaling Results\Pi\PythonPiWeakScaling.csv"
    out_file = open(path, "w")
    out_file.write("parallel_tasks_num,achieved_speedup,max_theoretical_speedup\n")

    number_of_simulations_n = 500000

    for number_of_processes_p in range(2, 14):
        increased_number_of_simulations = number_of_simulations_n * number_of_processes_p
        monte_carlo_simulation_pi = MonteCarloSimulationPi(number_of_processes_p)
        monte_carlo_simulation_pi.experiment_flag = True
        print("Approximation of Pi by using the Monte Carlo simulation serial version")
        serial_pi, serial_execution_time = monte_carlo_simulation_pi.mcs_pi_serial(number_of_simulations_n)
        print("Pi(n = {}, p = {}) = {}".format(number_of_simulations_n,1, serial_pi))
        print("Execution time (duration): {} seconds".format(serial_execution_time))
        print("Approximation of Pi by using the Monte Carlo simulation parallel version")
        parallel_pi, parallel_execution_time = monte_carlo_simulation_pi.mcs_pi_parallel(number_of_simulations_n)
        print("Pi(n = {}, p = {}) = {}".format(number_of_simulations_n, number_of_processes_p, parallel_pi))
        print("Execution time (duration): {} seconds".format(parallel_execution_time))
        achieved_speedup = serial_execution_time / parallel_execution_time
        max_speedup = calculate_gustafson_speedup(number_of_processes_p)
        print("Achieved speedup is: {} times.".format(achieved_speedup))
        print("Maximum speedup according to Gustafson’s law is: {} times.\n".format(max_speedup))
        out_file.write("{},{},{}\n".format(number_of_processes_p, achieved_speedup, max_speedup))
    out_file.close()



if __name__ == '__main__':
    print("Start strong scaling:")
    strong_scaling()
    print("End strong scaling:")
    print("Start weak scaling:")
    weak_scaling()
    print("End weak scaling:")
