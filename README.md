# Timescale DB - Tutorials

This repository contains tutorials and use cases examples on how to setup and working with time series data using [timescale database](https://docs.timescale.com/).

## Folder Structure :books:

- Each folder represents a specifc tutorial and have it's own environment.

### Get Started

The *[get-started](get-started)* folder contains a brief introduction on how to setup timescaledb `dockerized` environment showing up some features with a simple example with just SQL stuff.

### New York City Taxi Cab

The *[nyc-taxi](nyc-taxi)* example has historical data from New York's yellow taxi network provided by [NYC TLC](https://www1.nyc.gov/site/tlc/about/tlc-trip-record-data.page) which is tracked over 200.000 vehicles making about 1 milion trips each day.

### Cryptocurrency - Bitcoin blockchain

The *[cryptocurrency](cryptocurrency)* example has a updated dataset containing data from the last five days of bitcoin transactions. Including information about each transaction, also states if is the frist transaction in a block, a.k.a `coinbase transaction`, which includes the reward a coin miner receives for mining the coin.