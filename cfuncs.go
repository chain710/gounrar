package gounrar

/*
#include "dll.hpp"
int callbackInC(UINT msg,uintptr_t UserData,uintptr_t P1,uintptr_t P2) {
	int callbackInGo(UINT, uintptr_t, uintptr_t, uintptr_t);
	return callbackInGo(msg, UserData, P1, P2);
}
*/
import "C"
