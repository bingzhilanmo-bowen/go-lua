package utils

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"regexp"
	"runtime"
	"runtime/pprof"
	"time"
)



func ConsumeMem() uint64 {
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Sys
}

func DoResourceMonitor() {
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				m := pprof.Lookup("goroutine")
				memStats := ConsumeMem()
				log.Infof("Resource monitor: [%d goroutines] [%.3f kb]", m.Count(), float64(memStats)/1e3)
			}
		}
	}()
}

func Uuid() string {
   uid,err := uuid.NewUUID()

   if err == nil {
	   return uid.String()
   }

   return ""
}

func GetExprList(findStr ,exprStr string) []string {

	match,err := regexp.Compile(exprStr)
	if err != nil {
		log.Errorf("regexp.Compile error %s ", err.Error())
		return nil
	}

	return  match.FindAllString(findStr,-1)
}