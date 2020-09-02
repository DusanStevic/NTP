from MonteCarloSimulationPi import MonteCarloSimulationPi


def strong_scaling():
    for i in range(1, 11):
        monte_carlo_simulation_pi = MonteCarloSimulationPi(i)
        monte_carlo_simulation_pi.experiment_flag = False
        monte_carlo_simulation_pi.mcs_pi_parallel(10000000)


def weak_scaling():
    for i in range(1, 11):
        monte_carlo_simulation_pi = MonteCarloSimulationPi(i)
        monte_carlo_simulation_pi.experiment_flag = False
        monte_carlo_simulation_pi.mcs_pi_parallel(10000000*i)


if __name__ == '__main__':
    print("Start strong scaling:")
    strong_scaling()
    print("End strong scaling:")
    print("Start weak scaling:")
    weak_scaling()
    print("End weak scaling:")
