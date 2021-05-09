package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"
	"github.com/stretchr/testify/mock"
)

// StorageHandlerI is an autogenerated mock type for the StorageHandlerI type
type storageHandlerI struct {
	mock.Mock
}

// verifyAllocation provides a mock function with given fields: ctx, tx, readonly
func (_m *storageHandlerI) verifyAllocation(ctx context.Context, tx string, readonly bool) (*allocation.Allocation, error) {
	ret := _m.Called(ctx, tx, readonly)

	var r0 *allocation.Allocation
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) *allocation.Allocation); ok {
		r0 = rf(ctx, tx, readonly)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*allocation.Allocation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, tx, readonly)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// verifyAuthTicket provides a mock function with given fields: ctx, authTokenString, allocationObj, refRequested, clientID
func (_m *storageHandlerI) verifyAuthTicket(ctx context.Context, authTokenString string, allocationObj *allocation.Allocation, refRequested *reference.Ref, clientID string) (bool, error) {
	ret := _m.Called(ctx, authTokenString, allocationObj, refRequested, clientID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, *allocation.Allocation, *reference.Ref, string) bool); ok {
		r0 = rf(ctx, authTokenString, allocationObj, refRequested, clientID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *allocation.Allocation, *reference.Ref, string) error); ok {
		r1 = rf(ctx, authTokenString, allocationObj, refRequested, clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type TestDataController struct {
	db *gorm.DB
}

func NewTestDataController(db *gorm.DB) *TestDataController {
	return &TestDataController{db: db}
}

// ClearDatabase deletes all data from all tables
func (c *TestDataController) ClearDatabase() error {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil {
			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					log.Println(errRollback)
				}
			}
		}
	}()

	db, err := c.db.DB()
	if err != nil {
		return err
	}

	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec("truncate allocations cascade")
	if err != nil {
		return err
	}

	_, err = tx.Exec("truncate reference_objects cascade")
	if err != nil {
		return err
	}

	_, err = tx.Exec("truncate commit_meta_txns cascade")
	if err != nil {
		return err
	}

	_, err = tx.Exec("truncate collaborators cascade")
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *TestDataController) AddGetAllocationTestData() error {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil {
			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					log.Println(errRollback)
				}
			}
		}
	}()

	db, err := c.db.DB()
	if err != nil {
		return err
	}

	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	expTime := time.Now().Add(time.Hour * 100000).UnixNano()

	_, err = tx.Exec(`
INSERT INTO allocations (id, tx, owner_id, owner_public_key, expiration_date, payer_id)
VALUES ('exampleId' ,'exampleTransaction','exampleOwnerId','exampleOwnerPublicKey',` + fmt.Sprint(expTime) + `,'examplePayerId');
`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *TestDataController) AddGetFileMetaDataTestData() error {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil {
			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					log.Println(errRollback)
				}
			}
		}
	}()

	db, err := c.db.DB()
	if err != nil {
		return err
	}

	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	expTime := time.Now().Add(time.Hour * 100000).UnixNano()

	_, err = tx.Exec(`
INSERT INTO allocations (id, tx, owner_id, owner_public_key, expiration_date, payer_id)
VALUES ('exampleId' ,'exampleTransaction','exampleOwnerId','exampleOwnerPublicKey',` + fmt.Sprint(expTime) + `,'examplePayerId');
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO reference_objects (id, allocation_id, path_hash,lookup_hash,type,name,path,hash,custom_meta,content_hash,merkle_root,actual_file_hash,mimetype,write_marker,thumbnail_hash, actual_thumbnail_hash)
VALUES (1234,'exampleId','exampleId:examplePath','exampleId:examplePath','f','filename','examplePath','someHash','customMeta','contentHash','merkleRoot','actualFileHash','mimetype','writeMarker','thumbnailHash','actualThumbnailHash');
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO commit_meta_txns (ref_id,txn_id)
VALUES (1234,'someTxn');
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO collaborators (ref_id, client_id)
VALUES (1234, 'someClient');
`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *TestDataController) AddGetFileStatsTestData() error {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil {
			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					log.Println(errRollback)
				}
			}
		}
	}()

	db, err := c.db.DB()
	if err != nil {
		return err
	}

	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	expTime := time.Now().Add(time.Hour * 100000).UnixNano()

	_, err = tx.Exec(`
INSERT INTO allocations (id, tx, owner_id, owner_public_key, expiration_date, payer_id)
VALUES ('exampleId' ,'exampleTransaction','exampleOwnerId','exampleOwnerPublicKey',` + fmt.Sprint(expTime) + `,'examplePayerId');
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO reference_objects (id, allocation_id, path_hash,lookup_hash,type,name,path,hash,custom_meta,content_hash,merkle_root,actual_file_hash,mimetype,write_marker,thumbnail_hash, actual_thumbnail_hash)
VALUES (1234,'exampleId','exampleId:examplePath','exampleId:examplePath','f','filename','examplePath','someHash','customMeta','contentHash','merkleRoot','actualFileHash','mimetype','writeMarker','thumbnailHash','actualThumbnailHash');
`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *TestDataController) AddListEntitiesTestData() error {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil {
			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					log.Println(errRollback)
				}
			}
		}
	}()

	db, err := c.db.DB()
	if err != nil {
		return err
	}

	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	expTime := time.Now().Add(time.Hour * 100000).UnixNano()

	_, err = tx.Exec(`
INSERT INTO allocations (id, tx, owner_id, owner_public_key, expiration_date, payer_id)
VALUES ('exampleId' ,'exampleTransaction','exampleOwnerId','exampleOwnerPublicKey',` + fmt.Sprint(expTime) + `,'examplePayerId');
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
INSERT INTO reference_objects (id, allocation_id, path_hash,lookup_hash,type,name,path,hash,custom_meta,content_hash,merkle_root,actual_file_hash,mimetype,write_marker,thumbnail_hash, actual_thumbnail_hash)
VALUES (1234,'exampleId','exampleId:examplePath','exampleId:examplePath','f','filename','examplePath','someHash','customMeta','contentHash','merkleRoot','actualFileHash','mimetype','writeMarker','thumbnailHash','actualThumbnailHash');
`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *TestDataController) AddGetObjectPathTestData() error {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil {
			if tx != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					log.Println(errRollback)
				}
			}
		}
	}()

	db, err := c.db.DB()
	if err != nil {
		return err
	}

	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	expTime := time.Now().Add(time.Hour * 100000).UnixNano()

	_, err = tx.Exec(`
INSERT INTO allocations (id, tx, owner_id, owner_public_key, expiration_date, payer_id)
VALUES ('exampleId' ,'exampleTransaction','exampleOwnerId','exampleOwnerPublicKey',` + fmt.Sprint(expTime) + `,'examplePayerId');
`)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
