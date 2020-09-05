import copy
import math
import random
import time
from multiprocessing import Pool

import numpy as np
from pandas_datareader import data
from scipy.stats import norm


def calculate_execution_time(function):
    def calculate_duration(*args, **kwargs):
        start_time = time.time()
        executing_function = function(*args, **kwargs)
        end_time = time.time()
        execution_time = round(end_time - start_time, 7)
        return executing_function, execution_time

    return calculate_duration


class MonteCarloSimulationFinance:
    def __init__(self, start_date, end_date, ticker_symbol, number_of_processes):
        self.start_date = start_date
        self.end_date = end_date
        self.ticker_symbol = ticker_symbol
        self.data = None
        self.time_series = None
        self.number_of_processes = number_of_processes
        self.parallel_flag = False

    def data_acquisition(self):
        stock = data.DataReader(self.ticker_symbol, 'yahoo', self.start_date, self.end_date)
        stock = stock.dropna()
        self.time_series = stock['Close']

    # Differencing time series = Shifting and lagging time series
    def calculate_periodic_daily_return(self):
        # Differencing time series
        # The diff() function calculates the first differences of the time series.
        self.data = np.log(self.time_series).diff().dropna()
        # Shifting and lagging time series
        # self.data=np.log(self.time_series / self.time_series.shift(1)).dropna()

    # A z-table, also called the standard normal table, is a mathematical table that allows us to know
    # the percentage of values below (to the left) a z-score in a standard normal distribution (SND).
    # A z-score, also known as a standard score, indicates the number of standard deviations
    # a raw score lays above or below the mean. When the mean of the z-score is calculated it is always 0,
    # and the standard deviation (variance) is always in increments of 1.
    def calculate_z_score(self):
        return norm.ppf(random.random())

    def calculate_average_daily_return(self):
        return np.mean(self.data)

    def calculate_variance(self):
        return np.var(self.data)

    def calculate_standard_deviation(self):
        return np.std(self.data)

    def calculate_drift(self):
        return self.calculate_average_daily_return() - self.calculate_variance() / 2

    def calculate_random_value(self):
        return self.calculate_standard_deviation() * self.calculate_z_score()

    # prediction window size: number of prediction days per simulation
    def simulation_finance(self, number_of_simulations, prediction_window_size):
        predictions = []
        prediction = []
        for i in range(number_of_simulations):
            # today’s price
            prediction.append(self.time_series.iloc[-1])
            for j in range(prediction_window_size):
                # Next Day’s Price=Today’s Price × e^(Drift+Random Value)
                prediction.append(
                    prediction[-1] * pow(math.e, (self.calculate_drift() + self.calculate_random_value())))

            predictions.append(copy.deepcopy(prediction))
            prediction.clear()
        return predictions

    @calculate_execution_time
    def mcs_finance_serial(self, number_of_simulations, prediction_window_size):
        self.parallel_flag = False
        return self.simulation_finance(number_of_simulations, prediction_window_size)

    @calculate_execution_time
    def mcs_finance_parallel(self, number_of_simulations, prediction_window_size):
        self.parallel_flag = True
        pool = Pool(processes=self.number_of_processes)
        number_of_simulations_per_process = int(number_of_simulations / self.number_of_processes)
        simulations_per_process = []
        # Append the same value multiple times to a list
        # To add v, n times, to l:
        # l += n * [v]
        # Mapping a function with multiple arguments to a multiprocessing pool will distribute
        # the input data across processes to be run with the referenced function.
        simulations_per_process += self.number_of_processes * [
            (number_of_simulations_per_process, prediction_window_size)]
        predictions = pool.starmap(self.simulation_finance, simulations_per_process)
        return predictions

    def export_finance_file(self, predictions):
        serial_number = 1
        if self.parallel_flag == False:
            # r before string converts normal string to raw string
            path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
                   r"\Execution Results\Finance\PythonFinanceSerial.txt"
            out_file = open(path, "w")
            for i in range(len(predictions)):
                for j in range(len(predictions[i])):
                    if j == 0:
                        # Write a serial number of simulations at the beginning of a file.
                        out_file.write(str(serial_number) + "," + " ")
                    # A component in one row of the output file.
                    out_file.write(str(predictions[i][j]))
                    if j == (len(predictions[i]) - 1):
                        out_file.write("\r\n")
                        serial_number += 1
                    else:
                        out_file.write("," + " ")

            out_file.close()

        else:
            # r before string converts normal string to raw string
            path = r"C:\Users\Dule\Desktop\NAPREDNE TEHNIKE PROGRAMIRANJA\PROJEKAT\NTP" \
                   r"\Execution Results\Finance\PythonFinanceParallel.txt"
            out_file = open(path, "w")
            for i in range(len(predictions)):
                for j in range(len(predictions[i])):
                    for k in range(len(predictions[i][j])):
                        if k == 0:
                            # Write a serial number of simulations at the beginning of a file.
                            out_file.write(str(serial_number) + "," + " ")
                        # A component in one row of the output file.
                        out_file.write(str(predictions[i][j][k]))
                        if k == (len(predictions[i][j]) - 1):
                            out_file.write("\r\n")
                            serial_number += 1
                        else:
                            out_file.write("," + " ")

            out_file.close()


if __name__ == "__main__":
    number_of_simulations_serial = 10
    prediction_window_size_serial = 100
    number_of_processes_serial = 1
    monte_carlo_simulation_finance_serial = MonteCarloSimulationFinance(
        '1980-01-01', '2019-12-31', 'AAPL', number_of_processes_serial)
    monte_carlo_simulation_finance_serial.data_acquisition()
    monte_carlo_simulation_finance_serial.calculate_periodic_daily_return()
    serial_predictions, serial_execution_time = monte_carlo_simulation_finance_serial.mcs_finance_serial(number_of_simulations_serial,
                                                                                                  prediction_window_size_serial)
    print("Stock market price predictions using the Monte Carlo simulation serial version")
    print("Execution time(n = {}, p = {}, w = {}) = {} seconds".format(number_of_simulations_serial,
                                                                       number_of_processes_serial,
                                                                       prediction_window_size_serial, serial_execution_time))
    monte_carlo_simulation_finance_serial.export_finance_file(serial_predictions)

    number_of_simulations_parallel = 10
    prediction_window_size_parallel = 100
    number_of_processes_parallel = 4
    monte_carlo_simulation_finance_parallel = MonteCarloSimulationFinance('1980-01-01', '2019-12-31', 'AAPL',
                                                                 number_of_processes_parallel)
    monte_carlo_simulation_finance_parallel.data_acquisition()
    monte_carlo_simulation_finance_parallel.calculate_periodic_daily_return()
    parallel_predictions, parallel_execution_time = monte_carlo_simulation_finance_parallel.mcs_finance_parallel(
        number_of_simulations_parallel,
        prediction_window_size_parallel)
    print("Stock market price predictions using the Monte Carlo simulation parallel version")
    print("Execution time(n = {}, p = {}, w = {}) = {} seconds".format(number_of_simulations_parallel,
                                                                       number_of_processes_parallel,
                                                                       prediction_window_size_parallel, parallel_execution_time))
    monte_carlo_simulation_finance_parallel.export_finance_file(parallel_predictions)



















