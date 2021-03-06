package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/myutil"
)

type Board struct {
}

var articlesDir string = myutil.StoragePath() + "/articles/"

func (bd Board) List() []string {
	files, err := ioutil.ReadDir(articlesDir)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	var boardNames []string
	for _, file := range files {
		name, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		boardNames = append(boardNames, name)
	}
	return boardNames
}

func (bd Board) Exist(boardName string) bool {
	file := articlesDir + boardName + ".json"
	_, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}
	return true
}

func (bd Board) GetArticles(boardName string) article.Articles {
	file := articlesDir + boardName + ".json"
	articlesJSON, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    file,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
	}
	articles := make(article.Articles, 0)
	err = json.Unmarshal(articlesJSON, &articles)
	if err != nil {
		myutil.LogJSONDecode(err, articlesJSON)
	}
	return articles
}

func (bd Board) Create(boardName string) error {
	err := ioutil.WriteFile(articlesDir+boardName+".json", []byte("[]"), 664)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (bd Board) Save(boardName string, articles article.Articles) error {
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		myutil.LogJSONEncode(err, articles)
	}
	err = ioutil.WriteFile(articlesDir+boardName+".json", articlesJSON, 0644)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (bd Board) Delete(boardName string) error {
	err := os.Remove(articlesDir + boardName + ".json")
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
