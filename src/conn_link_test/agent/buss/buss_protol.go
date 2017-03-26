package main

import (
	"reflect"
	"unsafe"
)

/*
struct tagBUSS_HEADER
{
	K_INT32  PacketLen;
    K_INT16  Flag;
    //消息名称的长度,建议采用namespace.name 的方式
    K_INT16  NameLen;
    //消息名称,长度是NameLen个字节
    K_CHAR  MessageName[NameLen];
    //实际的消息数据,放置业务数据,长度是
    //PacketLen - NameLen - sizeof(Flag) - sizeof(NameLen) - sizeof(CheckSum)
    //K_CHAR  BufData[PacketLen-NameLen-sizeof(Flag)-sizeof(NameLen)- sizeof(CheckSum)];
    //K_INT32  CheckSum;
};
*/

func Int32ToBytes(value int32) []byte {
	var slice reflect.SliceHeader
	var data_buf []byte

	slice.Data = (uintptr)(unsafe.Pointer(&value))
	slice.Cap = int(unsafe.Sizeof(int32(0)))
	slice.Len = int(unsafe.Sizeof(int32(0)))
	data_buf = *(*[]byte)(unsafe.Pointer(&slice))

	return data_buf
}

func Int16ToBytes(value int16) []byte {
	var slice reflect.SliceHeader
	var data_buf []byte

	slice.Data = (uintptr)(unsafe.Pointer(&value))
	slice.Cap = int(unsafe.Sizeof(int16(0)))
	slice.Len = int(unsafe.Sizeof(int16(0)))
	data_buf = *(*[]byte)(unsafe.Pointer(&slice))

	return data_buf
}

func Pack(body []byte, message_name string, enctrypt_method int) []byte {
	buf := make([]byte, 0, 256)

	/********************************************************************/
	/**********************PacketLen*************************************/
	packet_len := 0
	packet_len += int(unsafe.Sizeof(int16(0))) // flag
	packet_len += int(unsafe.Sizeof(int16(0))) // name_len
	packet_len += len(message_name)
	packet_len += len(body)
	packet_len += int(unsafe.Sizeof(int32(0))) // check_sum

	data_buf := Int32ToBytes(int32(packet_len))
	buf = append(buf, data_buf...)

	/*********************************************************************/
	/****************************Flag*************************************/
	data_buf = Int16ToBytes(int16(enctrypt_method))
	buf = append(buf, data_buf...)

	/**********************************************************************/
	/******************************NameLen*********************************/
	data_buf = Int16ToBytes(int16(len(message_name)))
	buf = append(buf, data_buf...)

	/**********************************************************************/
	/******************************message_name****************************/
	buf = append(buf, ([]byte(message_name))...)

	/**********************************************************************/
	/******************************data***********************************/
	if body != nil {
		buf = append(buf, body...)
	}

	// 计算checksum
	check_sum := CheckSum(buf)
	data_buf = Int32ToBytes(check_sum)
	buf = append(buf, data_buf...)

	return buf
}

func Unpack(buf []byte) ([]byte, string, bool) {
	var enctrypt_method int

	offset := 0
	buf_len := len(buf)

	packet_len := int(*((*int32)(unsafe.Pointer(&buf[0]))))
	packet_data := buf[4:]
	if packet_len < offset+int(unsafe.Sizeof(int16(0))) {
		return nil, "", false
	}

	// 求出flag
	enctrypt_method = int(*(*int16)(unsafe.Pointer(&packet_data[offset])))
	// 默认为0，暂时不加密解密
	if enctrypt_method == 0 {

	}

	offset += int(unsafe.Sizeof(int16(0)))
	if packet_len < offset+int(unsafe.Sizeof(int16(0))) {
		return nil, "", false
	}

	name_len := int(*(*int16)(unsafe.Pointer(&packet_data[offset])))
	offset += int(unsafe.Sizeof(int16(0)))

	if packet_len < offset+name_len {
		return nil, "", false
	}
	message_name := string(packet_data[offset : offset+name_len])
	offset += name_len

	// packet_body
	if packet_len < offset+int(unsafe.Sizeof(int32(0))) {
		return nil, "", false // 没有checksum字段
	}

	var packet_body []byte = nil

	body_len := packet_len - offset - int(unsafe.Sizeof(int32(0)))
	if body_len > 0 {
		packet_body = make([]byte, body_len)
		copy(packet_body, packet_data[offset:offset+body_len])
	}
	offset += body_len

	// 求出checksum
	check_sum := int(*(*int32)(unsafe.Pointer(&packet_data[offset])))
	calc_check_sum := int(CheckSum(buf[0 : buf_len-int(unsafe.Sizeof(int32(0)))]))
	if calc_check_sum != check_sum {
		return nil, "", false
	}

	return packet_body, message_name, true
}

func CheckSum(buf []byte) int32 {
	var check_sum int32

	check_sum = 0
	for _, v := range buf {
		check_sum += int32(v)
	}

	return check_sum
}
