package init

import (
	"hcc/cello/lib/handler"
)

func getVolumeDB() error {

	err := handler.ReloadPoolObject()
	if err != nil {
		return err
	}
	// ReloadAllofVolInfo
	err = handler.ReloadAllOfVolInfo()
	if err != nil {
		return err
	}
	return nil
}
