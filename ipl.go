package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//these variables are used for basic authorization
var (
	username = "Supriya"
	password = "Project@123"
)

//used to get properties related to score
type Score struct {
	Match   string `json:"match"`
	Runs    int    `json:"runs"`
	Wickets int    `json:"wickets"`
}

//main structure contains player and score properties
type Player struct {
	Name   string  `json:"name"`
	ID     int     `json:"id"`
	Team   string  `json:"team"`
	Scores []Score `json:"scores"`
}

//used in displayplayer struct
type OnlyPlayer struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Team string `json:"team"`
}

//used to display player details in get player api o/p format
type DisplayPlayer struct {
	Players []OnlyPlayer `json:"players"`
}

//used in display score struct
type OnlyScores struct {
	ID     int     `json:"id"`
	Scores []Score `json:"scores"`
}

//used to display score details in get score api o/p format
type displayScores struct {
	PlayerScores []OnlyScores `json:"playerscores"`
}

//used for display fantasy score of individual player

type FantasyScore struct {
	Name   string `json:"name"`
	FScore int    `json:"fantasyscore"`
}

//used for display capholders

type CapHolder struct {
	PurpleCap string `json:"purpleCap"`
	OrangeCap string `json:"orangeCap"`
}

//this slice contains players and their scores-act as db(main slice)
var playerservice []Player

//this slice helps for adding scores to main slice
var TempScoreData []Score

//this slice contains fantasy scores of individual players
var FantasyScores []FantasyScore

// used to add the player to the main slice
func postPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, p, ok := r.BasicAuth()
	if ok && u == username && p == password {
		var newplayerdata Player
		_ = json.NewDecoder(r.Body).Decode(&newplayerdata)
		// check for empty values and eliminate them in the records
		if newplayerdata.ID != 0 && newplayerdata.Name != "" {
			newplayerdata.Scores = nil
			playerservice = append(playerservice, newplayerdata)
			fmt.Println("Player details is successfully added to the db")
		} else {
			response := "Invalid ID or Name is entered. Please provide details correct details"
			http.Error(w, response, http.StatusMethodNotAllowed)
			fmt.Println("Player details is not successfully added to db. Please provide vaild ID/Name")
		}
	}
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		fmt.Println("Invalid Credentials")
	}
}

// used to add the player score to the main slice
func postPlayerScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, p, ok := r.BasicAuth()
	if ok && u == username && p == password {
		params := mux.Vars(r)
		var newMatchScore Score
		_ = json.NewDecoder(r.Body).Decode(&newMatchScore)
		for index, item := range playerservice {
			id, _ := strconv.Atoi(params["id"])
			if item.ID == id {
				playerservice = append(playerservice[:index], playerservice[index+1:]...)
				item.Scores = append(item.Scores, newMatchScore)
				playerservice = append(playerservice, item)
				fmt.Println("Player score is successfully added to the db")
				break
			}
		}
	}
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		fmt.Println("Invalid Credentials")
	}
}

// Used to get player details from the main slice
func getPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, p, ok := r.BasicAuth()
	if ok && u == username && p == password {
		var playerdetails DisplayPlayer
		var playerdetail OnlyPlayer
		for _, item := range playerservice {
			playerdetail.ID = item.ID
			playerdetail.Name = item.Name
			playerdetail.Team = item.Team
			playerdetails.Players = append(playerdetails.Players, playerdetail)
		}
		json.NewEncoder(w).Encode(playerdetails)
	}
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		fmt.Println("Invalid Credentials")
	}
}

// Used to get player details along with there scores from the main slice
func getPlayerScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, p, ok := r.BasicAuth()
	if ok && u == username && p == password {
		var tempScores displayScores
		var playerdetails OnlyScores
		for _, item := range playerservice {
			playerdetails.ID = item.ID
			playerdetails.Scores = item.Scores
			tempScores.PlayerScores = append(tempScores.PlayerScores, playerdetails)
		}
		json.NewEncoder(w).Encode(tempScores)
		tempScores.PlayerScores = nil
	}
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		fmt.Println("Invalid Credentials")
	}
}

// used to calculate fantasy score based on the individual player score in main slice
func fantasyScoreCal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, p, ok := r.BasicAuth()
	if ok && u == username && p == password {

		for _, item := range playerservice {
			var singlePFS FantasyScore
			singlePFS.Name = item.Name
			singlePFS.FScore = 0
			for _, temp := range item.Scores {
				if temp.Wickets > 0 {
					singlePFS.FScore = singlePFS.FScore + 10*temp.Wickets
				}
				if temp.Wickets > 5 {
					singlePFS.FScore = singlePFS.FScore + 50
				}
				if temp.Runs >= 30 {
					singlePFS.FScore = singlePFS.FScore + 20
				}
				if temp.Runs >= 50 {
					singlePFS.FScore = singlePFS.FScore + 50
				}
				if temp.Runs >= 100 {
					singlePFS.FScore = singlePFS.FScore + 100
				}
			}
			FantasyScores = append(FantasyScores, singlePFS)
		}
		json.NewEncoder(w).Encode(FantasyScores)
		FantasyScores = nil
	}
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		fmt.Println("Invalid Credentials")
	}
}

// used to calculate capholders bassed on the overall score of the player present in the main slice
func capHolders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, p, ok := r.BasicAuth()
	if ok && u == username && p == password {
		var capholderdetails CapHolder
		for _, item := range playerservice {
			var PerformanceCal Score
			var totalwickets int = 0
			var totalruns int = 0
			for _, temp := range item.Scores {
				totalwickets = totalwickets + temp.Wickets
				totalruns = totalruns + temp.Runs
			}
			PerformanceCal.Match = item.Name
			PerformanceCal.Wickets = totalwickets
			PerformanceCal.Runs = totalruns
			TempScoreData = append(TempScoreData, PerformanceCal)
		}
		var maxruns int = 0
		var maxwickets int = 0
		for _, item := range TempScoreData {
			if maxwickets < item.Wickets {
				maxwickets = item.Wickets
				capholderdetails.PurpleCap = item.Match
			}
			if maxruns < item.Runs {
				maxruns = item.Runs
				capholderdetails.OrangeCap = item.Match
			}
		}
		json.NewEncoder(w).Encode(capholderdetails)
		TempScoreData = nil
	}
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		fmt.Println("Invalid Credentials")
	}
}

//this is the main func that create apis
func main() {
	r := mux.NewRouter()
	//adding two players detail to the main slice for demo
	TempScoreData = append(TempScoreData, Score{Match: "1", Wickets: 2, Runs: 150})
	playerservice = append(playerservice, Player{
		ID:     1,
		Name:   "Virat",
		Team:   "RCB",
		Scores: TempScoreData,
	})
	TempScoreData = nil
	TempScoreData = append(TempScoreData, Score{Match: "1", Wickets: 3, Runs: 50})
	playerservice = append(playerservice, Player{
		ID:     7,
		Name:   "Dhoni",
		Team:   "CSK",
		Scores: TempScoreData,
	})
	TempScoreData = nil
	//creating apis as per business requirement
	r.HandleFunc("/player", postPlayer).Methods("POST")
	r.HandleFunc("/player/{id}/score", postPlayerScore).Methods("POST")
	r.HandleFunc("/players", getPlayers).Methods("GET")
	r.HandleFunc("/players/scores", getPlayerScore).Methods("GET")
	r.HandleFunc("/fantasy-scores", fantasyScoreCal).Methods("GET")
	r.HandleFunc("/cap-holders", capHolders).Methods("GET")

	fmt.Print("Supriya starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
