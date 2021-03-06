package main

import (
	"log"
	"sync"

	cl "github.com/vestlog/nix/pkg/client"
	"github.com/vestlog/nix/pkg/models"
	"github.com/vestlog/nix/pkg/storage"
)

var (
	baseurl = "https://jsonplaceholder.typicode.com/"
	dsn     = "storage.db"
	userID  = 7
)

func main() {
	client, err := cl.CreateAPIClient(dsn, baseurl)
	if err != nil {
		log.Fatal("Error creating client:", err)
	}
	// db, err := nix.CreateSQLiteDatabase(dsn)
	db, err := storage.CreateGormDatabase(dsn)
	if err != nil {
		log.Fatal("Could not create DB connection:", err)
	}
	// defer db.Close()
	if err := db.CreateTables(); err != nil {
		log.Fatal(err)
	}
	posts, err := client.GetPosts(userID)
	if err != nil {
		log.Fatal("Error getting posts:", err)
	}
	wg := &sync.WaitGroup{}
	for _, post := range posts {
		if err := db.SavePost(&post); err != nil {
			log.Fatal("Could not save post:", err)
		}
		wg.Add(1)
		go func(post models.Post) {
			defer wg.Done()
			comments, err := client.GetComments(post.ID)
			if err != nil {
				log.Fatal(err)
			}
			for _, comment := range comments {
				wg.Add(1)
				go func(comment models.Comment) {
					defer wg.Done()
					if err := db.SaveComment(&comment); err != nil {
						log.Println(err)
						return
					}
				}(comment)
			}
		}(post)
	}
	wg.Wait()
}
