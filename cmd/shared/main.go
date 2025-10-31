package main

/*
#include <stdlib.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// 可以在这里添加 C++ 相关的声明

#ifdef __cplusplus
}
#endif
*/
import "C"
import (
	"github.com/chenwes/licensemodule/internal/license"
	"github.com/chenwes/licensemodule/pkg/utils"
)

//export VerifyLicense
func VerifyLicense(licenseFile, timestampFile, machineID, appID *C.char) *C.char {
	err := license.VerifyAndUpdate(
		C.GoString(licenseFile),
		C.GoString(timestampFile),
		C.GoString(machineID),
		C.GoString(appID),
	)
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("ok")
}

//export GenerateLicense
func GenerateLicense(machineID, appID *C.char, days C.int, outFile *C.char) *C.char {
	lic, err := license.NewLicense(
		C.GoString(machineID),
		C.GoString(appID),
		int(days),
		nil,
	)
	if err != nil {
		return C.CString(err.Error())
	}

	err = lic.Save(C.GoString(outFile))
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("ok")
}

//export GetMachineID
func GetMachineID(isContainer C.bool) *C.char {
	id, err := utils.GetMachineID()
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString(id)
}

func main() {}
