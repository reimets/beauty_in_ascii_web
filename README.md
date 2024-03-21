# beauty_in_ascii_web: README.md

## Overview

This project provides a command-line utility or web based interface written in Go for encoding and decoding ASCII art. It supports both single-line and multi-line modes for various operations.

## Features

- Encode ASCII art into a custom format.
- Decode text from the custom format to ASCII art form.
- Support for multi-line encoding and decoding stright from the web interface.
- Support for multi-line encoding and decoding stright from the terminal.
- If using terminal, you can read input directly from text files.

## Getting Started

### Prerequisites

- Go installed on your machine. Visit [Go's official site](https://golang.org/dl/) for download and installation instructions.

### Installation

Clone the repository to your local machine:

```bash
git clone https://gitea.kood.tech/reigoreimets/art.git
```

Navigate into the project directory:
```bash
cd art
```

## Requirements

Go (Golang) installed on your system.

## Usage

Files with suffix .encoded.txt have coded ASCII art pictures what you can transform into the ASCII art.
Files with suffix .art.txt have ASCII art pictures what you can transform into the coded version. 

You can find those files under catalog "./static/txt-files/"

Run the program using the Go command line:

```bash
go run main.go [options] [input]
```

### Getting usage instructions using -h or -help: 

```bash
go run main.go -h
```
or
```bash
go run main.go -help
```

### Usage instructions: 

#### To enable web server mode 
use \"go run main.go -w\"
Then go to http://localhost:8000

##### To decode from web interface:
Follow this patter => "[[number][single space][character(s)]]
For example: 
```bash
[5 #][5 -_]-[5 #]
```
And click "Decode" button

##### To encode from web interface:
Insert pattern you wish to encode
For example: 
``` bash
   *   *  
  *** *** 
  ******* 
   *****  
    ***   
     *  
```   
And click "Encode" button

Uou can see HTTP Status codes from terminal or use browser integrated developer tool toolkit (Using usually: F12 on your keyboard).


#### Using terminal
For decoding from the terminal
  For single line decoding:          Follow this patter => go run main.go "[[number][single space][character(s)]][same logic as in previous brackets][etc.]]" 
For decoding from file:            use file with the end ".encoded.txt".
For multiline decoding from input: type "go run main.go -m" and into the next lines insert coded lines you want to decode.

For encoding from the terminal:
For single line encoding:          add "-e" after main.go (For example: go run main.go -e "[pattern_you_wish_to_encode]") 
For decoding from file:            use file with the end ".art.txt". For example: go run main.go cats.art.txt 
For multiline encoding from input: add "-m" & "-e" (example: go run main.go -m -e)
and next lines enter for example:  
          
``` bash
   *   *  
  *** *** 
  ******* 
   *****  
    ***   
     *  
```   

 NB! After completing the multi-line input in the terminal, please push "enter" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. 

#### Example Commands for Decoding

##### For single line decoding from command line:
```bash
go run main.go "[5 #][5 -_]-[5 #]"
```
##### Decoding from file:
```bash
go run main.go cats.encoded.txt
```
##### Multi-line decoding from command line 
```bash
go run main.go -m
```
push "enter" and into the next lines insert coded lines you want to decode.
for example:                       
[5 |\---/|]
[5 | o_o |]
[5  \_^_/ ]
 NB! After completing the multi-line input in the terminal, please push "enter" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. 

#### Example Commands for Encoding

##### For single line encoding from command line:
```bash
go run main.go -e "#####-_-_-_-_-_-#####"
```
##### Encoding from file:
```bash
go run main.go cats.art.txt
```
##### Multi-line encoding from command line 
```bash
go run main.go -m -e
```
and into the next lines insert coded lines you want to encode.
for example:                       
``` bash
   *   *  
  *** *** 
  ******* 
   *****  
    ***   
     *  
```

 NB! After completing the multi-line input in the terminal, please push "enter" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. 

## Bonus Features
- Light/Dark mode button to change mode
- Has both Encode and Decode functions
- Has additional links to instagram, FB and LinkedIn webpages
- Has "Home" button 

## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Authors
Reigo Reimets
