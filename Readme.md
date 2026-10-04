# Advent of Code Solutions in Go

This repository contains my solutions to the [Advent of Code](https://adventofcode.com/) programming challenges, written in **Go**.

The main goal of this repository is to practice problem solving, improve my understanding of Go, and develop solutions that are **correct, readable, maintainable, and reasonably optimized**.

The solutions are not necessarily the most optimal implementations possible. However, I aim to make them efficient enough while keeping the code understandable and avoiding unnecessary complexity.

## Repository Structure

Each Advent of Code year has its own directory, and each day has its own folder containing the solution for that day's puzzle.

```text
.
├── 2015/
│   ├── day01/
│   ├── day02/
│   ├── day03/
│   ├── ...
│   └── run-all.ps1
├── 2016/
│   ├── day01/
│   ├── day02/
│   ├── ...
│   └── run-all.ps1
└── ...
```

The `run-all.ps1` script is located inside each year's directory and is used to run and benchmark all available solutions for that year.

## Requirements

* **Go 1.26.2**
* **PowerShell**

This repository uses **Go 1.26.2**. Make sure this version is installed before running the solutions.

Check your installed Go version with:

```powershell
go version
```

It should report a version beginning with:

```text
go1.26.2
```

## Running the Solutions

To run all solutions and benchmarks for a particular year, navigate to that year's directory and execute its `run-all.ps1` script.

For example, for **2015**:

```powershell
cd 2015
./run-all.ps1
```

The script runs the available solutions for that year and displays their benchmark results along with the answers for Part 1 and Part 2.

Example output:

```text
━━━ day01 ━━━
Total: 543611.20 microseconds
Average: 54.3611 microseconds
Total time: 543611200ns
part1 995
part2 5847
✔ day01 finished in 1607ms (outer)
```

### Running a Specific Day

To run a specific day's solution directly, navigate to its directory.

For example:

```powershell
cd 2015/day01
go run .
```

## Benchmarking

Each year's `run-all.ps1` script also benchmarks the solutions for that year.

The benchmark output provides information such as:

* Total benchmark time
* Average execution time
* Total execution time
* Part 1 answer
* Part 2 answer
* Overall execution time for the day

Benchmarking is used to measure the performance of the solutions and to help identify opportunities for improvement.

Performance is an important consideration when developing these solutions. The goal is not necessarily to find the absolute fastest implementation, but to write solutions that use appropriate algorithms and data structures and are **optimized enough to solve the problems efficiently without sacrificing readability or introducing unnecessary complexity**.

## Approach

For each problem, I aim to balance three things:

1. **Correctness** — The solution should produce the correct answers.
2. **Readability** — The implementation should be clear and easy to understand.
3. **Performance** — The solution should use an appropriate and reasonably efficient approach.

This means a solution may not be the theoretical or practical fastest possible implementation, but it should avoid obviously inefficient approaches and be optimized enough for the problem's requirements.

## Notes

These are my personal solutions and are primarily intended for learning and experimentation.

Different problems may use different approaches depending on their requirements. Benchmark results are included to help compare implementations and understand how algorithmic choices affect performance.

## About Advent of Code

[Advent of Code](https://adventofcode.com/) is an annual programming event created by **Eric Wastl**.

It features a series of programming puzzles released throughout December. Each puzzle consists of two parts and covers a variety of programming concepts, algorithms, data structures, mathematics, and problem-solving techniques.

You can find more information at [adventofcode.com](https://adventofcode.com/).
