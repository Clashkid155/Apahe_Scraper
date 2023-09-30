// Command pdf is a chromedp example demonstrating how to capture a pdf of a
// page.

// Was supposed to be me testing out chromedp
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func chrome() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	// create context
	ctx, cancel := chromedp.NewContext(ctx) //, chromedp.WithDebugf(log.Printf))
	defer cancel()

	// capture pdf
	var buf []byte
	var title = new(string)
	if err := chromedp.Run(ctx, printToPDF(`https://animepahe.ru`, &buf, title)); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s.pdf", *title), buf, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s.pdf\n", *title)
}

// print a specific pdf page.
func printToPDF(urlstr string, res *[]byte, title *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitReady("/html/body/section/article/div/div/div[2]/div/div[2]/div", chromedp.BySearch),
		chromedp.Sleep(time.Second * 2),
		chromedp.Title(title),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("Screenshot")
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
