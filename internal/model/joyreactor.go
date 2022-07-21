package model

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrInvalidUrl        = errors.New("ошибка получения страницы по заданному URL")
	ErrPageNotFound      = errors.New("страница не найдена")
	ErrNoImage           = errors.New("нет картинки")
	ErrInvalidPageNumber = errors.New("некорректный номер страницы")
	ErrInvalidPostNumber = errors.New("некорректный номер поста")
	ErrConnectionFailed  = errors.New("ошибка подключения")
	ErrNoMaxPages        = errors.New("ошибка получения максимальной страницы")
)

const (
	MaxPostsOnPage = 10

	urlPictureDraw = "http://joyreactor.cc/tag"
)

type JoyLoader struct {
	Image       string
	Description string
}

func GetDocumentFromURL(joyUrl string) (*goquery.Document, error) {
	joyUrl, err := url.QueryUnescape(joyUrl)
	if err != nil {
		return nil, ErrInvalidUrl
	}
	getPage, err := http.Get(joyUrl)
	if err != nil {
		return nil, ErrInvalidUrl
	}
	defer getPage.Body.Close()
	if getPage.StatusCode != http.StatusOK {
		return nil, ErrPageNotFound
	}

	page, err := goquery.NewDocumentFromReader(getPage.Body)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func GetRandomPictures(tag string) (*JoyLoader, error) {
	pages, err := GetPagesCount(tag)
	if err != nil {
		return nil, ErrNoMaxPages
	}

	for {
		page := rand.Intn(pages) + 1
		post := rand.Intn(9) + 1
		joy, err := GetJoy(page, post, tag, pages)
		if err != ErrNoImage && err != ErrPageNotFound && err != nil {
			return nil, err
		}
		if err != nil {
			continue
		}
		return joy, nil
	}
}

func GetJoy(page int, post int, tag string, pages int) (*JoyLoader, error) {
	if page < 1 || page > pages {
		return nil, ErrInvalidPageNumber
	}
	if post < 1 || post > MaxPostsOnPage {
		return nil, ErrInvalidPostNumber
	}

	joyUrl := fmt.Sprintf("%s/%s/%d", urlPictureDraw, tag, page)
	doc, err := GetDocumentFromURL(joyUrl)

	if err != nil {
		return nil, err
	}
	loader := &JoyLoader{}
	doc.Find(".postContainer").Each(func(i int, s *goquery.Selection) {
		if i != (post - 1) {
			return
		}
		findSelector := s.Find("a > img").Last()
		attr, exists := findSelector.Attr("src")
		if !exists {
			err = ErrNoImage
			return
		}
		loader.Image = fmt.Sprintf("https:%s", attr)
		if attr, exists = findSelector.Attr("alt"); exists {
			loader.Description = attr
		}

	})
	return loader, err
}

func DownloadFile(joyUrl string) (string, error) {
	resp, err := http.Get(joyUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, err := ioutil.TempFile("", "*.jpeg")
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return out.Name(), err
}

func GetPagesCount(tag string) (int, error) {
	joyUrl := fmt.Sprintf("http://joyreactor.cc/tag/%s", tag)
	doc, err := GetDocumentFromURL(joyUrl)
	if err != nil {
		return 0, ErrConnectionFailed
	}
	pages, err := strconv.Atoi(doc.Find(".pagination_expanded .current").Text())
	if err != nil {
		return 0, ErrNoMaxPages
	}
	return pages, nil
}
