package main

import (
    "fmt"
    "os"
    "encoding/json"
    "strings"
)

const (
	BEGIN = 1
	CHECK = 2
	EDIT  = 3
)

func help() string {
	return `
        todo [Flags] [Args]
    `
}

func is_flag() bool {
	if len(os.Args) < 3 {
		return false
	}
	switch {
	case os.Args[1] == "-a" && len(os.Args) <= 4:
		return true
	case os.Args[1] == "-b":
		return true
	case os.Args[1] == "-c":
		return true
	case os.Args[1] == "-d":
		return true
	case os.Args[1] == "-e" && len(os.Args) == 4:
		return true
	}
	return false
}

func is_proj_syntax(proj string) bool {
    if proj[0] != '@' {
        fmt.Printf("Wrong syntax ! try --help\n") 
        return false
    }
    return true
}

func is_proj_name(projn [][]byte, projA string) bool {
    exist := false 
    for y := 0; y < len(projn); y++ {
        if strstr(projA, string(projn[y])) == true {
            exist = true
            return true
        }
    }
    if !exist {
        fmt.Printf("the project doesn't exist !\n")
    }
    return false
}

func index_end_proj(arr [][]byte, begin int) int {
    for y := begin; y < len(arr); y++ {
        for x := 0; x < len(arr[y]); x++ {
            if arr[y][x] == ']' {
                return y
            }
        }
    }
    return 0
}

func index_begin_proj(arr [][]byte, projA string) int {
    begin := 0
    infos := infos_proj(arr)
    for i := 1; i < len(infos); i++ {
        if strstr(projA, string(arr[infos[i]])) {
            begin = infos[i] + 1
        }
    }
    return begin
} 

func size_proj_select(begin int, end int, arr [][]byte) int {
    size := 0

    for y := begin; y < end; y++ {
        for x := 0; x < len(arr[y]); x++ {
            size++
        } 
    } 
    return size
}

func select_proj(size int, arr [][]byte, begin int, end int) []byte {
    diff := end - begin
    p := make([]byte, size + diff + 3)
    w := 2

    p[0] = '['
    p[1] = '\n'
    for y := begin; y < end; y++ {
        for x := 0; x < len(arr[y]); x++ {
            p[w] = arr[y][x]
            w++
        }
        p[w] = '\n'
        w++
    }
    p[w] = ']'
    return p
}

func board_proj() (int, int, []byte, [][]byte) {
    json_file := get_file()
    arr := str_to_arr(json_file) 
    projn := proj_names(arr)
    projA := ""
    if len(os.Args) == 2 {
        projA = os.Args[1]
    } else {
        projA = os.Args[(len(os.Args)-1)]
    }
    if !is_proj_syntax(projA) {
        os.Exit(1)
    } 
    projA = projA[1:]
    if !is_proj_name(projn, projA) {
        os.Exit(1)
    }
    begin := index_begin_proj(arr, projA)
    end := index_end_proj(arr, begin)
    size := size_proj_select(begin, end, arr)
    p := select_proj(size, arr, begin, end)
    return begin, end, p, arr
}

func main() {
	if len(os.Args) == 1 {
		json_file := get_file()
        gest_proj(json_file)
		return
	} else if len(os.Args) == 2 {
        _, _, p, _ := board_proj()
        tab := json_to_array(p)
        print_tasks(tab)
    } else if len(os.Args) == 3 {
        switch (os.Args[1]) {
            case "-a" :
                json_file := get_file()
                a := str_to_arr(json_file) 
                projname := os.Args[2]
                task := "{\n\"Id\": \"0\",\n\"Status\": 0,\n\"Name\": \"task\"\n}\n],\n"
                nproj := strconc([]byte("{\n\""), []byte(projname))
                nproj = strconc(nproj, []byte("\""))
                nproj = strconc(nproj, []byte(": [\n"))
                nproj = strconc(nproj, []byte(task))
                nproj = strconc(nproj, arr2str(a,1,len(a)))
                home := os.Getenv("HOME")
                os.WriteFile((home + "/.todo.json"), nproj, 0777)
            case "-d" :
                b, e, _, a := board_proj()
                projn := os.Args[2]
                projn = projn[1:]
                c := strconc(arr2str(a,0,b), arr2str(a,e,len(a)))
                arr := str_to_arr(c)
                for i := 0; i < len(arr); i++ {
                    if strstr(projn, string(arr[i])) {
                        arr[i] = []byte(" ")
                        arr[i+1] = []byte(" ")
                        break;
                    }
                }
                c = arr2str(arr,0,len(arr))
                home := os.Getenv("HOME")
                os.WriteFile((home + "/.todo.json"), c, 0777)

        }
    } else if len(os.Args) >= 4 {
        if os.Args[2] == "all" {
            exclude := []string{"out", "a","h"}
            excludeDir := false
            excludeExt := true
            path, err := os.Getwd()
            if err != nil {
                
                fmt.Printf("Error obtaining current directory :", err)
            }
            files, err := os.ReadDir(path)
            if err != nil {
                fmt.Printf("Error reading directory :", err)
            }
            if (os.Args[3] == "-ex"){
            exclude = strings.Split(os.Args[4], ",")
            }else if (os.Args[3] == "-d") {
                excludeDir = true
            } else if (os.Args[3] == "-e") {
                excludeExt = false
            }
            if(os.Args[4] == "-ex" ){
                exclude = strings.Split(os.Args[5], ",")
            }else if (os.Args[4] == "-d") {
                excludeDir = true
            }  else if (os.Args[4] == "-e") {
                excludeExt = false
            }
            if(os.Args[5] == "-ex" ){
                exclude = strings.Split(os.Args[6], ",")
            }else if (os.Args[5] == "-d") {
                excludeDir = true
            } else if (os.Args[5] == "-e") {
                excludeExt = false
            }
            for _, file := range files {
                isExcluded := false
                filename := file.Name()
                for _, ext := range exclude {
                    if strings.HasSuffix(filename, ext) {
                        isExcluded = true
                    }
                }
                if isExcluded || excludeDir && file.IsDir() {
                    continue
                }else{
                    index := strings.LastIndex(filename, ".")
                    if index >= 0 {
                        filenames := filename[0:index]
                        if excludeExt {
                        filename = filenames
                        }
                        if filename == "" || string(filename[0]) == "." {
                            continue
                         }
                    }
                    b, e, p, a := board_proj() 
                    tab := json_to_array(p)
                    l := List{}
                    array2list(tab, &l)
                    t := flag(tab, &l, filename)
                    test, _ := json.MarshalIndent(t, "", " ")
                    c := strconc(arr2str(a,0,b), test[1:len(test)-1])
                    c = strconc(c, arr2str(a,e,len(a)))
                    home := os.Getenv("HOME")
                    os.WriteFile((home + "/.todo.json"), c, 0777)
                }
            }
        }else{
        b, e, p, a := board_proj() 
        tab := json_to_array(p)
        l := List{}
        array2list(tab, &l)
        t := flag(tab, &l, os.Args[2])
        test, _ := json.MarshalIndent(t, "", " ")
        c := strconc(arr2str(a,0,b), test[1:len(test)-1])
        c = strconc(c, arr2str(a,e,len(a)))
        home := os.Getenv("HOME")
        os.WriteFile((home + "/.todo.json"), c, 0777)
        }
    } else if !is_flag() {
        fmt.Println(help())
        os.Exit(1)
    } 
}