package subsystems

import(
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	println("MemorySubSystem Set " + cgroupPath + res.MemoryLimit)
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err != nil{
		return err
	} else {
		if res.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit),0644); err != nil {
				return fmt.Errorf("set cgroup memory failed %v", err)
			}
		}
		return nil
	}
}

func (s *MemorySubSystem) Remove(cgroupPath string) error {
        if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(subsysCgroupPath)
	} else { 
		return err
	}
}

func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	println("MemorySubSystem Apply " + cgroupPath + strconv.Itoa(pid))
        if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err != nil{
                return err
        } else {
                        if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)),0644); err != nil {
                                return fmt.Errorf("set cgroup memory failed %v", err)
                        }
                	return nil
        }
}

func (s *MemorySubSystem) Name() string {
	return "memory"
}

