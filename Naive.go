package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strings"
    "unicode"
	"math"
)

type Assets struct{
	freqs 	map[string]int{}
	HamW 	map[string]int{}
	vocab  	map[string]bool{}
	spamDocs  int
	hamDocs  	int
	totalSpamWords  int
	totalHamWords  int
	spam  int
	ham  int
	spamEmails  []string{}
	hamEmails  []string{}
	

}

//TOKENIZATION THE STRINGS
func tokenize(text string) []string {
    text = strings.ToLower(text)
    return strings.FieldsFunc(text, func(r rune) bool {
        return !unicode.IsLetter(r)
    })
}

func Opencsv(file string) string{
	
	Filename, err := os.Open(file)
    if err != nil {
        log.Fatalf("cannot open file: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(FileName)
    records, err := reader.ReadAll()
    if err != nil {
        log.Fatalf("cannot read file: %v", err)
    }

	return records
}
func train(file string) (bool, string){
	Records := Opencsv(file)
	assets := new(Assets)
	var logHam float64
	var logSpam float64

	for r, row := range Records{
		if r == 0 {
			continue
		}

		labels := row[1]
		text := row[2]
		tokens := tokenize(text)

		//hitung spam dan ham berdasarkan label
			if labels == "spam" {
				assets.spamDocs++
				for w := range freqs{
					assets.SpamW[w]++
					assets.totalSpamWords++
					assets.vocab[w] = true
				}
				} else {
					assets.hamDocs++
					for h := range freqs {
						assets.HamW[h]++
						assets.totalHamWords++
						assets.vocab[h] = true
					}
				}
        
		//hitung probabilitas dari kata
		totalCount := float64(assets.spamDocs + assets.hamDocs)
		//bagi total kata dengan spam
		spam_train := float64(assets.spamDocs)/totalCount
		ham_train := float64(assets.hamDocs)/totalCount

		logSpam = math.Log(spam_train)	
		logHam = math.Log(ham_train)
		vocabsize := float64(len(assets.vocab))


		
		for _, token := range tokens {
            assets.freqs[token]++
			spamWord := assets.SpamW[token]
			//laplace smoothing menghindari probabilitas 0
			likehood_s := float64(spamWord + 1) / (float64(assets.totalSpamWords) + vocabsize)
			logSpam += math.Log(likehood_s)

			//fmt.Println(freqs)

			hamWord:= assets.HamW[token]
			likehood_h := float64(hamWord + 1) / (float64(assets.totalHamWords) + vocabsize)
			logHam += math.Log(likehood_h)
			
			fmt.Printf("v %s", token)

			fmt.Printf("ham probability  %f  spam probability %f \n", likehood_s, likehood_h)
				
	}
	
}

	return logSpam > logHam ? "spam" : "ham"
}

func main() {
    

	spamEmails := []string{}
	hamEmails := []string{}
	
    //output the text in csv file
    for r, row := range records {
        if r == 0 {
            continue
        }



		labels := row[1]
        text := row[2]
        tokens := tokenize(text)

	       
		
		
		//count if logspam bigger then logham 
		if logSpam > logHam {
			spam++
			spamEmails = append(spamEmails, text)
            fmt.Printf("\n %v spam email \n", r)

            fmt.Printf("this email is spam %s", text)
			
			} else {
				ham++

			hamEmails = append(hamEmails, text)
            fmt.Printf("\n %v ham email \n", r)

            fmt.Printf("this email  %s", text)
			}
		
	}
   
		fmt.Printf("ham %d \n", ham)
  
}

