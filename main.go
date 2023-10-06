package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"paheScraper/config"
	"paheScraper/model"
	"regexp"
	"strings"
	"time"
)

var (
	browser playwright.Browser
	pw      *playwright.Playwright
	name    *string
	quality *string
)

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	name, quality = GetFlags()
	if *name == "" {
		fmt.Println("No anime name given.")
		os.Exit(1)
	}

	if !IsExists(config.FirstRun) {
		err = os.Mkdir(config.UserDocument, os.ModePerm)
		if err != nil {
			log.Printf("can't create directory. err: %v", err)
		}
		runOption := &playwright.RunOptions{
			DriverDirectory: "",
			//SkipInstallBrowsers: true,
			Browsers: []string{"chrome"},
			Verbose:  false,
		}
		err = playwright.Install(runOption)
		if err != nil {
			log.Printf("can't install browser %v", err)
		}
		_, err = os.Create(config.FirstRun)
		if err != nil {
			log.Printf("can't create %s\nErr: %v", config.FirstRun, err)
		}
	}
	pw, err = playwright.Run()
	if err != nil {
		log.Printf("could not start playwright: %v", err)
	}
	option := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		//Args:     []string{"--start-maximized"},
	}
	browser, err = pw.Chromium.Launch(option)
	if err != nil {
		log.Printf("could not launch browser: %v", err)
	}
}

func main() {

	playWrigh()
}

func playWrigh() {
	page, err := browser.NewPage()
	if err != nil {
		log.Printf("could not create page: %v", err)
	}

	defer func(pw *playwright.Playwright) {
		err = pw.Stop()
		if err != nil {
			log.Printf("could not stop Playwright: %v", err)
		}
	}(pw)
	defer func(browser playwright.Browser) {
		err = browser.Close()
		if err != nil {
			log.Printf("could not close browser: %v", err)
		}
	}(browser)
	gotoOptions := playwright.PageGotoOptions{Timeout: playwright.Float(15000)}

	if _, err = page.Goto("https://animepahe.ru/", gotoOptions); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	/// Select search box and type in anime name
	search := page.GetByPlaceholder("Search") //Locator(".input-search")

	err = search.Click(playwright.LocatorClickOptions{
		Delay: playwright.Float(1000),
	})
	if err != nil {
		log.Printf("couldn't click search bar %v", err)
	}
	err = search.Type(*name, playwright.LocatorTypeOptions{
		Delay: playwright.Float(200),
	})
	if err != nil {
		log.Printf("couldn't fill text field %v", err)
	}
	err = page.Locator(".search-results > li:nth-child(1)").Click()
	if err != nil {
		log.Printf("can't click search result %v", err)
	}

	pattern := `Episodes \((\d*[1-9]\d*)\)`
	re := regexp.MustCompile(pattern)
	assertions := playwright.NewPlaywrightAssertions(10000)                    //PlaywrightAssertions()
	err = assertions.Locator(page.Locator(".episode-count")).ToContainText(re) //Not().ToHaveCount(2)
	if err != nil {
		log.Printf("assertions error. err: %v", err)
	}
	html, err := page.Locator(".episode-count").InnerHTML()
	if err != nil {
		return
	}

	count, err := page.Locator(".episode-wrap").Count()
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Println("episodes:", count, html)
	//panic("Exiting")
	page.OnPopup(func(tab playwright.Page) {
		//fmt.Println("Popup url: ", tab.URL())
		go func() {
			err = tab.Close()
			if err != nil {
				log.Fatalf("can't close new tab %v", err)
			}
		}()
	})
	/// Click the first episode in the list
	err = page.Locator("div.episode-wrap:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Click()
	if err != nil {
		log.Printf("couldn't click first episode %v", err)
	}
	for i := count; i > 0; i-- {

		textContents, err := page.Locator(".theatre-info > h1:nth-child(2)").TextContent()
		if err != nil {
			log.Printf("%v", err)
		}
		animeDetail := &model.AnimeDetails{}
		animeDetail.Name, animeDetail.Episode = GetNameAndEpisode(textContents)
		fmt.Println("Anime:", animeDetail.Name, "Episode:", animeDetail.Episode, i)

		/// Click the dropdown button then click the last item in the dropdown
		err = page.Locator("div.col-12:nth-child(4) > div:nth-child(1)").Click()
		if err != nil {
			log.Printf("couldn't click first episode %v", err)
		}

		text := page.Locator("#pickDownload") //("SubsPlease")

		/// This enables us to get the contents of the dropdown
		/// Also try to get by attribute
		listLink, err := text.Locator(".dropdown-item").All() //text.All()

		fmt.Print("\n---------\n")

	list:
		for _, content := range listLink {
			links, err := content.GetAttribute("href")
			if err != nil {
				log.Printf("no such attribute %v", err)
			}
			linkName, err := content.InnerText()
			if err != nil {
				log.Printf("no text %v", err)
			}

			if strings.Contains(linkName, *quality) {
				//fmt.Println("Entered scope")
				getDownloadLink(animeDetail, links)
				break list
			}
		}
		//pages := page.Context().Pages()
		//fmt.Println(pages, len(pages))
		/// TODO: Track state of scraping, save to file after 10 episodes have been scraped
		/// TODO: Replace log.Fatalf calls

		nextEpisode, err := page.GetByTitle("Play Next Episode").GetAttribute("href")
		if err != nil {
			log.Printf("can't get link. err: %v", err)
		}
		_, err = page.Goto("https://animepahe.ru"+nextEpisode, gotoOptions)
		if err != nil {
			log.Printf("can't goto page %v", err)
		}
		if i%10 == 0 && i != count {
			time.Sleep(time.Minute * 3)
		}
	}
	err = page.Close()
	if err != nil {
		log.Println("can't close tab. err:", err)
	}

}

func getDownloadLink(details *model.AnimeDetails, links string) {
	newContext, err := browser.NewContext(
		playwright.BrowserNewContextOptions{
			UserAgent: playwright.String("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0")})
	if err != nil {
		log.Printf("can't create new browser context %v", err)
	}

	defer func(newContext playwright.BrowserContext) {
		err = newContext.Close()
		if err != nil {
			log.Printf("could not close browser context: %v", err)
		}
	}(newContext)
	newPage, err := newContext.NewPage()
	if err != nil {
		log.Printf("can't create new page %v", err)
	}

	fmt.Println("New page", links)
	if _, err = newPage.Goto(links, playwright.PageGotoOptions{
		Timeout: playwright.Float(10000)}); err != nil {
		log.Printf("can't goto page %v", err)
	}
	waitPage := newPage.GetByText("Continue")
	nextPageLink, err := waitPage.GetAttribute("href")
	if err != nil {
		log.Printf("can't get attribute %v", err)
	}
	//fmt.Println("About to move")
	_, err = newPage.Goto(nextPageLink)
	if err != nil {
		log.Printf("can't navigate to page %v", err)
	}
	fmt.Println("Moved")
	/// Click anywhere to trigger ads
	err = newPage.Mouse().Click(100, 100)
	if err != nil {
		log.Printf("can't click anywhere %v", err)
	}

	download, err := newPage.ExpectDownload(func() error {
		err = newPage.GetByText("Download").Click()
		if err != nil {
			log.Printf("can't click download button %v", err)
		}
		return err
	})
	if err != nil {
		log.Printf("can't see any download %v", err)
	}
	err = download.Cancel()
	if err != nil {
		log.Printf("can't cancel downlaod %v", err)
	}
	details.Url = download.URL()
	details.SetExpireTime()
	//fmt.Println(string(details.ToJson()))
	//fmt.Println("Reached here too", *details)
	details.SaveToFile()
}
