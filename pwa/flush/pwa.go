package flush

import (
	"log"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func Run() {
	app.Route("/", func() app.Composer {
		return &RootContainer{}
	})
	app.Route("/login", func() app.Composer { return &LoginContainer{} })
	app.Route("/new", func() app.Composer { return &NewFlushContainer{} })
	app.Route("/settings", func() app.Composer { return &SettingsContainer{} })
	app.RunWhenOnBrowser()

	if os.Getenv("BUILD_STATIC") == "true" {
		err := app.GenerateStaticWebsite(".", &app.Handler{
			Name:        "Flush-Log",
			Description: "bowel tracking app",
			Resources:   app.GitHubPages("flush-log"),
			Styles: []string{
				"/web/style.css",
			},
			Image: "/web/paper.png",
			Icon: app.Icon{
				Maskable: "/web/paper.png",
				Default:  "/web/paper.png",
				Large:    "/web/paper.png",
				SVG:      "/web/paper.svg",
			},
			LoadingLabel: "Unrolling paper...",
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	port := os.Getenv("GOAPP_PORT")
	if port == "" {
		log.Fatal("GOAPP_PORT not set")
	}
	apiUrlBytes, err := os.ReadFile("web/apiurl")
	if err != nil {
		log.Fatal(err)
	}
	apiUrl := string(apiUrlBytes)

	log.Println("Listening on " + port)
	log.Println("API url: ", apiUrl)
	http.Handle("/", &app.Handler{
		Name:        "Flush-Log",
		Description: "bowel tracking app",
		Scripts: []string{
			"https://cdn.tailwindcss.com",
		},
		Image: "/web/paper.png",
		Icon: app.Icon{
			Maskable: "/web/paper.png",
			Default:  "/web/paper.png",
			Large:    "/web/paper.png",
			SVG:      "/web/paper.svg",
		},
		LoadingLabel: "Unrolling paper...",
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
