package model

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

type AnimeDetails struct {
	Url        string `json:"url"`
	Name       string `json:"name"`
	LinkExpire string `json:"link_expires"`
}

func (detail *AnimeDetails) SetExpireTime() {
	exp := strings.Split(detail.Url, "expires=")
	expInt, err := strconv.Atoi(exp[len(exp)-1])
	if err != nil {
		log.Fatalf("can't convert string %v", err)
	}
	unix := time.Unix(int64(expInt), 0)
	/// day month year hour second
	/// 05-12-2023 03:45 pm
	/// fmt.Println("Unix: ", unix, "UTC: ", unix.Format("Mon 02 Jan 06 3:04 pm"))
	detail.LinkExpire = unix.Format("Mon 02 Jan 06 3:04 pm")
}

func (detail *AnimeDetails) ToJson() []byte {
	buffer := &bytes.Buffer{}

	jsonE := json.NewEncoder(buffer)
	jsonE.SetIndent("", "  ")
	jsonE.SetEscapeHTML(false)
	err := jsonE.Encode(detail)
	if err != nil {
		log.Fatalf("can't marshal %v", err)
	}
	return buffer.Bytes()
}
