package init

import (
	"hcc/cello/lib/handler"
)

func getVolumeDB() error {

	err := handler.ReloadPoolObject()
	if err != nil {
		return err
	}

	err = handler.ReloadAllofVolInfo()
	if err != nil {
		return err
	}
	return nil
}
