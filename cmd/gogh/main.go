package main

import (
	"flag"
	"fmt"
	"log"

	"gogh/internal/gogh"

	gh "github.com/j178/github-s3"
)

func main() {
	var repo = flag.String("repo", "", "target repo on github")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalln("Usage: gogh [file-path]")
	}

	gogh, err := gogh.New()
	if err != nil {
		log.Fatalln(err)
	}
	if gogh.Data.Settings.SessionToken == "" {
		var token string
		fmt.Println("enter user_session cookie:")
		if _, err := fmt.Scanln(&token); err != nil {
			log.Fatalln(err)
		}
		gogh.Data.Settings.SessionToken = token
		if err := gogh.SaveData(); err != nil {
			log.Fatalln(err)
		}
	}

	_gh := gh.New(gogh.Data.Settings.SessionToken, *repo)
	for _, path := range flag.Args() {
		res, err := _gh.UploadFromPath(path)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Println(res.GithubLink)
	}
}
