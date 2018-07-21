package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	desktop "gitlab.com/axet/desktop/go"

	"github.com/0xAX/notificator"
	"github.com/bradylove/concourse-monitor/pkg/assets"
	"github.com/bradylove/concourse-monitor/pkg/concourse"
	"github.com/bradylove/concourse-monitor/pkg/state"
)

var (
	refreshIntSeconds = flag.Int("refresh-interval", 15, "interval for pulling status from concourse")
	deamonize         = flag.Bool("d", false, "run concourse-monitor in the background")
)

func init() {
	flag.Parse()
}

func main() {
	if *deamonize {
		var args []string

		for _, a := range os.Args[1:] {
			if a != "-d" {
				args = append(args, a)
			}
		}

		fmt.Println("Starting concourse-monitor in the background")
		cmd := exec.Command(os.Args[0], args...)
		cmd.Start()

		os.Exit(0)
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

func openInBrowser(target *concourse.Target, path string) desktop.MenuAction {
	uri, err := url.Parse(target.API)
	if err != nil {
		return func(*desktop.Menu) {}
	}

	return func(*desktop.Menu) {
		uri.Path = path
		desktop.BrowserOpenURI(uri.String())
	}
}

func syncState(tray *desktop.DesktopSysTray, cache *Cache) {
	targets, err := concourse.LoadTargets(filepath.Join(os.Getenv("HOME"), ".flyrc"))
	if err != nil {
		log.Printf("Failed to load .flyrc: %s", err)
		return
	}

	client := concourse.NewClient(targets)
	pipelines, err := client.Pipelines()
	if err != nil {
		log.Printf("Failed to fetch pipelines: %s", err)
		return
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

func pipelineToMenu(p *concourse.Pipeline) desktop.Menu {
	return desktop.Menu{
		Type:    desktop.MenuItem,
		Enabled: true,
		Name:    p.DisplayName,
		Menu:    jobsToMenus(p.Target, p.Jobs),
		Action:  openInBrowser(p.Target, p.URL),
		Icon:    state.StatusIcon(state.PipelineStatus(p)),
	}
}

func jobsToMenus(target *concourse.Target, jobs []*concourse.Job) []desktop.Menu {
	var menu []desktop.Menu
	for _, j := range jobs {
		item := desktop.Menu{
			Type:    desktop.MenuItem,
			Icon:    state.StatusIcon(state.JobStatus(j)),
			Enabled: true,
			Name:    j.Name,
			Action:  openInBrowser(target, j.URL),
		}
		menu = append(menu, item)
	}

	return menu
}
