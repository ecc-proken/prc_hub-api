package mysql

import "encoding/json"

// null, null => jsonにキーが存在しない
// ポインター, null => jsonにキーが存在し、値がnull
// ポインター, ポインター => jsonにキーが存在し、値がnull以外
type PatchNullJSONString struct {
	String **string `validate:"omitempty"`
}

func (p *PatchNullJSONString) UnmarshalJSON(data []byte) error {
	// jsonにキーが存在する場合にこの関数が呼び出される
	var valueP *string = nil
	if string(data) == "null" {
		// jsonにキーが存在し、値がnull
		p.String = &valueP
		return nil
	}

	var tmp string
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// typeエラー
		return err
	}
	// valid value
	p.String = &tmpP
	return nil
}
