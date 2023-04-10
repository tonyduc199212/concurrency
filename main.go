package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
)

const path = "/Users/nguyenpc/src/input/input"

// strArray1 := [3]string{"Japan", "Australia", "Germany"}
func CreateFiles(num int) {
	content1 := [...]string{"aaa", "bbb"}

	content := strings.Join(content1[:], " ")
	// num := 2
	for i := 0; i < num; i++ {
		CreateFile(strconv.Itoa(i)+".txt", content)
	}
	// log.Println(content)

}
func CreateFile(fileName string, content string) {
	data := []byte(content)
	_path := filepath.Join(path, fileName)
	err := ioutil.WriteFile(_path, data, 0777)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfull:::", fileName)
}
func ReadFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
func ReadFolder(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
func CountAWord(word string, content string) int {
	arr := strings.Split(content, " ")
	count := 0
	for _, _word := range arr {
		if _word == word {
			count++
		}
	}
	return count

}
func main() {
	// CreateFiles(200000)
	files := ReadFolder(path)
	log.Println("len files ", len(files))
	tick0 := time.Now()
	Cach0(files)
	tick1 := time.Now()
	log.Println("time cach 0::: ", tick1.Sub(tick0))
	Cach1(files)

	tick2 := time.Now()
	log.Println("time cach 1::: ", tick2.Sub(tick1))
	numPool := 100
	Cach2(files, numPool)
	tick3 := time.Now()
	log.Println("time cach 2::: ", tick3.Sub(tick2))

}
func Cach0(files []fs.FileInfo) {
	countAll := 0

	for _, file := range files {
		_filePath := filepath.Join(path, file.Name())
		content := ReadFile(_filePath)
		count := CountAWord("aaa", content)
		countAll += count
	}
	// for range files {
	// 	count += <-result
	// 	// fmt.Println("count::: ", count)
	// }
	fmt.Println("countAll final::: ", countAll)
}
func Cach1(files []fs.FileInfo) {
	result := make(chan int, len(files))
	for _, file := range files {
		go func(file fs.FileInfo) {
			_filePath := filepath.Join(path, file.Name())
			content := ReadFile(_filePath)
			count := CountAWord("aaa", content)
			result <- count
		}(file)
	}
	count := 0
	for range files {
		count += <-result
		// fmt.Println("count::: ", count)
	}
	fmt.Println("count final::: ", count)
}
func Cach2(files []fs.FileInfo, num int) {
	result := make(chan int, len(files))
	wp := workerpool.New(num)
	for _, file := range files {
		file := file
		wp.Submit(func() {
			_filePath := filepath.Join(path, file.Name())
			content := ReadFile(_filePath)
			count := CountAWord("aaa", content)
			result <- count
		})
	}
	wp.StopWait()
	count := 0
	for range files {
		count += <-result
		// fmt.Println("count::: ", count)
	}
	fmt.Println("count final::: ", count)

}
