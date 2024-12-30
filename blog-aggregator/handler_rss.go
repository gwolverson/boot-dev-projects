package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gwolverson/go-courses/blog-aggregator/internal/database"
	"html"
	"io"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for index, item := range feed.Channel.Item {
		feed.Channel.Item[index].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[index].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func handlerRSS(state *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(state)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	post, err := db.CreatePost(context.Background(), database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		Title:       html.UnescapeString(feedData.Channel.Title),
		Url:         feed.Url,
		Description: html.UnescapeString(feedData.Channel.Description),
		PublishedAt: feed.CreatedAt,
		FeedID:      feed.ID,
	})
	if err != nil {
		log.Printf("Couldn't create post: %v", err)
	}
	log.Printf("Post stored: %v", post)
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func handlerAddFeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("add feed command expects two arguments; name and url")
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := state.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = state.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("New Feed: %s\n", feed)
	return nil
}

func handlerListFeeds(state *state, cmd command) error {
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		feedUser, err := state.db.GetUsername(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("Error getting user associated to feed: %s\n", err)
			return err
		}
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Feed user: %s\n", feedUser)
	}
	return nil
}

func handlerFollowFeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow feed command expects a single argument; the feed url")
	}

	feed, err := state.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := state.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Feed name: %s, user: %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(state *state, cmd command) error {
	currentUser, err := state.db.GetUser(context.Background(), state.appConfig.CurrentUserName)
	if err != nil {
		fmt.Printf("Error getting current user: %s\n", err)
		return err
	}
	followFeeds, err := state.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	for _, feed := range followFeeds {
		fmt.Printf("%s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollowFeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("unfollow command expects one argument; the feed url")
	}

	err := state.db.DeleteFeedFollowByUserAndUrl(context.Background(), database.DeleteFeedFollowByUserAndUrlParams{
		UserID: user.ID,
		Url:    cmd.args[0],
	})

	if err != nil {
		fmt.Printf("Error deleting feed follow: %s\n", err)
		return err
	}
	return nil
}
