package ssa_helpers

const jacoImport = `import jacodbInst.*
`

const jacoInstImport = `import jacodbInst.impl.location.GoInstLocationImpl
`

const jacoTypeImport = `import jacodbInst.GoType
`

const structDefinitionWithInterface = `class %s : %s {

`

const ssaToJacoExpr = `import jacodbInst.GoExpr

interface ssaToJacoExpr {
    fun createJacoDBExpr(): GoExpr
}
`

const ssaToJacoInst = `import jacodbInst.GoInst
import jacodbInst.GoMethod

interface ssaToJacoInst {
    fun createJacoDBInst(parent: GoMethod): GoInst
}
`

const ssaToJacoValue = `import jacodbInst.GoValue

interface ssaToJacoValue {
    fun createJacoDBValue(): GoValue
}
`

const ssaCallExpr = `import jacodbInst.*

class ssa_CallExpr(init: ssa_Call) : ssaToJacoExpr, ssaToJacoValue {
    val type = init.register!!.typ!! as GoType
    val value = (init.Call!!.Value!! as ssaToJacoValue).createJacoDBValue()
    val operands = init.Call!!.Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }

    override fun createJacoDBExpr(): GoCallExpr {
        return GoCallExpr(
            type,
            value,
            operands
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
}
`

const functionExtra = `
	fun createJacoDBMethod(): GoFunction {
        val returns = mutableListOf<GoType>()

        if (Signature!!.results!!.vars != null) {
            for (ret in Signature!!.results!!.vars!!) {
                returns.add(ret.Object!!.typ!! as GoType)
            }
        }

        val noBlocksFunction =
            GoFunction(
                Signature!!,
                Params!!.map { it.createJacoDBExpr() }, // TODO
                name!!,
                listOf(),
                returns, //TODO
                Pkg!!.Pkg!!.name!!
            )

        noBlocksFunction.blocks = Blocks!!.map { it.createJacoDBBasicBlock(noBlocksFunction) }

        return noBlocksFunction
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBMethod()
    }
`

const programExtra = `
	fun createJacoDBProject(): GoProject {
        val methods = mutableListOf<GoMethod>()
        for (pkg in packages!!) {
            for (member in pkg.value.Members!!) {
                if (member.value is ssa_Function) {
                    methods.add((member.value as ssa_Function).createJacoDBMethod())
                }
            }
        }

        return GoProject(
            methods.toList()
        )
    }
`

const basicBlockExtra = `
	fun createJacoDBBasicBlock(method: GoMethod): GoBasicBlock {
        val inst = mutableListOf<GoInst>()

        for (value in Instrs!!) {
            if (value is ssaToJacoInst) {
                inst.add(value.createJacoDBInst(method))
            }
        }

        return GoBasicBlock(
            Index!!.toInt(),
            Preds!!.map { it.Index!!.toInt() },
            Succs!!.map { it.Index!!.toInt() },
            inst
        )
    }
`

const jumpExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoJumpInst {
        return GoJumpInst(
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
    }
`

const ifExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoIfInst {
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

        return GoIfInst(
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
    }
`

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
			(Call!! as ssaToJacoExpr).createJacoDBExpr(),
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
			(Call!! as ssaToJacoExpr).createJacoDBExpr(),
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

const callExtra = `
	override fun createJacoDBInst(parent: GoMethod): GoCallInst {
        return GoCallInst(
            GoInstLocationImpl(
                register!!.anInstruction!!.block!!.Index!!.toInt(),
                Call!!.pos!!.toInt(),
                parent,
            ),
            parent,
            ssa_CallExpr(this).createJacoDBExpr()
        )
    }

    override fun createJacoDBValue(): GoValue {
        return ssa_CallExpr(this).createJacoDBValue()
    }
`

const freeVarExtra = `
	override fun createJacoDBExpr(): GoFreeVar {
        return GoFreeVar(
            pos!!.toInt(),
            name!!,
            typ!! as GoType
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
`

const parameterExtra = `
	override fun createJacoDBExpr(): GoParameter {
        return GoParameter(
            Object!!.Object!!.pos!!.toInt(),
            name!!,
            typ!! as GoType
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
`

const constExtra = `
	override fun createJacoDBExpr(): GoConst {
        val innerVal = Value!!
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
            else -> {
                name = GoNullConstant(
                    typ!! as GoType
                ).toString()
            }
        }

        return GoConst(
            0,
            name,
            typ!! as GoType
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
`

const globalExtra = `
	override fun createJacoDBExpr(): GoGlobal {
        return GoGlobal(
            pos!!.toInt(),
            name!!,
            typ!! as GoType
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
`

const builtinExtra = `
	override fun createJacoDBExpr(): GoBuiltin {
        return GoBuiltin(
            0,
            name!!,
            sig!!
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
`

const basicExtra = `
	override val typeName: String
        get() = name!!
`

const interfaceExtra = `
	override val typeName: String
        get() = "Any"
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

const allocExtra = `
	override fun createJacoDBExpr(): GoAllocExpr {
        return GoAllocExpr(
            register!!.typ!! as GoType,
        )
    }
`

const phiExtra = `
	override fun createJacoDBExpr(): GoPhiExpr {
        return GoPhiExpr(
            register!!.typ!! as GoType,
			Edges!!.map { (it as ssaToJacoValue).createJacoDBValue() }
        )
    }
`

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
`

const unOpExtra = `
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

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
`

const changeTypeExtra = `
	override fun createJacoDBExpr(): GoChangeTypeExpr {
        return GoChangeTypeExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const convertExtra = `
	override fun createJacoDBExpr(): GoConvertExpr {
        return GoConvertExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const multiConvertExtra = `
	override fun createJacoDBExpr(): GoMultiConvertExpr {
        return GoMultiConvertExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const changeInterfaceExtra = `
	override fun createJacoDBExpr(): GoChangeInterfaceExpr {
        return GoChangeInterfaceExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const sliceToArrayPointerExtra = `
	override fun createJacoDBExpr(): GoSliceToArrayPointerExpr {
        return GoSliceToArrayPointerExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const makeInterfaceExtra = `
	override fun createJacoDBExpr(): GoMakeInterfaceExpr {
        return GoMakeInterfaceExpr(
            register!!.typ!! as GoType,
			(X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const makeClosureExtra = `
	override fun createJacoDBExpr(): GoMakeClosureExpr {
        return GoMakeClosureExpr(
			register!!.typ!! as GoType,
            (Fn!! as ssa_Function).createJacoDBMethod(),
			Bindings!!.map { (it as ssaToJacoValue).createJacoDBValue() },
        )
    }
`

const makeMapExtra = `
	override fun createJacoDBExpr(): GoMakeMapExpr {
        return GoMakeMapExpr(
			register!!.typ!! as GoType,
            (Reserve!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const makeChanExtra = `
	override fun createJacoDBExpr(): GoMakeChanExpr {
        return GoMakeChanExpr(
			register!!.typ!! as GoType,
            (Size!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const makeSliceExtra = `
	override fun createJacoDBExpr(): GoMakeSliceExpr {
        return GoMakeSliceExpr(
			register!!.typ!! as GoType,
            (Len!! as ssaToJacoValue).createJacoDBValue(),
			(Cap!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const sliceExtra = `
	override fun createJacoDBExpr(): GoSliceExpr {
        return GoSliceExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Low!! as ssaToJacoValue).createJacoDBValue(),
			(High!! as ssaToJacoValue).createJacoDBValue(),
			(Max!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const fieldAddrExtra = `
	override fun createJacoDBExpr(): GoFieldAddrExpr {
        return GoFieldAddrExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			Field!!.toInt(),
        )
    }
`

const fieldExtra = `
	override fun createJacoDBExpr(): GoFieldExpr {
        return GoFieldExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			Field!!.toInt(),
        )
    }
`

const indexAddrExtra = `
	override fun createJacoDBExpr(): GoIndexAddrExpr {
        return GoIndexAddrExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Index!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const indexExtra = `
	override fun createJacoDBExpr(): GoIndexExpr {
        return GoIndexExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Index!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const lookupExtra = `
	override fun createJacoDBExpr(): GoLookupExpr {
        return GoLookupExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			(Index!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const selectExtra = `
	override fun createJacoDBExpr(): GoSelectExpr {
        return GoSelectExpr(
			register!!.typ!! as GoType,
            States!!.map { (it.Chan!! as ssaToJacoValue).createJacoDBValue() },
			States!!.map { (it.Send!! as ssaToJacoValue).createJacoDBValue() },
			Blocking!!,
        )
    }
`

const rangeExtra = `
	override fun createJacoDBExpr(): GoRangeExpr {
        return GoRangeExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const nextExtra = `
	override fun createJacoDBExpr(): GoNextExpr {
        return GoNextExpr(
			register!!.typ!! as GoType,
            (Iter!! as ssaToJacoValue).createJacoDBValue(),
        )
    }
`

const typeAssertExtra = `
	override fun createJacoDBExpr(): GoTypeAssertExpr {
        return GoTypeAssertExpr(
			register!!.typ!! as GoType,
            (X!! as ssaToJacoValue).createJacoDBValue(),
			AssertedType!! as GoType,
        )
    }
`

const extractExtra = `
	override fun createJacoDBExpr(): GoExtractExpr {
        return GoExtractExpr(
			register!!.typ!! as GoType,
            (Tuple!! as ssaToJacoValue).createJacoDBValue(),
			Index!!.toInt(),
        )
    }
`

const callCommonExtra = `
	override fun createJacoDBExpr(): GoExpr {
        return GoCallExpr(
            Method!!.Object!!.typ!! as GoType,
            (Value!! as ssaToJacoValue).createJacoDBValue(),
            Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }
        )
    }
`
