package randomx_test

import (
	"bytes"
	"encoding/hex"
	"runtime"
	"sync"
	"testing"
	"xdago/randomx"
)

var testPairs = [][][]byte{
	// randomX
	{
		[]byte("test key 000"),
		[]byte("This is a test"),
		[]byte("46b49051978dcce1cd534a4066035184afb16a0591b43522466e10cc2496717e"),
	},
}

func TestAllocCache(t *testing.T) {
	cache, _ := randomx.AllocCache(randomx.FlagDefault)
	randomx.InitCache(cache, []byte("123"))
	randomx.ReleaseCache(cache)
}

func TestAllocDataset(t *testing.T) {
	t.Log("warning: cannot use FlagDefault only, very slow!. After using FlagJIT, really fast!")

	ds, err := randomx.AllocDataset(randomx.FlagJIT)
	if err != nil {
		panic(err)
	}
	cache, err := randomx.AllocCache(randomx.FlagJIT)
	if err != nil {
		panic(err)
	}

	seed := make([]byte, 32)
	randomx.InitCache(cache, seed)
	t.Log("rxCache initialization finished")

	count := randomx.DatasetItemCount()
	t.Log("dataset count:", count/1024/1024, "mb")
	randomx.InitDataset(ds, cache, 0, count)
	t.Log(randomx.GetDatasetMemory(ds))

	randomx.ReleaseDataset(ds)
	randomx.ReleaseCache(cache)
}

func TestCreateVM(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var tp = testPairs[0]
	cache, _ := randomx.AllocCache(randomx.FlagDefault)
	t.Log("alloc cache mem finished")
	seed := tp[0]
	randomx.InitCache(cache, seed)
	t.Log("cache initialization finished")

	ds, _ := randomx.AllocDataset(randomx.FlagDefault)
	t.Log("alloc dataset mem finished")
	count := randomx.DatasetItemCount()
	t.Log("dataset count:", count)
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	t.Log("Here though using FlagDefault, but we use multi-core to accelerate it")
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	t.Log("dataset initialization finished") // too slow when one thread
	vm, _ := randomx.CreateVM(cache, ds, randomx.FlagJIT, randomx.FlagHardAES, randomx.FlagFullMEM)

	var hashCorrect = make([]byte, hex.DecodedLen(len(tp[2])))
	_, err := hex.Decode(hashCorrect, tp[2])
	if err != nil {
		t.Log(err)
	}

	hash := randomx.CalculateHash(vm, tp[1])
	if !bytes.Equal(hash, hashCorrect) {
		t.Logf("answer is incorrect: %x, %x", hash, hashCorrect)
		t.Fail()
	}
}

// go test -v -run=^$ -benchtime=1m  -timeout 20m -bench=.
func BenchmarkCalculateHashDefault(b *testing.B) {
	cache, _ := randomx.AllocCache(randomx.FlagDefault)
	ds, _ := randomx.AllocDataset(randomx.FlagDefault)
	randomx.InitCache(cache, []byte("123"))
	count := randomx.DatasetItemCount()
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	vm, _ := randomx.CreateVM(cache, ds, randomx.FlagDefault)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomx.CalculateHash(vm, []byte("123"))
	}

	randomx.DestroyVM(vm)
}

func BenchmarkCalculateHashJIT(b *testing.B) {
	cache, _ := randomx.AllocCache(randomx.FlagDefault, randomx.FlagJIT)
	ds, _ := randomx.AllocDataset(randomx.FlagDefault, randomx.FlagJIT)
	randomx.InitCache(cache, []byte("123"))
	count := randomx.DatasetItemCount()
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	vm, _ := randomx.CreateVM(cache, ds, randomx.FlagDefault)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomx.CalculateHash(vm, []byte("123"))
	}

	randomx.DestroyVM(vm)
}

func BenchmarkCalculateHashFullMEM(b *testing.B) {
	cache, _ := randomx.AllocCache(randomx.FlagDefault, randomx.FlagFullMEM)
	ds, _ := randomx.AllocDataset(randomx.FlagDefault, randomx.FlagFullMEM)
	randomx.InitCache(cache, []byte("123"))
	count := randomx.DatasetItemCount()
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	vm, _ := randomx.CreateVM(cache, ds, randomx.FlagDefault)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomx.CalculateHash(vm, []byte("123"))
	}

	randomx.DestroyVM(vm)
}

func BenchmarkCalculateHashHardAES(b *testing.B) {
	cache, _ := randomx.AllocCache(randomx.FlagDefault, randomx.FlagHardAES)
	ds, _ := randomx.AllocDataset(randomx.FlagDefault, randomx.FlagHardAES)
	randomx.InitCache(cache, []byte("123"))
	count := randomx.DatasetItemCount()
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	vm, _ := randomx.CreateVM(cache, ds, randomx.FlagDefault)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomx.CalculateHash(vm, []byte("123"))
	}

	randomx.DestroyVM(vm)
}

func BenchmarkCalculateHashAll(b *testing.B) {
	cache, _ := randomx.AllocCache(randomx.FlagDefault, randomx.FlagArgon2, randomx.FlagArgon2AVX2, randomx.FlagArgon2SSSE3, randomx.FlagFullMEM, randomx.FlagHardAES, randomx.FlagJIT) // without lagePage to avoid panic
	ds, _ := randomx.AllocDataset(randomx.FlagDefault, randomx.FlagArgon2, randomx.FlagArgon2AVX2, randomx.FlagArgon2SSSE3, randomx.FlagFullMEM, randomx.FlagHardAES, randomx.FlagJIT)
	randomx.InitCache(cache, []byte("123"))
	count := randomx.DatasetItemCount()
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	vm, _ := randomx.CreateVM(cache, ds, randomx.FlagDefault)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomx.CalculateHash(vm, []byte("123"))
	}

	randomx.DestroyVM(vm)
}
