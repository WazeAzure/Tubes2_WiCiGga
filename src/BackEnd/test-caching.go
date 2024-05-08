package main

import (
	"backend/caching"
	"backend/scraper"
)

func main() {
	caching.CheckCacheFolder()
	caching.InitCache()
	// base_url := "https://en.wikipedia.org/wiki/"
	// title := []string{"Donald_Trump", "Barack_Obama", "Jesus", "Kevin_Bacon", "Philosophy", "Elon_Musk", "Bitcoin", "Death_Grips", "Indonesia", "Hinduism", "Bali", "2004_Indian_Ocean_earthquake_and_tsunami", "Komodo_dragon", "Indonesia_national_football_team", "2024_Indonesian_general_election", "East_Timor", "G20", "Boeing_737_MAX_groundings", "Pacific_War", "Orangutan", "Jakarta", "Sama-Bajau", "Ferdinand_Magellan", "Prabowo_Subianto", "Comfort_women", "Durian"}

	// title := []string{"Nyepi", "Joe_Taslim", "Sukarno", "Suharto", "Sumatra", "Kalimantan", "Java", "Borneo", "Sulawesi", "Bandung", "Balikpapan", "Film"}

	// for idx := range title {
	// fmt.Println("=========== CURRENTLY :", title[idx])
	// url := base_url + title[idx]
	// fmt.Println(caching.GetKeyHash(url))

	// caching.DeleteCache(url)
	// list_url := scraper.ScrapWeb(url)
	// for idx := range list_url {
	// 	fmt.Println(list_url[idx])
	// }
	// caching.SetCacheUrl(url, list_url)
	// }

	url := "https://en.wikipedia.org/wiki/%C4%9E"
	// key := caching.GetKeyHash(url)
	caching.DeleteCache(url)
	// caching.CachedWebpage.Delete([]byte(key))

	list_url := scraper.ScrapWeb(url)
	caching.SetCacheUrl(url, list_url)

}
