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


//TOKENIZATION THE STRINGS
func tokenize(text string) []string {
    text = strings.ToLower(text)
    return strings.FieldsFunc(text, func(r rune) bool {
        return !unicode.IsLetter(r)
    })
}

func main() {
    file, err := os.Open("email_dataset.csv")
    if err != nil {
        log.Fatalf("cannot open file: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        log.Fatalf("cannot read file: %v", err)
    }

    freqs := map[string]int{}
	SpamW :=  map[string]int{}
	HamW :=  map[string]int{}
	vocab := map[string]bool{}
	spamDocs := 0
	hamDocs := 0
	totalSpamWords := 0
	totalHamWords := 0
	spam := 0
	ham := 0
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

	        //hitung spam dan ham berdasarkan label
			if labels == "spam" {
				spamDocs++
				for w := range freqs{
					SpamW[w]++
					totalSpamWords++
					vocab[w] = true
				}
				} else {
					hamDocs++
					for h := range freqs {
						HamW[h]++
						totalHamWords++
						vocab[h] = true
					}
				}
        
		//hitung probabilitas dari kata
		totalCount := float64(spamDocs + hamDocs)
		//bagi total kata dengan spam
		spam_train := float64(spamDocs)/totalCount
		ham_train := float64(hamDocs)/totalCount

		logSpam := math.Log(spam_train)	
		logHam:= math.Log(ham_train)
		vocabsize := float64(len(vocab))


		
		for _, token := range tokens {
            freqs[token]++
			spamWord := SpamW[token]
			//laplace smoothing menghindari probabilitas 0
			likehood_s := float64(spamWord + 1) / (float64(totalSpamWords) + vocabsize)
			logSpam += math.Log(likehood_s)

			//fmt.Println(freqs)

			hamWord:= HamW[token]
			likehood_h := float64(hamWord + 1) / (float64(totalHamWords) + vocabsize)
			logHam += math.Log(likehood_h)
			
			fmt.Printf("v %s", token)

			fmt.Printf("ham probability  %f  spam probability %f \n", likehood_s, likehood_h)
				
	}

		
		
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

