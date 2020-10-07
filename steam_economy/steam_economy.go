package steam_economy

import (
	"bytes"
	"encoding/json"
	"github.com/0x1a0b/hooked/config"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/d4l3k/messagediff.v1"
	"net/http"
	"regexp"
)

const (
	SteamApiEconomyBase = "https://api.steampowered.com/ISteamEconomy/GetAssetPrices/v1/"
	SteamRustGameId = "252490"
	SteamWebShop = "https://store.steampowered.com/itemstore/252490/browse/?filter=All"
)

var (
	currWebShopstate WebShopState
	currEconState = SteamApiResponse{Result: SteamApiResult{Success: false}}
)

func UpdateShop() () {
	newEconState, _ := GetEconomyResponse()
	newWebShopState := GetShopState()
	_, equal := messagediff.PrettyDiff(newWebShopState, currWebShopstate)
	if currEconState.Result.Success != true {
		log.Debugf("success is false, either there is a problem or a restart.. anyways...")
		currWebShopstate = newWebShopState
		currEconState = newEconState
		sendUpdate()
	} else if equal == true {
		log.Debugf("no change in steam econ")
	} else {
		log.Debugf("change in econ, updating")
		currWebShopstate = newWebShopState
		currEconState = newEconState
		sendUpdate()
	}
	return
}

func sendUpdate() () {
	for _, item := range currEconState.Result.Assets {
		id := item.Name
		var thisWebShopItem WebShopItem
		for _, wsItem := range currWebShopstate.Items {
			if wsItem.Id == id {
				thisWebShopItem = wsItem
			}
		}
		object := map[string]interface{}{
			"title": "New Skip: "+thisWebShopItem.Name,
			"url":  thisWebShopItem.Link,
			"color": 2724948,
			"fields": []interface{}{
				map[string]interface{}{
					"name": "Price Euro",
					"value": item.Prices["EUR"],
					"inline": true,
				},
				map[string]interface{}{
					"name": "Price CHF",
					"value": item.Prices["CHF"],
					"inline": true,
				},
			},
			"thumbnail": map[string]interface{}{
				"url": thisWebShopItem.Picture,
			},
			}

		o, _ := json.Marshal(object)
		http.Post(config.GetConf().Discord.Url, "application/json", bytes.NewBuffer(o))
		}
	return
	}

func GetShopState() (wss WebShopState) {
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
			var thisItem WebShopItem
			icon_container := itemhtml.Find(".item_def_icon_container").First()
				links := icon_container.Find("a")
				link, _ := links.First().Attr("href")
				thisItem.Link = link

			    re := regexp.MustCompile(`^https://store.steampowered.com/itemstore/252490/detail/(?P<id>[0-9]+)/$`)
			    id := re.FindStringSubmatch(link)[1]
			    thisItem.Id = id

				pics := icon_container.Find("img")
				pic, _ := pics.First().Attr("src")
				thisItem.Picture = pic

				textContainer := itemhtml.Find(".item_def_name").First()
				textLink := textContainer.Find("a").First()
				name := textLink.Contents().Text()
				thisItem.Name = name

				wss.Items = append(wss.Items, thisItem)
		})
	})
	return
}
type WebShopState struct {
	Items []WebShopItem
}
type WebShopItem struct {
	Id string
	Link string
	Name string
	Picture string
}

func GetEconomyResponse() (result SteamApiResponse, err error) {
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
	Classid string `json:"classid"`
}
type SteamApiAssetClass struct {
	Name string `json:"name"`
	Value string `json:"value"`
}