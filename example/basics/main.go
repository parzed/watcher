package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"../../../watcher"
)

func main() {
	w := watcher.New()

	// Uncomment to use SetMaxEvents set to 1 to allow at most 1 event to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	// w.SetMaxEvents(1)

	// Uncomment to only notify rename and move events.
	// w.FilterOps(watcher.Rename, watcher.Move)

	// Uncomment to filter files based on a regular expression.
	//
	// Only files that match the regular expression during file listing
	// will be watched.
	// r := regexp.MustCompile("^abc$")
	// w.AddFilterHook(watcher.RegexFilterHook(r, false))
	counter := 0
	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(counter) // Print the event's info.
				_ = event
				counter++
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
		fmt.Println(counter)
	}()

	// Watch this folder for changes.
	if err := w.Add("."); err != nil {
		log.Fatalln(err)
	}

	// Watch test_folder recursively for changes.
	dirPath, _ := os.Getwd()
	if err := w.AddRecursive(dirPath); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	//for path, f := range w.WatchedFiles() {
	//	fmt.Printf("%s: %s\n", path, f.Name())
	//}

	fmt.Println()

	// Trigger 2 events after watcher started.
	go func() {
		w.Wait()
		start := time.Now()
		for i := 0; i < 300000; i++ {
			fmt.Println(i)
			go w.TriggerEvent(watcher.Rename, nil)
		}
		duration := time.Since(start)
		fmt.Println(duration)
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
