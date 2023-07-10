package database

import (
	"context"

	"github.com/birdglove2/nitad-backend/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
)

// https://www.mongodb.com/developer/quickstart/golang-multi-document-acid-transactions/
func ExecTx(ctx context.Context, callback func(context.Context) errors.CustomError) errors.CustomError {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := client.StartSession()
	if err != nil {
		return errors.NewInternalServerError("error trying to start session - " + err.Error())
	}
	defer session.EndSession(ctx)

	transactionCallback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		return nil, callback(sessionContext)
	}

	for i := 0; i < 3; i++ {
		_, txErr := session.WithTransaction(ctx, transactionCallback, txnOpts)
		if txErr != nil {
			if cmdErr, ok := txErr.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				zap.S().Info("retry transaction")
				continue
			}
		}

		// if there is no error or the error occurred is not related to transaction
		// retry is not necessary
		break
	}

	if err != nil {
		return errors.NewInternalServerError("error executing transaction - " + err.Error() + ", transaction aborted.")
	}

	return nil
}
