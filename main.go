package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// Kui soovite veateateid kuvada samal viisil nagu edukaid tulemusi, saate kasutada sama malli ja lihtsalt muuta vastavat Result väärtust,
// et see sisaldaks veateadet. Selleks saate luua üldise struktuuri, mis võimaldab teil veateateid ja edukaid tulemusi ühtemoodi käsitleda.
type PageData struct {
	Result string
	Error  string
}

// Defines command-line flags with default values and descriptions.
// flag.Parse() // Processes the command-line flags and sets the corresponding variables.
var (
	helpFlag      = flag.Bool("h", false, "Enable help mode")
	helpFlagHelp  = flag.Bool("help", false, "Enable help mode")
	encodeFlag    = flag.Bool("e", false, "Enable encode mode")
	multiLineFlag = flag.Bool("m", false, "Enable multi-line mode")
	webFlag       = flag.Bool("w", false, "Enable web server mode")
)

func main() {
	log.Println()
	flag.Parse() // Processes the command-line flags and sets the corresponding variables.

	if *webFlag {
		startWebServer()
	} else {
		runCommandLine()
	}
}

func runCommandLine() {
	args := flag.Args() // Returns the command-line arguments that are not flags.

	// If the help flag is used, displays usage instructions and exits the program.
	if *helpFlag || *helpFlagHelp {
		displayTheUsage()
		return
	}

	var input string

	// Checks if an argument with .encoded.txt suffix is provided and reads its content.
	if len(args) == 1 && strings.HasSuffix(args[0], ".encoded.txt") {
		filePath := "./static/txt-files/" + args[0]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[41mError:\033[0m\n Error reading file \"%s\": %v\n", filePath, err)
			return
		}
		input = string(fileContent)
		result, err := decode(input, *multiLineFlag)
		if err != nil {
			fmt.Println("\033[41m Decode error: \033[0m\n", err)
			return
		}
		fmt.Println("Decoded result:")
		fmt.Println(result)
		return
	}

	// Checks if an argument with .art.txt suffix is provided and reads its content.
	if len(args) == 1 && strings.HasSuffix(args[0], ".art.txt") {
		filePath := "./static/txt-files/" + args[0]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[41mError:\033[0m\n Error reading file \"%s\": %v\n", filePath, err)
			return
		}
		input = string(fileContent)
		result, err := encode(input, *multiLineFlag)
		if err != nil {
			fmt.Println("Encode error:", err)
			return
		}
		fmt.Println("Encoded result:")
		fmt.Println(result)
		return
	}

	///////////////////////////////////////////////
	// If the encode flag is used and there is one argument provided.
	if *encodeFlag {
		if len(args) == 1 { // checks if one more argument is given besides flags, and for that it is args[0]
			result, err := encode(args[0], *multiLineFlag)
			if err != nil {
				fmt.Println("Encode error:", err)
				return
			}
			fmt.Println("Encoded result:")
			fmt.Println(result)
			return
		}
	}

	// Handles multi-line input according to the multi-line flag.
	if *multiLineFlag {
		var input string
		if *encodeFlag {
			// Encode the input
			fmt.Println("Enter multi-line input for encoding (Ctrl+D to finish):")
			input = handleMultiLineInput()
			result, err := encode(input, *encodeFlag)
			if err != nil {
				fmt.Println("Encode error:", err)
				return
			}
			fmt.Printf("Encoded result:\n")
			fmt.Println(result)
			return

		} else if len(args) == 0 {
			fmt.Println("Enter multi-line input for decoding (Ctrl+D to finish):")
			input = handleMultiLineInput()
			result, err := decode(input, *multiLineFlag)
			if err != nil {
				fmt.Println("\033[41m Decode error: \033[0m\n", err)
				return
			}
			fmt.Println("Decoded result:")
			fmt.Println(result)
			return

		} else {
			fmt.Println("\n\033[41mError:\033[0m\n Invalid usage with -m or -e flag.")
			// displayTheUsage()
			return
		}

	} else {
		if len(args) != 1 {
			fmt.Println("\n\033[41mError:\033[0m\nNo correct input provided or too many arguments.\nOr maybe you didn't used \"\"-s?\nUse \"go run main.go -h\" for help.")
			return
		}
		// If the -m flag is not used
		input = args[0]
		if len(args) == 1 {
			result, err := decode(input, *multiLineFlag)
			if err != nil {
				fmt.Println("\033[41m Decode error: \033[0m\n", err)
				return
			}
			fmt.Println("Decoded result:")
			fmt.Println(result)
			return

		} else {
			return
		}
	}
}

func startWebServer() {
	http.HandleFunc("/", serveHomepage)
	http.HandleFunc("/decoded-encoded", handleDecoderEncoder)

	// Seadista staatiliste failide marsruut
	fs := http.FileServer(http.Dir("static"))                 // Eeldab, et teil on "static" kataloog peakataloogi tasemel
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // Eemalda "/static" prefiks enne faili otsimist

	fmt.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleDecoderEncoder(w http.ResponseWriter, r *http.Request) {
	log.Println("Decoder/Encoder Started")

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing the form", http.StatusBadRequest)
		return
	}

	action := r.FormValue("buttonAction")
	inputString := r.FormValue("inputString")

	var pageData PageData
	var result string
	var err error

	tmpl, tmplErr := template.ParseFiles("template/decoded-encoded.html")
	if tmplErr != nil {
		log.Printf("Template error: %v", tmplErr) // Logi vea sõnum
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		return
	}

	switch action {
	case "decode":
		result, err = decode(inputString, true)
	case "encode":
		result, err = encode(inputString, true)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if inputString == "" {
		log.Println("HTTP400 Bad Request - No input inserted")
		errorInInput(w, "No input", tmpl)
		return
		// http.Error(w, "empty input", http.StatusBadRequest)
	}

	if err != nil {
		log.Println("HTTP400 Bad Request - Check README.md info for input information")
		errorInInput(w, "Wrong input - Check README.md info for input information", tmpl)
		pageData.Error = err.Error() // Vea sõnum
		return
	} else {
		log.Println("HTTP202 valid encoded strings")
		w.WriteHeader(http.StatusAccepted)
		pageData.Result = result // Edukas tulemus
	}

	// Käivitada mall pageData'ga, mis sisaldab nii tulemust kui ka võimalikku veateadet
	tmplErr = tmpl.Execute(w, pageData)
	if tmplErr != nil {
		// Kui malli käivitamisel tekib viga, logige see ja saatke kasutajale veateade
		log.Printf("Error executing template: %v", tmplErr)
		http.Error(w, "Error executing the template", http.StatusInternalServerError)
	}
}

func errorInInput(w http.ResponseWriter, errorMessage string, template *template.Template) {
	w.WriteHeader(http.StatusBadRequest)
	result := errorMessage

	p := PageData{
		Error: result,
	}
	if err := template.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveHomepage(w http.ResponseWriter, r *http.Request) {
	// Veenduge, et server ei serveeriks kogemata kataloogi sisu
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func decode(input string, multiLine bool) (string, error) {
	var resultBuilder strings.Builder
	if !isBracketsBalanced(input) {
		return "", fmt.Errorf("error: Square brackets are unbalanced")
	}

	if strings.Contains(input, "[]") {
		return "", fmt.Errorf("error: Square brackets are empty")
	}

	if multiLine {
		lines := strings.Split(input, "\n")
		for _, line := range lines {
			decodedLine, err := decodeString(line)
			if err != nil { // Kontrolli, kas `err` on mitte-nil, mis tähendab vea esinemist
				return "", fmt.Errorf("multiline decoding failed: %v", err)
			}
			resultBuilder.WriteString(decodedLine + "\n")
		}
	} else {
		decodedString, err := decodeString(input)
		if err != nil {
			return "", fmt.Errorf("decoding failed - check arguments between brackets [ ]: %v", err)
		}
		resultBuilder.WriteString(decodedString)
	}
	return resultBuilder.String(), nil
}

// Logic for decoding, using regular expressions and checking for balanced brackets.
// Additional checks and the implementation of the decoding process.
func decodeString(input string) (string, error) {
	// Implement decoding logic
	pattern := regexp.MustCompile(`\[([^\]]+)\]|([^[]+)`)
	matches := pattern.FindAllStringSubmatch(input, -1)

	var result string
	for _, match := range matches {
		if match[1] != "" {
			arguments := strings.SplitN(match[1], " ", 2)
			if len(arguments) != 2 {
				return "", fmt.Errorf("error: Wrong amout of arguments between square brackets. Use \"go run main.go -h\" for help")
			}

			number, err := strconv.Atoi(arguments[0])
			if err != nil || arguments[1] == "" {
				return "", fmt.Errorf("error: invalid number or missing argument in brackets. Use \"go run main.go -h\" for help")
			}

			result += strings.Repeat(arguments[1], number)

		} else if match[2] != "" {
			result += match[2]
		}
	}
	return result, nil
}

// Similar logic to the decode function but for encoding.
// Uses strings.Builder to create the encoded string.
func encode(input string, multiLine bool) (string, error) {
	var resultBuilder strings.Builder

	if multiLine {
		lines := strings.Split(input, "\n")
		for _, line := range lines {
			encodedLine, err := encodeString(line)
			if err != nil {
				return "", fmt.Errorf("multiline encoding failed")
			}
			resultBuilder.WriteString(encodedLine + "\n")
		}
	} else {
		encodedString, err := encodeString(input)
		if err != nil {
			return "", fmt.Errorf("encoding failed")
		}
		resultBuilder.WriteString(encodedString)
	}
	return resultBuilder.String(), nil
}

// Logic for encoding, handling character repetition and encoding accordingly.
func encodeString(input string) (string, error) {
	var encodedBuilder strings.Builder // Creates a new string builder for constructing the encoded string.
	i := 0                             // Initialize the index for traversing the input string.

	for i < len(input) {
		count := 1 // Initialize the count of consecutive characters.
		// Check if we are at the end or if the next symbol is the same as the current one.
		for i+count < len(input) && input[i] == input[i+count] {
			count++ // Increment count if the next character is the same.
		}
		if count > 1 {
			// If there are more than one of the same character in sequence, encode them.
			encodedBuilder.WriteString(fmt.Sprintf("[%d %c]", count, input[i])) // Add the encoded sequence to the builder.
			i += count                                                          // Move the index forward by the count.
			continue                                                            // Skip to the next iteration of the loop.
		}

		// Check for a two-character pattern only if there are enough characters following it.
		if i+1 < len(input) {
			nextPatternCount := 1 // Initialize the count for the two-character pattern.
			// Loop to find if the next two characters form a repeating pattern.
			for i+nextPatternCount*2 < len(input) && i+nextPatternCount*2+1 < len(input) && input[i] == input[i+nextPatternCount*2] && input[i+1] == input[i+nextPatternCount*2+1] {
				nextPatternCount++ // Increment pattern count if a repeating pattern is found.
			}
			if nextPatternCount > 1 {
				// If a repeating pattern is found, encode it.
				encodedBuilder.WriteString(fmt.Sprintf("[%d %s]", nextPatternCount, input[i:i+2])) // Add the encoded pattern to the builder.
				i += nextPatternCount * 2                                                          // Move the index forward by the pattern length times the count.
				continue                                                                           // Skip to the next iteration of the loop.
			}
		}

		// If no condition is met, simply add the character to the result.
		encodedBuilder.WriteString(fmt.Sprintf("%c", input[i])) // Add the single character to the builder.
		i++                                                     // Move the index forward by one.
	}
	return encodedBuilder.String(), nil // Return the encoded string or error
}

// Reads multi-line input from standard input using bufio.Scanner.
func handleMultiLineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// Checks if the square brackets in the input are balanced.
func isBracketsBalanced(input string) bool {
	var brackets int

	for _, character := range input {
		switch character {
		case '[':
			brackets++
		case ']':
			brackets--
			if brackets < 0 {
				return false
			}
		}
	}
	return brackets == 0
}

// Displays usage instructions, including examples of how to use the program.
func displayTheUsage() {
	fmt.Println("\n\033[41m Usage instructions: \033[0m")
	fmt.Println()

	fmt.Println("\033[45mFor decoding\033[0m")

	fmt.Println("\033[35mTo enable web server mode use \"go run main.go -w\" \033[0m")

	fmt.Println()
	fmt.Println("\033[35mFor single line decoding:          Follow this patter => go run main.go \"[\033[34m[number]\033[35m[single space]\033[34m[character(s)]\033[35m][same logic as in previous brackets][etc.]]\" \033[0m")
	fmt.Println("\033[35m             for example:          go run main.go \"[5 #][5 -_]-[5 #]\" \033[0m")
	fmt.Println("\033[34mFor decoding from file:            use file with the end \033[35m\".encoded.txt\"\033[34m. Example: go run main.go cats.encoded.txt\033[0m")
	fmt.Println("\033[35mFor multiline decoding from input: type \"go run main.go -m\" \033[0m")
	fmt.Println("\033[35mand into the next lines insert coded lines you want to decode.\033[0m")
	fmt.Println("\033[35mfor example:                       \n[5 |\\---/|]\n[5 | o_o |]\n[5  \\_^_/ ]\033[0m")
	fmt.Println("\033[45m\033[1m NB! After completing the multi-line input in the terminal, please push \"enter\" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. \033[0m\033[22m")

	fmt.Println()
	fmt.Println("\033[44mFor encoding\033[0m")
	fmt.Println("\033[34mFor single line encoding:          add \"-e\" after main.go (For example: go run main.go -e \"[pattern_you_wish_to_encode]\") \033[0m")
	fmt.Println("\033[34m             for example:          go run main.go -e \"#####-_-_-_-_-_-#####\" \033[0m")
	fmt.Println("\033[35mFor decoding from file:            use file with the end \033[34m\".art.txt\"\033[35m. For example: go run main.go cats.art.txt \033[0m")
	fmt.Println("\033[34mFor multiline encoding from input: add \"-m\" & \"-e\" (example: go run main.go -m -e)\033[0m")
	fmt.Println("\033[34mand next lines enter for example:  \n" +
		"          \n" +
		"   *   *  \n" +
		"  *** *** \n" +
		"  ******* \n" +
		"   *****  \n" +
		"    ***   \n" +
		"     *    \n\033[0m")
	fmt.Println("\033[44m\033[1m NB! After completing the multi-line input in the terminal, please push \"enter\" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. \033[0m\033[22m")
}
