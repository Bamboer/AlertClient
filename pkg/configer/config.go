package configer

import (
	"fmt"
	"grafana/pkg/log"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
        "bufio"
)

var (
	DefaultLevel = 1
)

type Config struct {
	filename       string
	data           map[string]string
	lastModifyTime int64
	rwLock         sync.RWMutex
	notifyList     []Notifyer
}

type Notifyer interface {
	Callback(*Config)
}

func (c *Config) AddObserver(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}

func NewConfig(file string) (conf *Config, err error) {
	conf = &Config{
		filename: file,
		data:     make(map[string]string, 1024),
	}
	m, err := conf.parse()
	if err != nil {
		log.Errorln("parse conf error: ", err)
	}
	conf.rwLock.Lock()
	conf.data = m
	conf.rwLock.Unlock()
	go conf.reload()
	return
}

func (c *Config) parse() (m map[string]string, err error) {
	m = make(map[string]string, 1024)
	f, err := os.Open(c.filename)
	if err != nil {
		log.Errorln("Failed opend file: ", err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var lineNo int
	for {
		line, errRet := reader.ReadString('\n')
		if errRet == io.EOF {
			lineParse(&lineNo, &line, &m)
			break
		}
		if errRet != nil {
			err = errRet
			return
		}
		lineParse(&lineNo, &line, &m)
	}
	return
}

func lineParse(lineNo *int, line *string, m *map[string]string) {
	*lineNo++
	l := strings.TrimSpace(*line)
	if len(l) == 0 || l[0] == '\n' || l[0] == '#' || l[0] == ';' {
		return
	}
	itemSlice := strings.Split(l, "=")
	if len(itemSlice) == 0 {
		log.Infoln("invalid config, line: ", lineNo)
		return
	}
	key := strings.TrimSpace(itemSlice[0])
	if len(key) == 0 {
		log.Infoln("invalid config, line: ", lineNo)
		return
	} else if len(key) == 1 {
		(*m)[key] = ""
		return
	}
	value := strings.TrimSpace(itemSlice[1])
	(*m)[key] = value
	return
}

func (c *Config) GetInt(key string) (value int, err error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	str, ok := c.data[key]
	if !ok {
		err = fmt.Errorf("key [%s] not found", key)
	}
	value, err = strconv.Atoi(str)
	return
}

func (c *Config) GetIntDefault(key string, defaultInt int) (value int) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	str, ok := c.data[key]
	if !ok {
		value = defaultInt
		return
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		value = defaultInt
	}
	return
}

func (c *Config) GetString(key string) (value string, err error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	value, ok := c.data[key]
	if !ok {
		err = fmt.Errorf("key [%s] not found", key)
	}
	return
}

func (c *Config) GetStringDefault(key string, defaultStr string) (value string) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	value, ok := c.data[key]
	if !ok {
		value = defaultStr
		return
	}
	return
}

func (c *Config) reload() {
	t := time.Tick(30 * time.Second)
	for _ = range t {
		func() {
			f, err := os.Open(c.filename)
			if err != nil {
				log.Infoln("open file error: ", err)
				return
			}
			defer f.Close()

			fileinfo, err := f.Stat()
			if err != nil {
				log.Infoln("stat file error: ", err)
				return
			}
			curModifyTime := fileinfo.ModTime().Unix()
			if curModifyTime > c.lastModifyTime {
				m, err := c.parse()
				if err != nil {
					log.Infoln("parse config error: ", err)
					return
				}
				c.rwLock.Lock()
				c.data = m
				c.rwLock.Unlock()
				c.lastModifyTime = curModifyTime

				for _, n := range c.notifyList {
					n.Callback(c)
				}
			}
		}()
	}
}
