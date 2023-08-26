# E-Commerce Campaign Tool

This command line tool simulates and manages campaigns for price manipulation based on demand in an e-commerce platform. It enables the creation of products, orders, and campaigns, and facilitates time-based simulation of their behavior.

## Table of Contents

- [E-Commerce Campaign Tool](#e-commerce-campaign-tool)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Project Structure](#project-structure)
  - [How to Run](#how-to-run)
  - [Testing](#testing)
  - [Example Usage](#example-usage)
  - [Contributing](#contributing)

## Overview

This tool is designed for an e-commerce platform to handle products, orders, and campaigns. It allows the creation of products with product codes, prices, and stock levels. Orders can be placed with product codes and quantities, and campaigns can be created with names, product codes, durations, price manipulation limits, and target sales counts.

Campaigns start after creation and last for a specified duration in hours. The tool also supports time simulation by allowing the user to increase time in hourly increments. Price manipulation within the specified limit is possible to influence demand. The ultimate goal is to reach the target sales count during the campaign duration.

## Project Structure

The project follows this folder structure:

- `app`: Contains the application's main logic.
- `cmd`: Entry point of the application.
- `domain`: Defines the domain-specific logic and repositories.
  - `campaign`: Handles campaign-related logic.
  - `order`: Manages order-related logic.
  - `product`: Contains product-related logic.
- `entity`: Defines the core entity structs for campaigns, orders, and products.
- `mock`: Provides mock implementations.
- `pkg`: Contains utility packages, such as in memory storage.
- `service`: Implements business logic for campaigns, orders, and products.
- `types`: Defines common type definitions used throughout the application.
- `valueobject`: Contains value objects for various attributes, like price and quantity.

## How to Run

Follow these steps to run the command line tool:

1. Clone this repository to your local machine.
2. Navigate to the project directory.
3. Make sure you have Go installed on your system.
4. Install dependencies `go mod download`
5. Run tests.
    To run only unit tests `go test ./...`         
    To run with integration tests `go test ./... -tags=integration`
6. Open a terminal and run the following command to build the tool and run binary output:

   ```sh
   go build ./cmd/
   ```
7. Or you can just run the following command to run without build
   ```sh
   go run ./cmd/
   ```
8. You can also run the tool with a scenario file. The tool will run the commands in the scenario file and print the output to the console. To run with a scenario file, run the following command:

   ```sh
   go run ./cmd/ --file <path-to-scenario-file>
   ```

## Testing

Run tests to ensure the tool's functionality:

1. Unit tests can be run with the following command:

   ```sh
    go test ./...
    ```
2. Integration tests can be run with the following command: 

   ```sh
    go test ./... -tags=integration
    ```

   
## Example Usage

Here's an example of how to use the tool with scenario inputs:

|Steps in Example Inputs |Output|
| :- | :-: |
|create\_product ABC 100 100|Product created; code ABC, price 100, stock 100|
|create\_campaign C1 ABC 5 20 50|Campaign created; name C1, product ABC, duration 10, limit 20, target sales count 100|
|create\_order ABC 10|Order created; product ABC, quantity 10|
|increase\_time 1|Time is 01:00|
|get\_product\_info ABC|Product ABC info; price 120, stock 90|
|get\_product\_info ABC|Product ABC info; price 120, stock 90|
|get\_product\_info ABC|Product ABC info; price 120, stock 90|
|increase\_time 1|Time is 02:00|
|get\_product\_info ABC|Product ABC info; price 110.8, stock 90|
|increase\_time 1|Time is 03:00|
|get\_product\_info ABC|Product ABC info; price 108.6, stock 90|
|increase\_time 1|Time is 04:00|
|get\_product\_info ABC|Product ABC info; price 106.7, stock 90|
|increase\_time 2|Time is 06:00|
|get\_product\_info ABC|Product ABC info; price 100.0, stock 90|
|get\_campaign\_info C1|Campaign C1 info; Status Ended, Target Sales 50, Total Sales 10, Turnover 1000.0, Average Item Price 100.0|



## Contributing

You can contribute to the E-Commerce Campaign Tool project and add new features or improve existing ones. If any existing repository or service has changed run `go generate ./...`. If a repository or service was created that needs a mock file, define `mockgen` command in the abstraction level and run `go generate ./...`. After all make sure to run test commands `go test ./...` for unit tests `go test ./... -tags=integration` for integration tests.
