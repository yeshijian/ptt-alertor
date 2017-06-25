package top

import (
	"strings"

	"strconv"

	log "github.com/meifamily/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/myutil"
)

const prefix string = "top:"

type BoardWord struct {
	Board, Word string
}

type WordOrder struct {
	BoardWord
	Count int
}

type WordOrders []WordOrder

func (wos WordOrders) SaveKeywords() error {
	return wos.save("keywords")
}

func (wos WordOrders) SaveAuthors() error {
	return wos.save("authors")
}

func (wos WordOrders) save(kind string) error {
	conn := connections.Redis()
	for _, wo := range wos {
		if _, err := conn.Do("ZADD", prefix+kind, wo.Count, wo.String()); err != nil {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
			return err
		}
	}
	return nil
}

func ListAuthors(num int) []string {
	return list("authors", num)
}

func ListKeywords(num int) []string {
	return list("keywords", num)
}

func list(kind string, num int) []string {
	conn := connections.Redis()
	lists, err := redis.Strings(conn.Do("ZREVRANGE", prefix+kind, 0, num-1))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return lists
}

func ListKeywordWithScore(num int) WordOrders {
	return listWithScore("keywords", num)
}

func ListAuthorWithScore(num int) WordOrders {
	return listWithScore("authors", num)
}

func listWithScore(kind string, num int) (wos WordOrders) {
	conn := connections.Redis()
	list, err := redis.Strings(conn.Do("ZREVRANGE", prefix+kind, 0, num-1, "WITHSCORES"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	bw := BoardWord{}
	for index, element := range list {
		if index%2 == 0 {
			bw = BoardWord{}
			strs := strings.Split(element, ":")
			bw.Board = strs[0]
			bw.Word = strs[1]
			continue
		}
		count, _ := strconv.Atoi(element)
		wo := WordOrder{bw, count}
		wos = append(wos, wo)
	}
	return wos
}

func (wo WordOrder) String() string {
	return wo.Board + ": " + wo.Word
}
