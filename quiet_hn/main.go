package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type storyCache struct {
	numStories     int
	cache          []item
	duration       time.Duration
	expirationTime time.Time
	mutex          sync.Mutex
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		numStories: numStories,
		duration:   5 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   5 * time.Second,
			}
			temp.refresh()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expirationTime = temp.expirationTime
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	t := time.Now()
	if t.Sub(sc.expirationTime) > 0 {
		err := sc.refresh()
		if err != nil {
			return nil, err
		}
	}
	res := sc.readCache()
	return res, nil
}

func (sc *storyCache) readCache() []item {
	return sc.cache
}

func (sc *storyCache) refresh() error {
	t := time.Now()
	var err error
	sc.cache, err = getTopStories(sc.numStories)
	sc.expirationTime = t.Add(sc.duration)
	if err != nil {
		return err
	}
	return nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		err = errors.New("Failed to load top stories")
		return nil, err
	}
	var stories []item
	offset := 0
	for len(stories) < numStories {
		want := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[offset:want])...)
		offset += want
	}
	return stories[:numStories], nil
}

type chanMessage struct {
	it  item
	idx int
	err error
}

func getStories(ids []int) []item {
	ch := make(chan chanMessage)
	// Create goroutines to fetch stories
	for i := 0; i < len(ids); i++ {
		go func(i int) {
			var client hn.Client
			id := ids[i]
			hnItem, err := client.GetItem(id)
			if err != nil {
				ch <- chanMessage{err: err}
			}
			item := parseHNItem(hnItem)
			ch <- chanMessage{it: item, idx: i}
		}(i)
	}
	// Add to slice
	var messages []chanMessage
	for i := 0; i < len(ids); i++ {
		m := <-ch
		if isStoryLink(m.it) && m.err == nil {
			messages = append(messages, m)
		}
	}
	// Sort the slice
	sort.Slice(messages, func(i, j int) bool { return messages[i].idx < messages[j].idx })
	// Transfer to stories
	var stories []item
	for _, m := range messages {
		stories = append(stories, m.it)
	}
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
