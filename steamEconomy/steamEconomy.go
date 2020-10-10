package steamEconomy

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/0x1a0b/hooked/discordSender"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/d4l3k/messagediff.v1"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	SteamApiEconomyBase = "https://api.steampowered.com/ISteamEconomy/GetAssetPrices/v1/"
	SteamRustGameId = "252490"
	SteamWebShop = "https://store.steampowered.com/itemstore/252490/browse/?filter=All"
)

type Instance struct {
	sender *discordSender.Sender
	logger *logrus.Logger
	currWebShopstate WebShopState
}

func Setup() (i *Instance) {

	i = &Instance{
		currWebShopstate: WebShopState{
			Items: []WebShopItem{},
		},
	}

	secret := config.GetConf().Hooks.RustNewItems
	i.sender = discordSender.New(secret)

	i.logger = logrus.New()
	i.logger.SetLevel(config.GetLogLevel())
	i.logger.SetReportCaller(true)

	return
}

func (i *Instance) Run() () {
	newWebShopState := i.GetShopState()
	_, equal := messagediff.PrettyDiff(newWebShopState, i.currWebShopstate)
	if len(i.currWebShopstate.Items) == 0 {
		i.logger.Debugf("items array empty, resuming from restart")
		i.currWebShopstate = newWebShopState
		i.FireAll()
	} else if equal == true {
		i.logger.Debugf("same econ, doing nothing")
	} else {
		i.logger.Debugf("detected delta")
		i.currWebShopstate = newWebShopState
		i.FireAll()
	}
	return
}


func (i *Instance) FireAll() () {
	econ, err := i.GetEconomyResponse()
	if err != nil {
		i.logger.Errorf("could not get econ status: %v", err)
		return
	}

	for _, item := range econ.Result.Assets {
		for _, wsItem := range i.currWebShopstate.Items {
			if wsItem.Id == item.Name {
				i.FireOne(item, wsItem)
				time.Sleep(2*time.Second)
			}
		}
	}
	return

}

func (i *Instance) FireOne(econItem SteamApiItemAsset, shopItem WebShopItem) {
	hook := i.SetHook(econItem, shopItem)
	if err := i.sender.Send(hook); err != nil {
		i.logger.WithField("hook", hook).Errorf("error sending hook: %v", err)
	} else {
		i.logger.WithField("hook", hook).Debugf("fired hook")
	}
	return
}

func (i *Instance) SetHook(econItem SteamApiItemAsset, shopItem WebShopItem) (h discordSender.Hook) {
	h = discordSender.Hook{
		Content: "New Skin: " + shopItem.Name,
		AvatarUrl: "https://rust.facepunch.com/dist/img/logo-face.svg",
		Username: "Mc Skinner",
		Embeds: []discordSender.Embed{
			{
				Title: "New Skin Available: " + shopItem.Name,
				Url: shopItem.Link,
				Color: 3092790,
				Thumbnail: discordSender.Thumbnail{
					Url: shopItem.Picture,
				},
				Fields: []discordSender.Field{
					{
						Name: "Price CHF",
						Value: strconv.Itoa(econItem.Prices["CHF"]),
						Inline: true,
					},
				},
				Footer: discordSender.Footer{
					Text: "Mc Skinner Bot",
				},
			},
		},
	}
	return
}

func (i *Instance) GetShopState() (wss WebShopState) {
	res, err := http.Get(SteamWebShop)
	if err != nil {
		i.logger.Errorf("error %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		i.logger.Errorf("error %v", err)
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

func (i *Instance) GetEconomyResponse() (result SteamApiResponse, err error) {
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
		i.logger.Errorf("no resty error, got status code %v from client", resp.StatusCode())
	} else {
		i.logger.Errorf("got resty error: %v", err)
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
