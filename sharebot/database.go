package sharebot

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"
	"strings"
)

type cardFilePath string
type ShopsDataBase struct {
	directoryPath string
	ShopList     []string
	allowedCards map[string]cardFilePath
}

func NewShopsDataBase(directoryPath string) ShopsDataBase {
	// update func (?)

	allowedCardsfile, err := os.Open(directoryPath + "/content.json")

	if err != nil {
		log.Fatal(err)
	}
	defer allowedCardsfile.Close()
	var buf []struct {
		ShopName string `json:"shopname"`
		CardPath string `json:"cardpath"`
	}

	shopList := make([]string, 0)
	allowedCards := make(map[string]cardFilePath)

	dec := json.NewDecoder(bufio.NewReader(allowedCardsfile))

	for dec.More() {
		err := dec.Decode(&buf)
		if err != nil {
			panic(err)
		}

		for _, elem := range buf {
			shopList = append(shopList, elem.ShopName)
			allowedCards[elem.ShopName] = cardFilePath(elem.CardPath)
		}
	}

	sort.Strings(shopList)

	return ShopsDataBase{
		directoryPath: directoryPath,
		ShopList:     shopList,
		allowedCards: allowedCards,
	}
}

type shopCard []byte

func (db ShopsDataBase) FindShopCard(name string) (shopCard, error) {
	cardFilename, ok := db.allowedCards[name]
	if ok {
		res, err := os.ReadFile(strings.Join([]string{db.directoryPath, string(cardFilename)}, "/"))
		if err != nil {
			log.Panic(cardFilename + " undefined")
		} 
		return res, nil
	} else {
		return nil, errors.New(name + " is undefined")
	}
}
