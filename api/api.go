package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "squancher:AHappyTedyBear@tcp(127.0.0.1:3306)/wumbo?parseTime=true")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	type Subtitle struct {
		Episode   string
		TimeStamp int
		//StartTime string
		//StopTime  string
		Text string
	}
	router := gin.Default()

	// GET a subtitle search
	router.GET("/search", func(c *gin.Context) { // search/search?q=search%20quer&s=S01E01
		var (
			subtitle  Subtitle
			subtitles []Subtitle
		)
		q := c.Query("q") + "*"
		if len(q) < 4 {
			return
		}
		s := c.Query("s") + "*"
		if len(s) != 1 {
			fmt.Print("the season has been set")
			rows, err := db.Query("SELECT episode, midTime, text FROM data WHERE (MATCH ( TEXT ) AGAINST (? IN BOOLEAN MODE)) AND (MATCH ( episode ) AGAINST (? IN BOOLEAN MODE))", q, s)
			if err != nil {
				fmt.Print(err.Error())
			}
			for rows.Next() {
				err = rows.Scan(&subtitle.Episode, &subtitle.TimeStamp, &subtitle.Text)
				subtitles = append(subtitles, subtitle)
				if err != nil {
					fmt.Print(err.Error())
				}
			}
			defer rows.Close()
		} else {
			// SELECT `episode` ,  `startTime`, `stopTime` ,  `text` , MATCH (TEXT) AGAINST (? IN NATURAL LANGUAGEMODE) AS relevance FROM data ORDER BY relevance DESC LIMIT 64
			//sorting by relevance is *okay* it yields slower load times though ~8-16ms compared to ~2-4ms without and also returns many non relevant results and is only limited to the 64 limit it will return all subtitles.
			rows, err := db.Query("SELECT episode, midTime, text FROM data WHERE (MATCH ( TEXT ) AGAINST (? IN BOOLEAN MODE))", q)
			if err != nil {
				fmt.Print(err.Error())
			}
			for rows.Next() {
				err = rows.Scan(&subtitle.Episode, &subtitle.TimeStamp, &subtitle.Text)
				subtitles = append(subtitles, subtitle)
				if err != nil {
					fmt.Print(err.Error())
				}
			}
			defer rows.Close()
		}
		c.JSON(http.StatusOK, gin.H{
			"subtitles": subtitles,
			"count":     len(subtitles),
		})
	})
	router.Run(":9002")
}
