import random
import time
from multiprocessing import Pool
import math
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from pandas_datareader import data
from scipy.stats import norm


class MonteCarloSimulationFinance:
    def __init__(self, number_of_processes):
        self.number_of_processes = number_of_processes

    def calculate_periodic_daily_return(self):
        return 4

    # A z-table, also called the standard normal table, is a mathematical table that allows us to know
    # the percentage of values below (to the left) a z-score in a standard normal distribution (SND).
    # A z-score, also known as a standard score, indicates the number of standard deviations
    # a raw score lays above or below the mean. When the mean of the z-score is calculated it is always 0,
    # and the standard deviation (variance) is always in increments of 1.
    def calculate_z_score(self):
        return norm.ppf(random.random())

    def calculate_average_daily_return(self,data):
        return np.mean(data)

    def calculate_variance(self,data):
        return np.var(data)

    def calculate_standard_deviation(self,data):
        return np.std(data)

    def calculate_drift(self,data):
        return self.calculate_average_daily_return(data) - self.calculate_variance(data) / 2

    def calculate_random_value(self,data):
        return self.calculate_standard_deviation(data) * self.calculate_z_score()

    # prediction window size: number of prediction days per simulation
    def simulation_finance(self,number_of_simulations,prediction_window_size):
        predictions = []
        for x in range(number_of_simulations):
            x = random.random()
            y = random.random()
            if x * x + y * y < 1:
                inside = inside + 1

        return inside

    def mcs_finance_serial(self,number_of_simulations):
        pass

    def mcs_finance_parallel(self, number_of_simulations):
        pass


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

    monte_carlo_simulation_finance = MonteCarloSimulationFinance(5)
    start = time.time()
    # print("Approximation of Pi by using the Monte Carlo simulation serial version:" + str(
    # monte_carlo_simulation_pi.mcs_pi_serial(10000000)))
    #print("Approximation of Pi by using the Monte Carlo simulation parallel version:" + str(
        #monte_carlo_simulation_pi.mcs_pi_parallel(10000000)))
    end = time.time()
    duration = end - start
    print(f"Duration {duration} seconds")

    start_date = '1980-01-01'
    end_date = '2019-12-31'
    stocks = data.DataReader('AAPL', 'yahoo', start_date, end_date)
    stocks = stocks.dropna()
    time_series = stocks['Close']
    # print(time_series)
    l = [10, 20, 30, 40, 50, 60, 70, 80, 90, 100]
    rez = []
    for index in range(1, len(l)):
        rez.append(math.log(l[index] / l[index - 1]))

    l = pd.DataFrame(l)
    log = np.log(time_series).diff()
    print(rez)
    print(log)

    #print(calculate_average_daily_return(log))
    #print('ovo je varijansa:' + str(calculate_variance(log)))
    #print('ovo je standardna devijacija:' + str(calculate_standard_deviation(log)))

    #drift = calculate_drift(log)
    #random_value = calculate_random_value(log)

    print(time_series.tail(30))

    # time_series.plot(figsize=(20, 10))
    # plt.title('Closing stock price')
    # plt.ylabel('Stock Price($)')
    # plt.xlabel('Date')
    # plt.legend(loc='upper right')
    # plt.show()

    # print(time_series.iloc[-1])
    # print(time_series.iloc[-2]*pow(math.e, (drift+random_value)))
    predictions = []
    predictions.append(time_series.iloc[-30])
    # 30 days prediction
    #for i in range(30):
     #   predictions.append(predictions[-1] * pow(math.e, (drift + random_value)))
    #for i in predictions:
     #   print(i)
