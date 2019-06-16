package medium

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	url            = "https://medium.com/feed/@romanyx90"
	defaultTimeout = time.Second * 5
)

// RSS used to unmarshal rss feed.
type RSS struct {
	Channel struct {
		Items []Story `xml:"item"`
	} `xml:"channel"`
}

// Story holds data of the user story.
type Story struct {
	Title      string   `xml:"title" json:"title"`
	Link       string   `xml:"link" json:"link"`
	Categories []string `xml:"category" json:"categories"`
}

// Stories returns list of user stories.
func Stories() ([]Story, error) {
	c := http.Client{
		Timeout: defaultTimeout,
	}

	res, err := c.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "client request failed")
	}
	defer res.Body.Close()

	var r RSS
	if err := xml.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	var result []Story
	for _, s := range r.Channel.Items {
		if len(s.Title) < 56 {
			s.Link = strings.Split(s.Link, "?source=")[0]
			result = append(result, s)
		}
	}

	return result, nil
}
