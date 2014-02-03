package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// put iperf specific logic here

//var findNamesRegex = regexp.MustCompile(`name="([^"]+)"`)
var findNamesRegex = regexp.MustCompile(`{{.([^}]+)}}`)

func GrepNamesFromTemplateFile(filepath string) ([]string, []byte) {
	t, err := ioutil.ReadFile(filepath)
	if err != nil {
		msg := fmt.Sprintf("GrepNamesFromTemplateFile() got bad path: '%v', err: %v", filepath, err)
		panic(msg)
	}
	nam := findNamesRegex.FindAllStringSubmatch(string(t), -1)

	fmt.Printf("nam is %v\n", nam)

	names := make([]string, len(nam))
	for i, v := range nam {
		names[i] = v[1]
	}

	//fmt.Printf("names is %v\n", names)
	return names, t
}

var findCheckboxesRegex = regexp.MustCompile(`type="checkbox"[^n]+name="([^"]+)"`)

func GrepCheckboxesFromTemplate(t []byte) []string {
	chk := findCheckboxesRegex.FindAllStringSubmatch(string(t), -1)

	//fmt.Printf("chk is %v\n", chk)

	checks := make([]string, len(chk))
	for i, v := range chk {
		checks[i] = v[1]
	}

	//fmt.Printf("checks is %v\n", checks)
	return checks
}

// In order to locate the radio buttons automagically, we require that
// all radio buttons be laid out in the template iperf.html file in
// this exact order: type, then name, then value, then {{.name_value}}
//
// For example:
//
// <input type="radio" name="trafprot" value="trafprot_tcp" {{.trafprot_tcp}}>TCP</input>
//
// The regex won't work otherwise.
//
var findRadioBtnRegex = regexp.MustCompile(`type="radio"[^n]+name="([^"]+)"\svalue="([^"]+)"`)

func GrepRadioBtnFromTemplate(t []byte) *map[string][]string {
	rad := findRadioBtnRegex.FindAllStringSubmatch(string(t), -1)

	fmt.Printf("rad is %v\n", rad)

	radios := make(map[string][]string, len(rad))
	for _, v := range rad {
		// example of what we are doing here:
		//  radios["trafprot] = append(radios["trafprot"], "tcp")
		//  radios["trafprot] = append(radios["trafprot"], "udp")
		radios[v[1]] = append(radios[v[1]], v[2])
	}

	fmt.Printf("radios is %v\n", radios)
	return &radios
}

func FormToIperfMap(form *map[string][]string) *map[string]string {
	// get a list of every value the map should contain,
	// to avoid rendering <no value> for missing values.
	values, tmplate := GrepNamesFromTemplateFile("templates/iperf.html")

	m := make(map[string]string)

	// first extract from form
	for k, v := range *form {
		m[k] = strings.Join(v, " ")
	}

	fmt.Printf("after form extraction, m is %#+v\n", m)

	// next add empty strings for any missing values...to
	// avoid producing "<no value>"
	for _, v := range values {
		_, ok := m[v]
		if !ok {
			m[v] = ""
		}
	}

	// and, preserve state of checkboxes, radios, and dropdowns
	PreserveCheckboxes(tmplate, &m)
	PreserveRadioBtn(tmplate, &m)
	PreserveDropDowns(&m, form)
	PreserveSelectedFiles(tmplate, &m)
	PreserveSelectedTab(tmplate, &m)

	return &m
}

func PreserveCheckboxes(tmplate []byte, m *map[string]string) {
	checkboxes := GrepCheckboxesFromTemplate(tmplate)

	for _, chk := range checkboxes {
		if (*m)[chk] == (chk + "on") {
			(*m)[chk] = `checked="checked"`
		} else {
			(*m)[chk] = ""
		}
	}
}

//
// PreserveRadioBtn() intent:
// ex: radios["trafprot"] will contain []string{"tcp", "udp"}
// and we need to set trafprot_tcp and trafprot_udp to `` or to `checked="checked"`
// depending on whether (*m)["trafprot"] == "udp" or "tcp"
//
func PreserveRadioBtn(tmplate []byte, m *map[string]string) {
	radios := GrepRadioBtnFromTemplate(tmplate)

	for rad, allOptions := range *radios {
		// ex: rad could be "trafprot", allOptions []string{"tcp", "udp"}
		//fmt.Printf("PreserveRadioBtn(): rad is %v,  allOptions is %v\n", rad, allOptions)

		for _, rad_option := range allOptions {
			// ex: rad_option could be "tcp"

			//fmt.Printf("PreserveRadioBtn(): next rad_option is %v\n", rad_option)

			if (*m)[rad] == rad_option {
				(*m)[rad+"_"+rad_option] = `checked="checked"`
			} else {
				(*m)[rad+"_"+rad_option] = ""
			}
		}
	}
}

var findFileSelectorRegex = regexp.MustCompile(`type="file"[^n]+name="([^"]+)"`)

func GrepSelectedFilesFromTemplate(t []byte) []string {
	re := findFileSelectorRegex.FindAllStringSubmatch(string(t), -1)

	fmt.Printf("re is %v\n", re)

	files := make([]string, len(re))
	for i, v := range re {
		files[i] = v[1]
	}

	fmt.Printf("files is %v\n", files)
	return files
}

func PreserveSelectedFiles(tmplate []byte, m *map[string]string) {

	selfiles := GrepSelectedFilesFromTemplate(tmplate)
	// ex: selfiles =

	for _, f := range selfiles {
		if (*m)[f] == "" {
		}
	}
}

func PreserveSelectedTab(tmplate []byte, m *map[string]string) {

	// make reload work for tab display even if no generate called.
	if (*m)["selectedTab"] == "" {
		(*m)["selectedTab"] = "0"
	}
}

// fragile/keep me in sync: the drop-downs in the iperf.html file
// need to be kept in sync with the names and lists below here in iperf.go
func PreserveDropDowns(m *map[string]string, umap *map[string][]string) {
	PreserveDropDownsKeyChoices(m, "rptformat", []string{"Kb", "Mb", "KB", "MB"})
	PreserveDropDownsKeyChoices(m, "congestionalgo", []string{"default", "reno", "yeah", "highspeed", "bic", "vegas"})
	PreserveDropDownsKeyChoices(m, "bandwidthUnits", []string{"Kb", "Mb", "Gb"})
	PreserveDropDownsKeyChoices(m, "TTtimeUnits", []string{"ttts", "tttm", "ttth"})

	PreserveMultipleDropDownsKeyChoices(m, "exclusions", []string{"con", "dat", "mca", "set", "srv"}, umap)
}

func PreserveMultipleDropDownsKeyChoices(m *map[string]string, key string, choices []string, umap *map[string][]string) {
	// null out all choices
	for _, ch := range choices {
		(*m)[key+"_"+ch] = ""
	}
	// multiple can be selected
	chosen, ok := (*umap)[key]
	if !ok {
		return
	}

	for _, ch := range chosen {
		(*m)[key+"_"+ch] = `selected="selected"`
	}
}

func PreserveDropDownsKeyChoices(m *map[string]string, key string, choices []string) {

	// null out all choices
	for _, ch := range choices {
		(*m)[key+"_"+ch] = ""
	}
	chosen, ok := (*m)[key]
	if !ok {
		return
	}
	// then select just the one
	(*m)[key+"_"+chosen] = `selected="selected"`
}

func reportExcludeCode(ex []string) string {
	res := ""

	for _, s := range ex {
		switch s {
		case "con":
			res += "C"
		case "dat":
			res += "D"
		case "mca":
			res += "M"
		case "set":
			res += "S"
		case "srv":
			res += "V"
		}
	}

	return res
}

func SetFormDefaults(v url.Values, m *map[string]string) {

	(*m)["trafprot_tcp"] = `checked="checked"`
	(*m)["ipv_ipv4"] = `checked="checked"`
	(*m)["bidiropt_sim"] = `checked="checked"`
	(*m)["congestionalgo_default"] = `checked="checked"`

	(*m)["Cv4o1"] = "10"
	(*m)["Cv4o2"] = "0"
	(*m)["Cv4o3"] = "1"
	(*m)["Cv4o4"] = "20"

	(*m)["Sv4o1"] = "10"
	(*m)["Sv4o2"] = "0"
	(*m)["Sv4o3"] = "1"
	(*m)["Sv4o4"] = "10"

	(*m)["Cv6o1"] = "fefe"
	(*m)["Cv6o2"] = "0"
	(*m)["Cv6o3"] = "0"
	(*m)["Cv6o4"] = "0"
	(*m)["Cv6o5"] = "10"
	(*m)["Cv6o6"] = "0"
	(*m)["Cv6o7"] = "1"
	(*m)["Cv6o8"] = "20"

	(*m)["Sv6o1"] = "fefe"
	(*m)["Sv6o2"] = "0"
	(*m)["Sv6o3"] = "0"
	(*m)["Sv6o4"] = "0"
	(*m)["Sv6o5"] = "10"
	(*m)["Sv6o6"] = "0"
	(*m)["Sv6o7"] = "1"
	(*m)["Sv6o8"] = "10"

	// set the default tab to display as the first one
	(*m)["selectedTab"] = "0"

}

func GenIperfCmd(v url.Values, m *map[string]string) string {

	Ccmd := "iperf --client "
	Scmd := "iperf --server --bind "
	Cip := ""
	Sip := ""

	v0 := ""
	_, ok := v["ipv"]
	if ok {
		if v["ipv"][0] == "ipv4" {
			Cip = v["Cv4o1"][0] + "." + v["Cv4o2"][0] + "." + v["Cv4o3"][0] + "." + v["Cv4o4"][0]
			Sip = v["Sv4o1"][0] + "." + v["Sv4o2"][0] + "." + v["Sv4o3"][0] + "." + v["Sv4o4"][0]
			Ccmd += " " + Sip
			Scmd += " " + Sip
		} else {
			Cip := v["Cv6o1"][0] + ":" + v["Cv6o2"][0] + ":" + v["Cv6o3"][0] + ":" + v["Cv6o4"][0] + ":" + v["Cv6o5"][0] + ":" + v["Cv6o6"][0] + ":" + v["Cv6o7"][0] + ":" + v["Cv6o8"][0]
			Sip := v["Sv6o1"][0] + ":" + v["Sv6o2"][0] + ":" + v["Sv6o3"][0] + ":" + v["Sv6o4"][0] + ":" + v["Sv6o5"][0] + ":" + v["Sv6o6"][0] + ":" + v["Sv6o7"][0] + ":" + v["Sv6o8"][0]
			Ccmd += " " + Cip + " --IPv6Version"
			Scmd += " " + Sip + " --IPv6Version"
		}
	}
	_, ok = v["Cport"]
	if ok && v["Cport"][0] != "" {
		Ccmd += " -p " + v["Cport"][0]
	}
	_, ok = v["Sport"]
	if ok && v["Sport"][0] != "" {
		Scmd += " -p " + v["Sport"][0]
	}

	for key, value := range v {

		if len(value) > 0 {
			v0 = value[0]
			if v0 == "" {
				continue
			}

			switch key {

			case "npthreads":
				Ccmd += (" --parallel " + v0)
				Scmd += (" --parallel " + v0)
				(*m)["npthreads"] = v0
			case "sbsize":
				Ccmd += (" --window " + v0)
				Scmd += (" --window " + v0)
				(*m)["sbsize"] = v0
			case "nodelay":
				Ccmd += " --no-delay"
				Scmd += " --no-delay"
			case "CLport":
				Ccmd += " -L " + v0
			case "mss":
				Ccmd += " -M " + v0
				Scmd += " -M " + v0
			case "trafprot":
				if v0 == "udp" {
					Ccmd += " --udp"
					Scmd += " --udp"

					// other udp-dependent stuff:
					if t, ok := v["sglthread"]; ok && t[0] == "sglthreadon" {
						Scmd += " --single_udp"
					}

					if t, ok := v["srvdaemon"]; ok && t[0] == "srvdaemonon" {
						Scmd += " --daemon"
					}

					if t, ok := v["bwidth"]; ok && t[0] != "" {
						Ccmd += " --bandwidth " + v0
					}

				}
			case "bidiropt":
				if bi, ok := v["bidir"]; ok && bi[0] == "bidiron" {
					if t, ok := v["bidiropt"]; ok && t[0] == "sim" {
						Ccmd += " --dualtest"
					} else if t, ok := v["bidiropt"]; ok && t[0] == "seq" {
						Ccmd += " --tradeoff"
					}
				}
			case "TTtime":
				multiplier := 1
				if t, ok := v["TTtimeUnits"]; ok && t[0] == "tttm" {
					multiplier = 60
				} else if t, ok := v["TTtimeUnits"]; ok && t[0] == "ttth" {
					multiplier = 60 * 60
				}
				num, err := strconv.Atoi(v0)
				if err != nil {
					panic(err)
				}
				num *= multiplier
				Ccmd += " --time " + fmt.Sprintf("%d", num)
			case "ttl":
				Ccmd += " --ttl " + v0
			case "rptformat":
				Ccmd += " --format " + v0
			case "exclusions":
				Ccmd += " --reportexclude " + reportExcludeCode(v["exclusions"])
			case "TotByt":
				Ccmd += " --num " + v0
			case "rptinterval":
				Ccmd += " --interval " + v0
				Scmd += " --interval " + v0

			default:
				//panic(fmt.Sprintf("unhandled option"))
			}
		} else {
			v0 = ""
		}
	}
	fmt.Printf("%#v", v)
	fmt.Printf("Ccmd: '%s'   Scmd:'%s'\n", Ccmd, Scmd)

	iperfCallCount++
	cmd := fmt.Sprintf("ssh %s %s ; ssh %s %s  #(call count: %d)", Sip, Scmd, Cip, Ccmd,
		iperfCallCount)

	(*m)["IperfCmd"] = cmd
	return cmd
}

var iperfCallCount int

func IperfHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request r is %+#v\n", r)

	err := r.ParseForm()
	if err != nil {
		log.Printf("IperfHandler() could not ParseForm on request '%#+v': returned err: %v\n", r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawFormInputMap := (map[string][]string)(r.Form)

	finishedOutputMap := FormToIperfMap(&rawFormInputMap)

	if iperfCallCount == 0 {
		SetFormDefaults(r.Form, finishedOutputMap)

	} else if len(rawFormInputMap) == 0 {
		fmt.Printf("rawFormInput map had length 0, using default values...\n")
		SetFormDefaults(r.Form, finishedOutputMap)
	}

	fmt.Printf("=================================\n")
	fmt.Printf("input rawFormInputMap is %#+v\n", rawFormInputMap)
	fmt.Printf("---------------------------------\n")
	fmt.Printf("output finishedOutputMap is %#+v\n", *finishedOutputMap)
	fmt.Printf("=================================\n")

	cmd := ""
	//	if len(rawFormInputMap) > 0 {
	cmd = GenIperfCmd(r.Form, finishedOutputMap)
	/*
		} else {
			fmt.Printf("empty rawFormInputMap detected... setting defaults\n")
			fmt.Printf("here is the rawFormInputMap: %#+v\n", rawFormInputMap)
		}
	*/
	log.Printf("p.IperfCmd is %v\n", cmd)

	// injects the p.IperfCmd string into the textarea (see
	// templates/iperf.html
	//	renderTemplate(w, "iperf", p)

	fnTmpl := "iperf"
	myT, err := template.ParseFiles("templates/" + fnTmpl + ".html")
	if err != nil {
		log.Printf("renderTemplate could not load file: %s.html", fnTmpl)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = myT.Execute(w, finishedOutputMap)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
