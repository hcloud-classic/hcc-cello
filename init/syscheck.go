package init

import "hcc/cello/lib/syscheck"

func syscheckInit() error {
	err := syscheck.CheckRoot()
	if err != nil {
		return err
	}

	return nil
}
