package Posts

import "time"

type Post struct {
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Author           Author    `json:"author"`
	Category         string    `json:"category"`
	Text             string    `json:"text"`
	URL              string    `json:"url"`
	Votes            []Votes   `json:"votes"`
	Comments         []Comment `json:"comments"`
	Created          time.Time `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`
	ID               string    `json:"id"`
}

type Comment struct {
	Created time.Time `json:"created"`
	Author  Author    `json:"author"`
	Body    string    `json:"body"`
	ID      string    `json:"id"`
}

type IncomingComment struct {
	Body string `json:"comment"`
}

type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Votes struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type Message struct {
	Message string `json:"message"`
}
