package tool

import (
	"math/rand"
	"strconv"
	"time"
)

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init() {
	// TODO 拆分微服务后，每个服务要设置自己的Seed
	//  crypto/rand.Read，调用系统getrandom(2)接口，实际上也是调用/dev/urandom
	//  而容器共用物理机内核，也就是说容器不存在自己的/dev/urandom，只有linux的/dev/random，因此随机性有一定保证，可以用作seed
	rand.Seed(time.Now().UnixNano())
}

func RandStr(n uint32) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTERS[rand.Int63()%int64(len(LETTERS))]
	}
	return string(b)
}

func IsNumber(str string) bool {
	if len(str) == 0 {
		return false
	}
	_, err := strconv.ParseUint(str, 10, 64)
	return err == nil
}

func FindMax(arr []int) int {
	maxVal := arr[0]

	for i := 1; i < len(arr); i++ {
		//从第二个 元素开始循环比较，如果发现有更大的，则交换
		if maxVal < arr[i] {
			maxVal = arr[i]
		}
	}
	return maxVal
}
