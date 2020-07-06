package main

import (
	"log"
	"sync"
	"syncTool/src/pkg/config"
	"syncTool/src/repositories"
	"syncTool/tools"
)

var queryLimit = 100
var maxThread = 5
var updateLimit = queryLimit / maxThread

func init() {
	repositories.InitPostgres()
}

func main() {
	queryLimit = config.GetInteger("config.queryLimit")
	maxThread = config.GetInteger("config.maxThread")
	updateLimit = queryLimit / maxThread

	idsTube := make(chan []string)
	exitTube := make(chan string)
	notifyTube := make(chan string)
	times := 0

	go func() {
		for {
			wg := sync.WaitGroup{}
			wg.Add(1)

			ids, err := repositories.GetLegacyIDs(queryLimit)
			if err != nil {
				log.Println("batch GetLegacyIDs terminate: err=", err)
				break
			}

			if len(ids) == 0 {
				log.Printf("batch GetLegacyIDs terminate: have no user to execute")
				break
			} else {
				times += queryLimit
				log.Printf("batch GetLegacyIDs running in %d times", times)
				idsTube <- ids

				select {
				case <- notifyTube:
					wg.Done()
				}
			}
			wg.Wait()
		}
	}()

	for {
		select {
		case <- exitTube:
			return
		default:
			ids := <- idsTube
			wg := sync.WaitGroup{}
			wg.Add(maxThread)

			startIndex := 0
			endIndex := updateLimit
			for i := 0; i < maxThread; i++ {
				idsForGoroutine := ids[startIndex:endIndex]
				go func(ids []string) {
					t1 := tools.GetTimestamp()
					userInfo, err := repositories.GetPortiereUserInfo(ids)
					if err != nil {
						log.Println("batch GetPortiereUserInfo terminate: err=", err)
					}
					if userInfo != nil {
						repositories.MigrationPortiereUsersInfoToHermes(userInfo)
					}
					t2 := tools.GetTimestamp()
					log.Println("UpdateUsers cost ",  t2-t1, " (ns)")
					wg.Done()
				}(idsForGoroutine)

				startIndex += updateLimit
				endIndex += updateLimit
			}

			wg.Wait()
			notifyTube <- ""
		}
	}
}