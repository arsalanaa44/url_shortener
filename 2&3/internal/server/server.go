package server 

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"


	"github.com/mehrdad3301/url-shortner/internal/config"
	"github.com/mehrdad3301/url-shortner/internal/redis"
)

var ( 
	redisCli *redis.RedisClient
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	LongURL string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
	IsCached bool `json:"isCached"`
	HostName string `json:"hostName"`
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {

	

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse the request JSON
	var req ShortenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Failed to parse request JSON", http.StatusBadRequest)
		return
	}

	finalRes := ShortenResponse{
		LongURL: req.URL, 
		HostName: r.Host, 
	}

	if short, err := redisCli.Get(req.URL) ; err != nil { 

		fmt.Println("URL not found in Redis cache")

		// Prepare the request payload
		payload := struct {
			Input string `json:"input"`
		}{
			Input: req.URL,
		}

		// Convert the payload to JSON
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "Failed to marshal request payload", http.StatusInternalServerError)
			return
		}

		// Send a POST request to gotiny.cc
		resp, err := http.Post("https://gotiny.cc/api", "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			http.Error(w, "Failed to call gotiny.cc API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read gotiny.cc response", http.StatusInternalServerError)
			return
		}

		type ShortenerResponse struct {
			Long     string `json:"long"`
			Code     string `json:"code"`
		}

		// Parse the response JSON
		var shortenerResp []ShortenerResponse
		err = json.Unmarshal(respBody, &shortenerResp)
		if err != nil {
			http.Error(w, "Failed to parse gotiny.cc response JSON", http.StatusInternalServerError)
			return
		}

		if len(shortenerResp) == 0 {
			http.Error(w, "Empty gotiny.cc response", http.StatusInternalServerError)
			return
		}
		finalRes.ShortURL= "https://gotiny.cc/" + shortenerResp[0].Code
		if rerr := redisCli.Set(finalRes.LongURL, finalRes.ShortURL) ; rerr != nil { 
			fmt.Printf("couldn't cache request to redis: %s\n", rerr)
		}
	} else { 
		finalRes.IsCached = true 
		finalRes.ShortURL = short 
	}
	
	respBytes, err := json.Marshal(finalRes) 
	// Convert the response to JSON

	if err != nil {
		http.Error(w, "Failed to marshal response JSON", http.StatusInternalServerError)
		return
	}

	// Set the response Content-Type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

func Start(cfg config.Config) {
	http.HandleFunc("/", shortenURLHandler)
	redisCli, _ = redis.NewRedisClient(cfg.Database) 
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}