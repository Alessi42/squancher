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
	router.GET("/api/search/:q/:s", func(c *gin.Context) { // search/search/search%20quer/S01E01 better for nginx caching
		var (
			subtitle  Subtitle
			subtitles []Subtitle
			score     int
			q         string
		)
		q = c.Param("q") + "*"
		if len(q) < 4 {
			return
		}
		if len(q) > 48 {
			q = q[:47] + "*"
		}
		s := c.Param("s") + "*"
		// old SELECT episode, midTime, text FROM data WHERE (MATCH ( TEXT ) AGAINST (? IN BOOLEAN MODE)) AND (MATCH ( episode ) AGAINST (? IN BOOLEAN MODE))
		rows, err := db.Query("SELECT episode, midTime, text, MATCH(text) AGAINST (? IN BOOLEAN MODE) AS relevance FROM data WHERE MATCH(text) AGAINST (? IN BOOLEAN MODE) AND (MATCH ( episode ) AGAINST (? IN BOOLEAN MODE)) ORDER BY relevance DESC LIMIT 32", q, q, s)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&subtitle.Episode, &subtitle.TimeStamp, &subtitle.Text, &score)
			subtitles = append(subtitles, subtitle)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"subtitles": subtitles,
			"count":     len(subtitles),
		})
	})
	router.GET("/api/search/:q", func(c *gin.Context) { // search/search/search%20quer/S01E01 better for nginx caching
		var (
			subtitle  Subtitle
			subtitles []Subtitle
			score     int
			q         string
		)
		q = c.Param("q") + "*"
		if len(q) < 4 {
			return
		}
		if len(q) > 48 {
			q = q[:47] + "*"
		}
		fmt.Print("length of query: ", len(q))
		// SELECT `episode` ,  `startTime`, `stopTime` ,  `text` , MATCH (TEXT) AGAINST (? IN NATURAL LANGUAGEMODE) AS relevance FROM data ORDER BY relevance DESC LIMIT 64
		//sorting by relevance is *okay* it yields slower load times though ~8-16ms compared to ~2-4ms without and also returns many non relevant results and is only limited to the 64 limit it will return all subtitles.
		rows, err := db.Query("SELECT episode, midTime, text, MATCH(text) AGAINST (? IN BOOLEAN MODE) AS relevance FROM data WHERE MATCH(text) AGAINST (? IN BOOLEAN MODE) ORDER BY relevance DESC LIMIT 32", q, q)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&subtitle.Episode, &subtitle.TimeStamp, &subtitle.Text, &score)
			if len(subtitles) < 32 {
				subtitles = append(subtitles, subtitle)
			}
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"subtitles": subtitles,
			"count":     len(subtitles),
		})
	})
	router.GET("/api/random/", func(c *gin.Context) { // search/search/search%20quer/S01E01 better for nginx caching
		var (
			subtitle  Subtitle
			subtitles []Subtitle
			stop      int
			start     int
			mid       int
		)
		// SELECT `episode`, `midTime`, `text` FROM `data` WHERE 1 ORDER BY RAND() LIMIT 0, 32
		// SELECT * FROM data AS r1 JOIN (SELECT CEIL(RAND() * (SELECT MAX(midTime) FROM data)) AS midTime) AS r2 WHERE r1.midTime >= r2.midTime ORDER BY r1.midTime ASC LIMIT 32
		rows, err := db.Query("SELECT * FROM data AS r1 JOIN (SELECT CEIL(RAND() * (SELECT MAX(midTime) FROM data)) AS midTime) AS r2 WHERE r1.midTime >= r2.midTime ORDER BY r1.midTime ASC LIMIT 16")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&subtitle.Episode, &stop, &start, &subtitle.TimeStamp, &subtitle.Text, &mid)
			if len(subtitles) < 32 {
				subtitles = append(subtitles, subtitle)
			}
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"subtitles": subtitles,
			"count":     len(subtitles),
		})
	})
	router.GET("/api/random/:s", func(c *gin.Context) { // search/search/search%20quer/S01E01 better for nginx caching
		var (
			subtitle  Subtitle
			subtitles []Subtitle
		)
		s := c.Param("s") + "*"
		// SELECT `episode`, `midTime`, `text` FROM `data` WHERE 1 ORDER BY RAND() LIMIT 0, 32
		// SELECT * FROM data AS r1 JOIN (SELECT CEIL(RAND() * (SELECT MAX(midTime) FROM data)) AS midTime) AS r2 WHERE r1.midTime >= r2.midTime ORDER BY r1.midTime ASC LIMIT 32
		rows, err := db.Query("SELECT episode, midTime, text FROM data WHERE (MATCH ( episode ) AGAINST (? IN BOOLEAN MODE)) ORDER BY RAND() LIMIT 16", s)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&subtitle.Episode, &subtitle.TimeStamp, &subtitle.Text)
			if len(subtitles) < 32 {
				subtitles = append(subtitles, subtitle)
			}
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"subtitles": subtitles,
			"count":     len(subtitles),
		})
	})
	router.Run(":9002")
}
