package main

import (
	"github.com/Scrimzay/4chansite/apihandlers"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)
	
func main() {
	r := gin.Default()
	// Create the FuncMap first
    funcMap := template.FuncMap{
        "add": func(a, b int) int { return a + b },
        "sub": func(a, b int) int { return a - b },
        "seq": func(n int) []int {
            var seq []int
            for i := 1; i <= n; i++ {
                seq = append(seq, i)
            }
            return seq
        },
    }
    
    // Set the FuncMap before loading templates
    r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("public/*.html")
	r.Static("/static", "./static")

	r.GET("/", indexHandler)
	r.GET("/boardslist", boardsListHandler)
	r.GET("/catalog", catalogIndexHandler)
	r.POST("/catalog/search", catalogSearchHandler)
	r.GET("/catalog/:id/info", catalogBoardHandler)

	err := r.Run(":4000")
	if err != nil {
		log.Fatal(err)
	}
}
	
func indexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func catalogIndexHandler(c *gin.Context) {
	c.HTML(200, "catalogindex.html", nil)
}

func boardsListHandler(c *gin.Context) {
	boards, err := apihandlers.CallGeneral4ChanBoardInformation()
	if err != nil {
		log.Printf("Error fetching 4chan boards: %v", err)
        c.HTML(http.StatusInternalServerError, "4chanboardsinfo.html", gin.H{
            "Error": "Failed to fetch 4chan boards info",
        })
	}

	c.HTML(200, "4chanboardsinfo.html", gin.H{
		"Boards": boards,
	})
}

func catalogBoardHandler(c *gin.Context) {
    id := c.Param("id")
    pageNum := 1 // Default to page 1
    
    // Parse page number from query parameter
    if p := c.Query("page"); p != "" {
        if pn, err := strconv.Atoi(p); err == nil && pn >= 1 {
            pageNum = pn
        }
    }

    // Get all catalog pages
    pages, err := apihandlers.CallGeneral4ChanCatalogInformation(id)
    if err != nil {
        log.Printf("Error fetching catalog for board %s: %v", id, err)
        c.HTML(http.StatusNotFound, "catalogindex.html", gin.H{
            "Error": "Board not found for ID: " + id,
        })
        return
    }

    // Validate page number
    if pageNum > len(pages) {
        pageNum = len(pages)
    }

    c.HTML(http.StatusOK, "catalogboard.html", gin.H{
        "BoardID":    id,
        "Page":       pageNum,
        "Threads":    pages[pageNum-1].Threads, // -1 because pages are 0-indexed
        "TotalPages": len(pages),
    })
}

func catalogSearchHandler(c *gin.Context) {
	id := c.PostForm("id")

    _, err := apihandlers.CallGeneral4ChanCatalogInformation(id)
    if err != nil {
        log.Printf("Error validating board ID %s: %v", id, err)
        c.HTML(http.StatusNotFound, "catalogindex.html", gin.H{
            "Error": "Board not found for ID: " + id,
        })
        return
    }

    redirectURL := "/catalog/" + id + "/info"
    c.Redirect(http.StatusFound, redirectURL)
}