package GoToJava

import java.io.BufferedReader
import jacodbInst.*
class ssa_CallCommon : ssaToJacoExpr {

	var Value: Any? = null
	var Method: types_Func? = null
	var Args: List<Any>? = null
	var pos: Long? = null

	override fun createJacoDBExpr(): GoExpr {
        return GoCallExpr(
            Method!!.Object!!.typ!! as GoType,
            (Value!! as ssaToJacoValue).createJacoDBValue(),
            Args!!.map { (it as ssaToJacoValue).createJacoDBValue() }
        )
    }
}

fun read_ssa_CallCommon(buffReader: BufferedReader, id: Int): ssa_CallCommon {
	val res = ssa_CallCommon()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_CallCommon
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
    res.Value = mapDec[readType]?.invoke(buffReader, id) as Any?

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
    res.Method = mapDec[readType]?.invoke(buffReader, id) as types_Func?

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
    res.Args = mapDec[readType]?.invoke(buffReader, id) as List<Any>?

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
    res.pos = mapDec[readType]?.invoke(buffReader, id) as Long?

	buffReader.readLine()
	return res
}
