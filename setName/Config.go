/*usage
* var cfg config
* cfg.Load("sample.cfg") //load config fields from a file
* cfg.Get("keyName") //get key value
* cfg.Set("keyName","keyValue") //set key value
* cfg.Remove("keyName") //remove a config field
* cfg.Write() //write config fields to file
*
* config file type
* key0=value0
* key1=1
* key2=true
*/
package main

import (
	"io/ioutil"
	"strings"
)

type config struct {
	file   string
	fields map[string]string
}

func (c *config) Load(fileName string) error {
	c.file = fileName
	c.fields = make(map[string]string)
	//读取配置字段
	fileByte, err := ioutil.ReadFile(c.file)

	if err != nil {
		return err
	}
	lines := strings.Split(string(fileByte), "\n")
	for _, line := range lines {
		field := strings.Split(line, "=")
		if len(field) >= 2 {
			c.fields[field[0]] = field[1]
		}
	}
	return nil
}
func (c *config) Set(key string, value string) {
	c.fields[key] = value
}
func (c *config) Get(key string) (string, bool) {
	value, ok := c.fields[key]
	return value, ok
}
func (c *config) Remove(key string) bool {
	_, ok := c.fields[key]
	if !ok {
		return false
	}
	delete(c.fields, key)
	return true
}
func (c *config) Write() error {
	text := ""
	for key, value := range c.fields {
		text += key + "=" + value + "\n"
	}
	return ioutil.WriteFile(c.file, []byte(text), 0777)
}
