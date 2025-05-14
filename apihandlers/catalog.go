package apihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FourChanCatalogAPIResponse []Page

type Page struct {
	Page int `json:"page"`
	Threads []Thread `json:"threads"`
}

type Thread struct {
	No int `json:"no"`
	Now string `json:"now"`
	Name string `json:"name"`
	Sub string `json:"sub"`
	Com string `json:"com"`
	SemanticURL string `json:"semantic_url"`
	Replies int `json:"replies"`
	Images int `json:"images"`
	OmittedPosts int `json:"omitted_posts"`
	OmittedImages int `json:"omitted_images"`
	LastModified int `json:"last_modified"`
}

func CallGeneral4ChanCatalogInformation(input string) ([]*Page, error) {
	baseURL := fmt.Sprintf("https://a.4cdn.org/%s/catalog.json", input)

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Error sending get req to 4chan catalog api: %v", err)
	}
	defer resp.Body.Close()

	var catalog []*Page
	if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return nil, fmt.Errorf("Failed to decode JSON response from 4chan catalog api: %w", err)
	}

	return catalog, nil
}