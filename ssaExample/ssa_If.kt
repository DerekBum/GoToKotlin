package GoToJava

import java.io.BufferedReader
import jacodbInst.*
import jacodbInst.impl.location.GoInstLocationImpl
class ssa_If : ssaToJacoInst {

	var anInstruction: generatedInlineStruct_000? = null
	var Cond: Any? = null

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
}

fun read_ssa_If(buffReader: BufferedReader, id: Int): ssa_If {
	val res = ssa_If()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_If
        }
        ptrMap[id] = res
    }
    var line: String
    var split: List<String>
    var id: Int
    var readType: String

	line = buffReader.readLine()
	if (line == "end") {
        return res
    }
    split = line.split(" ")
    readType = split[1]
    id = -1
    if (split.size > 2) {
        id = split[2].toInt()
    }
    res.anInstruction = mapDec[readType]?.invoke(buffReader, id) as generatedInlineStruct_000?

	line = buffReader.readLine()
	if (line == "end") {
        return res
    }
    split = line.split(" ")
    readType = split[1]
    id = -1
    if (split.size > 2) {
        id = split[2].toInt()
    }
    res.Cond = mapDec[readType]?.invoke(buffReader, id) as Any?

	buffReader.readLine()
	return res
}
