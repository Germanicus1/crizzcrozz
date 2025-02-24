v0.0.4

# CrizzCrozz - Crossword Puzzle Generator

## Introduction

CrizzCrozz is an crossword puzzle generator implemented in Go. It allows users to create custom crossword puzzles using a list of words and hints supplied via a CSV file. This project includes features such as asymmetric puzzle generation and flexible board dimensions, catering to various complexity levels and user preferences.

## Features

- **Dynamic Board Sizing:** Configure the crossword grid size to your preference.
- **Custom Word Lists:** Utilize your own word lists through CSV input to generate puzzles.
- **Asymmetrical Puzzles:** Generates puzzles without symmetry constraints for a unique challenge.
- **Retry Logic:** Attempts to place words multiple times to optimize board space usage.

## Getting Started

### Prerequisites

- Go 1.15 or higher
- A CSV file with words and hints formatted as `word,hint`

### Installation

Clone the repository and build the application:

```bash
git clone https://github.com/yourusername/CrizzCrozz.git
cd CrizzCrozz
go build .
```

### Running the Application

To start generating your crossword puzzles, run:

```bash
./CrizzCrozz -width=15 -height=15 -f=path/to/your/words.csv
```

Replace `path/to/your/words.csv` with the path to your CSV file containing the words and hints.

## Usage

CrizzCrozz accepts several command-line options to customize the crossword generation:

- `-width` specifies the width of the crossword grid.
- `-height` specifies the height (defaults to the width if not provided).
- `-r` specifies the number of retries for placing a word.
- `-f` specifies the file path to the CSV containing words and hints.

### Example

Generate a 23x23 crossword puzzle and try to place 39 words:

```bash
Words to place: 39
Words placed: 35
Failed to generate the crossword: Generate: failed to generate a complete puzzle
. g e t r e n n t . v i e l . . . d . . k . .
. . . h . . . . . . . . . a . . . u . . o . .
. . h u n g e r . . . . . m . r . r . . s . .
. g . n . . . . w . . . a p f e l s a f t . .
. e . f . . . . i . . . . e . s . t . . e . .
. t r i n k g e l d . . . . . t . . . . n . .
. r . s . . . . l . . e i n l a d u n g . . .
. ä . c e n t . k . . . . . . u . . . . . . .
. n . h . . . . o . . . . b . r . . r . . . .
. k . . z u s a m m e n . e . a . . e . . . .
. . . . . . . . m . . . . z . n . . c . . . .
. . . . w . s p e i s e k a r t e . h a u s .
. . . . a . c . n . c . . h . . . . n . . . .
. . k . s . h . . . h . f l u g z e u g . a .
. b e i s p i e l . ö . . e . . a . n . . u .
. . l . e . n . . . n . . n . . h . g a s t .
. . l . r . k . . . . . . . n . l . . . . o .
. . n . . . e . . k . . w i e d e r . . . . .
. . e . . a n a n a s . . . h . n . . . . . .
p a r t y . . . . r . . . . m . . . . . . . .
. . . . . . . v a t e r . . e s s e n . . . .
. . . . . . . . . e . . . . n . . . . . . . .
. . . . . . . . . . . . . . . . . . . . . . .
```

```bash
./CrizzCrozz -width=10 -height=10 -f=path/to/your/words.csv
```

## ToDo

- Saving the result to a file (json) for further usage in for example a frontend.
- Building a frontend.
- Publish Crossword puzzles.

## License

This project is licensed under the MIT License.

## Acknowledgments

Inspired by VG Ventures IO_Crossword https://github.com/VGVentures/io_crossword

Enjoy creating your unique crossword puzzles with CrizzCrozz!
