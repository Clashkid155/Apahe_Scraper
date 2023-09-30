package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"paheScraper/model"
	"strings"
)

var (
	browser playwright.Browser
	pw      *playwright.Playwright
)

const (
	firstRun = ".first_run"
)

func init() {
	var err error
	if !IsExists(firstRun) {
		runOption := &playwright.RunOptions{
			DriverDirectory: "",
			//SkipInstallBrowsers: true,
			//Browsers: []string{"chrome"},
			Verbose: false,
		}
		err = playwright.Install(runOption)
		if err != nil {
			log.Fatalf("could not install playwright dependencies: %v", err)
		}
		_, err = os.Create(firstRun)
		if err != nil {
			log.Fatalf("can't create %s %v", firstRun, err)
		}
	}
	pw, err = playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	option := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		SlowMo:   playwright.Float(0),
		//Args:     []string{"--start-maximized"},
	}
	browser, err = pw.Chromium.Launch(option)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
}

func main() {

	playWrigh()
}

func playWrigh() {
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	defer func(browser playwright.Browser) {
		err = browser.Close()
		if err != nil {
			log.Fatalf("could not close browser: %v", err)
		}
	}(browser)
	defer func(pw *playwright.Playwright) {
		err := pw.Stop()
		if err != nil {
			log.Fatalf("could not stop Playwright: %v", err)
		}
	}(pw)

	gotoOptions := playwright.PageGotoOptions{Timeout: playwright.Float(15000)}

	if _, err = page.Goto("https://animepahe.ru/", gotoOptions); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	/// Select search box and type in anime name
	search := page.Locator(".input-search")
	err = search.Click(playwright.LocatorClickOptions{
		Delay: playwright.Float(1000),
	})
	if err != nil {
		log.Fatalf("couldn't click search bar %v", err)
	}
	err = search.Type("Dark Gathering", playwright.LocatorTypeOptions{
		Delay: playwright.Float(200),
	})
	if err != nil {
		log.Fatalf("couldn't fill text field %v", err)
	}
	err = page.Locator(".search-results > li:nth-child(1)").Click()
	if err != nil {
		log.Fatalf("can't click search result %v", err)
	}

	/// Click the first episode in the list
	err = page.Locator("div.episode-wrap:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Click()
	if err != nil {
		log.Fatalf("couldn't click first episode %v", err)
	}
	textContents, err := page.Locator(".theatre-info > h1:nth-child(2)").TextContent()
	if err != nil {
		log.Fatalf("%v", err)
	}
	animeDetail := &model.AnimeDetails{}
	animeDetail.Name, animeDetail.Episode = GetNameAndEpisode(textContents)
	fmt.Println("Anime:", animeDetail.Name, "Episode:", animeDetail.Episode)

	/// Click the dropdown button then click the last item in the dropdown
	err = page.Locator("div.col-12:nth-child(4) > div:nth-child(1)").Click()
	if err != nil {
		log.Fatalf("couldn't click first episode %v", err)
	}
	/*/// Close popup tab
	page.OnPopup(func(tab playwright.Page) {
		//fmt.Println("Page url: ", tab.URL())
		err = tab.Close()
		if err != nil {
			log.Fatalf("can't close new tab %v", err)
		}
	})*/

	text := page.Locator("#pickDownload") //("SubsPlease")

	/// This enables us to get the contents of the dropdown
	/// Also try to get by attribute
	listLink, err := text.Locator(".dropdown-item").All() //text.All()

	fmt.Print("\n---------\n")

	for i, content := range listLink {
		links, err := content.GetAttribute("href")
		if err != nil {
			log.Fatalf("no such attribute %v", err)
		}
		linkName, err := content.InnerText()
		if err != nil {
			log.Fatalf("no text %v", err)
		}
		fmt.Printf("%d. %s %s\n", i+1, linkName, links)

		if strings.Contains(linkName, "1080p") {
			fmt.Println("Entered scope")
			getDownloadLink(animeDetail, links)

		}
	}

}

func getDownloadLink(details *model.AnimeDetails, links string) {
	newContext, err := browser.NewContext(
		playwright.BrowserNewContextOptions{
			UserAgent: playwright.String("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0")})
	if err != nil {
		log.Fatalf("can't create new browser context %v", err)
	}

	defer func(newContext playwright.BrowserContext) {
		err = newContext.Close()
		if err != nil {
			log.Fatalf("could not close browser: %v", err)
		}
	}(newContext)
	newPage, err := newContext.NewPage()
	if err != nil {
		log.Fatalf("can't create new page %v", err)
	}

	fmt.Println("New page", links)
	if _, err = newPage.Goto(links, playwright.PageGotoOptions{
		Timeout: playwright.Float(10000)}); err != nil {
		log.Fatalf("can't goto page %v", err)
	}
	waitPage := newPage.GetByText("Continue")
	nextPageLink, err := waitPage.GetAttribute("href")
	if err != nil {
		log.Fatalf("can't get attribute %v", err)
	}
	fmt.Println("About to move")
	_, err = newPage.Goto(nextPageLink)
	if err != nil {
		log.Fatalf("can't navigate to page %v", err)
	}
	fmt.Println("Moved")
	/// Click anywhere to trigger ads
	err = newPage.Mouse().Click(100, 100) //Locator(".").Click()
	if err != nil {
		log.Fatalf("can't click anywhere %v", err)
	}
	fmt.Println("Reached here", newPage.Context().Pages())
	download, err := newPage.ExpectDownload(func() error {
		fmt.Println("Yes expected exe")
		err = newPage.GetByText("Download").Click()
		if err != nil {
			log.Fatalf("can't click download button %v", err)
		}
		return err
	})
	if err != nil {
		log.Fatalf("can't see any download %v", err)
	}
	err = download.Cancel()
	if err != nil {
		log.Fatalf("can't cancel downlaod %v", err)
	}
	details.Url = download.URL()
	details.SetExpireTime()
	fmt.Println(string(details.ToJson()))
	fmt.Println("Reached here too", *details)
}
