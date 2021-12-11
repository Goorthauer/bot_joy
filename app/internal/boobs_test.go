package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
Тесты получения картинки по заданным параметрам
*/
func TestGetBoobs(t *testing.T) {
	tests := []struct {
		name string
		page int
		post int
		url  string
		err  error
	}{
		{
			name: "success",
			page: 1,
			post: 10,
			url:  "http://img10.joyreactor.cc/pics/post/%D0%AD%D1%80%D0%BE%D1%82%D0%B8%D0%BA%D0%B0-%D1%81%D0%B8%D1%81%D1%8C%D0%BA%D0%B8-81.jpeg",
			err:  nil,
		},
		{
			name: "pageError",
			page: 0,
			post: 10,
			url:  "",
			err:  ErrInvalidPageNumber,
		},
		{
			name: "postError",
			page: 1,
			post: 11,
			url:  "",
			err:  ErrInvalidPostNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := GetBoobs(tt.page, tt.post, "%D1%81%D0%B8%D1%81%D1%8C%D0%BA%D0%B8", 2000)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.url, gotURL)
		})
	}
}

/**
Тесты получения документа по URL
*/
func TestGetDocumentFromURL(t *testing.T) {
	tests := []struct {
		name   string
		joyUrl string
		err    error
	}{
		{
			name:   "success",
			joyUrl: "http://joyreactor.cc/tag/сиськи",
			err:    nil,
		},
		{
			name:   "invalidUrl",
			joyUrl: "http://test",
			err:    ErrInvalidUrl,
		},
		{
			name:   "pageNotFound",
			joyUrl: "http://joyreactor.cc/test",
			err:    ErrPageNotFound,
		},
	}

	for _, getDocumentTest := range tests {
		t.Run(getDocumentTest.name, func(t *testing.T) {
			_, err := GetDocumentFromURL(getDocumentTest.joyUrl)
			assert.Equal(t, getDocumentTest.err, err)
		})
	}
}

/**
Тесты получения случайного поста с сиськами
*/
func TestGetRandomBoobs(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		err  error
	}{
		{
			name: "success",
			tag:  "сиськи",
			err:  nil,
		},
		{
			name: "noMaxPage",
			tag:  "%!invalidTag%!",
			err:  ErrNoMaxPages,
		},
	}

	for _, getRandomBoobs := range tests {
		t.Run(getRandomBoobs.name, func(t *testing.T) {
			_, err := GetRandomBoobs(getRandomBoobs.tag)
			assert.Equal(t, getRandomBoobs.err, err)
		})
	}
}

/**
Тесты получения количества страниц
*/
func TestGetPagesCount(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		err  error
	}{
		{
			name: "success",
			tag:  "сиськи",
			err:  nil,
		},
		{
			name: "connectionFailed",
			tag:  "%!invalidTag%!",
			err:  ErrConnectionFailed,
		},
		{
			name: "noMaxPage",
			tag:  "joyre",
			err:  ErrNoMaxPages,
		},
	}

	for _, getPagesCount := range tests {
		t.Run(getPagesCount.name, func(t *testing.T) {
			_, err := GetPagesCount(getPagesCount.tag)
			assert.Equal(t, getPagesCount.err, err)
		})
	}
}
