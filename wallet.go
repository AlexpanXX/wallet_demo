package main

import (
	"github.com/elastos/Elastos.ELA.SPV/sdk"
	"github.com/elastos/Elastos.ELA.SPV/util"
	"github.com/elastos/Elastos.ELA.Utility/common"
	"github.com/elastos/Elastos.ELA.Utility/crypto"
	"github.com/elastos/Elastos.ELA/core"
	"time"
	"wallet/database"
)

const (
	defaultMaxPeers = 125
)

var (
	foundation, _ = common.
		Uint168FromAddress("8ZNizBf4KhhPjeJRGpox6rPcHE5Np6tFx3")
)

type wallet struct {
	// SPV 服务
	sdk.IService
}

// 同步区块从而接收和自己有关的交易
// 将接收到的交易存储为UTXO
// 将花费掉的UTXO删除
// 使用存储的UTXO生成交易

func newWallet(cfg *config) (*wallet, error) {
	w := wallet{}

	service, err := sdk.NewService(&sdk.Config{
		Magic:         cfg.Magic,
		SeedList:      cfg.SeedList,
		DefaultPort:   cfg.DefaultPort,
		MaxPeers:      defaultMaxPeers,
		GenesisHeader: GenesisHeader(),
		ChainStore:    database.NewDatabase(),
		NewTransaction: func() util.Transaction {
			return &core.Transaction{}
		},
		NewBlockHeader: func() util.BlockHeader {
			return util.NewElaHeader(&core.Header{})
		},
		GetFilterData: w.getFilterData,
		StateNotifier: &w,
	})
	if err != nil {
		return nil, err
	}
	w.IService = service

	return &w, nil
}

func (w *wallet) getFilterData() ([]*common.Uint168, []*util.OutPoint) {
	return nil, nil
}

// TransactionAnnounce will be invoked when received a new announced transaction.
func (w *wallet) TransactionAnnounce(tx util.Transaction) {

}

// TransactionAccepted will be invoked after a transaction sent by
// SendTransaction() method has been accepted.  Notice: this method needs at
// lest two connected peers to work.
func (w *wallet) TransactionAccepted(tx util.Transaction) {

}

// TransactionRejected will be invoked if a transaction sent by SendTransaction()
// method has been rejected.
func (w *wallet) TransactionRejected(tx util.Transaction) {

}

// TransactionConfirmed will be invoked after a transaction sent by
// SendTransaction() method has been packed into a block.
func (w *wallet) TransactionConfirmed(tx *util.Tx) {

}

// BlockCommitted will be invoked when a block and transactions within it are
// successfully committed into database.
func (w *wallet) BlockCommitted(block *util.Block) {

}

// GenesisHeader creates a specific genesis header by the given
// foundation address.
func GenesisHeader() util.BlockHeader {
	// Genesis time
	genesisTime := time.Date(2017, time.December, 22, 10, 0, 0, 0, time.UTC)

	// header
	header := core.Header{
		Version:    core.BlockVersion,
		Previous:   common.EmptyHash,
		MerkleRoot: common.EmptyHash,
		Timestamp:  uint32(genesisTime.Unix()),
		Bits:       0x1d03ffff,
		Nonce:      core.GenesisNonce,
		Height:     uint32(0),
	}

	// ELA coin
	elaCoin := &core.Transaction{
		TxType:         core.RegisterAsset,
		PayloadVersion: 0,
		Payload: &core.PayloadRegisterAsset{
			Asset: core.Asset{
				Name:      "ELA",
				Precision: 0x08,
				AssetType: 0x00,
			},
			Amount:     0 * 100000000,
			Controller: common.Uint168{},
		},
		Attributes: []*core.Attribute{},
		Inputs:     []*core.Input{},
		Outputs:    []*core.Output{},
		Programs:   []*core.Program{},
	}

	coinBase := &core.Transaction{
		TxType:         core.CoinBase,
		PayloadVersion: core.PayloadCoinBaseVersion,
		Payload:        new(core.PayloadCoinBase),
		Inputs: []*core.Input{
			{
				Previous: core.OutPoint{
					TxID:  common.EmptyHash,
					Index: 0x0000,
				},
				Sequence: 0x00000000,
			},
		},
		Attributes: []*core.Attribute{},
		LockTime:   0,
		Programs:   []*core.Program{},
	}

	coinBase.Outputs = []*core.Output{
		{
			AssetID:     elaCoin.Hash(),
			Value:       3300 * 10000 * 100000000,
			ProgramHash: *foundation,
		},
	}

	nonce := []byte{0x4d, 0x65, 0x82, 0x21, 0x07, 0xfc, 0xfd, 0x52}
	txAttr := core.NewAttribute(core.Nonce, nonce)
	coinBase.Attributes = append(coinBase.Attributes, &txAttr)

	transactions := []*core.Transaction{coinBase, elaCoin}
	hashes := make([]common.Uint256, 0, len(transactions))
	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash())
	}
	header.MerkleRoot, _ = crypto.ComputeRoot(hashes)

	return util.NewElaHeader(&header)
}
