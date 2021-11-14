// THIS CODE HAS BEEN GENERATED, DO NOT TOUCH THIS
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

const url = "{{.Url}}"
const (
	statusOK           = 0
	statusTransportErr = 1
	statusRpcErr       = 2
	statusInputErr     = 3
)

func main(){
    {{- range $idx, $methodDef := .Def.Methods }}
    flags{{$idx}} := flag.NewFlagSet("{{.Name}}", flag.ContinueOnError)
    {{ range $name, $type := $methodDef.Inputs}}
    {{- if eq $type "integer"}}{{$name}}{{$idx}} := flags{{$idx}}.Int("{{$name}}", 0, "-{{$name}}=42"){{end}}
    {{- if eq $type "number"}}{{$name}}{{$idx}}  := flags{{$idx}}.Float64("{{$name}}", 0.0, "-{{$name}}=4.2"){{end}}
    {{- if eq $type "string" }}{{$name}}{{$idx}} := flags{{$idx}}.String("{{$name}}", "", "-{{$name}}=fourty two"){{end}}
    {{end}}
    {{end}}

    if len(os.Args) < 2 {
        errSubCommandMissing()
    }

    switch os.Args[1] {
    {{range $idx, $methodDef := .Def.Methods }}
    case "{{$methodDef.Name}}":
        err := flags{{$idx}}.Parse(os.Args[2:])
        if err!=nil {
            if errors.Is(err,flag.ErrHelp){
                flags0.PrintDefaults()
                os.Exit(statusInputErr)
            }
            stdErrf("bad input: %v", err)
            os.Exit(statusInputErr)
        }
        post("{{$methodDef.Name}}",struct{
        {{range $name, $type := $methodDef.Inputs}}
        P{{$name}} {{GoType $type}} `json:"{{$name}}"`
        {{- end}}
        }{
            {{- range $name, $type := $methodDef.Inputs}}
            P{{$name}} : *{{$name}}{{$idx}},
            {{- end}}
        })
    {{- end}}
    	default:
    		errSubCommandMissing()
    }
}

func errSubCommandMissing() {
	fmt.Println("expected subcommand, one of: split add ")
	os.Exit(statusInputErr)
}

func post(method string, data interface{}) {
	type(
        base struct {
            Version string `json:"jsonrpc"`
            Id      int    `json:"id"`
        }
        req struct {
            base
            Method string `json:"method"`
            Params interface{} `json:"params"`
        }
        result struct {
            Value interface{} `json:"value"`
        }
        rpcErr struct{
            Code    int    `json:"code"`
            Message string `json:"message"`
            Data    string `json:"data"`
        }
        res struct {
            base
            Error rpcErr `json:"error"`
            Result  result `json:"result"`
        }
    )

	jsonData, err := json.Marshal(req{
		base: base{
			Version: "2.0",
			Id:      1,
		},
		Method: method,
		Params: data,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(statusTransportErr)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		stdErrf("TransportErr HTTP-POST: %v", err)
		os.Exit(statusTransportErr)
	}

	if resp.StatusCode != 200 {
		var err = errors.New("unexpected status code")
		stdErrf("TransportErr HTTP-STATUS %v: %s", resp.StatusCode, err)
		os.Exit(statusTransportErr)
	}

	var response res
	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		stdErrf("TransportErr HTTP-STATUS %v: %v", resp.StatusCode, err)
		os.Exit(statusTransportErr)
	}

	err = json.Unmarshal(resBytes, &response)
	if err != nil {
		stdErrf("TransportErr HTTP-STATUS %v: %v", resp.StatusCode, err)
		os.Exit(statusTransportErr)
	}

	if response.Error.Code != 0 {
		stdErrf("RpcErr %v: %v: %v", response.Error.Code, response.Error.Message, response.Error.Data)
		os.Exit(statusRpcErr)
	}

	if response.Result.Value == nil {
		stdErrf("RpcErr 0:  empty result (response bytes: %s)", string(resBytes))
		os.Exit(statusRpcErr)
	} else {
		outputValue(response.Result.Value)
	}
}

func outputValue(value interface{}){
	var v = reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			stdout(v.Index(i).Interface())
		}
	default:
		stdout(value)
	}

	os.Exit(statusOK)
}

func stdErrf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format+"\n", a...)
	if err != nil {
		panic(fmt.Sprintf("cannot write to stderr: %v", err))
	}
}

func stdout(msg interface{}) {
	_, err := fmt.Fprintln(os.Stdout, msg)
	if err != nil {
		panic(fmt.Sprintf("cannot write to stdout: %v", err))
	}
}