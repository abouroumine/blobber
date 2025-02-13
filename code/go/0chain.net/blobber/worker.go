package main

import (
	"time"

	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/allocation"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/challenge"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/config"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/handler"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/readmarker"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/writemarker"
	"github.com/0chain/blobber/code/go/0chain.net/core/common"
	"github.com/0chain/blobber/code/go/0chain.net/core/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func setupWorkers() {
	var root = common.GetRootContext()
	handler.SetupWorkers(root)
	challenge.SetupWorkers(root)
	readmarker.SetupWorkers(root)
	writemarker.SetupWorkers(root)
	allocation.StartUpdateWorker(root,
		config.Configuration.UpdateAllocationsInterval)
}

func refreshPriceOnChain() {
	var REPEAT_DELAY = 60 * 60 * time.Duration(viper.GetInt("price_worker_in_hours")) // 12 hours with default settings
	for {
		time.Sleep(REPEAT_DELAY * time.Second)
		if err := registerBlobberOnChain(); err != nil {
			logging.Logger.Error("refresh price on chain ", zap.Error(err))
		}
	}
}

func healthCheckOnChain() {
	const REPEAT_DELAY = 60 * 15 // 15 minutes

	for {
		time.Sleep(REPEAT_DELAY * time.Second)
		txnHash, err := handler.BlobberHealthCheck()
		if err != nil {
			handler.SetBlobberHealthError(err)
		} else {
			t, err := handler.TransactionVerify(txnHash)
			if err != nil {
				logging.Logger.Error("Failed to verify blobber health check", zap.Any("err", err), zap.String("txn.Hash", txnHash))
			} else {
				logging.Logger.Info("Verified blobber health check", zap.String("txn_hash", t.Hash), zap.Any("txn_output", t.TransactionOutput))
			}

			handler.SetBlobberHealthError(err)
		}
	}
}
