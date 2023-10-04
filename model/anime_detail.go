package model

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"paheScraper/config"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AnimeDetails struct {
	Url        string `json:"url"`
	Name       string `json:"name"`
	Episode    string `json:"episode"`
	LinkExpire string `json:"link_expires"`
}

func (detail *AnimeDetails) SetExpireTime() {
	exp := strings.Split(detail.Url, "expires=")
	expInt, err := strconv.Atoi(exp[len(exp)-1])
	if err != nil {
		log.Printf("can't convert string %v", err)
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
		log.Printf("can't marshal %v", err)
	}
	return buffer.Bytes()
}

func (detail *AnimeDetails) SaveToFile() {
	filePath := filepath.Join(config.UserDocument, strings.ToValidUTF8(detail.Name,
		"_")+".json")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("can't open file %s %v", detail.Name, err)
	}
	_, err = file.Write(detail.ToJson())
	if err != nil {
		log.Printf("can't write details to file %v", err)
	}
}
