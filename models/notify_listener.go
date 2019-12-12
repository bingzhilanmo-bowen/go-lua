package models

import (
	"bingzhilanmo/go-lua/pkg/utils"
	"github.com/go-pg/pg/v9"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const (
	KEYWORD_CHAN = "nlg_keywords_chan"
	TEMPLATE_CHAN = "nlg_template_chan"
	BIND_CHAN = "nlg_bind_chan"
)

func ListenKeywordChange() {
	go func(c <-chan *pg.Notification) {
		for notify := range c {
			log.Infof("keyword change listener %s", notify.Payload)
			changeId := notify.Payload
			changeKeyword,_ := QueryByKeywordId(changeId)

			if changeKeyword != nil {
				kyCacheKey := utils.FormatCacheKey(changeKeyword.KeywordName, utils.NLG_KEYWORD)
				utils.Set(kyCacheKey, changeKeyword, cache.NoExpiration)
				utils.WarmLuaVm(changeKeyword.KeywordName, changeKeyword.ScriptText)
			}
		}
	}(DB.Listen(KEYWORD_CHAN).Channel())
}

