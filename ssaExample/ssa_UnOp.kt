package GoToJava

import java.io.BufferedReader
import jacodbInst.*
class ssa_UnOp : ssaToJacoExpr, ssaToJacoValue {

	var register: ssa_register? = null
	var Op: Long? = null
	var X: Any? = null
	var CommaOk: Boolean? = null

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
}

fun read_ssa_UnOp(buffReader: BufferedReader, id: Int): ssa_UnOp {
	val res = ssa_UnOp()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_UnOp
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
    res.register = mapDec[readType]?.invoke(buffReader, id) as ssa_register?

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
    res.Op = mapDec[readType]?.invoke(buffReader, id) as Long?

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
    res.X = mapDec[readType]?.invoke(buffReader, id) as Any?

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
    res.CommaOk = mapDec[readType]?.invoke(buffReader, id) as Boolean?

	buffReader.readLine()
	return res
}
