package classfile

import "fmt"

type ClassFile struct {
	// magic        uint32
	minorVersion uint16
	majorVersion uint16
	constantPool ConstantPool
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	fields       []*MemberInfo
	methods      []*MemberInfo
	attributes   []AttributeInfo
}

func Parse(classData []byte) (cf *ClassFile, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	cr := &ClassReader{ClassData}
	cf = &ClassFile{}
	cf.read(cr)
	return
}

func (self *ClassFile) read(reader *ClassReader) {
	self.readAndCheckMagic(reader)
	self.readAndCheckVersion(reader)
	self.ConstantPool = readConstantPool(reader)
	self.accessFlags = reader.readUint16()
	self.thisClass = reader.readUint16()
	self.superClass = reader.readUint16()
	self.interfaces = reader.readUint16()
	self.fields = readMembers(reader, self.constantPool)
	self.methods = readMembers(reader, self.constantPool)
	self.attribute = readMembers(reader, self.constantPool)
}

func (self *ClassFile) readAndCheckMagic(reader *ClassReader) {
	magic := reader.readUint32()
	if magic != 0xCAFEBABE { // the head of balid class file.
		// There will throw a exception, but we have not implement exception.
		panic("java.lang.ClassFormatError: magic")
	}
}

func (self *ClassFile) readAndCheckVersion(reader *ClassReader) {
	self.minorVersion = reader.readUint16()
	self.majorVersion = reader.readUint16()

	// java-se-8->52, java-se-7->51, java-se-6->50, java-se-5->49,
	// j2se-1.4->48, j2se-1.3->47, j2se-1.2->46, jdk1.1->45~45.65535,
	// jdk1.0.2->45~45.3
	switch self.majorVersion {
	case 45:
		return
	case 46, 47, 48, 49, 50, 51, 52:
		if self.MinorVersion == 0 { // why it? I think that it is meaningless.
			return
		}
	}
	panic("java.lang.UnsupportedClassVersionError!")
}

func (self *ClassFile) MinorVersion() uint16 { //getter
	return self.minorVersion
}

func (self *ClassFile) MajorVersion() uint16 { //getter
	return self.majorVersion
}

func (self *ClassFile) ConstantPool() ConstantPool { //getter
	return self.constantPool
}

func (self *ClassFile) AccessFlags() uint16 { //getter
	return self.AccessFlags
}

func (self *ClassFile) Fields() []*MemberInfo { //getter
	return self.fields
}

func (self *ClassFile) Methods() []*MemberInfo { //getter
	return self.methods
}

func (self *ClassFile) ClassName() string {
	return self.constantPool.getClassName(self.thisClass)
}

func (self *ClassFile) SuperClassName() string {
	if self.superClass > 0 {
		return self.constantPool.getClassName(self.superClass)
	}
	return ""
}

func (self *ClassFile) InterFaceNames() []string {
	interfaceNames := make([]string, len(self.interfaces))
	for i, cpIndex := range self.interfaces {
		interfaceName[i] = self.constantPool.getClassName(cpIndex)
	}
	return interfaceNames
}
