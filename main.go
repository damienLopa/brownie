package main

import (
    "fmt"
    "os"
    "strings"
    "encoding/json"
    "io/ioutil"
    "path/filepath"
    "text/template"
    "reflect"
    "github.com/manifoldco/promptui"
    "github.com/mkideal/cli"
    "github.com/go-git/go-git"
    "github.com/gosimple/slug"
)

type argT struct {
	cli.Helper
	Template string  `cli:"t,template" usage:"short and long format flags both are supported"`
}


func getProjectName()string {
    return promptUi("What's the name of your new project ?")
}

func selectUi(items []string, label string)string {  
        index := -1
        var result string
        var err error

        for index < 0 {
            prompt := promptui.SelectWithAdd{
                Label:    label,
                Items:    items,
            }
            index, result, err = prompt.Run()
            if index == -1 {
                items = append(items, result)
            }
        }
        if err != nil {
            fmt.Printf("Prompt failed %v\n", err)
            return ""
        }

    return result
}

func promptUi(label string)string { 
    var result string
    var err error
    
    prompt := promptui.Prompt{
        Label: label,
    }
    result, err = prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        return ""
    }
    return result
}

func cloneRepo(url string, pathTemplate string) {
    
    _ = os.Mkdir(os.Getenv("HOME") + "/.brownie", 0755)

    _, err := git.PlainClone(pathTemplate, false, &git.CloneOptions{
        URL:      url,
    })
    if err != nil {
        panic(err)
    }
}

func selectTemplate(onlineTemplate string, projectName string)string {
    pathTemplate := os.Getenv("HOME") + "/.brownie/" + slug.Make(projectName)

    if onlineTemplate != "" {
        cloneRepo(onlineTemplate, pathTemplate)
    }else {
        items := []string{"Django", "Flask"}
        label := "What's your favorite framework ?"
        selectedTemplate := selectUi(items, label)
        onlineTemplate := "https://github.com/damienLopa/" + slug.Make(selectedTemplate) + "-brownie-template"
        cloneRepo(onlineTemplate, pathTemplate)
    }
    return pathTemplate
}

func selectProjectName()string {
    items := []string{"Django", "Flask"}
    label := "What's your favorite framework ?"

    return selectUi(items, label)
}

func generateProject(project_name string,
                     path_template string,
                     jsonData map[string]interface{})error {
    var err error

    os.Mkdir(project_name, 0755)

    err = filepath.Walk(path_template,

    func(path string, info os.FileInfo, err error)error {
        if err != nil {
            return err
        }
        if info.IsDir(){
            a := strings.ReplaceAll(path, path_template, project_name)
            os.Mkdir(a, 0755)
        }else{
            a := strings.ReplaceAll(path, path_template, project_name)
            tmpl, err := template.ParseFiles(path)
            if err != nil {
	        	return err
            }
            f, err := os.Create(a)
            if err != nil {
	        	return err
            }
            err = tmpl.Execute(f, &jsonData)
            if err != nil {
	        	return err
            }
        }
        return err  
    })
    return err
}


func getSelectedOption(jsonData map[string]interface{})map[string]interface {}{
    selectedOption := make(map[string]interface{})

    for key, value := range jsonData {
        t := reflect.TypeOf(value)
        switch t.Kind(){
            case reflect.Slice:
                optionInterface := reflect.ValueOf(value)
                optionList := make([]string, optionInterface.Len())
                for i := 0; i < optionInterface.Len(); i++ {
                    optionList[i] = fmt.Sprint(optionInterface.Index(i))
                }
                selectedOption[key] = selectUi(optionList, key)
            default:
                selectedOption[key] = promptUi(key)
            }
    }
    return selectedOption
}

func getJsonData(template string)map[string]interface {} {
    var err error
    jsonData := make(map[string]interface{})

    jsonFile, err := os.Open("templates/" + template + "/.brownie.json")
    if err != nil {
        return jsonData
    }

    byteValueFile, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        panic(err)
    }

    err = json.Unmarshal([]byte(byteValueFile), &jsonData);
    if err != nil {
		panic(err)
    }

    return jsonData
}

func main() {
    
    os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
        argv := ctx.Argv().(*argT)
        
        projectName := getProjectName()
        path_template := selectTemplate(argv.Template, projectName)

        jsonData := getJsonData(argv.Template)

        selectedOption := getSelectedOption(jsonData)

        err := generateProject(projectName, path_template, selectedOption)
        fmt.Printf("Someting went wrong %v\n", err)
        return nil
    }))
}