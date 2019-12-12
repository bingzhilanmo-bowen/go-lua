package models

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func QueryByKeyword(kd string) (*Keyword, error) {
	var keyword Keyword
	err := DB.Model(&keyword).Where("keyword_name = ? ", kd).Select()
	return &keyword, err
}

func QueryByKeywordId(id string) (*Keyword, error) {
	var keyword Keyword
	err := DB.Model(&keyword).Where("id = ? ", id).Select()
	return &keyword, err
}

func CreateNewkeyword(newKeyword *Keyword) error {
	_, err := DB.Model(newKeyword).Insert()
	if err == nil {
		log.Infof("send notify channel %s , payload  %s", KEYWORD_CHAN, newKeyword.Id)
		sendNotify(newKeyword.Id)
	}
	return err
}

func Updatekeyword(newKeyword *Keyword) error {
	err := DB.Update(newKeyword)
	if err == nil {
		log.Infof("send notify channel %s , payload  %s", KEYWORD_CHAN, newKeyword.Id)
		sendNotify(newKeyword.Id)
	}
	return err
}

func ExistsByKeyword(kd string) (bool, error) {
	var keyword Keyword
	bol, err := DB.Model(&keyword).Where("keyword_name = ? ", kd).Exists()
	return bol, err
}

func QueryAllKeywords() ([]Keyword, error) {
	var keys []Keyword
	err := DB.Model(&keys).Select()
	return keys, err
}

func QueryPageKeyword(pageNo, pageSize int, keyword string) ([]Keyword, int, error) {
	var keys []Keyword

	q := DB.Model(&keys)

	if keyword != "" {
		q = q.Where("keyword_name LIKE ?", BuildLike(keyword))
	}

	count, err := q.Count()

	if err != nil {
		log.Errorf("query keyword page list error %s", err.Error())
		return nil, 0, err
	}

	if count <= 0 {
		return nil, 0, nil
	}

	if pageNo > 0 && pageSize > 0 {
		q = q.Limit(pageSize).Offset((pageNo - 1) * pageSize)
	}
	q = q.OrderExpr("created_time DESC")

	err = q.Select()

	if err != nil {
		log.Errorf("query keyword page list error %s", err.Error())
		return nil, 0, err
	}

	return keys, count, nil
}

func sendNotify(id string) {
	execInfo := fmt.Sprintf("NOTIFY %s, ?", KEYWORD_CHAN)
	log.Infof("send notify channel %s , payload  %s", KEYWORD_CHAN, id)
	res, err := DB.Exec(execInfo, id)
	if err != nil {
		log.Errorf("send notify id [%s], error %s", id, err.Error())
	}
	log.Info(res)
}
