package watcher

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"time"
)

type sStruct struct {
	npm   map[uintptr]bool
	exNum int
}

func SizeOf(data interface{}) int {
	var npm = &sStruct{make(map[uintptr]bool), 0}
	num := npm.sizeof(reflect.ValueOf(data))
	return num //+ npm.exNum
}

func SizeTOf(data interface{}) int {
	var npm = &sStruct{make(map[uintptr]bool), 0}
	num := npm.sizeof(reflect.ValueOf(data))
	return num + npm.exNum
}

func (s *sStruct) sizeof(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Map:
		sum := 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			num := s.sizeof(mapkey)
			if num < 0 {
				return -1
			}
			sum += num
			num = s.sizeof(v.MapIndex(mapkey))
			if num < 0 {
				return -1
			}
			sum += num
		}
		s.exNum += int(v.Type().Size())
		return sum
	case reflect.Slice:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			num := s.sizeof(v.Index(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		s.exNum += int(v.Type().Size())
		return sum

	case reflect.Array:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			num := s.sizeof(v.Index(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		return sum

	case reflect.String:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			num := s.sizeof(v.Index(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		s.exNum += int(v.Type().Size())
		return sum

	case reflect.Ptr:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		//fmt.Println(v.Pointer())
		if _, ok := s.npm[v.Pointer()]; ok {
			return 0
		} else {
			s.npm[v.Pointer()] = true
		}
		return s.sizeof(v.Elem())

	case reflect.Interface:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		return s.sizeof(v.Elem())

	case reflect.Uintptr: //Don't think it's Pointer 不认为是指针
		return int(v.Type().Size())

	case reflect.UnsafePointer: //Don't think it's Pointer 不认为是指针
		return int(v.Type().Size())

	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			if v.Type().Field(i).Tag.Get("ss") == "-" {
				continue
			}
			num := s.sizeof(v.Field(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		return sum

	case reflect.Func, reflect.Chan:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		return 0 //Temporary non handling func,chan.
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int, reflect.Uint:
		return int(v.Type().Size())
	case reflect.Bool:
		return int(v.Type().Size())
	default:
		fmt.Println("t.Kind() no found:", v.Kind())
	}

	return -1
}

func getSize(v interface{}) int {
	return SizeTOf(v)
}

func getSize1(v interface{}) int {
	size := int(reflect.TypeOf(v).Size())
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			size += getSize(s.Index(i).Interface())
		}
	case reflect.Map:
		s := reflect.ValueOf(v)
		keys := s.MapKeys()
		size += int(float64(len(keys)) * 10.79) // approximation from https://golang.org/src/runtime/hashmap.go
		for i := range keys {
			size += getSize(keys[i].Interface()) + getSize(s.MapIndex(keys[i]).Interface())
		}
	case reflect.String:
		size += reflect.ValueOf(v).Len()
	case reflect.Struct:
		s := reflect.ValueOf(v)
		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).CanInterface() {
				size += getSize(s.Field(i).Interface())
			}
		}
	}
	return size
}

func StaticMemory() {
	for {
		printmem()
		time.Sleep(time.Second * 10)
	}
}

func printmem() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	asize := mem.Alloc
	mb := float64(asize) / (1024 * 1024)
	if mb > float64(watcherMut*1000) {
		dumpMemPprof(int(mb))
	}
	sys := float64(mem.HeapSys) / (1024 * 1024)
	fmt.Printf("******lyh****** Alloc %.2f MB, Sys %.2f \n", mb, sys)
}

func dumpMemPprof(size int) error {
	fileName := fmt.Sprintf("watchdb_pprof_%s.%d.mem.bin", time.Now().Format("20060102150405"), size)
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create mem pprof file %s error: %w", fileName, err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err = pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("could not write memory profile: %w", err)
	}
	return nil
}
