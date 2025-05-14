package apihandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FourChanBoardAPIResponse struct {
	Boards []Boards `json:"boards"`
}

type Boards struct {
	Board string `json:"board"`
	Title string `json:"title"`
	MaxFilesize int `json:"max_filesize"`
	MaxWebmFilesize int `json:"max_webm_filesize"`
	BumpLimit int `json:"bump_limit"`
	ImageLimit int `json:"image_limit"`
	Description string `json:"meta_description"`
}

func CallGeneral4ChanBoardInformation() ([]Boards, error) {
	baseURL := "https://a.4cdn.org/boards.json"

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Error sending get req to 4chan api boards: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading resp body from 4chan api boards: %v", err)
	}

	var boards FourChanBoardAPIResponse
	err = json.Unmarshal(body, &boards)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling json response from 4chan api boards: %v", err)
	}

	return boards.Boards, nil
}