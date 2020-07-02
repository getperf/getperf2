package agent

import (
    "time"
)

type Schedule struct {
    DiskCapacity int /**< Disk free threshold(%) */
    SaveHour     int /**< Metric data retention(H) */
    RecoveryHour int /**< Metric data retransmission(H) */
    MaxErrorLog  int /**< Max rows of error output */

    Pid    int /**< Scheduler process id */
    Status int /**< Process status */

    LogLevel     int  /**< Log level */
    DebugConsole bool /**< Console log enabled */
    LogSize      int  /**< Log size */
    LogRotation  int  /**< Number of log rotation */
    LogLocalize  bool /**< Flag of Japanese log */

    HanodeEnable bool   /**< HA node check flag */
    HanodeCmd    string /**< HA node check script */

    PostEnable bool   /**< Post command enabled */
    PostCmd    string /**< Post command */

    RemhostEnable bool   /**< Remote transfer enabled */
    UrlCM         string /**< Web service URL (Configuration manager) */
    UrlPM         string /**< Web service URL (Performance manager) */
    SoapTimeout   int    /**< Web service timeout */
    SiteKey       string /**< Site key */

    WebServiceEnable bool   // Web service enable
    WebServiceUrl    string // Web service url

    ProxyEnable bool   /**< HTTP proxy enabled */
    ProxyHost   string /**< Proxy host */
    ProxyPort   int    /**< Proxy port */

    LastUpdate  time.Time /**< Last update of parameter file */
    ParseFailed bool      /**< Set true if parser failed */

    License    *License
    Collectors map[string]*Collector // Collector pids
}

func NewSchedule() *Schedule {
    var schedule Schedule
    schedule.DiskCapacity = DEFAULT_DISK_CAPACITY
    schedule.SaveHour = DEFAULT_SAVE_HOUR
    schedule.RecoveryHour = DEFAULT_RECOVERY_HOUR
    schedule.MaxErrorLog = DEFAULT_MAX_ERROR_LOG
    schedule.LogLevel = DEFAULT_LOG_LEVEL
    schedule.LogSize = DEFAULT_LOG_SIZE
    schedule.LogRotation = DEFAULT_LOG_ROTATION

    schedule.WebServiceEnable = true
    schedule.WebServiceUrl = DEFAULT_SERVICE_URL

    schedule.LogLocalize = true
    schedule.LastUpdate = time.Now()
    schedule.ParseFailed = false
    schedule.License = NewLicense()
    schedule.Collectors = make(map[string]*Collector)
    return &schedule
}

func (schedule *Schedule) AddCollector(collector *Collector) {
    schedule.Collectors[collector.StatName] = collector
}

func (config *Config) GetCollector(statName string) *Collector {
    if schedule := config.Schedule; schedule != nil {
        return schedule.GetCollector(statName)
    }
    return nil
}

func (schedule *Schedule) GetCollector(statName string) *Collector {
    if val, ok := schedule.Collectors[statName]; ok {
        return val
    }
    return nil
}

func (schedule *Schedule) FindOrCreateCollector(statName string) *Collector {
    if collector, ok := schedule.Collectors[statName]; ok {
        return collector
    } else {
        id := len(schedule.Collectors) + 1
        newCollector := NewCollector(statName)
        newCollector.Id = id
        schedule.AddCollector(newCollector)
        return newCollector
    }
}
