package main

import (
	"fmt"
	"github.com/vaberof/yadro-test-task/internal/app/entrypoint/event/eventhandler"
	"github.com/vaberof/yadro-test-task/internal/app/entrypoint/file/filehandler"
	"github.com/vaberof/yadro-test-task/internal/domain/computerclub"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("filename must be provided as second argument")
	}

	filename := os.Args[1]

	computerClubConfig := computerclub.Config{}

	invalidLine, err := filehandler.ProcessComputerClubConfig(filename, &computerClubConfig)
	if err != nil {
		if invalidLine != nil {
			fmt.Println(*invalidLine)
		} else {
			panic(err.Error())
		}
		return
	}

	computerClubService := computerclub.NewComputerClub(&computerClubConfig)

	eventHandler := eventhandler.NewHandler(computerClubService)

	fileHandler := filehandler.NewHandler(eventHandler)

	workingDayReport, invalidLine, err := fileHandler.GetWorkingDayReport(filename, computerClubConfig.TablesCount)
	if err != nil {
		if invalidLine != nil {
			fmt.Println(*invalidLine)
		} else {
			panic(err.Error())
		}
		return
	}

	fmt.Print(workingDayReport)
}
