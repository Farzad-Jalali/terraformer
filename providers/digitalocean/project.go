// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package digitalocean

import (
	"context"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/digitalocean/godo"
)

type ProjectGenerator struct {
	DigitalOceanService
}

func (g ProjectGenerator) listProjects(ctx context.Context, client *godo.Client) ([]godo.Project, error) {
	list := []godo.Project{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		projects, resp, err := client.Projects.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			list = append(list, project)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.Pages == nil || resp.Links.Pages.Next == "" {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

func (g ProjectGenerator) createResources(projectList []godo.Project) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, project := range projectList {
		resources = append(resources, terraform_utils.NewSimpleResource(
			project.ID,
			project.Name,
			"digitalocean_project",
			"digitalocean",
			[]string{}))
	}
	return resources
}

func (g *ProjectGenerator) InitResources() error {
	client := g.generateClient()
	output, err := g.listProjects(context.TODO(), client)
	if err != nil {
		return err
	}
	g.Resources = g.createResources(output)
	return nil
}
