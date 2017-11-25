package mkl
var mkl_versions = make(map[string]string)
var mkl_licenses = make(map[string]string)

func Version(n string,v string){
    mkl_versions[n] = v
}

func Lic(n string,l string){
    mkl_licenses[n] = l
}
