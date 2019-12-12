package service

import (
	"bingzhilanmo/go-lua/models"
	"bingzhilanmo/go-lua/pkg/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"time"
)

type KeywordDto struct {
	Id          string `json:"id"`
	Keyword     string `json:"keyword"`
	ScriptText  string `json:"script_text"`
	KeywordType int    `json:"keyword_type"`
}

type KeywordQueryDto struct {
	PageQueryBase
	Keyword string `json:"keyword"`
}

func CreateNewKeyword(dto *KeywordDto) (*models.Keyword, error) {

	exists, err := models.ExistsByKeyword(dto.Keyword)

	if err != nil {
		log.Errorf("query keyword script by keyword %s error %s", dto.Keyword, err.Error())
		return nil, err
	}

	if exists {
		return nil, errors.New(fmt.Sprintf("keyword has %s exists !!!", dto.Keyword))
	}

	newScript := &models.Keyword{
		KeywordName: dto.Keyword,
		ScriptText:  dto.ScriptText,
		CreateTime:  time.Now().Unix(),
		KeywordType: dto.KeywordType,
		Id:          utils.Uuid(),
	}

	err_create := models.CreateNewkeyword(newScript)
	if err_create == nil {
		return newScript, nil
	}

	log.Errorf("create new keyword  %v ,error %s", dto, err.Error())
	return nil, err_create
}

func UpdateKeyword(dto *KeywordDto) (*models.Keyword, error) {

	old, err := models.QueryByKeywordId(dto.Id)

	if err != nil {
		log.Errorf("query keyword script by id %s error %s", dto.Keyword, err.Error())
		return nil, err
	}

	if old == nil {
		return nil, errors.New(fmt.Sprintf("old keyword %s not exists !!!", dto.Keyword))
	}

	if dto.Keyword != old.KeywordName {
		//说明改了 name 需要检查
		exists, err := models.ExistsByKeyword(dto.Keyword)

		if err != nil {
			log.Errorf("query keyword script by keyword %s error %s", dto.Keyword, err.Error())
			return nil, err
		}

		if exists {
			return nil, errors.New(fmt.Sprintf("keyword has %s exists !!!", dto.Keyword))
		}
	}

	old.KeywordName = dto.Keyword
	old.ScriptText = dto.ScriptText
	old.KeywordType = dto.KeywordType
	old.UpdateTime = time.Now().Unix()

	err_update := models.Updatekeyword(old)
	if err_update == nil {
		return old, nil
	}

	log.Errorf("update keyword  %v ,error %s", old, err.Error())
	return nil, err_update
}

func QueryKeywordPage(query *KeywordQueryDto) ([]KeywordDto, int) {

	list, count, err := models.QueryPageKeyword(query.PageNo, query.PageSize, query.Keyword)

	if err != nil {
		log.Errorf("query keyword page error %s", err.Error())
		return nil, 0
	}

	if count <= 0 {
		return nil, count
	}

	le := len(list)

	dtoList := make([]KeywordDto, le)

	for i, va := range list {
		dto := &KeywordDto{
			Id:          va.Id,
			Keyword:     va.KeywordName,
			ScriptText:  va.ScriptText,
			KeywordType: va.KeywordType,
		}
		dtoList[i] = *dto
	}

	return dtoList, count
}

func QueryKeywordDetail(id string) (*KeywordDto, error) {
	detail, err := models.QueryByKeywordId(id)
	if err != nil {
		log.Errorf("query keyword detail error %s", err.Error())
		return nil, err
	}

	dto := &KeywordDto{
		Id:          detail.Id,
		Keyword:     detail.KeywordName,
		ScriptText:  detail.ScriptText,
		KeywordType: detail.KeywordType,
	}

	return dto, nil
}

func WarmLuaVm() {
	allKeys, err := models.QueryAllKeywords()
	if err != nil {
		log.Errorf("warm lua vm error %s", err.Error())
	}

	log.Infof("warm lua vm cache %d", len(allKeys))

	for _, va := range allKeys {
		utils.WarmLuaVm(va.KeywordName, va.ScriptText)
	}
}

func RunKeywordScript(keyword string, params map[string]string) (string, error) {

	k := CacheHitKeyword(keyword)

	if k != nil {
		va, err := utils.CacheExecuteWithMap(k.KeywordName, k.ScriptText, "run", params)
		res, _ := va.(lua.LString)
		return res.String(),err

	}

	return "", errors.New("keyword not found")
}



