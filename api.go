package main

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	_"io/ioutil"
	"os"
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

type Data struct {
	UserInput float64 `json:"feature"` 
	UserResponse int `json:"predict`
}

var R *Node

func main(){
	x := []float64{2,4,5,6,8,10,12,14,16,18,20,22,24,26,28,30}
	y := []int {0,0,0,0,1,1,1,1,1,1,1,1,1,1,1,1}
	_, err := train("email_dataset.csv")
	if err != nil {
		log.Fatal(err)
	}

	R = TREE(y,x,0,3)
	mux := http.NewServeMux()
	data:= Readfile("user.json")
//	if err != nil {
		//log.Println(err)
	//}
	fmt.Println(data)

	mux.HandleFunc("POST /tree", makeTree)
	mux.HandleFunc("POST /Cekspam", Spamcek)
	mux.HandleFunc("POST /Write", Userwrite)
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
	m := &Assets{}
	
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "invalid message", http.StatusBadRequest)
		return
	}

	result := m.PredictMsg(msg.Chat)
	response := MsgResponse{
		ChatResponses: result,
	}


	w.Header().Set("Conten-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)


}

func checkfile(file string)error{
	_, err := os.Stat(file)
	if os.IsNotExist(err){
		_, err := os.Create(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func Readfile(file string) (Data){
	err := checkfile(file)
	if err != nil {
		log.Println(err)
	}

	filename, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)

	}
	var data Data
	err = json.Unmarshal(filename, &data)
	  if err != nil {
		log.Println(err)
	}
	 return data
}

func Userwrite(w http.ResponseWriter, r *http.Request) {
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)

	jsonStr, err := json.MarshallIndent(data, "", " ")
	if err != nil {
		fmt.Println("failed encode json", err)
	}

	err = os.Writefile("user.json", jsonStr, 0644)
	if err != nil {
		fmt.Println("failed write the json", err)
	}

	w.Header().Set("Conten-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "data berhasil disimpan"}`))


}

