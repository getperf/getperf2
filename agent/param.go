package agent

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ParamType int

const (
	INT ParamType = iota
	BOOL
	STRING
)

type ParamMode int

const (
	BASE ParamMode = iota
	COLLECTOR
	JOB
)

func paramLineTrim(line string) string {
	matches := regexp.MustCompile("^(\\s*);").FindStringSubmatch(line)
	if len(matches) == 2 {
		return ""
	}
	return strings.TrimSpace(line)
}

func searchParam(line string, matchWord string, paramMode ParamMode) (string, string) {
	var search string
	switch paramMode {
	case BASE:
		search = fmt.Sprintf("%s\\s*=\\s*(.+)", matchWord)
	case COLLECTOR:
		search = fmt.Sprintf("%s\\.(\\w+)\\s*=\\s*(.+)", matchWord)
	case JOB:
		search = fmt.Sprintf("%s\\.(\\w+)\\s*=\\s*(.+)", matchWord)
	}

	matches := regexp.MustCompile(search).FindStringSubmatch(line)
	if len(matches) > 0 {
		if paramMode == BASE {
			return matches[1], ""
		} else {
			return matches[2], matches[1]
		}
	} else {
		return "", ""
	}
}

func convertInt(word string, line string) (int, error) {
	valueInt, err := strconv.Atoi(word)
	if err != nil {
		log.Error("parse error : ", line, ",invalid number : ", word)
	}
	return valueInt, err
}

func convertBool(word string, line string) (bool, error) {
	switch word {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		log.Error("parse error : ", line, ",invalid bool : ", word)
		return false, fmt.Errorf("bool convert error")
	}
}

// convertCommandは "cmd", 'cmd' の文字列から cmd を抽出します。
func convertCommand(word string) string {
	cmd := strings.TrimSpace(word)
	m := regexp.MustCompile("^(['\"])(.+)(['\"])$").FindStringSubmatch(cmd)
	if len(m) != 4 || m[1] != m[3] {
		log.Error("parse error cmd : ", word)
		return ""
	}
	return m[2]
}

func stringToJob(body string) *Job {
	words := strings.Split(body, ",")
	cmd := convertCommand(strings.TrimSpace(words[0]))
	if cmd == "" {
		return nil
	}
	job := Job{Cmd: cmd}
	if len(words) >= 2 {
		job.Ofile = strings.TrimSpace(words[1])
	}
	if len(words) == 4 {
		label1 := strings.TrimSpace(words[2])
		if val, err := convertInt(label1, body); err == nil {
			job.Cycle = val
		} else {
			return nil
		}
		label2 := strings.TrimSpace(words[3])
		if val, err := convertInt(label2, body); err == nil {
			job.Step = val
		} else {
			return nil
		}
	}
	return &job
}

func (config *Config) ParseConfigDir(configDir string) {
	if config.Schedule != nil {
		config.Schedule.ParseConfigDir(configDir)
	}
}

func (schedule *Schedule) ParseConfigDir(configDir string) error {
	configFiles, err := ioutil.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("parse config dir %s : %s", configDir, err)
	}
	for _, configFile := range configFiles {
		configPath := filepath.Join(configDir, configFile.Name())
		err = schedule.ParseConfigFile(configPath)
	}
	return err
}

func (config *Config) LoadLicense() error {
	schedule := config.Schedule
	// if schedule != nil {
	// 	return fmt.Errorf("load license schedule is nil")
	// }
	log.Info("load license ", config.LicenseFile)
	schedule.ParseConfigFile(config.LicenseFile)
	return nil
}

func (config *Config) ParseConfigFile(configFile string) {
	// if config.Schedule != nil {
	config.Schedule.ParseConfigFile(configFile)
	// }
}

func (schedule *Schedule) ParseConfigFile(configFile string) error {
	if !strings.HasSuffix(configFile, ".ini") &&
		!strings.HasSuffix(configFile, ".txt") {
		return nil
	}
	baseDir, err := GetParentAbsPath(configFile, 1)
	log.Debug("load ", configFile)
	if err != nil {
		return fmt.Errorf("parse config file %s : %s", configFile, err)
	}
	fp, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("parse config file %s : %s", configFile, err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		schedule.ParseConfigLine(line)
		m := regexp.MustCompile("Include\\s+(.+?)$").FindStringSubmatch(line)
		if len(m) > 0 {
			source := filepath.Join(baseDir, strings.TrimSpace(m[1]))
			if ok, _ := CheckDirectory(source); ok {
				err = schedule.ParseConfigDir(source)
			}
			if ok, _ := CheckFile(source); ok {
				err = schedule.ParseConfigFile(source)
			}
		}
	}
	return err
}

func (config *Config) ParseConfigLine(line string) {
	if config.Schedule != nil {
		config.Schedule.ParseConfigLine(line)
	}
}

func (schedule *Schedule) ParseConfigLine(line string) {
	// コメント行は読み飛ばす
	if line = paramLineTrim(line); line == "" {
		return
	}
	// ベースパラメータの解析
	if body, _ := searchParam(line, "DISK_CAPACITY", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.DiskCapacity = value
		}
	}
	if body, _ := searchParam(line, "SAVE_HOUR", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.SaveHour = value
		}
	}
	if body, _ := searchParam(line, "RECOVERY_HOUR", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.RecoveryHour = value
		}
	}
	if body, _ := searchParam(line, "MAX_ERROR_LOG", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.MaxErrorLog = value
		}
	}
	if body, _ := searchParam(line, "LOG_LEVEL", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.LogLevel = value
		}
	}
	if body, _ := searchParam(line, "DEBUG_CONSOLE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.DebugConsole = value
		}
	}
	if body, _ := searchParam(line, "LOG_SIZE", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.LogSize = value
		}
	}
	if body, _ := searchParam(line, "LOG_ROTATION", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.LogRotation = value
		}
	}
	if body, _ := searchParam(line, "LOG_LOCALIZE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.LogLocalize = value
		}
	}
	if body, _ := searchParam(line, "HANODE_ENABLE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.HanodeEnable = value
		}
	}
	if body, _ := searchParam(line, "HANODE_CMD", BASE); body != "" {
		schedule.HanodeCmd = body
	}
	if body, _ := searchParam(line, "POST_ENABLE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.PostEnable = value
		}
	}
	if body, _ := searchParam(line, "POST_CMD", BASE); body != "" {
		schedule.PostCmd = body
	}
	if body, _ := searchParam(line, "REMHOST_ENABLE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.RemhostEnable = value
		}
	}
	if body, _ := searchParam(line, "URL_CM", BASE); body != "" {
		schedule.UrlCM = body
	}
	if body, _ := searchParam(line, "URL_PM", BASE); body != "" {
		schedule.UrlPM = body
	}
	if body, _ := searchParam(line, "SOAP_TIMEOUT", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.SoapTimeout = value
		}
	}
	if body, _ := searchParam(line, "SITE_KEY", BASE); body != "" {
		schedule.SiteKey = body
	}
	if body, _ := searchParam(line, "PROXY_ENABLE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.ProxyEnable = value
		}
	}
	if body, _ := searchParam(line, "PROXY_HOST", BASE); body != "" {
		schedule.ProxyHost = body
	}
	if body, _ := searchParam(line, "PROXY_PORT", BASE); body != "" {
		if value, err := convertInt(body, line); err == nil {
			schedule.ProxyPort = value
		}
	}
	if body, _ := searchParam(line, "WEB_SERVICE_ENABLE", BASE); body != "" {
		if value, err := convertBool(body, line); err == nil {
			schedule.WebServiceEnable = value
		}
	}
	if body, _ := searchParam(line, "WEB_SERVICE_URL", BASE); body != "" {
		schedule.WebServiceUrl = body
	}
	// ライセンス用パラメータの解析
	if body, _ := searchParam(line, "HOSTNAME", BASE); body != "" {
		schedule.License.Hostname = body
	}
	if body, _ := searchParam(line, "EXPIRE", BASE); body != "" {
		schedule.License.Expired = body
	}
	if body, _ := searchParam(line, "CODE", BASE); body != "" {
		schedule.License.Code = body
	}

	// コレクター用パラメータの解析
	if body, stat := searchParam(line, "STAT_ENABLE", COLLECTOR); body != "" {
		if value, err := convertBool(body, line); err == nil {
			collector := schedule.FindOrCreateCollector(stat)
			collector.StatEnable = value
		}
	}
	if body, stat := searchParam(line, "BUILD", COLLECTOR); body != "" {
		if value, err := convertInt(body, line); err == nil {
			collector := schedule.FindOrCreateCollector(stat)
			collector.Build = value
		}
	}
	if body, stat := searchParam(line, "STAT_STDOUTLOG", COLLECTOR); body != "" {
		if value, err := convertBool(body, line); err == nil {
			collector := schedule.FindOrCreateCollector(stat)
			collector.StatStdoutLog = value
		}
	}
	if body, stat := searchParam(line, "STAT_INTERVAL", COLLECTOR); body != "" {
		if value, err := convertInt(body, line); err == nil {
			collector := schedule.FindOrCreateCollector(stat)
			collector.StatInterval = value
		}
	}
	if body, stat := searchParam(line, "STAT_TIMEOUT", COLLECTOR); body != "" {
		if value, err := convertInt(body, line); err == nil {
			collector := schedule.FindOrCreateCollector(stat)
			collector.StatTimeout = value
		}
	}
	if body, stat := searchParam(line, "STAT_MODE", COLLECTOR); body != "" {
		collector := schedule.FindOrCreateCollector(stat)
		collector.StatMode = body
	}
	// ワーカー用パラメータの解析
	if body, stat := searchParam(line, "STAT_CMD", JOB); body != "" {
		collector := schedule.FindOrCreateCollector(stat)
		if job := stringToJob(body); job != nil {
			collector.Jobs = append(collector.Jobs, job)
		} else {
			schedule.ParseFailed = true
		}
	}
}
