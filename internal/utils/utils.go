package utils

import (
	"fmt"
	"github.com/biter777/countries"
	"remnawave-tg-shop-bot/internal/config"
	"sort"
	"strings"
)

func BuildAvailableCountriesLists(langCode string) string {
	var locationsText strings.Builder
	countriesMap := config.Countries()

	keys := make([]string, 0, len(countriesMap))
	for k := range countriesMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		country := strings.Split(countriesMap[k], " ")
		if i == len(keys)-1 {
			if langCode == "ru" {
				locationsText.WriteString(fmt.Sprintf("└ %s%s\n", country[0], countries.ByName(country[1]).StringRus()))
			} else {
				locationsText.WriteString(fmt.Sprintf("└ %s%s\n", country[0], countries.ByName(country[1]).String()))
			}
		} else {
			if langCode == "ru" {
				locationsText.WriteString(fmt.Sprintf("├ %s%s\n", country[0], countries.ByName(country[1]).StringRus()))
			} else {
				locationsText.WriteString(fmt.Sprintf("├ %s%s\n", country[0], countries.ByName(country[1]).String()))
			}
		}
	}

	return locationsText.String()
}
