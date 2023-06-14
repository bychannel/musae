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
func LoadWordCfg(fileName string) (int, error) {
	// 读取文件内容
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return 0, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Warn(err)
		}
	}()

	rows, err := f.GetRows("dirtyWord")
	if err != nil {
		return 0, err
	}

	temp := make([]string, 0)
	var counter int
	var total int
	for i, row := range rows {
		if i == 0 {
			continue
		}
		temp = append(temp, row[0])
		counter++
		total++

		// 大于5000存档一次
		if counter >= 5000 {
			GetIllegalWordMgr().SetKeywords(temp)
			// 清空
			counter = 0
			temp = nil
			logger.Debugf("加载屏蔽词 counter: %d", total)
		}
	}
	GetIllegalWordMgr().SetKeywords(temp)
	return total, nil
}
