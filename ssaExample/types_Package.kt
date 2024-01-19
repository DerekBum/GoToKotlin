class types_Package {

	var path: String? = null
	var name: String? = null
	var scope: types_Scope? = null
	var imports: Array<types_Package>? = null
	var complete: Boolean? = null
	var fake: Boolean? = null
	var cgo: Boolean? = null
	var goVersion: String? = null
}
