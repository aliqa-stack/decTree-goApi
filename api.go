package main

import(
	_"fmt"
	"log"
	"net/http"
	"encoding/json"
)

type Branch struct {
	Feature float64 `json:"feature"`


}

type Responses struct {
	Prediction int `json:"prediction"`
}

type Message struct {
	Chat string `json:"chat"`
}

type MsgResponse struct {
	ChatResponses string `json:"chatresponses`
}

var R *Node

func main(){
	x := []float64{2,4,5,6,8,10,12,14,16,18,20,22,24,26,28,30}
	y := []int {0,0,0,0,1,1,1,1,1,1,1,1,1,1,1,1}
	NaiveModel, err := train("email_dataset.csv")
	if err != nil {
		log.Fatal(err)
	}

	R = TREE(y,x,0,3)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tree", makeTree)
	mux.HandleFunc("POST /Cekspam", Spamcek)
	log.Println("server running on 8080")
	log.Fatal(http.ListenAndServe(":8080",mux))
}

func makeTree(w http.ResponseWriter, r *http.Request){
	var branch Branch

	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		log.Println("error: ", err)
		return
	}

	result := Predict(R,branch.Feature)
	response := Responses{
		Prediction: result,
	}

	w.Header().Set("Conten-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func Spamcek(w http.ResponseWriter, r *http.Request){
	var msg Message
	
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "invalid message", http.StatusBadRequest)
		return
	}

	result := Predictmsg(msg.Chat)
	response := MsgResponse{
		ChatResponses: result
	}


	w.Header().Set("Conten-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)


}

