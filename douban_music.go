package main

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"regexp"
	"os"
	"strings"
	id3 "github.com/mikkyang/id3-go"
	v2 "github.com/mikkyang/id3-go/v2"
//	iconv "github.com/djimenez/iconv-go"
)

func main() {
	body := get("http://site.douban.com/Ceekay/")
	os.Mkdir("music", 0755)
	
	re0 := regexp.MustCompile(`(?s)<div class="sp-nav">.*?<div class="sp-logo">.*?<img.*?alt="(.*?)"`)
	matches0 := re0.FindStringSubmatch(body)
	name := matches0[1]
	fmt.Println(name)
	os.Mkdir("music/" + name, 0755)
	
	re := regexp.MustCompile(`(?s)<div class="mod" id="playlist-(\d+)">.*?<span>(.*?)</span>`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, m := range matches {
		fmt.Println("\t" + m[2])
		os.Mkdir("music/" + name + "/" + m[2], 0755)
		re2 := regexp.MustCompile(`(?s)var widget = PlaylistWidget.findOrCreate\(` + m[1] + `\),.*?song_records = .*?is_login = false;`)
		match := re2.FindString(body)
		re3 := regexp.MustCompile(`(?s){"name":"(.*?)",.*?"rawUrl":"(.*?)",`)
		matches3 := re3.FindAllStringSubmatch(match, -1)
		for _, m3 := range matches3 {
			fmt.Printf("\t\t%s %s\n", m3[1], m3[2])
			filename := "music/" + name + "/" + m[2] + "/" + m3[1] + ".mp3"
			if !fileExist(filename) {
				getstore(strings.Replace(m3[2], "\\", "", -1), filename)
				
				tagIt(filename, m3[1], m[2], name, "hippop")
			}
		}
	}
}

func changeStr(str string) string {
	str2 := ""
	for i:=0;i<len(str);i++ {
		str2 += string(str[i])
	}
	return str2
}

func get(url string) string {
	resp, err := http.Get(url)
	if err != nil { panic(err) }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { panic(err) }
	return string(body)
}

func getstore(url string, path string) bool {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	//stat, err := f.Stat() //获取文件状态
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil { panic(err) }
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil { panic(err) }
	return true
}

func fileExist(filename string) bool {
	f, err := os.Open(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	return true
}

func tagIt(filename string, title string, album string, artist string, genre string) {
	mp3File, err := id3.Open(filename)
	if err != nil {
		panic(err)
	}
	defer mp3File.Close()
	fmt.Println("version " + mp3File.Version())
	if mp3File.Version() == "1.0" {
//		title = changeStr(title)
//		mp3File.SetTitle(title)
//		mp3File.SetAlbum(album)
//		mp3File.SetArtist(artist)
	} else {
		setFrame(mp3File, "TIT2", title)
		setFrame(mp3File, "TALB", album)
		setFrame(mp3File, "TPE1", artist)
	}
	mp3File.SetGenre("Hip-Hop")
}

func setFrame(tag *id3.File, frameName string, value string) bool {
	frame := tag.Frame(frameName)
	if frame != nil {
//		fmt.Println("changing frame")
		if textFramer, ok := frame.(v2.TextFramer); ok {
			textFramer.SetEncoding("UTF-8")
			textFramer.SetText(value)
			return true
		}
	} else {
//		fmt.Println("adding frame")
		ft := v2.V23FrameTypeMap[frameName]
		textFrame := v2.NewTextFrame(ft, "")
		textFrame.SetEncoding("UTF-8")
		textFrame.SetText(value)
		tag.AddFrames(textFrame)
	}
	return false
}

