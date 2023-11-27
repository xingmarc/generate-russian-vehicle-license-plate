package main

import (
    "os"
    "fmt"
    "bufio"
    "unicode/utf8"
    "sort"
)

// As of January 2014, in very rare cases, California has extended custom license plates to allow more than seven digits, but not to exceed nine characters. However, most plates are limited to seven-and-a-half characters (the half-character is a half-space).
const CaliforniaLiscenseLength int = 7

// UNUSED For debugging purpose: allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func main() {
    helpMsg := `Usage:

    go run main.go

    -h or --help          print this message
    --includeWeird        use those not exactly matching mappings. e.g. ч as 4
    --california          use maxLength=7, otherwise the maxLength is 8
    -f <InputFilePath>    provide the input file path
    -o <OutputFileName>   provide the outout file name. Default is output.txt`
    includeWeird := false
    maxLength := 8
    minLength := 3
    fileLocation := ""
    outputFileName := "output.txt"

    argsWithoutProg := os.Args[1:]

    i := 0
    for ; i < len(argsWithoutProg); {

        arg := argsWithoutProg[i]

        if arg == "--help" || arg == "-h" {
            fmt.Println(helpMsg)
            os.Exit(0)
        }

        if arg == "--includeWeird" {
            includeWeird = true
            i++
        }

        if arg == "--california" {
            maxLength = CaliforniaLiscenseLength
            i++
        }

        if arg == "-f" {
            if i + 1 >= len(argsWithoutProg) {
                panic("Invalid input: -f needs to be followed by inputFilePath")
            }
            fileLocation = argsWithoutProg[i + 1]
            i += 2
        }

        if arg == "-o" {
            if i + 1 >= len(argsWithoutProg) {
                panic("Invalid input, -o needs to be followed by outputFileName")
            }
            outputFileName = argsWithoutProg[i + 1]
            i += 2
        }
    }

    russianLetters := []string{
        // twenty consonants
        "б", "в", "г", "д", "ж", "з", "к", "л", "м", "н", "п", "р", "с", "т", "ф", "х", "ц", "ч", "ш", "щ",
        // ten vowels
        "а", "е", "ё", "и", "о", "у", "ы", "э", "ю", "я",
        // a semivowel / consonant
        "й",
        // and two modifier letters or "signs"
        "ъ", "ь",
    }
    russianLettersMap := map[string]struct{}{}

    for _, russianLetter := range russianLetters {
        russianLettersMap[russianLetter] = struct{}{}
    }

    /* Things with a little bit weird */
    convertMapButALittleBitWeird := map[string]string{
        "б": "6", /* Weirdness: This will make upper and Lower-case mix together.*/
        "ч": "4", /* Weirdness: Not very similar */
        "д": "D", /* Weirdness: д's uppercase's handwriting is similar to D */
        "ш": "W", /* Weirdness: Not very similar */
        "и": "N", /* Weirdness: Not very similar */
    }

    /* 99.9999% Identical */
    convertMap := map[string]string{
        "в": "B",
        "з": "3",
        "к": "K",
        "м": "M",
        "н": "H",
        "р": "P",
        "с": "C",
        "т": "T",
        "х": "X",
        "а": "A",
        "е": "E",
        "о": "O",
        "у": "Y",
    }

    if includeWeird {
        // combine two maps
        for k, v := range convertMapButALittleBitWeird {
            convertMap[k] = v
        }
    }

    // Return empty string if it is not convertable
    convert := func(word string) string {
        if utf8.RuneCountInString(word) < minLength || utf8.RuneCountInString(word) > maxLength {
            return ""
        }
        result := ""
        for _, s := range word {
            _, isValidLetter := russianLettersMap[string(s)]
            if !isValidLetter {
                fmt.Println("Error in parsing this word: " + word)
                return ""
            }
            val, ok := convertMap[string(s)]
            if !ok {
                return ""
            }
            result += val
        }
        return result
    }

    file, err := os.Open(fileLocation)
    // Writ to file
    out, outFileErr := os.Create(outputFileName)
    if err != nil {
        panic(err)
    }
    if outFileErr != nil {
        panic(outFileErr)
    }
    defer file.Close()
    defer out.Close()

    outputStringList := []string{}

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        s := scanner.Text()
        converted := convert(s)
        if converted != "" {
            outputStringList = append(outputStringList, s + " " + converted + "\n")
            fmt.Println(s, converted)
        }
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
    sort.Strings(outputStringList)

    for _, row := range outputStringList {
        _, writeStringErr := out.WriteString(row)
        if writeStringErr != nil {
            panic(writeStringErr)
        }
    }
}
