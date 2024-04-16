package ssa_helpers

var instructions = []string{
	"ssa_Jump",
	"ssa_If",
	"ssa_Return",
	"ssa_RunDefers",
	"ssa_Panic",
	"ssa_Go",
	"ssa_Defer",
	"ssa_Send",
	"ssa_Store",
	"ssa_MapUpdate",
	"ssa_DebugRef",
	"ssa_Call",
}

var variables = []string{
	"ssa_FreeVar",
	"ssa_Parameter",
	"ssa_Const",
	"ssa_Global",
	"ssa_Builtin",
}

var expressions = append(variables, []string{
	"ssa_Alloc",
	"ssa_Phi",
	"ssa_BinOp",
	"ssa_UnOp",
	"ssa_ChangeType",
	"ssa_Convert",
	"ssa_MultiConvert",
	"ssa_ChangeInterface",
	"ssa_SliceToArrayPointer",
	"ssa_MakeInterface",
	"ssa_MakeClosure",
	"ssa_MakeMap",
	"ssa_MakeChan",
	"ssa_MakeSlice",
	"ssa_Slice",
	"ssa_FieldAddr",
	"ssa_Field",
	"ssa_IndexAddr",
	"ssa_Index",
	"ssa_Lookup",
	"ssa_Select",
	"ssa_Range",
	"ssa_Next",
	"ssa_TypeAssert",
	"ssa_Extract",
	"ssa_Call",
	"ssa_Function",
}...)

var values = append(variables, expressions...)

var types = []string{
	"types_Array",
	"types_Basic",
	"types_Chan",
	"types_Interface",
	"types_Map",
	"types_Named",
	"types_Pointer",
	"types_Signature",
	"types_Slice",
	"types_Struct",
	"types_Tuple",
	"types_TypeParam",
	"types_Union",
	"ssa_opaqueType",
}

var ifaceToStringMap = map[int]string{
	instruction_:      "ssaToJacoInst",
	expression_:       "ssaToJacoExpr",
	value_:            "ssaToJacoValue",
	type_:             "ssaToJacoType",
	expressionValue_:  "ssaToJacoExpr, ssaToJacoValue",
	instructionValue_: "ssaToJacoInst, ssaToJacoValue",
	call_:             "ssaToJacoInst, ssaToJacoValue, ssaToJacoExpr",
}

var nameToExtra = map[string]string{
	"ssa_Program":    programExtra,
	"ssa_BasicBlock": basicBlockExtra,
	"ssa_Function":   functionExtra,

	"ssa_Jump":      jumpExtra,
	"ssa_If":        ifExtra,
	"ssa_Return":    returnExtra,
	"ssa_RunDefers": runDefersExtra,
	"ssa_Panic":     panicExtra,
	"ssa_Go":        goExtra,
	"ssa_Defer":     deferExtra,
	"ssa_Send":      sendExtra,
	"ssa_Store":     storeExtra,
	"ssa_MapUpdate": mapUpdateExtra,
	"ssa_DebugRef":  debugRefExtra,
	"ssa_Call":      callExtra,

	"ssa_FreeVar":   freeVarExtra,
	"ssa_Parameter": parameterExtra,
	"ssa_Const":     constExtra,
	"ssa_Global":    globalExtra,
	"ssa_Builtin":   builtinExtra,

	"types_Array":     arrayExtra,
	"types_Basic":     basicExtra,
	"types_Chan":      chanExtra,
	"types_Interface": interfaceExtra,
	"types_Map":       mapExtra,
	"types_Named":     namedExtra,
	"types_Pointer":   pointerExtra,
	"types_Signature": signatureExtra,
	"types_Slice":     sliceTypeExtra,
	"types_Struct":    structExtra,
	"types_Tuple":     tupleExtra,
	"types_TypeParam": typeParamExtra,
	"types_Union":     unionExtra,
	"ssa_opaqueType":  opaqueTypeExtra,

	"ssa_Alloc":               allocExtra,
	"ssa_Phi":                 phiExtra,
	"ssa_BinOp":               binOpExtra,
	"ssa_UnOp":                unOpExtra,
	"ssa_ChangeType":          changeTypeExtra,
	"ssa_Convert":             convertExtra,
	"ssa_MultiConvert":        multiConvertExtra,
	"ssa_ChangeInterface":     changeInterfaceExtra,
	"ssa_SliceToArrayPointer": sliceToArrayPointerExtra,
	"ssa_MakeInterface":       makeInterfaceExtra,
	"ssa_MakeClosure":         makeClosureExtra,
	"ssa_MakeMap":             makeMapExtra,
	"ssa_MakeChan":            makeChanExtra,
	"ssa_MakeSlice":           makeSliceExtra,
	"ssa_Slice":               sliceExtra,
	"ssa_FieldAddr":           fieldAddrExtra,
	"ssa_Field":               fieldExtra,
	"ssa_IndexAddr":           indexAddrExtra,
	"ssa_Index":               indexExtra,
	"ssa_Lookup":              lookupExtra,
	"ssa_Select":              selectExtra,
	"ssa_Range":               rangeExtra,
	"ssa_Next":                nextExtra,
	"ssa_TypeAssert":          typeAssertExtra,
	"ssa_Extract":             extractExtra,

	"constant_intVal":     intValExtra,
	"constant_stringVal":  stringValExtra,
	"constant_ratVal":     ratValExtra,
	"constant_floatVal":   floatValExtra,
	"constant_complexVal": complexValExtra,
}

var nameToStub = map[string]string{
	"constant_intVal":     intValStub,
	"constant_stringVal":  stringValStub,
	"constant_ratVal":     ratValStub,
	"constant_floatVal":   floatValStub,
	"constant_complexVal": complexValStub,
}
