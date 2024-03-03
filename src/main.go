package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/* Constants */
const DTD string = "<!DOCTYPE service_bundle SYSTEM '/usr/share/lib/xml/dtd/service_bundle.dtd.1'>"

/* Boolean settings */
var autoinstall bool = false
var custom_name bool = false
var quiet_mode bool = false
var validate_on_exit bool = false
var version float32 = 1.1

/* Main */
func main() {

	/*
		CMD Args:
		-f [WIP](filepath): allow for input file full of name/value pairs to be submitted instead of command line args
		-i : automatically install generated file
		-o (string) : specify output
		-q : quiet mode, suppressed output
		-s (name=value) : specify a name / value pair; needs a list of acceptable values pre-made
		-V : validate on exit
		version : print version and exit
	*/

	//Initialize map for -s name/value pairs
	var s_args map[string]string = make(map[string]string)
	s_args["start-method"] = ":true"
	s_args["stop-method"] = ":true"
	s_args["restart-method"] = ":true"
	s_args["service-name"] = "MyService"
	s_args["service-description"] = "MyService Description."
	s_args["timeout-seconds"] = "60"

	//Convert os.Args to list, exclude binary name
	var args map[int]string = make(map[int]string)
	for index, value := range os.Args[1:] {
		args[index] = value
	}

	//Set the file name
	var output_file string = fmt.Sprintf("./%s.xml", s_args["service-name"])

	//Determine boolean args
	for i := 0; i <= len(args); i++ {
		switch args[i] {
		case "-q":
			quiet_mode = true
			delete(args, i)
			break
		case "-V":
			validate_on_exit = true
			delete(args, i)
			break
		case "-i":
			autoinstall = true
			delete(args, i)
			break
		case "-o":
			custom_name = true
			output_file = fmt.Sprintf("./%s.xml", args[i+1])
			delete(args, i)
			break
		case "version":
			fmt.Printf("svcbundle version %.2f %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
			break
		}
	}

	//Check for exclusive args
	if autoinstall && custom_name {
		fmt.Println("error: '-i' and '-o' are exclusive arguments!")
		os.Exit(2)
	}

	//Iterate over and store -s name/value pairs
	for i := 0; i < len(args); i++ {
		//Ensure first arg is "-s" and next argument exists
		if args[i] == "-s" && len(args) >= i+1 && strings.Contains(args[i+1], "=") {
			//Split args
			tmp := strings.Split(args[i+1], "=")
			//Ensure arg is actually in premade map
			if _, ok := s_args[tmp[0]]; !ok {
				fmt.Println("BAD KEY: " + tmp[1])
				os.Exit(1)
			} else {
				//update args if so
				s_args[tmp[0]] = tmp[1]
				//Update output file if new service name is provided, will be overwritten if "-o" argument is present
				if tmp[0] == "service-name" && !custom_name {
					output_file = fmt.Sprintf("./%s.xml", s_args["service-name"])
				}
				//Determine if quiet mode is enabled or not
				if !quiet_mode {
					fmt.Printf("Setting '%s' to '%s'.\n", tmp[0], tmp[1])
				}
				continue
			}
		}
	}

	//Update filepath if autoinstall is true
	if autoinstall {
		output_file = "/lib/svc/manifest/system/" + output_file
		if !quiet_mode {
			fmt.Println("Set program to automatically install manifest after completion.")
		}
	}

	//Create service_bundle instance and add a service
	var svcbundle *service_bundle = &service_bundle{Name: s_args["service-name"], Type: "manifest"}
	svcbundle.Service = add_service(s_args)

	//Marshall XML
	out, err := xml.MarshalIndent(svcbundle, " ", "  ")
	err_check(err)

	//Open a file for writing
	file, err := os.Create(output_file)
	err_check(err)
	defer file.Close()

	//Add loctext content, write to file, return absolute path
	var tmpString string = strings.Replace(string(out), `<loctext xml:lang="C"></loctext>`, `<loctext xml:lang="C">`+s_args["service-name"]+`</loctext>`, 1)
	_, err = file.WriteString(xml.Header + DTD + generate_timestamp() + strings.Replace(tmpString, `<loctext xml:lang="C"></loctext>`, `<loctext xml:lang="C">`+s_args["service-description"]+`</loctext>`, 1))
	err_check(err)
	abspath, err := filepath.Abs(output_file)
	err_check(err)
	if !quiet_mode {
		fmt.Println("Wrote file to : " + abspath)
	}

	//Validate on exit if -V flag supplied
	if validate_on_exit {
		validate(abspath)
	}

}

/* Function to check errors */
func err_check(e error) {
	if e != nil {
		panic(e)
	}
}
