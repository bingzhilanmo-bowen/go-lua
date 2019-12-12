package service

import (
	"bingzhilanmo/go-lua/models"
	"bingzhilanmo/go-lua/pkg/utils"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type PageQueryBase struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}




func CacheHitKeyword(keyword string) *models.Keyword {
	kyCacheKey := utils.FormatCacheKey(keyword, utils.NLG_KEYWORD)

	ky, _ := utils.Get(kyCacheKey)

	if ky != nil {
		return ky.(*models.Keyword)
	}

	k, err := models.QueryByKeyword(keyword)
	if err != nil {
		log.Errorf(" %s keyword found  error %s .....", keyword, err.Error())
	}

	utils.Set(kyCacheKey, k, cache.NoExpiration)
	return k
}
