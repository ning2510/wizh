#!/bin/bash

Green="\033[32m"
Red="\033[31m"
GreenBG="\033[42;37m"
RedBG="\033[41;37m"
Font="\033[0m"

OK="${Green}[OK]${Font}"
Error="${Red}[错误]${Font}"

check_screen() {
    if ! command -v screen &> /dev/null; then
        echo -e "${Error} ${RedBG} 未安装 screen！${Font}"
        exit 1
    fi
}

check_screen

start_screen() {
    cmd="cd cmd/$1 && go run main.go\n"
    screen -S $1 -X quit
    screen -dmS $1
    screen -S $1 -p 0 -X stuff "$cmd"
    screen -ls
    echo -e "${OK} ${GreenBG} $1 服务已启动！ ${Font}"
}

stop_screen() {
    screen -S $1 -X quit
    screen -ls
    echo -e "${OK} ${GreenBG} $1 服务已暂停！ ${Font}"
}

if [ $# -gt 0 ]; then
    services=("comment" "favorite" "message" "relation" "user" "video" "api")

    if [[ "$1" == "init" ]] || [[ "$1" == "install" ]]; then
        go mod tidy
    elif [[ "$1" == "start" ]]; then
        shift 1
        if [[ "$1" == "" ]]; then
            for service in "${services[@]}"; do
                start_screen "$service"
            done
        else
            start_screen "$1"
        fi
    elif [[ "$1" == "stop" ]]; then
        shift 1
        if [[ "$1" == "" ]]; then
            for service in "${services[@]}"; do
                stop_screen "$service"
            done
        else
            stop_screen "$1"
        fi
    elif [[ "$1" == "info" ]]; then
        shift 1
        screen -r $1
    elif [[ "$1" == "ls" ]]; then
        screen -ls
    fi
fi