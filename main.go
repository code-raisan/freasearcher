package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrRequest
	exitCodeErrFuzzyFinder
)

var (
	isURL bool
	word  string
)

func init() {
	flag.BoolVar(&isURL, "u", false, "ブラウザを開かずURLを表示するか指定します。")
	flag.StringVar(&word, "w", "", "検索するキーワードを指定します。")
}

func main() {
	flag.Parse()

	if word == "" {
		fmt.Println("検索キーワードを指定して下さい。\n詳しくはfreasearcher --helpを実行してください。")
		os.Exit(int(exitCodeErrArgs))
	}

	param := Param{
		Language:   "ja",
		SafeSearch: 0,
		Query:      word,
	}
	ctns, err := getResp(param)
	if err != nil {
		fmt.Fprintf(os.Stderr, "検索結果の取得に失敗しました: %s\n", err)
		os.Exit(int(exitCodeErrRequest))
	}

	idx, err := fuzzyfinder.FindMulti(
		ctns,
		func(i int) string {
			return ctns[i].title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			return runewidth.Wrap(ctns[i].content, setWidth())
		}))
	if err != nil {
		if errors.Is(fuzzyfinder.ErrAbort, err) {
			// 正常終了
			os.Exit(int(exitCodeOK))
		}
		fmt.Fprintf(os.Stderr, "fuzzyfinderにてエラーが発生しました: %s\n", err)
		os.Exit(int(exitCodeErrFuzzyFinder))
	}
	url := ctns[idx[0]].url

	if isURL {
		fmt.Println(url)
	} else {
		openBrowser(url)
	}

	os.Exit(int(exitCodeOK))
}
