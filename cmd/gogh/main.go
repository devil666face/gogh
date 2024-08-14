package main

import (
	"gogh/internal/gogh"
	"gogh/internal/view"
	"log"
)

func main() {
	_gogh, err := gogh.New()
	if err != nil {
		log.Fatalln(err)
	}

	_view := view.New(
		_gogh,
	)
	_view.Run()

	// flag.Parse()
	// if len(flag.Args()) == 0 {
	// 	log.Fatalln("Usage: gogh [path]")
	// }

	// if gogh.Data.Settings.SessionToken == "" {
	// 	var token string
	// 	fmt.Println("enter user_session cookie:")
	// 	if _, err := fmt.Scanln(&token); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// 	gogh.SetToken(token)
	// 	if err := gogh.SaveData(); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	// for _, path := range flag.Args() {
	// 	if err := gogh.Upload(path); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// if err := gogh.UploadParalel(path); err != nil {
	// 	log.Fatalln(err)
	// }
	// }
}
