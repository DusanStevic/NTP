import copy
import math
import random
import time
from multiprocessing import Pool

import numpy as np
from pandas_datareader import data
from scipy.stats import norm


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

    def mcs_finance_serial(self, number_of_simulations, prediction_window_size):
        self.parallel_flag = False
        return self.simulation_finance(number_of_simulations, prediction_window_size)

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

    def export_finance_file(self, number_of_simulations, prediction_window_size):
        serial_number = 1
        if self.parallel_flag == False:
            predictions = self.mcs_finance_serial(number_of_simulations, prediction_window_size)
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
            predictions = self.mcs_finance_parallel(number_of_simulations, prediction_window_size)
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
    monte_carlo_simulation_finance = MonteCarloSimulationFinance('1980-01-01', '2019-12-31', 'AAPL', 5)
    monte_carlo_simulation_finance.data_acquisition()
    monte_carlo_simulation_finance.calculate_periodic_daily_return()
    start = time.time()
    #print("Stock market price predictions using the Monte Carlo simulation serial version:" + str(
    #monte_carlo_simulation_finance.mcs_finance_serial(10, 5)))
    print("Stock market price predictions using the Monte Carlo simulation parallel version:" + str(
        monte_carlo_simulation_finance.mcs_finance_parallel(10, 5)))
    end = time.time()
    duration = end - start
    print(f"Duration {duration} seconds")

    monte_carlo_simulation_finance.export_finance_file(10, 8)

