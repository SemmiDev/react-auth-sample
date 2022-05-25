package main

import (
	"log"
	"testing"
	"time"
)

func Test_timezone(t *testing.T) {
	now := time.Now().UTC().Unix()
	loc, _ := time.LoadLocation("Asia/Jakarta")

	log.Println(time.Unix(now, 0).In(loc))
	loc, _ = time.LoadLocation("Asia/Seoul")
	log.Println(time.Unix(now, 0).In(loc))
}

// list of timezone
// https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
