# imdb_chart_fetcher

### Step 1: Clone Repository
  git clone https://github.com/saravase/imdb_chart_fetcher.git

### Step 2: Navigate to imdb_chart_fetcher folder
  cd imdb_chart_fetcher
  
### Step 3: Execute imdb_chart_fetcher
  ./imdb_chart_fetcher 'https://www.imdb.com/india/top-rated-indian-movies' 1
  
### Time difference between normal single vs multiple http request at a time

#### make start-loop: [Single http request]
	go run loop/main.go 'https://www.imdb.com/india/top-rated-indian-movies' 100
  Using For loop - 100 http request  took 5m10.6330182s
  
#### make start-routine:[Multiple http request]
	go run routine/main.go 'https://www.imdb.com/india/top-rated-indian-movies' 100
  Using Goroutine - 100 http request  took 26.511292424s
  
