/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random dad joke",
	Long:  `This command fetches a random dad joke from the icanhazdadjoke api`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getRandomJokeWithTerm((jokeTerm))
		} else {
			getRandomJoke()
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().String("term", "", "A search term for a dad joke.")
}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResult struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}
	err := json.Unmarshal(responseBytes, &joke)

	if err != nil {
		log.Println("Could not unmarshall joke: %v", err)
	}

	fmt.Println(string(joke.Joke))
}

func getRandomJokeWithTerm(jokeTerm string) {
	length, results := getJokeDataWithTerm(jokeTerm)
	randomizeJokeList(length, results)
	// fmt.Println(results)
}

func randomizeJokeList(length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())
	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("No jokes found with this term")
		fmt.Println(err.Error())
	} else {
		randomNum := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNum])
	}

}

func getJokeDataWithTerm(jokeTerm string) (totalJokes int, jokeList []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)

	jokeListRaw := SearchResult{}
	err := json.Unmarshal(responseBytes, &jokeListRaw)
	if err != nil {
		log.Println("Could not unmarshall joke with term: %v", err)
	}

	jokes := []Joke{}
	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		log.Println("Could not unmarshall joke list results: %v", err)

	}
	return jokeListRaw.TotalJokes, jokes
}

func getJokeData(baseAPIURL string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		baseAPIURL,
		nil,
	)
	if err != nil {
		log.Println("Could not request a dad joke: %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "DadJoke CLI (github.com/thetajwar2003/Dad-Joke-CLI)")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("Could not make a request: %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("Could not read response body: %v", err)
	}

	return responseBytes
}
