package main

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"github.com/rivo/tview"
	"bufio"
	"io"
	"./cli"
	"./views"
)

func main() {

	containers, err := cli.Cli.ContainerList(context.Background(), types.ContainerListOptions{})


	if err != nil {
		panic(err)
	}

	containersArray := make([]string, 5)

	for i, container := range containers {
		containersArray[i] = container.Image
	}

	app := tview.NewApplication()
	var textView = views.TextView.SetChangedFunc(func() {
		app.Draw()
	})


	app.SetRoot(textView, true)
	dropdown := tview.NewDropDown().
		SetLabel("Select a container (hit Enter): ").
		SetOptions(containersArray,updateLogsOnSelect)
	app.SetRoot(dropdown, true)

	flex := tview.NewFlex().AddItem(textView, 0, 1, true)
	flex.AddItem(dropdown, 0, 1, true);
	app.SetFocus(dropdown)
	app.SetRoot(flex, true).SetFocus(dropdown).Run()

}

/**
*function to update the logs view when it's selected from the dropdown
 */
func updateLogsOnSelect(name string, index int) {
	var textView = views.TextView
	id := getIdByName(name)
	go func() {
		out, err := cli.Cli.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Timestamps: false,
		})

		if err != nil {
			panic(err)
		}
		defer out.Close()

		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			io.Copy(textView, out)
		}
	}()
}

/**
Get a container ID by it's image name
 */
func getIdByName(name string) string {

	containers, err := cli.Cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}


	for _, container := range containers {
		if(container.Image == name){
			return container.ID
		}

	}

	return ""
}




