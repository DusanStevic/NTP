import random
import time
from multiprocessing import Process, Queue,Pool


# mcs stands for Monte Carlo Simulation
# pi=3.1415926535
def simulation(number_of_simulations):
    inside = 0
    for x in range(number_of_simulations):
        x = random.random()
        y = random.random()
        if x * x + y * y < 1:
            inside = inside + 1

    return inside

def mcs_pi_serial(number_of_simulations):
    return 4 * simulation(number_of_simulations) / number_of_simulations



def mcs_pi_parallel(number_of_simulations,number_of_processes):
    pool = Pool(processes=number_of_processes)
    number_of_simulations_per_process = int(number_of_simulations/number_of_processes)
    simulations_per_process = []
    simulations_per_process += number_of_processes * [number_of_simulations_per_process]
    inside_sum = pool.map(simulation, simulations_per_process)
    return 4*sum(inside_sum) / number_of_simulations


# outFileName=""C:\Users\Dule\Desktop\dule.txt""
# outFile=open(outFileName, "w")
# outFile.write("""Hello my name is ABCD""")
# outFile.close()

# def save(self, path, parallel, iter):
#     with open(f'{path}/{"parallel" if parallel else "serial"}{iter}.mesh', 'w') as file:
#         for i in range(self.rows):
#             for j in range(self.columns):
#                 file.write(f"{self.mesh[i * self.columns + j]} ")
#             file.write("\n")
#
# def upis_test(test):
#     csvfile = open('test_previewA.csv', 'w', newline='',encoding='utf-8')
#     obj = csv.writer(csvfile)
#     for row in test:
#         obj.writerow(row)
#         csvfile.flush()
#     csvfile.close()
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
    #print(mcs_pi_serial(100_000_000))
    print(mcs_pi_parallel(100_000_000,5))
    end = time.time()
    print("vreme izvrsenja:" + str(end - start) + " sekundi")
