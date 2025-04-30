package libvlcPurego

import "unsafe"

func stringSliceToPtrPtrByte(args []string) (**byte, func()) {
	mem := make([][]byte, len(args))
	ptrs := make([]*byte, len(args))

	for i, s := range args {
		mem[i] = append([]byte(s), 0)
		ptrs[i] = &mem[i][0]
	}

	return (**byte)(unsafe.Pointer(&ptrs[0])), func() {}
}

func goStringFromCString(cstr *byte) string {
	var b []byte
	for p := uintptr(unsafe.Pointer(cstr)); ; p++ {
		ch := *(*byte)(unsafe.Pointer(p))
		if ch == 0 {
			break
		}
		b = append(b, ch)
	}
	return string(b)
}
