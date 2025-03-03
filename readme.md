v0.1.0

# CrizzCrozz - Crossword Puzzle Generator

## Introduction

CrizzCrozz is a **crossword puzzle generator** implemented in Go. It allows users to create custom crossword puzzles using a list of words and hints supplied via a CSV file. The generator features **asymmetric puzzle generation, automatic board sizing, and optimized placement strategies**, ensuring the best possible fit for the given word list.

## Features

- **Automatic Board Resizing**
  Dynamically finds the smallest board size that fits all words.
- **Custom Word Lists**
  Import your own word lists via CSV input.
- **Asymmetrical Puzzles**
  No symmetry constraints, allowing for unique crossword layouts.
- **Best Fit Optimization**
  If all words can't fit, CrizzCrozz will return the best possible solution.

## Getting Started

### Prerequisites

- Go 1.15 or higher
- A CSV file with words and hints formatted as:

```csv
word,hint
```

### Installation

Clone the repository and build the application:

```bash
git clone https://github.com/yourusername/CrizzCrozz.git
cd CrizzCrozz
go build .
```

### Running the Application

To generate a crossword puzzle, run:

```bash
./CrizzCrozz -f=path/to/your/words.csv
```

### Example Runs

#### **Generate a crossword with auto-sized board**

```bash
❯ go run . -f="vocabulary.csv"

Successfully unmarshalled 39 words and hints.
Trying board size: 16x16
Trying board size: 20x20
Trying board size: 18x18
Trying board size: 17x17
Optimal board size found: 18x18
Best Solution Found:
W O H E R . . . A P F E L S A F T . P
. . . S . . W . . . . . . P . . . . A
H A U S . . I . . . . . K E L L N E R
. . . E I N L A D U N G . I . . . . T
V . . N . . L . . . . . . S . G . . Y
I . . . . . K . . N . R . E . A . C .
E . . . B . O . . E . E . K O S T E N
L A M P E . M . . H . S . A . T . N .
. . . . I . M . . M . T . R . . . T .
. . A . S P E I S E K A R T E . B . .
. . N . P . N . . N . U . E . . E . .
. . A . I . . . G . . R . . Z . Z . .
H U N G E R . . E . . A . . U . A . .
. . A . L . . . T H U N F I S C H . .
. . S . . . . . R . . T . . A . L . .
. . . F L U G Z E U G . . . M . E . .
. . . . . . . . N . . . . . M . N . .
. . . . D U Z E N . . . . . E . . . .
. . . . . . . . T . S C H I N K E N .

Words placed: 39
```

## ToDo

- **Save the result to a file (JSON)** for usage in a frontend.
- **Build a frontend** to display generated crosswords.
- **Publish crossword puzzles**.

## License

This project is licensed under the MIT License.

## Acknowledgments

Inspired by VG Ventures IO_Crossword https://github.com/VGVentures/io_crossword

---

Enjoy creating your unique crossword puzzles with **CrizzCrozz**!
