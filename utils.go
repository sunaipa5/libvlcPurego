package libvlcPurego

import "unsafe"

func stringSliceToPtrPtrByte(args []string) (**byte, func()) {
	mem := make([][]byte, len(args))
	ptrs := make([]*byte, len(args))

	for i, s := range args {
		b := append([]byte(s), 0)
		mem[i] = b
		ptrs[i] = &b[0]
	}

	return (**byte)(unsafe.Pointer(&ptrs[0])), func() {
	}
}

