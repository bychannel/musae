package wordfilter

var SensitiveWordMgr *Filter

func GetSensitiveWordMgr() *Filter {
	if SensitiveWordMgr == nil {
		SensitiveWordMgr = New()
	}
	return SensitiveWordMgr
}

// 加载屏蔽词
func LoadSensitiveWordCfg(fileName string) (int, error) {
	total, err := GetSensitiveWordMgr().LoadWordDict(fileName)
	if err != nil {
		return 0, err
	}

	return total, nil
}
