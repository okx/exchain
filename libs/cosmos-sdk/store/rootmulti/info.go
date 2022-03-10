package rootmulti

import (
	"fmt"
	"sort"
	"strings"
)

type Infos []Info

func (f Infos) Len() int {
	return len(f)
}

func (f Infos) Less(i, j int) bool {
	return f[i].storeName < f[j].storeName
}

func (f Infos) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Info struct {
	storeName string
	keys      []string
	values    []string
}

type Hashes []HashInfo

func (this Hashes) String() string {
	sort.Sort(this)
	sb := strings.Builder{}
	for _, v := range this {
		sb.WriteString(fmt.Sprintf("%s:%s\n", v.storeName, v.hash))
	}
	return sb.String() + "\n"
}

func (h Hashes) Len() int {
	return len(h)
}

func (h Hashes) Less(i, j int) bool {
	return h[i].storeName < h[j].storeName
}

func (h Hashes) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type HashInfo struct {
	storeName string
	hash      string
}
