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
		StartTime string
		StopTime  string
		Text      string
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
			rows, err := db.Query("SELECT episode, startTime, stopTime, TEXT FROM data WHERE (MATCH ( TEXT ) AGAINST (? IN BOOLEAN MODE)) AND (MATCH ( episode ) AGAINST (? IN BOOLEAN MODE))", q, s)
			if err != nil {
				fmt.Print(err.Error())
			}
			for rows.Next() {
				err = rows.Scan(&subtitle.Episode, &subtitle.StartTime, &subtitle.StopTime, &subtitle.Text)
				subtitles = append(subtitles, subtitle)
				if err != nil {
					fmt.Print(err.Error())
				}
			}
			defer rows.Close()
		} else {
			// SELECT  `episode` ,  `stopTime` ,  `startTime` ,  `text` , MATCH ( `text` ) AGAINST ( 'I am krombopolis micheal and i just love killing' IN BOOLEAN MODE) AS score FROM data ORDER BY score DESC LIMIT 12120 , 30
			rows, err := db.Query("SELECT episode, startTime, stopTime, TEXT FROM data WHERE (MATCH ( TEXT ) AGAINST (? IN BOOLEAN MODE))", q)
			if err != nil {
				fmt.Print(err.Error())
			}
			for rows.Next() {
				err = rows.Scan(&subtitle.Episode, &subtitle.StartTime, &subtitle.StopTime, &subtitle.Text)
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
