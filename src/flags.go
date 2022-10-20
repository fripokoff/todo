package main

import (
	. "os"
    . "strconv"
)

func flag(arr []data, l *List) []data {
    switch Args[1] {
        case "-a":
            add_flag(l)
            arr = list2array(l)
        case "-d":
            del_flag(l)
            arr = list2array(l)
        default:
            arr = status_flag(arr)
    }
    return arr
}

func add_flag(l *List) {
    tmp := l.head
    l.head = tmp
    i := 0
    for ; tmp.next != nil; tmp = tmp.next {
        i++
    }  
    node := new(Node)
    node.Id = Itoa(i)
    node.Status = 0
    node.Name = Args[2]
    tmp.next = node
}

func del_flag(l *List) {
    i := 0
    taille := len(Args)
    check := 0
    for y := 2; y < taille; y++ { 
        for tmp := l.head; tmp.next != nil; tmp = tmp.next {
            if i == 0 && Args[y] == "0" {
                tmp = tmp.next
                l.head = tmp
                check++
                break
            } else if tmp.next.Id == Args[y] && tmp.next.next == nil {
                check++
                tmp.next = nil 
                break
            } else if tmp.next.Id == Args[y] {
                check++
                tmp.next = tmp.next.next
            } 
            i++
        }
    }
    l.len = l.len - (check + 1)
}

func status_flag(arr []data) []data {
    taille := len(Args)
    for i := 0; i < len(arr); i++ {
        for y := 2; y < taille; y++ {
            if Args[y] == arr[i].Id {
                switch Args[1] {
                case "-b":
                    arr[i].Status = BEGIN
                case "-c":
                    arr[i].Status = CHECK
                case "-e":
                    arr[i].Name = Args[3]
                }
            }
        }
    }
    return arr
}
