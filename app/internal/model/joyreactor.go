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

const MaxPostsOnPage = 10

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

func GetBoobs(page int, post int, tag string, pages int) (string, error) {
	if page < 1 || page > pages {
		return "", ErrInvalidPageNumber
	}
	if post < 1 || post > MaxPostsOnPage {
		return "", ErrInvalidPostNumber
	}

	joyUrl := fmt.Sprintf("http://joyreactor.cc/tag/%s/%d", tag, page)
	doc, err := GetDocumentFromURL(joyUrl)

	if err != nil {
		return "", err
	}

	var result string
	doc.Find(".postContainer").Each(func(i int, s *goquery.Selection) {
		if i != (post - 1) {
			return
		}
		attr, exists := s.Find("a > img").Last().Attr("src")
		if !exists {
			err = ErrNoImage
			return
		}

		result = attr
	})
	return result, err
}

func GetRandomBoobs(tag string) (string, error) {
	pages, err := GetPagesCount(tag)
	if err != nil {
		return "", ErrNoMaxPages
	}

	for {
		page := rand.Intn(pages) + 1
		post := rand.Intn(9) + 1
		joyUrl, err := GetBoobs(page, post, tag, pages)
		if err != ErrNoImage && err != ErrPageNotFound && err != nil {
			return "", err
		}
		if err != nil {
			continue
		}
		return joyUrl, nil
	}
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
