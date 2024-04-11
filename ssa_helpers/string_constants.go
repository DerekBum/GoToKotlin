package ssa_helpers

import "fmt"

const jacoImport = `import jacodb.*
`

const jacoInstImport = `import jacodb.GoInstLocationImpl
`

const jacoTypeImport = `import jacodb.GoType
`

const structDefinitionWithInterface = `class %s : %s {

`

const ssaToJacoExpr = `import jacodb.GoExpr

interface ssaToJacoExpr {
    fun createJacoDBExpr(): GoExpr
}
`

const ssaToJacoInst = `import jacodb.GoInst
import jacodb.GoMethod

interface ssaToJacoInst {
    fun createJacoDBInst(parent: GoMethod): GoInst
}
`

const ssaToJacoValue = `import jacodb.GoValue

interface ssaToJacoValue {
    fun createJacoDBValue(): GoValue
}
`

const createValueFunc = `override fun createJacoDBValue(): GoValue {
		if (structToPtrMap.containsKey(this) && ptrToJacoMap.containsKey(structToPtrMap[this])) {
            return ptrToJacoMap[structToPtrMap[this]] as %s
        }
        return createJacoDBExpr()
    }
`

const checkUsed = `if (structToPtrMap.containsKey(this) && ptrToJacoMap.containsKey(structToPtrMap[this])) {
            return ptrToJacoMap[structToPtrMap[this]] as %s
        }
`

const markUsed = `if (structToPtrMap.containsKey(this)) {
            ptrToJacoMap[structToPtrMap[this]!!] = res
        }
        return res`

var ssaCallExpr = fmt.Sprintf(`import jacodb.*

class ssa_CallExpr(init: ssa_Call) : ssaToJacoExpr, ssaToJacoValue {
    val type = init.register!!.typ!! as GoType
    val value = (init.Call!!.Value!! as ssaToJacoValue).createJacoDBValue()
    val operands = init.Call!!.Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }

    override fun createJacoDBExpr(): GoCallExpr {
		%s
        val res = GoCallExpr(
            type,
            value,
            operands
        )
		%s
    }
	%s
}
`, fmt.Sprintf(checkUsed, "GoCallExpr"), markUsed, fmt.Sprintf(createValueFunc, "GoCallExpr"))

var functionExtra = fmt.Sprintf(`
	fun createJacoDBMethod(): GoFunction {
		%s

        val returns = mutableListOf<GoType>()

        if (Signature!!.results!!.vars != null) {
            for (ret in Signature!!.results!!.vars!!) {
                returns.add(ret.Object!!.typ!! as GoType)
            }
        }

        val res =
            GoFunction(
                Signature!!,
                Params!!.map { it.createJacoDBExpr() }, //TODO
                name!!,
                listOf(),
                returns, //TODO
                Pkg?.Pkg?.name ?: "null"
            )

		if (structToPtrMap.containsKey(this)) {
            ptrToJacoMap[structToPtrMap[this]!!] = res
        }

        res.blocks = Blocks!!.map { it.createJacoDBBasicBlock(res) }

		return res
    }
	
	override fun createJacoDBValue(): GoValue {
		if (structToPtrMap.containsKey(this) && ptrToJacoMap.containsKey(structToPtrMap[this])) {
            return ptrToJacoMap[structToPtrMap[this]] as GoFunction
        }
        return createJacoDBMethod()
    }

	override fun createJacoDBExpr(): GoExpr {
        return createJacoDBValue()
    }
`, fmt.Sprintf(checkUsed, "GoFunction"))

var programExtra = fmt.Sprintf(`
	fun createJacoDBProject(): GoProject {
		%s

        val methods = mutableListOf<GoMethod>()
        for (pkg in packages!!) {
            for (member in pkg.value.Members!!) {
                if (member.value is ssa_Function) {
                    methods.add((member.value as ssa_Function).createJacoDBMethod())
                }
            }
        }

        val res = GoProject(
            methods.toList()
        )
		%s
    }
`, fmt.Sprintf(checkUsed, "GoProject"), markUsed)

var basicBlockExtra = fmt.Sprintf(`
	fun createJacoDBBasicBlock(method: GoMethod): GoBasicBlock {
		%s

        val inst = mutableListOf<GoInst>()

        for (value in Instrs!!) {
            if (value is ssaToJacoInst) {
                inst.add(value.createJacoDBInst(method))
            }
        }

        val res = GoBasicBlock(
            Index!!.toInt(),
            Preds!!.map { it.Index!!.toInt() },
            Succs!!.map { it.Index!!.toInt() },
            inst
        )
		%s
    }
`, fmt.Sprintf(checkUsed, "GoBasicBlock"), markUsed)

var jumpExtra = fmt.Sprintf(`
	override fun createJacoDBInst(parent: GoMethod): GoJumpInst {
		%s

        val res = GoJumpInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                0,
                parent,
            ),
            parent,
            GoInstRef(
                anInstruction!!.block!!.Index!!.toInt()
            )
        )
		%s
    }
`, fmt.Sprintf(checkUsed, "GoJumpInst"), markUsed)

var ifExtra = fmt.Sprintf(`
	override fun createJacoDBInst(parent: GoMethod): GoIfInst {
		%s

        var cond: GoConditionExpr

        val trueConst = ssa_Const()
        trueConst.Value = true
        val type = types_Basic()
        type.kind = 1
        type.info = 1
        type.name = "bool"
        trueConst.typ = type

        if (Cond!! is ssa_BinOp) {
            val parsed = (Cond!! as ssa_BinOp).createJacoDBExpr()
            if (parsed is GoConditionExpr) {
                cond = parsed
            } else {
                cond = GoEqlExpr(
                    lhv = trueConst.createJacoDBExpr(),
                    rhv = parsed as GoValue,
                    type = type,
                )
            }
        } else {
            cond = GoEqlExpr(
                lhv = trueConst.createJacoDBExpr(),
                rhv = (Cond!! as ssaToJacoValue).createJacoDBValue(),
                type = type,
            )
        }

        val res = GoIfInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                0,
                parent,
            ),
            parent,
            cond,
            GoInstRef(
                anInstruction!!.block!!.Succs!![0].Index!!.toInt()
            ),
            GoInstRef(
                anInstruction!!.block!!.Succs!![1].Index!!.toInt()
            ),
        )
		%s
    }
`, fmt.Sprintf(checkUsed, "GoIfInst"), markUsed)

const returnExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoReturnInst {
        return GoReturnInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
            Results!!.map { (it as ssaToJacoValue).createJacoDBValue() },
        )
    }
`

const runDefersExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoRunDefersInst {
        return GoRunDefersInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                0,
                parent,
            ),
            parent,
        )
    }
`

const panicExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoPanicInst {
        return GoPanicInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const goExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoGoInst {
        return GoGoInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
			(Call!!.Value!! as ssaToJacoValue).createJacoDBValue(),
            Call!!.Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }
        )
    }
`

const deferExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoDeferInst {
        return GoDeferInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
			(Call!!.Value!! as ssaToJacoValue).createJacoDBValue(),
            Call!!.Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }
        )
    }
`

const sendExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoSendInst {
        return GoSendInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
			(Chan!! as ssaToJacoValue).createJacoDBValue(),
			(X!! as ssaToJacoExpr).createJacoDBExpr(),
        )
    }
`

const storeExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoStoreInst {
        return GoStoreInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
            (Addr!! as ssaToJacoValue).createJacoDBValue(),
            (Val!! as ssaToJacoValue).createJacoDBValue()
        )
    }
`

const mapUpdateExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoMapUpdateInst {
        return GoMapUpdateInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
			(Map!! as ssaToJacoValue).createJacoDBValue(),
			(Key!! as ssaToJacoExpr).createJacoDBExpr(),
			(Value!! as ssaToJacoExpr).createJacoDBExpr(),
        )
    }
`

const debugRefExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoDebugRefInst {
        return GoDebugRefInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                pos!!.toInt(),
                parent,
            ),
            parent,
        )
    }
`

var callExtra = fmt.Sprintf(`
	var CallExpr: GoCallExpr? = null
	override fun createJacoDBInst(parent: GoMethod): GoCallInst {
        if (CallExpr == null) {
            CallExpr = ssa_CallExpr(this).createJacoDBExpr()
        }
        return GoCallInst(
            GoInstLocationImpl(
                register!!.anInstruction!!.block!!.Index!!.toInt(),
                Call!!.pos!!.toInt(),
                parent,
            ),
            parent,
            CallExpr!!
        )
    }
	
	override fun createJacoDBValue(): GoValue {
        if (CallExpr == null) {
            CallExpr = ssa_CallExpr(this).createJacoDBExpr()
        }
        return CallExpr!!
    }

	override fun createJacoDBExpr(): GoExpr {
        return createJacoDBValue()
    }
`)

var freeVarExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoFreeVar {
        return GoFreeVar(
            pos!!.toInt(),
            name!!,
            typ!! as GoType
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoFreeVar"))

var parameterExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoParameter {
        return GoParameter(
            Object!!.Object!!.pos!!.toInt(),
            name!!,
            typ!! as GoType
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoParameter"))

// TODO()
var constExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoConst {
        val innerVal = Value
        val name: String

        when (innerVal) {
            is Long -> {
                name = GoLong(
                    innerVal,
                    typ!! as GoType
                ).toString()
            }
            is Boolean -> {
                name = GoBool(
                    innerVal,
                    typ!! as GoType
                ).toString()
            }
            is Double -> {
                name = GoDouble(
                    innerVal,
                    typ!! as GoType
                ).toString()
            }
            is String -> {
                name = GoStringConstant(
                    innerVal,
                    typ!! as GoType
                ).toString()
            }
            is constant_intVal -> {
                name = GoStringConstant(
                    innerVal.toString(),
                    typ!! as GoType
                ).toString()
            }
            is constant_stringVal -> {
                name = GoStringConstant(
                    innerVal.toString(),
                    typ!! as GoType
                ).toString()
            }
            is constant_ratVal -> {
                name = GoStringConstant(
                    innerVal.toString(),
                    typ!! as GoType
                ).toString()
            }
			is constant_floatVal -> {
                name = GoStringConstant(
                    innerVal.toString(),
                    typ!! as GoType
                ).toString()
            }
			is constant_complexVal -> {
                name = GoStringConstant(
                    innerVal.toString(),
                    typ!! as GoType
                ).toString()
            }
            else -> {
                name = GoNullConstant().toString()
            }
        }

        return GoConst(
            0,
            name,
            typ!! as GoType
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoConst"))

const intValStub = `package GoToJava

class constant_intVal {}
`

const stringValStub = `package GoToJava

class constant_stringVal {}
`

const ratValStub = `package GoToJava

class constant_ratVal {}
`

const floatValStub = `package GoToJava

class constant_floatVal {}
`

const complexValStub = `package GoToJava

class constant_complexVal {}
`

const intValExtra = `
	override fun toString(): String {
        var num = Val!!.abs!!.joinToString { it.toString() }
        if (Val!!.neg!!) {
            num = "-$num"
        }
        return num
    }
`

const stringValExtra = `
    override fun toString(): String {
        var str = ""
        if (s != null) {
            str += s!!
        }
        if (l != null) {
            str += l!!.toString()
            str += r!!.toString()
        }
        return str
    }
`

const ratValExtra = `
    override fun toString(): String {
        return "${Val!!.a}/${Val!!.b}"
    }
`

const floatValExtra = `
    override fun toString(): String {
		var str = "2^${Val!!.exp!!}"
		val temp = ""
		for (w in Val!!.mant!!) {
			temp += w.toString()
		}
		str = "$temp * $str"
		if (Val!!.neg!!) {
			str = "-$str"
		}
        return str
    }
`

const complexValExtra = `
    override fun toString(): String {
        return "(${re!!} + ${im!!}i)"
    }
`

var globalExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoGlobal {
        return GoGlobal(
            pos!!.toInt(),
            name!!,
            typ!! as GoType
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoGlobal"))

var builtinExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoBuiltin {
        return GoBuiltin(
            0,
            name!!,
            sig!!
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoBuiltin"))

const arrayExtra = `
	override val typeName: String
        get() = "[${len!!}]${(elem!! as GoType).typeName}"
`

const basicExtra = `
	override val typeName: String
        get() = name!!
`

const chanExtra = `
	override val typeName: String
        get(): String {
            var res = (elem!! as GoType).typeName
            if (dir!! == 0L) {
                res = "chan $res"
            }
            else if (dir!! == 1L) {
                res = "<-chan $res"
            }
            else if (dir!! == 2L) {
                res = "chan <-$res"
            }
            return res
        }
`

const interfaceExtra = `
	override val typeName: String
        get() = "Any"
`

const mapExtra = `
	override val typeName: String
        get() = "map[${(key!! as GoType).typeName}]${(elem!! as GoType).typeName}"
`

const namedExtra = `
	override val typeName: String
        get() = (underlying!! as GoType).typeName
`

const pointerExtra = `
	override val typeName: String
        get() = "*${(base!! as GoType).typeName}"
`

const signatureExtra = `
	override val typeName: String
        get(): String {
            var res = "func ("
            var paramsString = ""
            for (p in params!!.vars!!) {
                paramsString += p.Object!!.name + ", "
            }
            res += paramsString.removeSuffix(", ") + ") ("
            var resultsString = ""
            for (r in results!!.vars!!) {
                resultsString += r.Object!!.name + ", "
            }
            res += resultsString.removeSuffix(", ") + ")"
            return res
        }
`

const sliceTypeExtra = `
	override val typeName: String
        get() = "[]${(elem!! as GoType).typeName}"
`

const structExtra = `
	override val typeName: String
        get(): String {
            var res = "struct {\n"
            fields!!.forEachIndexed { ind, elem ->
                res += (elem.Object!!.typ!! as GoType).typeName
                if (tags != null && tags!!.size > ind) {
                    res += " " + tags!![ind]
                }
                res += "\n"
            }
            res += "}"
            return res
        }
`

const tupleExtra = `
	override val typeName: String
        get(): String {
            var res = "["
            for (i in vars!!) {
                res += i.Object!!.name + ", "
            }
            return res.removeSuffix(", ") + "]"
        }
`

const typeParamExtra = `
	override val typeName: String
        get() = obj!!.Object!!.name!!
`

const unionExtra = `
	override val typeName: String
        get(): String {
            var res = "enum {\n"
            for (t in terms!!) {
                res += (t.typ!! as GoType).typeName + ",\n"
            }
            return "$res}"
        }
`

const opaqueTypeExtra = `
	override val typeName: String
        get() = name!!
`

var allocExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoAllocExpr {
        return GoAllocExpr(
            register!!.typ!! as GoType,
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoAllocExpr"))

var phiExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoPhiExpr {
		%s

        val res = GoPhiExpr(
            register!!.typ!! as GoType,
			listOf()
        )
        if (structToPtrMap.containsKey(this)) {
            ptrToJacoMap[structToPtrMap[this]!!] = res
        }
        
        res.edges = Edges!!.map { (it as ssaToJacoValue).createJacoDBValue() }
		return res
    }
	%s
`, fmt.Sprintf(checkUsed, "GoPhiExpr"), fmt.Sprintf(createValueFunc, "GoPhiExpr"))

const binOpExtra = `
	override fun createJacoDBExpr(): GoBinaryExpr {
        val type = register!!.typ!! as GoType

        when (Op!!) {
            12L -> return GoAddExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            13L -> return GoSubExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            14L -> return GoMulExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            15L -> return GoDivExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            16L -> return GoModExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            17L -> return GoAndExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            18L -> return GoOrExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            19L -> return GoXorExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            20L -> return GoShlExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            21L -> return GoShrExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            22L -> return GoAndNotExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            39L -> return GoEqlExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            44L -> return GoNeqExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            40L -> return GoLssExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            45L -> return GoLeqExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            41L -> return GoGtrExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            46L -> return GoGeqExpr(
                lhv = (X!! as ssaToJacoValue).createJacoDBValue(),
                rhv = (Y!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            else -> error("unexpected BinOp ${Op!!}")
        }
    }

	override fun createJacoDBValue(): GoValue {
        val res = createJacoDBExpr()
        if (res is GoValue) {
            return res
        }
        error("unexpected cast to Value $res")
    }
`

var unOpExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoUnaryExpr {
        val type = register!!.typ!! as GoType

        when (Op!!) {
            43L -> return GoUnNotExpr(
                value = (X!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            13L -> return GoUnSubExpr(
                value = (X!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            36L -> return GoUnArrowExpr(
                value = (X!! as ssaToJacoValue).createJacoDBValue(),
                type = type,
                commaOk = CommaOk ?: false
            )
            14L -> return GoUnMulExpr(
                value = (X!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            19L -> return GoUnXorExpr(
                value = (X!! as ssaToJacoValue).createJacoDBValue(),
                type = type
            )
            else -> error("unexpected UnOp ${Op!!}")
        }
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoUnaryExpr"))

var changeTypeExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoChangeTypeExpr {
        return GoChangeTypeExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoChangeTypeExpr"))

var convertExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoConvertExpr {
        return GoConvertExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoConvertExpr"))

var multiConvertExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoMultiConvertExpr {
        return GoMultiConvertExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoMultiConvertExpr"))

var changeInterfaceExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoChangeInterfaceExpr {
        return GoChangeInterfaceExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoChangeInterfaceExpr"))

var sliceToArrayPointerExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoSliceToArrayPointerExpr {
        return GoSliceToArrayPointerExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoSliceToArrayPointerExpr"))

var makeInterfaceExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoMakeInterfaceExpr {
        return GoMakeInterfaceExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoMakeInterfaceExpr"))

var makeClosureExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoMakeClosureExpr {
        return GoMakeClosureExpr(
			register!!.typ!! as GoType,
            (Fn!! as ssa_Function).createJacoDBMethod(),
			Bindings!!.map { (it as ssaToJacoValue).createJacoDBValue() },
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoMakeClosureExpr"))

var makeMapExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoMakeMapExpr {
		val reserve = if (Reserve == null) {
            GoLong(0, LongType())
        } else {
            (Reserve!! as ssaToJacoValue).createJacoDBValue()
        }

        return GoMakeMapExpr(
			register!!.typ!! as GoType,
            reserve,
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoMakeMapExpr"))

var makeChanExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoMakeChanExpr {
        return GoMakeChanExpr(
			register!!.typ!! as GoType,
            (Size!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoMakeChanExpr"))

var makeSliceExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoMakeSliceExpr {
        return GoMakeSliceExpr(
			register!!.typ!! as GoType,
            (Len!! as ssaToJacoValue).createJacoDBValue(),
			(Cap!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoMakeSliceExpr"))

var sliceExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoSliceExpr {
		%s

		val low: GoValue = if (Low == null) {
            GoNullConstant()
        } else {
            (Low!! as ssaToJacoValue).createJacoDBValue()
        }
        val high: GoValue = if (High == null) {
            GoNullConstant()
        } else {
            (High!! as ssaToJacoValue).createJacoDBValue()
        }
        val max: GoValue = if (Max == null) {
            GoNullConstant()
        } else {
            (Max!! as ssaToJacoValue).createJacoDBValue()
        }

		val res = GoSliceExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			low,
			high,
			max,
        )
		%s
    }
	%s
`, fmt.Sprintf(checkUsed, "GoSliceExpr"), markUsed, fmt.Sprintf(createValueFunc, "GoSliceExpr"))

var fieldAddrExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoFieldAddrExpr {
        return GoFieldAddrExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			Field!!.toInt(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoFieldAddrExpr"))

var fieldExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoFieldExpr {
        return GoFieldExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			Field!!.toInt(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoFieldExpr"))

var indexAddrExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoIndexAddrExpr {
        return GoIndexAddrExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Index!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoIndexAddrExpr"))

var indexExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoIndexExpr {
        return GoIndexExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Index!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoIndexExpr"))

var lookupExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoLookupExpr {
        return GoLookupExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Index!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoLookupExpr"))

var selectExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoSelectExpr {
        val chan = mutableListOf<GoValue>()
        val send = mutableListOf<GoValue>()
        if (States != null) {
            States!!.map {
                chan.add((it.Chan!! as ssaToJacoValue).createJacoDBValue())
                send.add(
                    if (it.Send == null) {
                        GoNullConstant()
                    } else {
                        (it.Send!! as ssaToJacoValue).createJacoDBValue()
                    }
                )
            }
        }

        return GoSelectExpr(
			register!!.typ!! as GoType,
            chan,
			send,
			Blocking!!,
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoSelectExpr"))

var rangeExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoRangeExpr {
        return GoRangeExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoRangeExpr"))

var nextExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoNextExpr {
        return GoNextExpr(
			register!!.typ!! as GoType,
            (Iter!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoNextExpr"))

var typeAssertExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoTypeAssertExpr {
        return GoTypeAssertExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			AssertedType!! as GoType,
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoTypeAssertExpr"))

var extractExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoExtractExpr {
        return GoExtractExpr(
			register!!.typ!! as GoType,
            (Tuple!! as ssaToJacoValue).createJacoDBValue(),
			Index!!.toInt(),
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoExtractExpr"))

/*var callCommonExtra = fmt.Sprintf(`
	override fun createJacoDBExpr(): GoCallExpr {
        return GoCallExpr(
            Method!!.Object!!.typ!! as GoType,
            (Value!! as ssaToJacoValue).createJacoDBValue(),
            Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }
        )
    }
	%s
`, fmt.Sprintf(createValueFunc, "GoCallExpr"))*/
