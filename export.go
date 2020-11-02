package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"process_exporter/psutil"
	"strconv"
	"strings"
	"text/template"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type export struct {
	config *Config
	descs  map[string]*prometheus.Desc
}

type Config struct {
	Global struct {
		Port string `yaml:"port"`
		Path string `yaml:"path"`
	} `yaml:"global"`
	Metrics []Metric `yaml:"metrics"`
}

type Metric struct {
	Type           string               `yaml:"type"`
	Filter         string               `yaml:"filter"`
	Name           string               `yaml:"name"`
	Help           string               `yaml:"help"`
	Value          string               `yaml:"value"`
	ValueType      prometheus.ValueType `yaml:"value_type"`
	VariableLabels []LabelValue         `yaml:"variable_labels"`
	ConstLabels    []LabelValue         `yaml:"const_labels"`
}

type LabelValue struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

type Info struct {
	Processes *psutil.Processes

	Process *psutil.Process
}

func newExport(path string) *export {
	// 记载配置文件
	cfg := new(Config)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	if err := yaml.Unmarshal(buf, cfg); err != nil {
		log.Panicln(err)
	}

	ept := new(export)
	ept.config = cfg
	ept.descs = make(map[string]*prometheus.Desc, len(cfg.Metrics))

	// 从配置生成 export metric
	// 遍历配置文件的每个 Metrics
	for _, metric := range cfg.Metrics {
		// 可变标签的生成
		variableLabels := make([]string, 0)
		for _, v := range metric.VariableLabels {
			variableLabels = append(variableLabels, v.Label)
		}
		// 固定标签的生成
		constLabels := make(map[string]string, 0)
		for _, c := range metric.ConstLabels {
			// 固定标签的值的生成
			// 可调用的变量
			info := new(Info)
			info.Processes = psutil.NewProcesses()
			info.Process = new(psutil.Process)

			constLabels[c.Label] = execute(c.Value, info)
		}

		// metric 的生成
		ept.descs[metric.Name] = prometheus.NewDesc(
			metric.Name,
			metric.Help,
			variableLabels,
			constLabels,
		)
	}
	return ept
}

func (e *export) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range e.descs {
		ch <- desc
	}
}

func (e *export) Collect(ch chan<- prometheus.Metric) {
	processes := psutil.NewProcesses()

	// info 用来存放供配置文件调用的变量
	info := new(Info)

	for _, metric := range e.config.Metrics {
		if metric.Type == "each" {
			processes.EachProcess(func(p *psutil.Process) bool {
				info.Processes = processes
				info.Process = p

				if execute(metric.Filter, info) != "true" {
					return true
				}

				putMetric(e, info, metric, ch)
				return true
			})
		} else if metric.Type == "end" {
			pse := psutil.NewEmptyProcesses()

			processes.EachProcess(func(p *psutil.Process) bool {
				info.Processes = processes
				info.Process = p

				if execute(metric.Filter, info) == "true" {
					pse.AppendProcess(p)
					return true
				}

				return true
			})

			info.Process = new(psutil.Process)
			info.Processes = pse
			putMetric(e, info, metric, ch)
		}
	}
}

//执行文本代码
func execute(text string, itf interface{}) string {
	if tpl, err := template.New("").Parse(text); err != nil {
		log.Panicln(err)
	} else {
		buf := new(bytes.Buffer)
		if err := tpl.Execute(buf, itf); err != nil {
			return ""
		}
		return buf.String()
	}

	return ""
}

// 从配置产生指标, 注入到 export
func putMetric(e *export, itf interface{}, metric Metric, ch chan<- prometheus.Metric) {
	// 从配置产生指标值
	t := strings.TrimSpace(execute(metric.Value, itf))
	var value float64
	var err error
	if t == "" {
		value = 0
	} else {
		value, err = strconv.ParseFloat(t, 64)
		if err != nil {
			log.Println(err)
		}
	}

	// 从配置产生指标标签和标签值
	values := make([]string, 0)
	for _, v := range metric.VariableLabels {
		values = append(values, execute(v.Value, itf))
	}

	//注入
	ch <- prometheus.MustNewConstMetric(
		e.descs[metric.Name],
		metric.ValueType,
		value,
		values...,
	)
}
