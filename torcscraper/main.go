package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"

	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/chromedp"
	_ "github.com/lib/pq"
)

func kategori(text string) string {
	t := strings.ToLower(text)
	if strings.Contains(t, "cvv") {
		return "cc"
	}
	if strings.Contains(t, "exploit") {
		return "exploit"
	}
	if strings.Contains(t, "malware") {
		return "malware"
	}
	return "genel"
}

func Kritiklik(content string, pageURL string) string {
	if strings.Contains(strings.ToLower(content), "cvv") {
		return "yüksek"
	}
	if strings.Contains(strings.ToLower(content), "malware") {
		return "orta"
	}
	if strings.Contains(strings.ToLower(content), "combolist") {
		return "orta"
	}
	if strings.Contains(strings.ToLower(content), "database") {
		return "orta"
	}
	return "düşük"
}

func main() {

	logfile, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(logfile)
		defer logfile.Close()
	}

	outputDir := "scans"
	os.MkdirAll(outputDir, 0755)

	connStr := "host=findings-postgres port=5432 user=postgres password=sifreniz dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS findings (
    id SERIAL PRIMARY KEY,
    url TEXT,
    title TEXT,
    source_title TEXT,
    category TEXT,
    criticality TEXT,
    content TEXT,
    screenshot TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		log.Fatal("Tablo oluşturulamadı: ", err)
	}
	defer db.Close()
	go StartWebServer()

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("/usr/bin/chromium-browser"),
			chromedp.ProxyServer("socks5://tor:9050"),
			chromedp.NoSandbox,
			chromedp.DisableGPU,
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Headless)...)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	file, err := os.Open("targets.yaml")
	if err != nil {
		log.Fatal("targets.yaml eksik")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		fmt.Printf("Tarama Başladı: %s\n", url)
		log.Printf("İşlem başlatıldı: %s", url)

		var buf []byte
		var hamhtml, kaynak string
		var baslik string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitReady("body"),
			chromedp.FullScreenshot(&buf, 90),
			chromedp.Evaluate(`document.title`, &kaynak),
			chromedp.Evaluate(`document.body.innerText`, &hamhtml),
			chromedp.Evaluate(`document.querySelector("h1")?.innerText`, &baslik),
		)

		if err != nil {
			log.Printf("HATA (%s): %v", url, err)
			continue
		}

		cleanName := strings.NewReplacer("http://", "", "https://", "", "/", "_", ".", "_").Replace(url)
		imgPath := filepath.Join(outputDir, cleanName+".png")
		htmlPath := filepath.Join(outputDir, cleanName+".html")

		os.WriteFile(imgPath, buf, 0644)
		os.WriteFile(htmlPath, []byte(hamhtml), 0644)

		record := Record{
			URL: url, Title: baslik, SourceTitle: kaynak, Category: kategori(hamhtml),
			Criticality: Kritiklik(hamhtml, url), Content: hamhtml, Screenshot: imgPath,
		}
		if err := insertRecord(db, record); err != nil {
			log.Printf("DB Hatası: %v", err)
		} else {
			log.Printf("Başarıyla kaydedildi: %s", url)
		}
	}
	select {}
}
