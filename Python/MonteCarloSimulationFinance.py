import random
import time
from multiprocessing import Pool
import math
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from pandas_datareader import data
from scipy.stats import norm
import copy


class MonteCarloSimulationFinance:
    def __init__(self, start_date, end_date, ticker_symbol, number_of_processes):
        self.start_date = start_date
        self.end_date = end_date
        self.ticker_symbol = ticker_symbol
        self.data = None
        self.time_series = None
        self.number_of_processes = number_of_processes

    def data_acquisition(self):
        stocks = data.DataReader(self.ticker_symbol, 'yahoo', self.start_date, self.end_date)
        stocks = stocks.dropna()
        self.time_series = stocks['Close']


    def calculate_periodic_daily_return(self):
        self.data= np.log(self.time_series).diff().dropna()

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
                #Next Day’s Price=Today’s Price × e^(Drift+Random Value)
                prediction.append(prediction[-1] * pow(math.e, (self.calculate_drift() + self.calculate_random_value())))

            predictions.append(copy.deepcopy(prediction))
            prediction.clear()
        return predictions

    def mcs_finance_serial(self, number_of_simulations):
        pass

    def mcs_finance_parallel(self, number_of_simulations):
        pass


if __name__ == "__main__":

    m = MonteCarloSimulationFinance('1980-01-01', '2019-12-31', 'AAPL', 5)
    m.data_acquisition()
    m.calculate_periodic_daily_return()
    print(m.simulation_finance(3, 5))




