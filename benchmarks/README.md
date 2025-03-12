# 🚀 Quick Benchmarking

## Introduction
Welcome to our benchmarking project! Here, we are conducting performance tests to evaluate different HTTP frameworks in Go. The goal is to better understand how each framework performs under heavy loads and identify which one offers the best performance for different scenarios.

## Tools Used
To ensure accurate and reliable benchmarking, we use two powerful load-testing tools:

- **[k6](./k6/README.md)**: An open-source tool designed for performance and load testing, known for its ease of use and advanced metric analysis.
- **[Vegeta](./vegeta/README.md)**: A highly efficient and flexible load tester that allows custom attack configurations to measure server capacity.

## What Are We Measuring?
Our benchmarking evaluates essential metrics, including:
- **Total HTTP Requests**: The total number of processed requests.
- **Requests per Second**: The number of requests handled each second.
- **Average Response Time**: How long the server takes to respond to each request.
- **Data Received and Sent**: The total volume of data transmitted.
- **Error Rate**: The percentage of failed requests.

## Project Structure
The project is organized into directories, each containing its own tests and configurations:

📁 **k6/** - Load tests using k6. See the [k6 README](./k6/README.md) for details.

📁 **vegeta/** - Load tests using Vegeta. See the [Vegeta README](./vegeta/README.md) for details.

## How to Run the Tests?
Each directory contains detailed instructions on how to execute the tests. Simply navigate to the corresponding directory and follow the guidelines.

## Contribution
If you want to contribute with new tests or improvements, feel free to open a PR or share your ideas! 🚀

---
Let's discover together which framework delivers the best performance! 🎯

