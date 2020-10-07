package main

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
)

const (
	SteamApiEconomyBase = "https://api.steampowered.com/ISteamEconomy/GetAssetPrices/v1/"
	SteamRustGameId = "252490"
	SteamWebShop = "https://store.steampowered.com/itemstore/252490/browse/?filter=All"
)

func main() ( ) {
	res, err := http.Get(SteamWebShop)
	if err != nil {
		log.Errorf("error %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("error %v", err)
	}

	doc.Find("#ItemDefsRows").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find(".item_def_grid_item").Each(func(indextr int, itemhtml *goquery.Selection) {
			icon_container := itemhtml.Find(".item_def_icon_container").First()
				links := icon_container.Find("a")
				link, _ := links.First().Attr("href")
						log.Errorf("asf %v", link)
			    re := regexp.MustCompile(`^https://store.steampowered.com/itemstore/252490/detail/(?P<id>[0-9]+)/$`)
			    id := re.FindStringSubmatch(link)[1]
			    log.Errorf("if %v", id)
				pics := icon_container.Find("img")
				pic, _ := pics.First().Attr("src")
				log.Errorf("asf %v", pic)

				textContainer := itemhtml.Find(".item_def_name").First()
				textLink := textContainer.Find("a").First()
				name := textLink.Contents().Text()
				log.Errorf("%v", name)

		})
	})

}

func GetEconomyResponse(key string, appid string) (result SteamApiResponse, err error) {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"key": config.GetConf().Steam.ApiKey,
			"appid": SteamRustGameId,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(SteamApiEconomyBase)
	if err == nil {
		log.Debugf("no resty error, got status code %v from client", resp.StatusCode())
	} else {
		log.Errorf("got resty error: %v", err)
	}
	return
}
type SteamApiResponse struct {
	Result SteamApiResult `json:"result"`
}
type SteamApiResult struct {
	Success bool `json:"success"`
	Assets []SteamApiItemAsset `json:"assets"`
}
type SteamApiItemAsset struct {
	Prices map[string]int `json:"prices"`
	OriginalPrices map[string]int `json:"original_prices"`
	Name string `json:"name"`
	Class []SteamApiAssetClass `json:"class"`
	Classid string
}
type SteamApiAssetClass struct {
	Name string `json:"name"`
	Value string `json:"value"`
}