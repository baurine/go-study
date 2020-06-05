package model

import (
	"encoding/json"
	"fmt"
)

type ComicResponse struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Link       string `json:"link"`
	News       string `json:"news"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
}

type Comic struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Date        string `json:"date"`
}

func (cr ComicResponse) FormattedDate() string {
	return fmt.Sprintf("%s-%s-%s", cr.Day, cr.Month, cr.Year)
}

func (cr ComicResponse) Comic() Comic {
	return Comic{
		Title:       cr.Title,
		Number:      cr.Num,
		Date:        cr.FormattedDate(),
		Description: cr.Alt,
		Image:       cr.Img,
	}
}

func (c Comic) PrettyString() string {
	p := fmt.Sprintf(
		"Title: %s\nComic No: %d\nDate: %s\nDescription: %s\nImage: %s\n",
		c.Title, c.Number, c.Date, c.Description, c.Image)
	return p
}

func (c Comic) JSON() string {
	cJSON, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(cJSON)
}
