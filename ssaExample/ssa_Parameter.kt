package GoToJava

import java.io.BufferedReader
import jacodbInst.*
class ssa_Parameter : ssaToJacoExpr, ssaToJacoValue {

	var name: String? = null
	var Object: types_Var? = null
	var typ: Any? = null
	var parent: ssa_Function? = null
	var referrers: List<Any>? = null

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
}

fun read_ssa_Parameter(buffReader: BufferedReader, id: Int): ssa_Parameter {
	val res = ssa_Parameter()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_Parameter
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
    res.name = mapDec[readType]?.invoke(buffReader, id) as String?

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
    res.Object = mapDec[readType]?.invoke(buffReader, id) as types_Var?

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
    res.typ = mapDec[readType]?.invoke(buffReader, id) as Any?

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
    res.parent = mapDec[readType]?.invoke(buffReader, id) as ssa_Function?

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
    res.referrers = mapDec[readType]?.invoke(buffReader, id) as List<Any>?

	buffReader.readLine()
	return res
}
