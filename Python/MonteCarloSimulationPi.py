import random
import time
from multiprocessing import Pool


class MonteCarloSimulationPi:
    def __init__(self, number_of_processes):
        self.number_of_processes = number_of_processes

    # mcs stands for Monte Carlo Simulation
    # pi=3.1415926535
    def simulation_pi(self, number_of_simulations):
        path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP\Pharo\MonteCarloSimulationPi.txt"
        out_file = open(path, "w")
        inside = 0
        for x in range(number_of_simulations):
            x = random.random()
            y = random.random()
            out_file.write(str(int(x * 500)) + ' ' + str(int(y * 500)) + '\n')
            if x * x + y * y < 1:
                inside = inside + 1
        out_file.close()
        return inside

    def mcs_pi_serial(self, number_of_simulations):
        return 4 * self.simulation_pi(number_of_simulations) / number_of_simulations

    def mcs_pi_parallel(self, number_of_simulations):
        pool = Pool(processes=self.number_of_processes)
        number_of_simulations_per_process = int(number_of_simulations / self.number_of_processes)
        simulations_per_process = []
        simulations_per_process += self.number_of_processes * [number_of_simulations_per_process]
        inside_sum = pool.map(self.simulation_pi, simulations_per_process)
        return 4 * sum(inside_sum) / number_of_simulations


if __name__ == "__main__":
    # OVO SLJAKA
    # outFileName=r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP\Pharo\MonteCarloSimulationPi.txt"
    # outFile=open(outFileName, "w")
    # for i in range(100):
    #     x = random.randint(0, 500)
    #     y = random.randint(0, 500)
    #     print(str(x)+' '+str(y))
    #     outFile.write(str(x)+' '+str(y)+'\n')
    # outFile.close()

    start = time.time()
    monte_carlo_simulation_pi = MonteCarloSimulationPi(5)

    print(monte_carlo_simulation_pi.mcs_pi_serial(100000))
    # print(monte_carlo_simulation_pi.mcs_pi_parallel(10000000))
    end = time.time()
    print("vreme izvrsenja:" + str(end - start) + " sekundi")
