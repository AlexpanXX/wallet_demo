package database

import (
	"testing"

	"github.com/elastos/Elastos.ELA.SPV/util"
	"github.com/elastos/Elastos.ELA.Utility/common"
	"github.com/elastos/Elastos.ELA/core"
)

func TestNew(t *testing.T) {
	db := &headers{
		store: make(map[common.Uint256]*util.Header),
	}
	headers := make([]*util.Header, 100)
	previous := common.Uint256{}
	for i := range headers {
		headers[i] = new(util.Header)
		if i > 0 {
			previous = headers[i-1].Hash()
		}
		headers[i].BlockHeader = util.NewElaHeader(&core.Header{
			Previous: previous,
		})
		headers[i].Height = uint32(i)
		db.Put(headers[i], true)
	}

	for i := range headers {
		hash := headers[i].Hash()
		header, _ := db.Get(&hash)
		if header == nil {
			t.Errorf("can not find header %s", hash)
		}
	}

	for i := 99; i >= 0; i-- {
		previous, err := db.GetPrevious(headers[i])
		if err != nil {
			t.Error(err)
		}
		if previous == nil {
			t.Errorf("can not find %d previous header",
				headers[i].Height)
		}
	}
}
