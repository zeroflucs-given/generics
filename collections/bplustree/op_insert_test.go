package bplustree

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics/collections"
)

const (
	DefaultTestPreAlloc = 100
)

// TestInsertsBasic tests we can do a few inserts
func TestInsertsBasic(t *testing.T) {
	inserts := []struct {
		K int
		V string
	}{
		{16, "maybe"},
		{5, "should"},
		{15, "but"},
		{32, "about"},
		{25, "you"},
		{3, "sentence"},
		{8, "some"},
		{12, "eventually"},
		{7, "make"},
		{42, "it"},
		{1, "This"},
		{22, "unless"},
		{20, "not"},
		{11, "sense"},
		{21, "immediately"},
		{43, "."},
		{44, "Sometimes"},
		{45, "keys"},
		{46, "come"},
		{47, "in"},
		{48, "order"},
		{30, "think"},
	}
	for i := 2; i < len(inserts); i++ {
		t.Run(fmt.Sprintf("WithOrder_%d", i), func(t *testing.T) {
			tree, err := New[int, string](i, DefaultTestPreAlloc)
			require.NoError(t, err, "Should be able to initialize")
			require.NotNil(t, tree, "Should have a tree result")

			for i, op := range inserts {
				tree.Insert(op.K, op.V)

				prev := -1
				count := 0
				for kvp := range tree.Scan() {
					require.LessOrEqual(t, prev, kvp.Key, "Keys should be ascending")
					prev = kvp.Key
					// t.Logf("  K=%v, V=%v", kvp.Key, kvp.Value)
					count++
					require.LessOrEqual(t, count-1, i, "Count shouldn't overrun during scan")
				}
				require.Equal(t, i+1, count, "Should have right number of values in tree")

			}
		})
	}
}

// TestBulkInsertWithVerify tests what happens as we add more data. After each insert
// we verify the contents of the data, which can take a considerable amount of time.
func TestBulkInsertWithVerify(t *testing.T) {
	// Arrange
	order := 5
	dataSize := 10000
	maxKey := int32(1000000)
	rs := rand.NewSource(int64(order * dataSize))
	rnd := rand.New(rs)

	tree, err := New[int, int](order, DefaultTestPreAlloc)
	require.NoError(t, err, "Should not error while starting")
	require.NotNil(t, tree, "Should have a tree")

	for i := 0; i < dataSize; i++ {
		if i%1000 == 0 {
			t.Logf("Progress report [Index: %d]", i)
		}

		// oldState := bytes.NewBuffer(nil)
		// tree.(collections.Diagnosable).Dump(oldState)

		randKey := int(rnd.Int31n(maxKey))
		tree.Insert(randKey, i)

		count := tree.Count()
		if count != i+1 {
			newState := bytes.NewBuffer(nil)
			tree.(collections.Diagnosable).Dump(newState)

			// os.WriteFile("before.tmp", oldState.Bytes(), 0644)
			err := os.WriteFile("after.tmp", newState.Bytes(), 0644)
			if err != nil {
				t.Logf("Failed to write after-dump: %v", err)
			}
			t.Logf("Got %v records but expected %v. Failed inserting %v", count, i+1, randKey)
			t.FailNow()
		}
	}
}

func BenchmarkSimpleTreeInsert(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(10000)
	b.StartTimer()
	i := 0

	for i < b.N {
		tr, _ := New[int, int](27, DefaultTestPreAlloc)
		for keyIndex, key := range insertP {
			tr.Insert(keyIndex, key)
			i++

			if i >= b.N {
				return
			}
		}
	}
}

func TestTreeInsertStress(t *testing.T) {
	order := 27
	reportInterval := 1000 * 1000 // 1M

	rs := rand.NewSource(0)
	rnd := rand.New(rs)
	data := rnd.Perm(5000000)

	tr, _ := New[int, int](order, DefaultTestPreAlloc)

	i := 0
	intervalStart := time.Now()

	for {
		if i > 0 && i%reportInterval == 0 {
			elapsed := time.Since(intervalStart)
			rate := float64(reportInterval) / elapsed.Seconds()
			t.Logf("Inserted %d items in %.2f seconds (Rate=%.2f M/sec, Total=%dM)", reportInterval, elapsed.Seconds(), rate/1000000.0, i/1000000)
			intervalStart = time.Now()

			if elapsed.Seconds() > 1 {
				t.Log("Test completing as this cycle took longer than the threshold. This is not an error.")
				break
			}
		}

		tr.Insert(data[i%len(data)], i)

		i++
	}
}
