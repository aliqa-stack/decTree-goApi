
package main

import (
    "encoding/csv"
    "fmt"
    _"log"
    "math"
    "os"
    "strings"
    "unicode"
)

type Assets struct {
    freqs map[string]int
    HamW  map[string]int
    SpamW map[string]int
    vocab map[string]bool

    spamDocs       int
    hamDocs        int
    totalSpamWords int
    totalHamWords int
}

// TOKENIZATION
func tokenize(text string) []string {
    text = strings.ToLower(text)

    return strings.FieldsFunc(text, func(r rune) bool {
        return !unicode.IsLetter(r)
    })
}

// OPEN CSV
func OpenCSV(file string) ([][]string, error) {
    filename, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer filename.Close()

    reader := csv.NewReader(filename)

    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    return records, nil
}

// TRAINING
func train(file string) (*Assets, error) {
    records, err := OpenCSV(file)
    if err != nil {
        return nil, err
    }

    assets := &Assets{
        freqs: make(map[string]int),
        HamW:  make(map[string]int),
        SpamW: make(map[string]int),
        vocab: make(map[string]bool),
    }

    for r, row := range records {
        if r == 0 {
            continue
        }

        if len(row) < 3 {
            continue
        }

        label := strings.ToLower(strings.TrimSpace(row[1]))
        text := row[2]

        tokens := tokenize(text)

        if label != "spam" && label != "ham" {
            continue
        }

        if label == "spam" {
            assets.spamDocs++

            for _, word := range tokens {
                assets.SpamW[word]++
                assets.totalSpamWords++
                assets.vocab[word] = true
            }

        } else {
            assets.hamDocs++

            for _, word := range tokens {
                assets.HamW[word]++
                assets.totalHamWords++
                assets.vocab[word] = true
            }
        }
    }

    if assets.spamDocs + assets.hamDocs == 0 {
        return nil, fmt.Errorf("no valid training data")
    }

    return assets, nil
}

// PREDICTION
func (a *Assets) PredictMsg(text string) string {
    tokens := tokenize(text)

    totalDocs := float64(a.spamDocs + a.hamDocs)

    logSpam := math.Log(float64(a.spamDocs) / totalDocs)
    logHam := math.Log(float64(a.hamDocs) / totalDocs)

    vocabSize := float64(len(a.vocab))

    for _, token := range tokens {

        if !a.vocab[token] {
            continue
        }

        spamWord := float64(a.SpamW[token])

        likelihoodSpam := (spamWord + 1) /
            (float64(a.totalSpamWords) + vocabSize)

        logSpam += math.Log(likelihoodSpam)

        hamWord := float64(a.HamW[token])

        likelihoodHam := (hamWord + 1) /
            (float64(a.totalHamWords) + vocabSize)

        logHam += math.Log(likelihoodHam)
    }

    if logSpam > logHam {
        return "spam"
    }

    return "ham"
}

//func main() {
    //model, err := train("email_dataset.csv")
    //if err != nil {
        //log.Fatal(err)
    //}

    //email := "Congratulations you won a free prize"

    //prediction := model.Predict(email)

    //fmt.Println("Prediction:", prediction)
//}
