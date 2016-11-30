package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func main() {
	// Do an initial request to determine how many pages of recipes there are
	// available to download.
	pageCount, err := getRecipePageCount()
	if err != nil {
		panic(err)
	}

	// Initialize the input (page count) channel and a cancelable request context.
	var (
		pc          = make(chan int)
		ctx, cancel = context.WithCancel(context.Background())
	)

	// Fill the input channel in a separate goroutine.
	go func() {
		for i := 1; i <= pageCount; i++ {
			pc <- i
		}
		close(pc)
	}()

	// Kick off the job. This returns two channels; one for receiving downloaded
	// recipes, and one for receiving errors.
	const numGoroutines = 50
	rc, errc := downloadRecipes(numGoroutines, ctx, pc)

	// Demonstration of pipeline cancellation. After 30 seconds, every goroutine
	// will be told to drop what it's doing and exit.
	go func() {
		<-time.After(30 * time.Second)
		fmt.Println("Telling everyone to clean up.")
		cancel()
	}()

	// Wait on input values from the spawned goroutines.
	for {
		// This block uses some clever channel tricks. If `ok` is false, it
		// means that the channel has been closed. A receive operation on a
		// closed channel immediately returns its type's zero value, but a
		// receive operation on a nil channel will never return. Once the
		// channel has been closed, we want to set it to nil to ensure that
		// that arm of the switch statement is never executed again.
		select {
		case recipe, ok := <-rc:
			if !ok {
				rc = nil
				break
			}

		// Do something with recipe. This will probably involve
		// saving it somewhere, unmarshaling it into a struct,
		// or both. Could be a good opportunity for another
		// pipeline!
			_ = recipe
			fmt.Println("Got a recipe.")

		case err, ok := <-errc:
			if !ok {
				errc = nil
				break
			}

			fmt.Println("Error: " + err.Error())
		}

		// Once both channels have been closed and set to nil, we need to
		// break out of the loop to avoid hanging indefinitely on no input.
		if rc == nil && errc == nil {
			break
		}
	}
}

// downloadRecipes spawns numGoroutines goroutines to download recipes from Brewtoad.
// The provided context can be used to cancel in-flight requests. The input channel
// provides the page numbers that should be downloaded. This method returns two channels:
// one for downloaded recipes in XML format, and one for any errors encountered.
func downloadRecipes(numGoroutines int, ctx context.Context, pc <-chan int) (<-chan string, <-chan error) {
	var (
		wg   sync.WaitGroup
		rc   = make(chan string)
		errc = make(chan error)
	)

	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func() {
			running := true
			for running {
				select {
				case page, ok := <-pc:
					if !ok {
						running = false
						break
					}

					err := getRecipesForPage(ctx, rc, page)
				// Don't send anything on the error channel if it was nil,
				// or if cancellation was requested, since we're trying to
				// abort everything anyway.
					if err != nil && ctx.Err() != context.Canceled {
						errc <- err
					}

				// A receive event on this channel means that we're cancelled,
				// so we should stop what we're doing and exit the loop.
				case <-ctx.Done():
					running = false
				}
			}
			wg.Done()
		}()
	}

	// Once all goroutines have finished, close the returned channels.
	go func() {
		wg.Wait()
		close(rc)
		close(errc)
	}()

	return rc, errc
}

// getRecipesForPage downloads a recipe page and sends each recipe found
// there along the provided channel.
func getRecipesForPage(ctx context.Context, rc chan<- string, page int) error {
	doc, err := downloadPage(ctx, page)
	if err != nil {
		return err
	}

	for _, link := range findRecipeLinks(doc) {
		r, err := downloadBeerXml(ctx, link)
		if err != nil {
			return err
		}

		beerXml, err := ioutil.ReadAll(r)
		r.Close()
		if err != nil {
			return err
		}

		rc <- string(beerXml)
	}

	return nil
}

// downloadPage downloads a recipe page from Brewtoad and parses it into
// an HTML document.
func downloadPage(ctx context.Context, page int) (*html.Node, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(
		"https://www.brewtoad.com/recipes?page=%d&sort=rank",
		page,
	), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return html.Parse(resp.Body)
}

// findRecipeLinks traverses an HTML document looking for beer recipe links.
func findRecipeLinks(doc *html.Node) (links []string) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		defer func() {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}()

		if n.Type == html.ElementNode && n.Data == "a" {
			if classes, ok := getAttr(n.Attr, "class"); ok {
				for _, class := range strings.Fields(classes) {
					if class == "recipe-link" {
						href, _ := getAttr(n.Attr, "href")
						links = append(links, href)
						return
					}
				}
			}
		}
	}
	f(doc)
	return
}

// downloadBeerXml downloads the XML for the provided recipe link. If no error
// is returned, then the io.ReadCloser must be closed by the caller in order to
// prevent a resource leak.
func downloadBeerXml(ctx context.Context, link string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", "https://www.brewtoad.com"+link+".xml", nil)
	if err != nil {
		return nil, errors.New("failed to build beer xml request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.New("failed to get beer xml: " + err.Error())
	}
	return resp.Body, nil
}

// getRecipePageCount downloads the first page and checks the pagination to determine
// how many pages of recipes there are.
func getRecipePageCount() (int, error) {
	resp, err := http.Get("https://www.brewtoad.com/recipes?page=1&sort=rank")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return 0, err
	}

	var pageCount int

	var f func(*html.Node)
	f = func(n *html.Node) {
		defer func() {
			if pageCount == 0 {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		}()

		if n.Type == html.ElementNode && n.Data == "a" {
			if classes, ok := getAttr(n.Attr, "class"); ok {
				for _, class := range strings.Fields(classes) {
					if class == "next_page" {
						pageCount, err = strconv.Atoi(n.PrevSibling.PrevSibling.FirstChild.Data)
						return
					}
				}
			}
		}
	}
	f(doc)

	return pageCount, nil
}

// getAttr is a utility method for looking up an HTML element attribute.
func getAttr(attrs []html.Attribute, name string) (string, bool) {
	for _, attr := range attrs {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
}



