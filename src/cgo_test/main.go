package main

import (
	"fmt"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdlib.h>
void hello()
{
	printf("Cgo hello\n");
}

typedef struct
{
	int nums;
	char name[64];
}Student;

void Set(Student **pstudent)
{
	*pstudent = (Student *)malloc(sizeof(Student));
	(*pstudent)->nums = 1;
	snprintf((*pstudent)->name, sizeof((*pstudent)->name), "hello");
}
*/
import "C"

func MergeSort(array []int){
	
}


func main() {
	var student *C.Student
	C.Set(&student)

	name := C.GoString(&(student.name[0]))
	fmt.Printf("name=%s,nums=%d\n", name, int(student.nums))
	C.free(unsafe.Pointer(student))
}
