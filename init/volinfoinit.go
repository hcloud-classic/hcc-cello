package init

import (
	"hcc/cello/lib/handler"
)

func getVolumeDB() error {

	err := handler.ReloadAllofVolInfo()
	if err != nil {
		return err
	}
	return nil
}
