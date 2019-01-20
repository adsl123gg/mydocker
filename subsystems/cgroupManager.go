package subsystems

import("fmt")

type CgroupManager struct {
	Path string
	Resource *ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Set(res *ResourceConfig) error {
	for _, subSysIns := range(SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range(SubsystemsIns) {
                subSysIns.Apply(c.Path, pid)
        }
        return nil
}

func (c *CgroupManager) Destory() error {
	for _, subSysIns := range(SubsystemsIns) {
        	if err := subSysIns.Remove(c.Path); err != nil {
			return fmt.Errorf("Destory cgroup failed %v", err)
		}
	}
        return nil
}


