package models

import (
	"bingzhilanmo/go-lua/config"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
)

const (
	KEYWORD_NEW_ASR = iota
	KEYWORD_ADD_OLD
	KEYWORD_REPLACE
)


var DB *pg.DB

func InitDb() error{

	db := pg.Connect(&pg.Options{
		Addr: config.GetGlobalConfig().DB.Addr,
		User:     config.GetGlobalConfig().DB.User,
		Password: config.GetGlobalConfig().DB.Password,
		Database: config.GetGlobalConfig().DB.Database,
	})

	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err == nil {
		log.Debug("init pg db success .......")
		DB = db

		if config.GetGlobalConfig().DB.Init {
			err := CreatedTables()
			if err != nil {
				log.Debugf("created table error %s", err.Error())
			}
		}

	}

	return err
}

func OpenListener() {
	log.Info("open listener ......")
	ListenKeywordChange()
}

func CloseDb() error {
	return	DB.Close()
}

type Keyword struct {
	Id	 string	`pg:"id,pk,type:varchar(64)"`
	KeywordType int `pg:"keyword_type,type:int2"`
	KeywordName string `pg:"keyword_name,type:varchar(64)"`
	ScriptText string `pg:"script_text,type:text"`
	CreateTime int64 `pg:"created_time,type:int8"`
	CreateUser string `pg:"created_user,type:varchar(64)"`
	UpdateTime int64 `pg:"update_time,type:int8"`
	UpdateUser string `pg:"updated_user,type:varchar(64)"`
}

func (t Keyword) IsNotNull() bool {
	return t.Id != "" && t.KeywordName != "" && t.ScriptText != ""
}

func (t Keyword) IsKeywordReplace() bool  {
	return t.KeywordType == KEYWORD_REPLACE
}

func (t Keyword) IsKeywordNewAsr() bool  {
	return t.KeywordType == KEYWORD_NEW_ASR
}

func (t Keyword) IsAddOld() bool  {
	return t.KeywordType == KEYWORD_ADD_OLD
}


func BuildLike(in string) string {
	return "%"+in+"%"
}


func CreatedTables() error {
	for _, model := range []interface{}{&Keyword{}} {
		err := DB.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:true,
			FKConstraints:false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}