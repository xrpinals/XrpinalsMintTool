package bitcoin

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const (
	GenesisNBits = uint32(0x1d00ffff)
)

func NBits2Target(nBits uint32) *big.Int {
	nBits = nBits & 0xffffffff
	nSize := nBits >> 24

	fmt.Println(nSize)

	nWord := big.NewInt(int64(nBits & 0x007fffff))
	if nSize <= 3 {
		nWord = new(big.Int).Mul(nWord, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(8*(3-nSize))), nil))
	} else {
		nWord = new(big.Int).Mul(nWord, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(8*(nSize-3))), nil))
	}
	return nWord
}

func GetTargetWork(target *big.Int) (float64, error) {
	targetMax := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(256)), nil)

	targetMaxFloat, err := strconv.ParseFloat(targetMax.Text(10), 64)
	if err != nil {
		return 0.0, err
	}

	targetFloat, err := strconv.ParseFloat(target.Text(10), 64)
	if err != nil {
		return 0.0, err
	}

	return targetMaxFloat / targetFloat, nil
}

func GetGenesisTargetWork() (float64, error) {
	targetGenesis := NBits2Target(GenesisNBits)
	return GetTargetWork(targetGenesis)
}

func GetNBitsDiff(nBits uint32) float64 {
	shift := (nBits >> 24) & 0xff
	diff := float64(0xffff) / float64(nBits&0xffffff)
	for shift < 29 {
		diff *= 256.0
		shift = shift + 1
	}
	for shift > 29 {
		diff /= 256.0
		shift = shift - 1
	}
	return diff
}

func GetTargetDiff(target *big.Int) (float64, error) {
	targetGenesis := NBits2Target(GenesisNBits)

	targetGenesisFloat, err := strconv.ParseFloat(targetGenesis.Text(10), 64)
	if err != nil {
		return 0.0, err
	}

	targetFloat, err := strconv.ParseFloat(target.Text(10), 64)
	if err != nil {
		return 0.0, err
	}

	return targetGenesisFloat / targetFloat, nil
}

func GetDiffWork(diff float64) (float64, error) {
	genesisWork, err := GetGenesisTargetWork()
	if err != nil {
		return 0.0, err
	}

	return genesisWork * diff, nil
}

func GetHashRateByWork(work float64, secs int64, unit string) float64 {
	hashRate := work / float64(secs)
	if unit == "k" || unit == "K" {
		return hashRate / math.Pow(10.0, 3.0)
	} else if unit == "m" || unit == "M" {
		return hashRate / math.Pow(10.0, 6.0)
	} else if unit == "g" || unit == "G" {
		return hashRate / math.Pow(10.0, 9.0)
	} else if unit == "t" || unit == "T" {
		return hashRate / math.Pow(10.0, 12.0)
	} else if unit == "p" || unit == "P" {
		return hashRate / math.Pow(10.0, 15.0)
	} else if unit == "e" || unit == "E" {
		return hashRate / math.Pow(10.0, 18.0)
	} else {
		return hashRate
	}
}

func GetHashRateByDiff(diff float64, secs int64, unit string) (float64, error) {
	work, err := GetDiffWork(diff)
	if err != nil {
		return 0.0, err
	}
	return GetHashRateByWork(work, secs, unit), nil
}

func GetHashRateByNBits(nBits uint32, secs int64, unit string) (float64, error) {
	work, err := GetTargetWork(NBits2Target(nBits))
	if err != nil {
		return 0.0, err
	}
	return GetHashRateByWork(work, secs, unit), nil
}
