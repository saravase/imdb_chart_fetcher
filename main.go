package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const DEFAULT_YEAR int = 0000
const DEFAULT_RATING float64 = 0.0

type Movie struct {
	Title            string  `json:"title"`
	MovieReleaseYear int     `json:"movie_release_year"`
	IMDBRating       float64 `json:"imdb_rating"`
	Summary          string  `json:"summary"`
	Duration         string  `json:"duration"`
	Genre            string  `json:"genre"`
}

type Doc struct {
	doc *goquery.Document
}

// ParseInt function parses the string value to int value and returns it.
func ParseInt(num string) (int, error) {
	return strconv.Atoi(num)
}

// Trim function trims the space present before and after the string and returns it.
func Trim(content string) string {
	return strings.TrimSpace(content)
}

// GetNewDocument function returns the html page content of corresponing URL.
func GetNewDocument(url string) *Doc {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalf("[ERROR] chartUrl Document Creation Failed. Reason :%v", err)
		os.Exit(1)
	}
	return &Doc{
		doc: doc,
	}
}

// GetURLProps function returns the correcsponding URL properties.
func GetURLProps(chartUrl string) *url.URL {
	urlProps, err := url.Parse(chartUrl)
	if err != nil {
		log.Fatal("[ERROR] URL Parse Failed. Reason: %v", err)
		os.Exit(1)
	}
	return urlProps
}

// GetMovieLinks functions fetches all the movie links from corresponding document.
func GetMovieLinks(docs *Doc, url string) []string {
	var movieLinks []string
	urlProps := GetURLProps(url)
	docs.doc.Find(".titleColumn a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		movieLinks = append(movieLinks, urlProps.Scheme+"://"+urlProps.Host+link)
	})
	return movieLinks
}

// GetTitleAndYear function returns the movie title and movie released year separately
// from the corresponding movie.
func GetTitleAndYear(docs *Doc) (string, int) {

	titleWithYear := Trim(docs.doc.Find("div .title_wrapper h1").Contents().Text()) // Title (YYYY)
	titleList := strings.Split(titleWithYear, "(")

	if len(titleList) == 2 {
		title := Trim(titleList[0])
		year, err := ParseInt(Trim(titleList[1][:len(titleList[1])-1]))
		if err != nil {
			year = DEFAULT_YEAR
		}
		return title, year
	}

	return "", DEFAULT_YEAR
}

// GetIMDBRating function returns IMDB rating of corresponding movie.
func GetIMDBRating(docs *Doc) float64 {
	rating, err := strconv.ParseFloat(Trim(docs.doc.Find("div .ratingValue strong span").Contents().Text()), 64)
	if err != nil {
		rating = float64(DEFAULT_RATING)
	}
	return rating
}

// GetSummary function returns the summary of the corresponding movie.
func GetSummary(docs *Doc) string {
	return Trim(docs.doc.Find("div .summary_text").Contents().Text())
}

// GetDuration function returns the duration of the corresponding movie.
func GetDuration(docs *Doc) string {
	return Trim(docs.doc.Find("div .subtext time").Contents().Text())
}

// GetGenre function returns the genre of the corresponding movie.
func GetGenre(docs *Doc) string {
	var genreList []string
	genreTags := docs.doc.Find("div .subtext a")
	count := len(genreTags.Nodes) - 1
	genreTags.Each(func(index int, item *goquery.Selection) {
		linkTag := item
		if count != index {
			genreList = append(genreList, linkTag.Text())
		}
	})

	genre := strings.Join(genreList, ", ")
	return genre
}

// GetMovieList function returns movies based on corresponding URL and itemCount.
func GetMovieList() []*Movie {

	if len(os.Args) != 3 {
		log.Fatal("[ERROR] Arguments count mismatch.")
		os.Exit(1)
	}

	chartUrl := os.Args[1]
	doc := GetNewDocument(chartUrl)
	movieLinks := GetMovieLinks(doc, chartUrl)

	itemsCount, err := ParseInt(os.Args[2])
	if err != nil || itemsCount < 0 {
		log.Fatal("[ERROR] Invalid itemsCount.")
		os.Exit(1)
	}

	var movieList []*Movie
	for index, movieLink := range movieLinks {
		if index+1 > itemsCount {
			break
		}

		doc = GetNewDocument(movieLink)
		title, year := GetTitleAndYear(doc)
		movie := &Movie{
			Title:            title,
			MovieReleaseYear: year,
			IMDBRating:       GetIMDBRating(doc),
			Summary:          GetSummary(doc),
			Duration:         GetDuration(doc),
			Genre:            GetGenre(doc),
		}
		movieList = append(movieList, movie)
	}
	return movieList
}

func main() {

	movies := GetMovieList()
	if len(movies) == 0 {
		log.Println("[INFO] No movie record.")
		os.Exit(1)
	}
	movieList, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("[ERROR] JSON serialization Failed. Reason: %v", err)
		os.Exit(1)
	}

	fmt.Println(string(movieList))
}
