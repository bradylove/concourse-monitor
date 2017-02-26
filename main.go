package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/0xAX/notificator"
	desktop "github.com/axet/desktop/go"
	"github.com/bradylove/cc-monitor/assets"
)

var (
	teamName          = flag.String("team-name", "main", "team name to monitor")
	concourseURL      = flag.String("concourse-url", "", "url for concourse instance")
	refreshIntSeconds = flag.Int("refresh-interval", 15, "interval for pulling status from concourse")
)

func init() {
	flag.Parse()
}

func main() {
	if *concourseURL == "" {
		log.Fatalf("concourse-url cannot be empty")
	}

	refreshInt := time.Duration(*refreshIntSeconds) * time.Second

	notify := notificator.New(notificator.Options{
		DefaultIcon: assets.CCIconPath,
		AppName:     "Concourse Monitor",
	})
	cache := NewCache(notify)

	tray := desktop.DesktopSysTrayNew()
	icon := assets.Image("icons/cc_icon.png")

	menu := []desktop.Menu{
		desktop.Menu{Type: desktop.MenuItem, Enabled: false, Name: "Loading..."},
	}

	tray.SetIcon(icon)
	tray.SetTitle("Concourse")
	tray.SetMenu(menu)
	tray.Show()

	go func() {
		syncState(tray, cache)

		for range time.Tick(refreshInt) {
			syncState(tray, cache)
		}
	}()

	desktop.Main()
}

func openInBrowser(path string) desktop.MenuAction {
	return func(*desktop.Menu) {
		desktop.BrowserOpenURI(fmt.Sprint(*concourseURL, path))
	}
}

func syncState(tray *desktop.DesktopSysTray, cache *Cache) {
	client := NewConcourseClient(*concourseURL)

	pipelines, err := client.GetPipelines(*teamName)
	if err != nil {
		log.Println("Failed to fetch pipelines: %s", err)
	}

	cache.Update(pipelines)

	if len(pipelines) < 1 {
		menu := []desktop.Menu{
			desktop.Menu{Type: desktop.MenuItem, Enabled: false, Name: "No pipelines configured..."},
		}

		tray.SetMenu(menu)
		tray.Update()
		return
	}

	var menu []desktop.Menu
	for _, p := range pipelines {
		menu = append(menu, pipelineToMenu(p))
	}

	tray.SetMenu(menu)
	tray.Update()
}

func pipelineToMenu(p *Pipeline) desktop.Menu {
	return desktop.Menu{
		Type:    desktop.MenuItem,
		Enabled: true,
		Name:    p.Name,
		Menu:    jobsToMenus(p.Jobs),
		Action:  openInBrowser(p.URL),
		Icon:    p.StatusIcon(),
	}
}

func jobsToMenus(jobs []*Job) []desktop.Menu {
	var menu []desktop.Menu
	for _, j := range jobs {
		item := desktop.Menu{
			Type:    desktop.MenuItem,
			Icon:    j.StatusIcon(),
			Enabled: true,
			Name:    j.Name,
			Action:  openInBrowser(j.URL),
		}
		menu = append(menu, item)
	}

	return menu
}
