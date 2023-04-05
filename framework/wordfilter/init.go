package wordfilter

import (
	"github.com/xuri/excelize/v2"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
)

var IllegalWordMgr *IllegalWordsSearch

func GetIllegalWordMgr() *IllegalWordsSearch {
	if IllegalWordMgr != nil {
		return IllegalWordMgr
	}

	IllegalWordMgr = NewIllegalWordsSearch()
	return IllegalWordMgr
}

// 加载屏蔽词
func LoadWordCfg(fileName string) {
	// 读取文件内容
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Warn(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		panic(err)
	}

	temp := make([]string, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		temp = append(temp, row[3])
	}

	// 初始化屏蔽词mgr
	GetIllegalWordMgr().SetKeywords(temp)
}
