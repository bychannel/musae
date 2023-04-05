package wordfilter

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func Test_WorldFilter(t *testing.T) {
	Test_StringSearch(t)
	Test_WordsSearch(t)
	Test_StringSearchEx(t)
	Test_WordsSearchEx(t)
	Test_IllegalWordsSearch(t)
	Test_Save_Load(t)
	Test_Save_Load2(t)
	Test_time(t)
}

func Test_StringSearch(t *testing.T) {
	fmt.Println("test_StringSearch")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search := NewStringSearch()
	search.SetKeywords(list)

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0] != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1] != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}
func Test_WordsSearch(t *testing.T) {
	fmt.Println("test_WordsSearch")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search := NewWordsSearch()
	search.SetKeywords(list)

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f.Keyword != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0].Keyword != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1].Keyword != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}
func Test_StringSearchEx(t *testing.T) {
	fmt.Println("test_StringSearchEx")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search := NewStringSearchEx()
	search.SetKeywords(list)

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0] != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1] != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}
func Test_WordsSearchEx(t *testing.T) {
	fmt.Println("test_WordsSearchEx")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search := NewWordsSearchEx()
	search.SetKeywords(list)

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f.Keyword != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0].Keyword != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1].Keyword != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}
func Test_IllegalWordsSearch(t *testing.T) {
	fmt.Println("test_IllegalWordsSearch")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search := NewIllegalWordsSearch()
	search.SetKeywords(list)

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f.Keyword != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0].Keyword != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1].Keyword != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}

func Test_Save_Load(t *testing.T) {
	fmt.Println("text_Save_Load")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search2 := NewStringSearchEx()
	search2.SetKeywords(list)
	search2.Save("1.dat")

	search := NewStringSearchEx()
	search.Load("1.dat")

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0] != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1] != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}

func Test_Save_Load2(t *testing.T) {
	fmt.Println("test_Save_Load2")

	test := "我是中国人"
	list := []string{"中国", "国人", "zg人"}

	search2 := NewIllegalWordsSearch()
	search2.SetKeywords(list)
	search2.Save("2.dat")

	search := NewIllegalWordsSearch()
	search.Load("2.dat")

	b := search.ContainsAny(test)
	if b == false {
		fmt.Println("ContainsAny is Error.")
	}

	f := search.FindFirst(test)
	if f.Keyword != "中国" {
		fmt.Println("FindFirst is Error.")
	}

	all := search.FindAll(test)
	if all[0].Keyword != "中国" {
		fmt.Println("FindAll is Error.")
	}
	if all[1].Keyword != "国人" {
		fmt.Println("FindAll is Error.")
	}
	if len(all) != 2 {
		fmt.Println("FindAll is Error.")
	}
	str := search.Replace(test, '*')
	if str != "我是***" {
		fmt.Println("Replace is Error.")
	}
}

func Test_time(t *testing.T) {
	bs, _ := ioutil.ReadFile("BadWord.txt")
	s := string(bs)
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\r", "\n", -1)
	sp := strings.Split(s, "\r")
	list := make([]string, 0)
	for _, item := range sp {
		list = append(list, item)
	}
	bs2, _ := ioutil.ReadFile("Talk.txt")
	words := string(bs2)

	search := NewStringSearchEx()
	search.SetKeywords(list)

	dt := time.Now()
	for i := 0; i < 100000; i++ {
		search.FindAll(words)
	}
	dt2 := time.Now()

	fmt.Println(dt2.Sub(dt))

}
