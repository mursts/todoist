package main

import (
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	projectsCount := len(client.Store.Projects)
	projectIds := make([]int, projectsCount, projectsCount)
	for i, project := range client.Store.Projects {
		projectIds[i] = project.GetID()
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)
	ex := Filter(c.String("filter"))

	itemList := make([][]string, 0, len(client.Store.ItemOrders))
	for _, itemOrder := range client.Store.ItemOrders {
		item := itemOrder.Data.(todoist.Item)
		r, err := Eval(ex, item)
		if err != nil {
			return err
		}
		if !r || item.Checked == 1 {
			continue
		}
		itemList = append(itemList, []string{
			IdFormat(item),
			PriorityFormat(item.Priority),
			DueDateFormat(item.DateTime(), item.AllDay),
			ProjectFormat(item.ProjectID, client.Store.Projects, projectColorHash, c),
			item.LabelsString(client.Store.Labels),
			ContentPrefix(client.Store.Items, item, c) + ContentFormat(item),
		})
	}

	defer writer.Flush()

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
