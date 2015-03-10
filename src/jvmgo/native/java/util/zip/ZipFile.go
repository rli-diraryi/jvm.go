package zip

import (
	. "jvmgo/any"
	"jvmgo/jvm/rtda"
	rtc "jvmgo/jvm/rtda/class"
	"jvmgo/util"
)

const (
	JZENTRY_NAME    = 0
	JZENTRY_EXTRA   = 1
	JZENTRY_COMMENT = 2
)

func init() {
	_zf(initIDs, "initIDs", "()V")
	_zf(getEntryBytes, "getEntryBytes", "(JI)[B")
	_zf(getEntryFlag, "getEntryFlag", "(J)I")
	_zf(getEntryTime, "getEntryTime", "(J)J")
	_zf(getNextEntry, "getNextEntry", "(JI)J")
	_zf(getTotal, "getTotal", "(J)I")
	_zf(open, "open", "(Ljava/lang/String;IJZ)J")
	_zf(startsWithLOC, "startsWithLOC", "(J)Z")
}

func _zf(method Any, name, desc string) {
	rtc.RegisterNativeMethod("java/util/zip/ZipFile", name, desc, method)
}

// private static native void initIDs();
// ()V
func initIDs(frame *rtda.Frame) {
	// todo
}

// private static native long open(String name, int mode, long lastModified,
//                                 boolean usemmap) throws IOException;
// (Ljava/lang/String;IJZ)J
func open(frame *rtda.Frame) {
	vars := frame.LocalVars()
	nameObj := vars.GetRef(0)

	name := rtda.GoString(nameObj)
	jzfile, err := openZip(name)
	if err != nil {
		// todo
		panic("IOException")
	}

	stack := frame.OperandStack()
	stack.PushLong(jzfile)
}

// private static native boolean startsWithLOC(long jzfile);
// (J)Z
func startsWithLOC(frame *rtda.Frame) {
	// todo
	stack := frame.OperandStack()
	stack.PushBoolean(true)
}

// private static native int getTotal(long jzfile);
// (J)I
func getTotal(frame *rtda.Frame) {
	vars := frame.LocalVars()
	jzfile := vars.GetLong(0)

	total := getEntryCount(jzfile)

	stack := frame.OperandStack()
	stack.PushInt(total)
}

// private static native long getNextEntry(long jzfile, int i);
// (JI)J
func getNextEntry(frame *rtda.Frame) {
	vars := frame.LocalVars()
	jzfile := vars.GetLong(0)
	i := vars.GetInt(2)

	jzentry := getJzentry(jzfile, i)

	stack := frame.OperandStack()
	stack.PushLong(jzentry)
}

// private static native int getEntryFlag(long jzentry);
// (J)I
func getEntryFlag(frame *rtda.Frame) {
	vars := frame.LocalVars()
	jzentry := vars.GetLong(0)

	entry := getEntry(jzentry)
	flag := int32(entry.Flags)

	stack := frame.OperandStack()
	stack.PushInt(flag)
}

// private static native byte[] getEntryBytes(long jzentry, int type);
// (JI)[B
func getEntryBytes(frame *rtda.Frame) {
	vars := frame.LocalVars()
	jzentry := vars.GetLong(0)
	_type := vars.GetInt(2)

	goBytes := _getEntryBytes(jzentry, _type)
	jBytes := util.CastUint8sToInt8s(goBytes)
	byteArr := rtc.NewByteArray(jBytes, frame.ClassLoader())

	stack := frame.OperandStack()
	stack.PushRef(byteArr)
}

func _getEntryBytes(jzentry int64, _type int32) []byte {
	entry := getEntry(jzentry)
	switch _type {
	case JZENTRY_NAME:
		return []byte(entry.Name)
	case JZENTRY_EXTRA:
		return entry.Extra
	case JZENTRY_COMMENT:
		return []byte(entry.Comment)
	}
	util.Panicf("BAD type: %v", _type)
	return nil
}

// private static native long getEntryTime(long jzentry);
// (J)J
func getEntryTime(frame *rtda.Frame) {
	vars := frame.LocalVars()
	jzentry := vars.GetLong(0)

	entry := getEntry(jzentry)
	modDate := entry.ModifiedDate
	modTime := entry.ModifiedTime
	time := int64(modDate)<<16 | int64(modTime)

	stack := frame.OperandStack()
	stack.PushLong(time)
}