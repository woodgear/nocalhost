/*
Copyright 2020 The Nocalhost Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"github.com/spf13/cobra"

	"nocalhost/internal/nhctl/app"
	"nocalhost/internal/nhctl/nocalhost"
	"nocalhost/pkg/nhctl/log"
)

func init() {
	pvcCleanCmd.Flags().StringVar(&pvcFlags.App, "app", "", "Clean up PVCs of specified application")
	pvcCleanCmd.Flags().StringVar(&pvcFlags.Svc, "svc", "", "Clean up PVCs of specified service")
	pvcCleanCmd.Flags().StringVar(&pvcFlags.Name, "name", "", "Clean up specified PVC")
	pvcCmd.AddCommand(pvcCleanCmd)
}

var pvcCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up PersistVolumeClaims",
	Long:  `Clean up PersistVolumeClaims`,
	Run: func(cmd *cobra.Command, args []string) {
		if pvcFlags.App == "" {
			log.Fatal("--app mush be specified")
		}

		if !nocalhost.CheckIfApplicationExist(pvcFlags.App) {
			log.Fatalf("Application %s not found", pvcFlags.App)
		}
		nhApp, err := app.NewApplication(pvcFlags.App)
		if err != nil {
			log.Fatalf("Failed to create application %s", pvcFlags.App)
		}

		// Clean up specified pvc
		if pvcFlags.Name != "" {
			err = nhApp.CleanUpPVC(pvcFlags.Name)
			if err != nil {
				log.FatalE(err, "Failed to clean up pvc: "+pvcFlags.Name)
			} else {
				log.Infof("%s cleaned up", pvcFlags.Name)
				return
			}
		}

		// Clean up PVCs of specified service
		if pvcFlags.Svc != "" {
			exist, err := nhApp.CheckIfSvcExist(pvcFlags.Svc, app.Deployment)
			if err != nil {
				log.FatalE(err, "failed to check if svc exists")
			} else if !exist {
				log.Fatalf("\"%s\" not found", pvcFlags.Svc)
			}
		}

		err = nhApp.CleanUpPVCs(pvcFlags.Svc, true)
		if err != nil {
			log.FatalE(err, "Cleaning up pvcs failed")
		}
	},
}
