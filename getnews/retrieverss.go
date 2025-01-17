package rss

import (
	"fmt"
	"github.com/ungerik/go-rss"
	"rssfeed/repositories"
	"time"
)

// The is the method responsible for retrieving rss feeds form different sources. it returns a channel of type rss.Channel
func GetRss(c chan rss.Channel, url string) (chan rss.Channel, error) {
	channel, err := rss.Read(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("I have received from: " + channel.Title)

	c <- *channel
	return c, err
}

// This method receives from the channel and returns feeds to be saved in database
func ReceiveFromChannel(c <-chan rss.Channel) []interface{} {

	channel := <-c
	//creating an interface for saving in database
	feeds := make([]interface{}, len(channel.Item))
	for i, v := range channel.Item {
		feeds[i] = v
	}

	return feeds

}

// This function uses the GetNew and Receive from channel to loop through the urls
func Spider() bool {

	var err1 error
	c := make(chan rss.Channel, 100)

	urls := []string{
		"http://rss.cnn.com/rss/edition_world.rss",
		"http://feeds.bbci.co.uk/news/world/rss.xml",
	}

	// Loop for getting feeds, passing to receive  and saving in database
	for i := 0; i < len(urls); i++ {
		GetRss(c, urls[i])
		feed := ReceiveFromChannel(c)
		_, err1 = repositories.SaveToDb(feed)

		time.Sleep(15 * time.Second)
	}

	if err1 != nil {
		return false
	}

	return true
}

//This method keeps the execution of the spider method at interval
func StartSpider() {
	/*timestamp := time.Now().Local()*/

	for _ = range time.Tick(2 * time.Minute) {
		Spider()
		/*fmt.Println("data at " + timestamp.String())
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		log.Print("Allocated memory: %fMB", float32(mem.Alloc)/1024.0/1024.0)*/
	}

}
