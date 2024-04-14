/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package shardedsql

import (
	"context"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BatchLookup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	migrations := NewStatementSequence("batchlookup")
	migrations.Insert(1, "CREATE TABLE batchlookup (k INTEGER, PRIMARY KEY (k))")
	err := testingDB.MigrateSchema(ctx, migrations)
	assert.NoError(t, err)

	// Insert some records
	n := 1011
	var keys []int
	shard1 := testingDB.Shard(1)
	for i := 0; i < n; i++ {
		_, err := shard1.ExecContext(ctx, "INSERT INTO batchlookup (k) VALUES (?)", i)
		assert.NoError(t, err)
		if i != 5 { // Exclude record 5
			keys = append(keys, i)
		}
	}

	// Batch lookup should return all records (excluding 5) in order they are requested
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	var readBack []int
	batches := 0
	blu := NewBatchLookup(shard1, "SELECT k FROM batchlookup WHERE k IN (?)", keys)
	for blu.Next() {
		batches++
		rows, err := blu.QueryContext(ctx)
		assert.NoError(t, err)
		for rows.Next() {
			var k int
			err := rows.Scan(&k)
			assert.NoError(t, err)
			readBack = append(readBack, k)
		}
	}
	assert.Equal(t, len(keys), len(readBack))
	sort.Ints(readBack)
	for i, k := range readBack {
		if i < 5 {
			assert.Equal(t, i, k)
		} else {
			assert.Equal(t, i+1, k)
		}
	}
	assert.Equal(t, (n+blu.batchSize-1)/blu.batchSize, batches)
}
